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
	EmbeddingDimension   int64  `json:"embedding_dimension"`
	OpenAIEndpoint       string `json:"openai_endpoint"`
	OpenAIKey            string `json:"openai_key"`
	GeminiAPIKey         string `json:"gemini_api_key"`
	MilvusAddress        string `json:"milvus_address"`
	MilvusCollectionName string `json:"milvus_collection_name"`
	Env                  string `json:"env"`
}

type S3 struct {
	Method           string `json:"method"`
	Endpoint         string `json:"endpoint"`
	ExternalEndpoint string `json:"external_endpoint"`
	AccessKey        string `json:"access_key"`
	SecretAccessKey  string `json:"secret_access_key"`
	BucketName       string `json:"bucket_name"`
	FileDuration     int    `json:"file_duration"`
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
	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = "config.json"
	}
	var data []byte

	log.Info("start init config profile: ", profile)
	bytes, err := files.ReadFile(configFile)
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
