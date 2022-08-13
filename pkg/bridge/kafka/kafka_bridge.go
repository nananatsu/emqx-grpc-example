package bridge

import (
	"bytes"
	"context"
	"log"
	"regexp"
	"strconv"

	pb "emqx-grpc/api/protobuf/exhook"

	"emqx-grpc/pkg/entity"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHookProviderServer
	kafkaWriter *kafka.Writer
	topicFilter map[*regexp.Regexp]string
	multiSend   bool
}

func (s *server) OnProviderLoaded(ctx context.Context, in *pb.ProviderLoadedRequest) (*pb.LoadedResponse, error) {
	hooks := []*pb.HookSpec{
		// {Name: "client.connect"},
		// {Name: "client.connack"},
		// {Name: "client.connected"},
		// {Name: "client.disconnected"},
		// {Name: "client.authenticate"},
		// {Name: "client.check_acl"},
		// {Name: "client.subscribe"},
		// {Name: "client.unsubscribe"},
		// {Name: "session.created"},
		// {Name: "session.subscribed"},
		// {Name: "session.unsubscribed"},
		// {Name: "session.resumed"},
		// {Name: "session.discarded"},
		// {Name: "session.takeovered"},
		// {Name: "session.terminated"},
		// {Name: "message.publish"},
		{Name: "message.delivered"},
		// {Name: "message.acked"},
		{Name: "message.dropped"},
	}
	return &pb.LoadedResponse{Hooks: hooks}, nil
}

func (s *server) OnProviderUnloaded(ctx context.Context, in *pb.ProviderUnloadedRequest) (*pb.EmptySuccess, error) {
	return &pb.EmptySuccess{}, nil
}

func (s *server) OnMessagePublish(ctx context.Context, in *pb.MessagePublishRequest) (*pb.ValuedResponse, error) {

	in.Message.Payload = []byte("hello world!!!")
	reply := &pb.ValuedResponse{}
	reply.Type = pb.ValuedResponse_STOP_AND_RETURN
	reply.Value = &pb.ValuedResponse_Message{Message: in.Message}
	return reply, nil
}

func (s *server) OnMessageDelivered(ctx context.Context, in *pb.MessageDeliveredRequest) (*pb.EmptySuccess, error) {

	log.Println("收到mqtt消息", in.Message)

	var buffer bytes.Buffer

	for r, v := range s.topicFilter {
		if r.MatchString(in.Message.Topic) {
			if buffer.Len() == 0 {
				buffer.WriteString(in.Message.Topic)
				buffer.WriteString("###")
				buffer.Write(in.Message.Payload)
				buffer.WriteString("###")
				buffer.WriteString(strconv.FormatUint(in.Message.Timestamp, 10))
			}
			s.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
				Topic: v,
				Key:   []byte(in.Message.From),
				Value: buffer.Bytes(),
			})

			if !s.multiSend {
				break
			}
		}
	}

	return &pb.EmptySuccess{}, nil
}

func (s *server) OnMessageDropped(ctx context.Context, in *pb.MessageDroppedRequest) (*pb.EmptySuccess, error) {

	log.Println("丢弃mqtt消息", in.Message)
	return &pb.EmptySuccess{}, nil
}

func (s *server) Close() error {

	return s.kafkaWriter.Close()
}

func RegisterGrpc(grpcServer *grpc.Server, kafkaConfig *entity.KafkaConfig) *server {

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaConfig.Server...),
		Balancer: &kafka.ReferenceHash{},
	}

	topicFilter := make(map[*regexp.Regexp]string, 4)

	for k, v := range kafkaConfig.Topic {
		reg, err := regexp.Compile(k)

		if err != nil {
			log.Printf("mqtt主题正则 %v无法编译: %v", k, err)
			continue
		}
		topicFilter[reg] = v
	}

	s := server{kafkaWriter: kafkaWriter, topicFilter: topicFilter, multiSend: kafkaConfig.MultiSend}

	pb.RegisterHookProviderServer(grpcServer, &s)

	return &s
}
