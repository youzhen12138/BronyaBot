package config

type Config struct {
	MySql  MySql  `yaml:"mysql"`
	Logger Logger `yaml:"logger"`
}
