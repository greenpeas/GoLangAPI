package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Psql struct {
			Dsn           string `yaml:"dsn"`
			MigrationPath string `yaml:"migration_path""`
		} `yaml:"psql"`
		Redis struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"database"`
	Salt string `yaml:"salt"`
	Jwt  struct {
		Secret     string `yaml:"secret"`
		Ttl        int64  `yaml:"ttl"`
		TtlRefresh int64  `yaml:"ttl_refresh"`
	} `yaml:"jwt"`
	Logger struct {
		FilePath   string `yaml:"file_path"`
		Level      string `yaml:"level"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Compress   bool   `yaml:"compress"`
	} `yuml:"logger"`
	Router struct {
		UseCorsMiddleware bool `yaml:"use_cors_middleware"`
		Port              int  `yaml:"port"`
	} `yaml:"router"`
	GRPCCommands struct {
		Addr    string `yaml:"addr"`
		Timeout uint8  `yaml:"timeout"`
	} `yaml:"grpc_commands"`
	Test struct {
		User struct {
			Login    string `yaml:"login"`
			Password string `yaml:"password"`
		}
	} `json:"test"`
	ShutdownTimeout   string `yaml:"shutdown_timeout"`
	ShippingFilesPath string `yaml:"shipping_files_path"`
	RunCheckTelemetry bool   `yaml:"run_check_telemetry"`
}

var instance *Config
var once sync.Once

func Get(path string) *Config {
	once.Do(func() {
		instance = &Config{}

		configFile, err := os.Open(path)

		if err != nil {
			log.Fatalln("Error read config file")
		}

		defer configFile.Close()

		decoder := yaml.NewDecoder(configFile)

		if err = decoder.Decode(instance); err != nil {
			log.Fatalln("Error parse config file")
		}

		setDefaults()

		log.Println("Config loaded: " + configFile.Name())
	})

	return instance
}

func setDefaults() {
	if instance.ShutdownTimeout == "" {
		instance.ShutdownTimeout = "5s"
	}

	if instance.Router.Port == 0 {
		instance.Router.Port = 8080
	}
}
