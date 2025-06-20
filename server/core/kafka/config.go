package kafka

type Config struct {
	Brokers          string `yaml:"brokers"`
	GroupId          string `yaml:"group-id"`
	TopicPrefix      string `yaml:"topic-prefix"`
	SecurityProtocol string `yaml:"security-protocol"`
	SaslMechanism    string `yaml:"sasl-mechanism"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	SchemaRegistry   string `yaml:"schema-registry"`
}
