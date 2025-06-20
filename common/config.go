package common

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/toolkits/file"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	config     *globalConfig
	configLock = new(sync.RWMutex)
)

func Config() *globalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

type Provider interface {
	Dsn() string
	TableNamePrefix() string
	success(cfg string)
}

type globalConfig struct {
	configType string
	Yaml       yamlConfig
	ini        iniConfig
	Global     Provider
}

func LoadConfig(cfg, configType string) (err error) {
	c := &globalConfig{configType: configType}
	switch configType {
	case "ini":
		err = c.loadEnvWithFile(cfg, &c.ini)
		c.Global = &c.ini
	default:
		err = c.loadFile(cfg, &c.Yaml)
		c.Global = &c.Yaml
	}
	if err != nil {
		log.Fatalln(err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = c
	c.Global.success(cfg)

	return
}

func (c *globalConfig) loadEnvWithFile(cfg string, rawVal interface{}) (err error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile != "" {
		err = c.loadFile(configFile, rawVal)
	} else {
		err = c.loadFile(cfg, rawVal)
	}
	return
}

func (c *globalConfig) loadFile(cfg string, rawVal interface{}) (err error) {
	if cfg == "" {
		err = fmt.Errorf("use --config to specify configuration file")
		return
	}
	if !file.IsExist(cfg) {
		err = fmt.Errorf("config file: %s is not existent", cfg)
		return
	}
	v := viper.New()
	v.AddConfigPath(filepath.Dir(cfg))
	v.SetConfigName(strings.TrimSuffix(filepath.Base(cfg), file.Ext(cfg)))
	v.SetConfigType(c.configType)

	if err = v.ReadInConfig(); err != nil {
		err = fmt.Errorf("read config file fail: %s", err.Error())
		return
	}

	if err = v.Unmarshal(rawVal); err != nil {
		err = fmt.Errorf("parse config file fail: %s", err.Error())
	}
	return
}

type dataBase struct {
	Username    string `mapstructure:"username"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	TablePrefix string `mapstructure:"table_prefix"`
}

type yamlConfig struct {
	SecretKey       string `mapstructure:"secret_key" json:"secret_key"`
	OauthEncryptKey string `mapstructure:"oauth_encrypt_key" yaml:"oauth_encrypt_key"`
	DbHost          string `mapstructure:"db_host" yaml:"db_host"`
	DbUsername      string `mapstructure:"db_username" yaml:"db_username"`
	DbPassword      string `mapstructure:"db_password" yaml:"db_password"`
	DbPort          int    `mapstructure:"db_port" yaml:"db_port"`
	DbDatabase      string `mapstructure:"db_database" yaml:"db_database"`
	TablePrefix     string `mapstructure:"table_prefix" yaml:"table_prefix"`
}

func (c *yamlConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbDatabase)
}

func (c *yamlConfig) TableNamePrefix() string {
	return c.TablePrefix
}

func (c *yamlConfig) success(cfg string) {
	log.Println("read config file:", cfg, "successfully")
}

type iniConfig struct {
	Mysql *dataBase `json:"mysql"`
}

func (c *iniConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.Username, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database)
}

func (c *iniConfig) TableNamePrefix() string {
	return c.Mysql.TablePrefix
}

func (c *iniConfig) success(cfg string) {
}
