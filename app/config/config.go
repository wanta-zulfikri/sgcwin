package config

import "github.com/spf13/viper"

type Server struct {
	Port string `mapstructure:"PORT"`
}
type DatabaseConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Name     string `mapstructure:"NAME"`
}
type RedisConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

type GCPConfig struct {
	Credential string `mapstructure:"CREDEN"`
	PRJID      string `mapstructure:"PROJECTID"`
	BCKNM      string `mapstructure:"BUCKETNAME"`
	Path       string `mapstructure:"PATH"`
}

type MidtransConfig struct {
	ServerKey      string `mapstructure:"SERVERKEY"`
	ClientKey      string `mapstructure:"CLIENTKEY"`
	Env            int    `mapstructure:"ENV"`
	URLHandler     string `mapstructure:"URL"`
	ExpiryDuration int    `mapstructure:"EXP"`
	Unit           string `mapstructure:"UNIT"`
}
type NSQConfig struct {
	Host    string `mapstructure:"HOST"`
	Port    string `mapstructure:"PORT"`
	Topic   string `mapstructure:"TOPIC"`
	Topic2  string `mapstructure:"TOPIC2"`
	Topic3  string `mapstructure:"TOPIC3"`
	Topic4  string `mapstructure:"TOPIC4"`
	Topic5  string `mapstructure:"TOPIC5"`
	Topic6  string `mapstructure:"TOPIC6"`
	Topic7  string `mapstructure:"TOPIC7"`
	Topic8  string `mapstructure:"TOPIC8"`
	Topic9  string `mapstructure:"TOPIC9"`
	Topic10 string `mapstructure:"TOPIC10"`
	Topic11 string `mapstructure:"TOPIC11"`
	Topic12 string `mapstructure:"TOPIC12"`
	Topic13 string `mapstructure:"TOPIC13"`
	Topic14 string `mapstructure:"TOPIC14"`
}
type PusherConfig struct {
	AppId   string `mapstructure:"APPID"`
	Key     string `mapstructure:"KEY"`
	Secret  string `mapstructure:"SECRET"`
	Cluster string `mapstructure:"CLUSTER"`
	Secure  bool   `mapstructure:"SECURE"`
	Channel string `mapstructure:"CHANNEL"`
	Event1  string `mapstructure:"EVENT1"`
	Event2  string `mapstructure:"EVENT2"`
	Event3  string `mapstructure:"EVENT3"`
}
type Config struct {
	Server     Server         `mapstructure:"SERVER"`
	Database   DatabaseConfig `mapstructure:"DATABASE"`
	Midtrans   MidtransConfig `mapstructure:"MIDTRANS"`
	JwtSecret  string         `mapstructure:"JWTSECRET"`
	GmapsKey   string         `mapstructure:"GMAPS"`
	Redis      RedisConfig    `mapstructure:"REDIS"`
	CSRFLength int            `mapstructure:"CSRFLENGTH"`
	CSRFMode   string         `mapstructure:"CSRFMODE"`
	NSQ        NSQConfig      `mapstructure:"NSQ"`
	GCP        GCPConfig      `mapstructure:"GCP"`
	Pusher     PusherConfig   `mapstructure:"PUSHER"`
	QuizAuth   string         `mapstructure:"QUIZ"`
}

func InitConfiguration() (*Config, error) {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
