package entity

type Config struct {
	Grpc  *GrpcConfig  `yaml:"grpc"`
	Emqx  *EmqxConfig  `yaml:"emqx"`
	Kafka *KafkaConfig `yaml:"kafka"`
}

type GrpcConfig struct {
	BindAddress string `yaml:"bindAddress"`
}

type EmqxConfig struct {
	AdapterServer string `yaml:"adapterServer"`
}

type KafkaConfig struct {
	Server      []string          `yaml:"server"`
	MsgTemplate string            `yaml:"msgTemplate"`
	Topic       map[string]string `yaml:"topic"`
	MultiSend   bool              `yaml:"multiSend"`
}
