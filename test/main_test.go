package main

import (
	"bytes"
	"html/template"
	"math/rand"
	"strings"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func RandByteBytesRmndr(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return b
}

type Msg struct {
	Topic     []byte
	Payload   []byte
	Timestamp []byte
}

func genMsgByTemplate(m *Msg, msg *template.Template) []byte {

	var value bytes.Buffer
	msg.Execute(&value, &Msg{
		Topic:     RandByteBytesRmndr(8),
		Payload:   RandByteBytesRmndr(8),
		Timestamp: RandByteBytesRmndr(8),
	})
	return value.Bytes()
}

func genMsgByStringBuilder() []byte {
	var builder strings.Builder
	builder.WriteString(RandStringBytesRmndr(8))
	builder.WriteString("###")
	builder.WriteString(RandStringBytesRmndr(8))
	builder.WriteString("###")
	builder.WriteString(RandStringBytesRmndr(8))
	return []byte(builder.String())
}

func genMsgByPlus() []byte {
	return []byte(RandStringBytesRmndr(8) + "###" + RandStringBytesRmndr(8) + "###" + RandStringBytesRmndr(8))
}

func stringConcat() string {

	str := RandStringBytesRmndr(8)
	for i := 0; i < 10000; i++ {
		str += RandStringBytesRmndr(8)
	}
	return str
}

func stringBuilderConcat() string {

	var builder strings.Builder

	builder.WriteString(RandStringBytesRmndr(8))
	for i := 0; i < 10000; i++ {
		builder.WriteString(RandStringBytesRmndr(8))
	}
	return builder.String()
}

func BenchmarkTemplate(b *testing.B) {

	m := new(Msg)
	msg, _ := template.New("kafka-message").Parse("{{ .Topic }}###{{ .Payload }}##{{ .Timestamp }}")

	for n := 0; n < b.N; n++ {
		genMsgByTemplate(m, msg)
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		genMsgByStringBuilder()
	}
}

func BenchmarkPlus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		genMsgByPlus()
	}
}

func BenchmarkConcatStringPlus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stringConcat()
	}
}

func BenchmarkConcatStringBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stringBuilderConcat()
	}
}
