package common

import (
	"github.com/mitchellh/mapstructure"
	"github.com/urfave/cli/v2"
)

var flagsGlobalConfig *flagsConfig

func Flags() *flagsConfig {
	configLock.Lock()
	defer configLock.Unlock()

	return flagsGlobalConfig
}

type flagsConfig struct {
	Prometheus   prometheusConfig
	Alertmanager alertmanagerConfig
}

type prometheusConfig struct {
	ListenAddress           string `mapstructure:"listen-address"`
	SecretKeyFile           string `mapstructure:"secret-key-file"`
	PrometheusRulesStoreDir string `mapstructure:"prometheus-rules-store-dir"`
	PrometheusRulesSuffix   string `mapstructure:"prometheus-rules-suffix"`
	PrometheusRulesPrefix   string `mapstructure:"prometheus-rules-prefix"`
	PrometheusListenAddress string `mapstructure:"prometheus-listen-address"`
	PrometheusBasicUsername string `mapstructure:"prometheus-basic-username"`
	PrometheusBasicPassword string `mapstructure:"prometheus-basic-password"`
	PrometheusReload        bool   `mapstructure:"prometheus-reload"`
}

type alertmanagerConfig struct {
	ListenAddress             string `mapstructure:"listen-address"`
	SecretKeyFile             string `mapstructure:"secret-key-file"`
	AlertmanagerBaseDir       string `mapstructure:"alertmanager-base-dir"`
	AlertmanagerConfigName    string `mapstructure:"alertmanager-config-name"`
	AlertmanagerListenAddress string `mapstructure:"alertmanager-listen-address"`
	AlertmanagerBasicUsername string `mapstructure:"alertmanager-basic-username"`
	AlertmanagerBasicPassword string `mapstructure:"alertmanager-basic-password"`
	AlertmanagerReload        bool   `mapstructure:"alertmanager-reload"`
}

func loadFlags(c *cli.Context, output interface{}) (err error) {
	newFlags := make(map[string]interface{})
	for _, v := range c.App.Flags {
		if len(v.Names()) == 1 {
			name := v.Names()[0]
			newFlags[name] = c.Value(name)
		}
	}
	err = mapstructure.Decode(newFlags, output)
	return
}

func LoadFlagsPrometheus(c *cli.Context) (err error) {
	var cfg prometheusConfig
	if err = loadFlags(c, &cfg); err != nil {
		return
	}

	configLock.Lock()
	defer configLock.Unlock()

	flagsGlobalConfig = &flagsConfig{Prometheus: cfg}
	return
}

func LoadFlagsAlertmanager(c *cli.Context) (err error) {
	var cfg alertmanagerConfig
	if err = loadFlags(c, &cfg); err != nil {
		return
	}

	configLock.Lock()
	defer configLock.Unlock()

	flagsGlobalConfig = &flagsConfig{Alertmanager: cfg}
	return
}
