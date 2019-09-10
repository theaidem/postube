package main

import (
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type cfg struct{}

var (
	cfgOnce     sync.Once
	cfgInstance cfg
)

func config() *cfg {
	cfgOnce.Do(func() {
		log.Println("Init Configuration")
		cfgInstance = cfg{}
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.WatchConfig()
		viper.OnConfigChange(cfgInstance.onConfigChange)
		if err := viper.ReadInConfig(); err != nil {
			log.Panic(err)
		}
	})

	return &cfgInstance
}

func (c *cfg) onConfigChange(e fsnotify.Event) {
	log.Println("Config file changed:", e.Name)
}

func (c *cfg) proxyEnabled() bool {
	return viper.GetBool("proxy.enabled")
}

func (c *cfg) proxyURL() string {
	return viper.GetString("proxy.url")
}

func (c *cfg) parserIsEnabled() bool {
	return viper.GetBool("parser.enabled")
}

func (c *cfg) isSkipMode() bool {
	return viper.GetBool("parser.skipMode")
}

func (c *cfg) parserInterval() time.Duration {
	return time.Second * viper.GetDuration("parser.interval")
}

func (c *cfg) channels() []string {
	return viper.GetStringSlice("channels")
}

func (c *cfg) vkAccessToken() string {
	return viper.GetString("vk.accessToken")
}

func (c *cfg) vkGroupID() string {
	return viper.GetString("vk.groupID")
}

func (c *cfg) telegramToken() string {
	return viper.GetString("telegram.token")
}

func (c *cfg) telegramChannel() string {
	return viper.GetString("telegram.channel")
}
