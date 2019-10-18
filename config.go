package main

import (
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const configPath = "./config"

type cfg struct {
	config *viper.Viper
}

var (
	cfgOnce     sync.Once
	cfgInstance cfg
)

func config() *cfg {
	cfgOnce.Do(func() {
		log.Println("Init Configuration")
		cfgInstance = cfg{
			config: viper.New(),
		}
		cfgInstance.config.AddConfigPath(configPath)
		cfgInstance.config.SetConfigName("config")
		cfgInstance.config.WatchConfig()
		cfgInstance.config.OnConfigChange(cfgInstance.onConfigChange)
		if err := cfgInstance.config.ReadInConfig(); err != nil {
			log.Panic(err)
		}
	})

	return &cfgInstance
}

func (c *cfg) onConfigChange(e fsnotify.Event) {
	log.Println("Config file changed:", e.Name)
}

func (c *cfg) proxyEnabled() bool {
	return c.config.GetBool("proxy.enabled")
}

func (c *cfg) proxyURL() string {
	return c.config.GetString("proxy.url")
}

func (c *cfg) parserIsEnabled() bool {
	return c.config.GetBool("parser.enabled")
}

func (c *cfg) isSkipMode() bool {
	return c.config.GetBool("parser.skipMode")
}

func (c *cfg) parserInterval() time.Duration {
	return time.Second * c.config.GetDuration("parser.interval")
}

func (c *cfg) tubes() []string {
	return c.config.GetStringSlice("tubes")
}
