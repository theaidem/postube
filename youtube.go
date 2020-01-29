package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/recoilme/pudge"
	"github.com/spf13/viper"
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

func parseRSSFeed(tubeName string, tubeConfig *viper.Viper) {

	channels := tubeConfig.GetStringSlice("channels")
	tubeDataPath := fmt.Sprintf("./data/%s", tubeName)
	tubeDB, err := pudge.Open(tubeDataPath, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer tubeDB.Close()

	log.Printf("(%s) Start RSS parsing... %d channels", tubeName, len(channels))
	posted := 0

	for _, id := range channels {

		isChanExists, err := tubeDB.Has(id)
		if err != nil {
			log.Println(err)
			continue
		}

		if !isChanExists {
			log.Printf("(%s) Detect new channel: %s\n", tubeName, id)
			if err := tubeDB.Set(id, true); err != nil {
				log.Println(err)
				continue
			}
		}

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

			exists, err := tubeDB.Has(entry.ID)
			if err != nil {
				log.Println("DB check record ERROR:", err.Error())
				continue
			}

			if !exists {
				if config().isSkipMode() || !isChanExists {
					tubeDB.Set(entry.ID, time.Now().Unix())
					log.Println("Stored with skip mode:", entry.Title)
				} else {
					if err = postingEntry(&entry, tubeConfig); err == nil {
						tubeDB.Set(entry.ID, time.Now().Unix())
						posted++
						log.Printf("(%s) [%d] Posted: %s\n", tubeName, posted, entry.Title)
						// TODO: post limitation?!!
					} else {
						log.Println(err)
						return
					}
				}
			}

		}
	}
	return
}

func postingEntry(entry *rssEntry, tubeConfig *viper.Viper) error {

	accessToken := tubeConfig.GetString("vk.accessToken")
	groupID := tubeConfig.GetString("vk.groupID")

	video, err := saveVideo(entry.Link.Href, accessToken, groupID)
	if err != nil {
		return errors.Wrap(err, "VK Save Video ERROR:")
	}

	time.Sleep(time.Second) // for rate limit!

	if err = addPost(video.Response.OwnerID, video.Response.VideoID, entry, accessToken); err != nil {
		if delErr := deleteVideo(video.Response.OwnerID, video.Response.VideoID, accessToken); delErr != nil {
			return errors.Wrap(delErr, "VK Delete video ERROR:")
		}
		return errors.Wrap(err, "VK Add Post ERROR:")
	}

	telegramToken := tubeConfig.GetString("telegram.token")
	telegramChannel := tubeConfig.GetString("telegram.channel")

	if err = sendMessage(entry.Link.Href, entry.Group.Title, telegramToken, telegramChannel); err != nil {
		log.Println("Telegram Add Post ERROR:", err)
		// Ignore telegram error returns...
	}

	return nil
}
