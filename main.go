package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

func main() {
	run()
	for {
		select {
		case <-time.Tick(config().parserInterval()):
			run() // boom
		}
	}
}

func run() {
	if config().parserIsEnabled() {
		tubes := config().tubes()
		log.Printf("Start processing for %d tubes..", len(tubes))
		for _, tubeName := range shuffleSlice(tubes) {

			tubeConfig := viper.New()
			configName := fmt.Sprintf("%s.config", tubeName)
			tubeConfig.AddConfigPath(configPath)
			tubeConfig.SetConfigName(configName)

			// Check config file
			if err := tubeConfig.ReadInConfig(); err != nil {
				log.Println(err)
				continue
			}

			parseRSSFeed(tubeName, tubeConfig)
		}
		log.Println("**** **** **** ****")
	}
}
