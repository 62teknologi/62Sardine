package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var LoadedConfig map[string]any

type Config struct {
	HTTPServerAddress        string `mapstructure:"HTTP_SERVER_ADDRESS"`
	FILESYSTEM_DISK          string `mapstructure:"FILESYSTEM_DISK"`
	FILESYSTEM_FOLDER        string `mapstructure:"FILESYSTEM_FOLDER"`
	APP_URL                  string `mapstructure:"APP_URL"`
	AWS_ACCESS_KEY_ID        string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_ACCESS_KEY_SECRET    string `mapstructure:"AWS_ACCESS_KEY_SECRET"`
	AWS_REGION               string `mapstructure:"AWS_REGION"`
	AWS_BUCKET               string `mapstructure:"AWS_BUCKET"`
	AWS_URL                  string `mapstructure:"AWS_URL"`
	AWS_ENDPOINT             string `mapstructure:"AWS_ENDPOINT"`
	DO_ACCESS_KEY_ID         string `mapstructure:"DO_ACCESS_KEY_ID"`
	DO_ACCESS_KEY_SECRET     string `mapstructure:"DO_ACCESS_KEY_SECRET"`
	DO_REGION                string `mapstructure:"DO_REGION"`
	DO_BUCKET                string `mapstructure:"DO_BUCKET"`
	DO_URL                   string `mapstructure:"DO_URL"`
	DO_ENDPOINT              string `mapstructure:"DO_ENDPOINT"`
	ALIYUN_ACCESS_KEY_ID     string `mapstructure:"ALIYUN_ACCESS_KEY_ID"`
	ALIYUN_ACCESS_KEY_SECRET string `mapstructure:"ALIYUN_ACCESS_KEY_SECRET"`
	ALIYUN_BUCKET            string `mapstructure:"ALIYUN_BUCKET"`
	ALIYUN_URL               string `mapstructure:"ALIYUN_URL"`
	ALIYUN_ENDPOINT          string `mapstructure:"ALIYUN_ENDPOINT"`
	MINIO_ACCESS_KEY_ID      string `mapstructure:"MINIO_ACCESS_KEY_ID"`
	MINIO_ACCESS_KEY_SECRET  string `mapstructure:"MINIO_ACCESS_KEY_SECRET"`
	MINIO_REGION             string `mapstructure:"MINIO_REGION"`
	MINIO_BUCKET             string `mapstructure:"MINIO_BUCKET"`
	MINIO_URL                string `mapstructure:"MINIO_URL"`
	MINIO_ENDPOINT           string `mapstructure:"MINIO_ENDPOINT"`
	MINIO_SSL                string `mapstructure:"MINIO_SSL"`
	GOOGLE_GCS_PATH          string `mapstructure:"GOOGLE_GCS_PATH"`
	GOOGLE_GCS_BUCKET        string `mapstructure:"GOOGLE_GCS_BUCKET"`
	GOOGLE_GCS_URL           string `mapstructure:"GOOGLE_GCS_URL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("HTTP_SERVER_ADDRESS", "0.0.0.0:10082")
	viper.SetDefault("FILESYSTEM_DISK", "s3")
	viper.SetDefault("MINIO_SSL", false)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func ReadConfig(name string) (string, error) {
	if LoadedConfig == nil {
		config, _ := LoadConfig()
		LoadedConfig = map[string]any{
			"app": map[string]any{
				"name": "Saldine",
			},
			"filesystems": map[string]any{

				"default":        config.FILESYSTEM_DISK,
				"default_folder": config.FILESYSTEM_FOLDER,

				// Filesystem Disks

				"disks": map[string]any{
					"local": map[string]any{
						"driver": "local",
						"root":   "storage/app",
						"url":    config.APP_URL + "/storage",
					},
					"s3": map[string]any{
						"driver":   "s3",
						"key":      config.AWS_ACCESS_KEY_ID,
						"secret":   config.AWS_ACCESS_KEY_SECRET,
						"region":   config.AWS_REGION,
						"bucket":   config.AWS_BUCKET,
						"url":      config.AWS_URL,
						"endpoint": config.AWS_ENDPOINT,
					},
					"do": map[string]any{
						"driver":   "s3",
						"key":      config.DO_ACCESS_KEY_ID,
						"secret":   config.DO_ACCESS_KEY_SECRET,
						"region":   config.DO_REGION,
						"bucket":   config.DO_BUCKET,
						"url":      config.DO_URL,
						"endpoint": config.DO_ENDPOINT,
					},
					"gcs": map[string]any{
						"driver": "gcs",
						"path":   config.GOOGLE_GCS_PATH,
						"bucket": config.GOOGLE_GCS_BUCKET,
						"url":    config.GOOGLE_GCS_URL,
					},
					"oss": map[string]any{
						"driver":   "oss",
						"key":      config.ALIYUN_ACCESS_KEY_ID,
						"secret":   config.ALIYUN_ACCESS_KEY_SECRET,
						"bucket":   config.ALIYUN_BUCKET,
						"url":      config.ALIYUN_URL,
						"endpoint": config.ALIYUN_ENDPOINT,
					},
					// "cos": map[string]any{
					// 	"driver": "cos",
					// 	"key":    config.Env("TENCENT_ACCESS_KEY_ID"),
					// 	"secret": config.Env("TENCENT_ACCESS_KEY_SECRET"),
					// 	"bucket": config.Env("TENCENT_BUCKET"),
					// 	"url":    config.Env("TENCENT_URL"),
					// },
					"minio": map[string]any{
						"driver":   "minio",
						"key":      config.MINIO_ACCESS_KEY_ID,
						"secret":   config.MINIO_ACCESS_KEY_SECRET,
						"region":   config.MINIO_REGION,
						"bucket":   config.MINIO_BUCKET,
						"url":      config.MINIO_URL,
						"endpoint": config.MINIO_ENDPOINT,
						"ssl":      config.MINIO_SSL,
					},
				},
			},
		}
	}

	keys := strings.Split(name, ".")
	var val interface{}
	var err error
	for _, k := range keys {
		switch v := val.(type) {
		case nil:
			val, err = GetConfigValue(LoadedConfig, k)
		case map[string]interface{}:
			val, err = GetConfigValue(v, k)
		default:
			err = fmt.Errorf("invalid key %s", name)
		}
		if err != nil {
			break
		}
	}
	val_key := fmt.Sprintf("%v", val)

	if val == nil {
		val_key = ""
	}

	return val_key, err

}

func GetConfigValue(data map[string]interface{}, key string) (interface{}, error) {
	val, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}
