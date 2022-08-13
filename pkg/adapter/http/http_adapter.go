package adapter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	pb "emqx-grpc/api/protobuf/exproto"

	"emqx-grpc/pkg/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedConnectionHandlerServer
	client pb.ConnectionAdapterClient
	conn   *grpc.ClientConn
}

func (*server) OnSocketCreated(in pb.ConnectionHandler_OnSocketCreatedServer) error {

	req, err := in.Recv()
	if err != nil {
		log.Println("处理socket连接失败", err)
	}
	log.Println("新连接", req.Conninfo)
	return nil
}

func (*server) OnSocketClosed(in pb.ConnectionHandler_OnSocketClosedServer) error {
	req, err := in.Recv()
	if err != nil {
		log.Println("处理socket关闭失败", err)
	}
	log.Println("连接关闭", req)
	return nil
}

func (s *server) OnReceivedBytes(in pb.ConnectionHandler_OnReceivedBytesServer) error {

	req, err := in.Recv()
	if err != nil {
		log.Println("接受socket数据失败", err)
	}

	reader := BytesReader{content: req.Bytes}
	bufReader := bufio.NewReader(&reader)
	httpReq, err := http.ReadRequest(bufReader)
	if err != nil {
		log.Println("解析http数据失败", err)
	}

	s.forwardRequest(httpReq, req.Conn)

	respBody, err := json.Marshal(entity.GBTResult{Success: 1})

	if err != nil {
		log.Println("json序列化http响应内容", err)
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")

	resp := http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         httpReq.Proto,
		ProtoMajor:    httpReq.ProtoMajor,
		ProtoMinor:    httpReq.ProtoMinor,
		Request:       httpReq,
		ContentLength: int64(len(respBody)),
		Header:        header,
		Body:          &BytesReader{content: respBody},
	}

	writer := BytesWriter{content: *bytes.NewBuffer(make([]byte, 0))}
	resp.Write(&writer)

	s.client.Send(context.Background(), &pb.SendBytesRequest{Conn: req.Conn, Bytes: writer.content.Bytes()})

	return nil
}

func (s *server) forwardRequest(httpReq *http.Request, conn string) {

	body, err := ioutil.ReadAll(httpReq.Body)

	if err != nil {
		log.Println("获取http body失败", err)
	}

	if httpReq.URL.Path == "/status_info" {
		var statusInfo entity.GBTStatusInfo
		json.Unmarshal(body, &statusInfo)

		_, err := s.client.Authenticate(context.Background(), &pb.AuthenticateRequest{
			Conn: conn,
			Clientinfo: &pb.ClientInfo{
				ProtoName: "mqtt",
				ProtoVer:  "3.1.1",
				Username:  "test",
			},
			Password: "test",
		})

		if err != nil {
			log.Panicln("登入mqtt服务失败", err)
		}

		deviceState, err := json.Marshal(&entity.DeviceState{
			Id: statusInfo.Serail,
			State: []entity.State{
				{
					Tag:   "RSSI",
					Value: strconv.Itoa(statusInfo.SignalStrength),
				},
				{
					Tag:   "Battery",
					Value: strconv.Itoa(statusInfo.BatteryPercetage),
				},
				{
					Tag:   "Emg",
					Value: strconv.Itoa(statusInfo.Status),
				},
			},
		})

		if err != nil {
			log.Println("序列化mqtt消息失败", err)
		}

		_, err = s.client.Publish(context.Background(), &pb.PublishRequest{
			Conn:    conn,
			Topic:   "http/thing/event/property/sync/post",
			Qos:     2,
			Payload: deviceState,
		})

		if err != nil {
			log.Println("发送mqtt消息失败", err)
		}
	} else if httpReq.URL.Path == "/install_info" {
		var installInfo entity.GBTInstallInfo
		json.Unmarshal(body, &installInfo)
	}

}

func (*server) OnTimerTimeout(in pb.ConnectionHandler_OnTimerTimeoutServer) error {
	log.Println("定时器超时处理")
	return nil
}

func (*server) OnReceivedMessages(in pb.ConnectionHandler_OnReceivedMessagesServer) error {
	log.Println("设备订阅消息转发")
	return nil
}

func (s *server) Close() error {
	return s.conn.Close()
}

func RegisterGrpc(grpcServer *grpc.Server, emqxConfig *entity.EmqxConfig) *server {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(emqxConfig.AdapterServer, opts...)
	if err != nil {
		log.Fatalf("exproto ConnectionAdapter服务器连接失败: %v", err)
	}

	log.Println("exproto ConnectionAdapter服务器连接成功", emqxConfig.AdapterServer)

	s := server{client: pb.NewConnectionAdapterClient(conn), conn: conn}

	pb.RegisterConnectionHandlerServer(grpcServer, &s)

	return &s
}
