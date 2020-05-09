package api

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port      int               `mapstructure:"port"`
	Schema    map[string]string `mapstructure:"schema"`
	Cassandra struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Keyspace string `mapstructure:"keyspace"`
	} `mapstructure:"cassandra"`
	Timeout int `mapstructure:"timeout"`
}

func (c Config) SchemaPath(key string) (string, error) {
	if filename, ok := c.Schema[key]; ok {
		path, err := filepath.Abs(filename)
		if err != nil {
			return "", fmt.Errorf("can't find schema %s path", filename)
		}

		return path, nil
	}

	return "", fmt.Errorf("schema name %s doesn't exist", key)
}

func CreateConfig() (Config, error) {
	var config Config
	absPath, err := filepath.Abs("")
	if nil != err {
		return config, err
	}

	v := viper.New()
	v.AddConfigPath(absPath)
	v.SetConfigName("config")
	if err = v.ReadInConfig(); nil != err {
		return config, err
	}

	sv := v.Sub("api")
	sv.AutomaticEnv()
	sv.SetEnvPrefix("PROJECTX")
	sv.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if e := sv.Unmarshal(&config); nil != e {
		return config, err
	}

	return config, nil
}
