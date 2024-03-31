package configuration

type ENV string

const (
	Development ENV = "development"
	Production  ENV = "production"
)

type Config struct {
	config AppConfig
	env    ENV
}

type AppConfig struct {
	Auth  AuthConfig  `yaml:"auth"`
	Chart ChartConfig `yaml:"chart"`
	DB    DBConfig    `yaml:"db"`
	MQTT  MQTTConfig  `yaml:"mqtt"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
}

type ChartConfig struct {
	SoilMoisture SensorValueConfig `yaml:"soil_moisture"`
}

type SensorValueConfig struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
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
