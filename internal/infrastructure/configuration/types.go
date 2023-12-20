package configuration

type Config struct {
	config AppConfig
}

type AppConfig struct {
	DB   DBConfig   `yaml:"db"`
	MQTT MQTTConfig `yaml:"mqtt"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type MQTTConfig struct {
	Broker   string `yaml:"broker"`
	ClientID string `yaml:"client_id"`
}
