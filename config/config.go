package config

import (
	"context"
	"embed"
	"encoding/json"
	"os"
	"sync"
	"tts-poc-service/lib/baselogger"
	pkgUtil "tts-poc-service/pkg/common/utils"
)

var Config *Cfg

type Cfg struct {
	Server   Server   `json:"server"`
	General  General  `json:"general"`
	Storage  S3       `json:"storage"`
	Database Database `json:"database"`

	mutex sync.RWMutex
}

type Server struct {
	Port           int              `json:"port"`
	ReadTimeout    pkgUtil.Duration `json:"read_timeout"`
	WriteTimeout   pkgUtil.Duration `json:"write_timeout"`
	MaxHeaderBytes int              `json:"max_header_bytes"`
}

type General struct {
	Env string `json:"env"`
}

type S3 struct {
	Method          string `json:"method"`
	Endpoint        string `json:"endpoint"`
	AccessKey       string `json:"access_key"`
	SecretAccessKey string `json:"secret_access_key"`
	BucketName      string `json:"bucket_name"`
	FileDuration    int    `json:"file_duration"`
}

type Database struct {
	Host        string           `json:"host"`
	DbName      string           `json:"db_name"`
	User        string           `json:"user"`
	Password    string           `json:"password"`
	Port        string           `json:"port"`
	MaxConn     int              `json:"max_connection"`
	MaxIdle     int              `json:"max_idle"`
	MaxLifetime pkgUtil.Duration `json:"max_lifetime"`
	MaxIdletime pkgUtil.Duration `json:"max_idletime"`
}

//go:embed *
var files embed.FS

func InitConfig(ctx context.Context, log *baselogger.Logger) {
	profile := os.Getenv("APP_ENV")
	var data []byte

	log.Info("start init config profile: ", profile)
	bytes, err := files.ReadFile("config.json")
	if err != nil {
		log.Fatal("error when load config: ", err)
	}
	data = bytes

	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal("error when unmarshal config: ", err)
	}

	Config.General.Env = os.Getenv("APP_ENV")
}

func (c *Cfg) Reload(ctx context.Context, log *baselogger.Logger) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	InitConfig(ctx, log)
}
