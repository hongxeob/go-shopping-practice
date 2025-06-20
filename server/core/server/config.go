package server

type Config struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Reflection bool   `yaml:"reflection"` // Enable or disable reflection for the server

}
