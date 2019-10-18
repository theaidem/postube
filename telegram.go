package main

import (
	"fmt"
)

const (
	sendMessageEndpoint = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=markdown"
)

func sendMessage(link, title, telegramToken, telegramChannel string) error {
	text := fmt.Sprintf("*%s* [%s]", title, link)
	endpoint := fmt.Sprintf(sendMessageEndpoint, telegramToken, telegramChannel, text)

	request, err := newRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	return nil
}
