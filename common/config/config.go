package config

import (
	"sync"
	"time"

	"github.com/jinzhu/configor"
)

var (
	Cfg Configuration
	mu  sync.RWMutex
)

type (
	AppConfig struct {
		Secret  string `json:"secret" default:"secret"`
		Env     string `json:"env" default:"dev"`
		ChainId int    `json:"chain_id" default:"56"`
	}

	HttpConfig struct {
		ListenAddr         string        `json:"listen_addr"`
		LimitConnection    int           `json:"limit_connection"`
		ReadTimeout        time.Duration `json:"read_timeout"`
		WriteTimeout       time.Duration `json:"write_timeout"`
		IdleTimeout        time.Duration `json:"idle_timeout"`
		MaxHeaderBytes     int           `json:"max_header_bytes"`
		MaxMultipartMemory int64         `json:"max_multipart_memory"`
	}

	LoggerConfig struct {
		Level        string        `json:"level"`
		Write        bool          `json:"write"`
		Path         string        `json:"path"`
		FileName     string        `json:"file_name"`
		MaxAge       time.Duration `json:"max_age"`
		RotationTime time.Duration `json:"rotation_time"`
	}

	MysqlConfig struct {
		Driver   string `json:"driver"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
	}

	S3Config struct {
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
		Bucket    string `json:"bucket"`
		BaseUrl   string `json:"base_url"`
	}

	RedisConfig struct {
		Host      string `json:"host" env:"REDIS_HOST"`
		Port      int    `json:"port" env:"REDIS_PORT"`
		Auth      string `json:"auth" env:"REDIS_AUTH"`
		MaxIdle   int    `json:"max_idle" env:"REDIS_MAX_IDLE"`
		MaxActive int    `json:"max_active" env:"REDIS_MAX_ACTIVE"`
		Db        int    `json:"db" env:"REDIS_DB"`
	}

	ElasticConfig struct {
		Url      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RabbitMqConfig struct {
		Username    string `json:"userName"`
		Password    string `json:"password"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		VirtualHost string `json:"virtualHost"`
	}
	RemoteApi struct {
		SsoApi    string `json:"sso_api"`
		SsoToken  string `json:"sso_token"`
		NativeWeb string `json:"native_web"`
	}

	UnisatConfig struct {
		ApiKey      string `json:"apikey"`
		StartHeight int    `json:"start_height"`
	}
	WalletAddress struct {
		BTC   string `json:"btc"`
		BRC20 string `json:"brc20"`
	}

	Configuration struct {
		App           AppConfig      `json:"app"`
		Http          HttpConfig     `json:"server"`
		Mysql         MysqlConfig    `json:"mysql"`
		MysqlAsset    MysqlConfig    `json:"mysql_asset"`
		BscPsql       MysqlConfig    `json:"bsc_psql"`
		RabbitMq      RabbitMqConfig `json:"rabbit_mq"`
		Logger        LoggerConfig   `json:"logger"`
		S3            S3Config       `json:"s3"`
		Redis         RedisConfig    `json:"redis"`
		Elastic       ElasticConfig  `json:"elastic"`
		RemoteApi     RemoteApi      `json:"remote_api"`
		Unisat        UnisatConfig   `json:"unisat"`
		WalletAddress WalletAddress  `json:"wallet_address"`
	}
)

func Load(file *string) (Configuration, error) {
	mu.Lock()
	defer mu.Unlock()

	err := configor.Load(&Cfg, *file)
	if err != nil {
		return Configuration{}, err
	}
	return Cfg, err
}

func GetConfig() Configuration {
	mu.Lock()
	defer mu.Unlock()
	return Cfg
}
