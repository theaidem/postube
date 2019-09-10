package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	rssURL = "https://www.youtube.com/feeds/videos.xml"
)

type rssFeed struct {
	Title     string `xml:"title"`
	ID        string `xml:"id"`
	ChannelID string `xml:"channelId"`
	Published string `xml:"published"`
	Author    struct {
		Name string `xml:"name"`
		URI  string `xml:"uri"`
	} `xml:"author"`
	Entries []rssEntry `xml:"entry"`
}

type rssEntry struct {
	ID    string `xml:"id"`
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Author struct {
		Name string `xml:"name"`
		URI  string `xml:"uri"`
	} `xml:"author"`
	Group struct {
		Title       string `xml:"title"`
		Content     string `xml:"content"`
		Thumbnail   string `xml:"thumbnail"`
		Description string `xml:"description"`
	} `xml:"group"`
}

func parseRSSFeed() {
	for _, id := range config().channels() {

		endpoint := fmt.Sprintf("%s?channel_id=%s", rssURL, id)
		request, err := newRequest("GET", endpoint, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		client, err := newClient()
		if err != nil {
			log.Println(err)
			continue
		}

		response, err := client.Do(request)
		if err != nil {
			log.Println(fmt.Sprintf("Read RSS feed ERROR: channel ID: %s", id), err.Error())
			continue
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			err := errors.New(fmt.Sprintf("Read RSS feed ERROR: Unexpected status code %d, channel ID: %s", response.StatusCode, id))
			log.Println(err)
			continue
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("Read RSS body ERROR:", err.Error())
			continue
		}

		var feed rssFeed
		xml.Unmarshal(body, &feed)

		for i := range feed.Entries {
			entry := feed.Entries[len(feed.Entries)-1-i]

			exists, err := db().has(entry.ID)
			if err != nil {
				log.Println("DB check record ERROR:", err.Error())
				continue
			}

			if !exists {
				if config().isSkipMode() {
					db().set(entry.ID, time.Now().Unix())
					log.Println("Stored with skip mode:", entry.Title)
				} else {
					if err = postingEntry(&entry); err == nil {
						db().set(entry.ID, time.Now().Unix())
						log.Println("Posted", entry.Title)
					} else {
						log.Println(err)
					}
				}
			}

		}

	}
}

func postingEntry(entry *rssEntry) error {

	video, err := saveVideo(entry.Link.Href)
	if err != nil {
		return errors.Wrap(err, "VK Save Video ERROR:")
	}

	time.Sleep(time.Second) // for rate limit!

	if err = addPost(video.Response.OwnerID, video.Response.VideoID, entry); err != nil {
		return errors.Wrap(err, "VK Add Post ERROR:")
	}

	if err = sendMessage(entry.Link.Href, entry.Group.Title); err != nil {
		return errors.Wrap(err, "Telegram Add Post ERROR:")
	}

	return nil
}
