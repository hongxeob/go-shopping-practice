package endpoint

type Config struct {
	targetEnv string `yaml:"target-env"`
}

func (c *Config) IsProd() bool {
	return IsProd(c.targetEnv)
}

func IsProd(targetEnv string) bool {
	return targetEnv == "prod"
}
