package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

const (
	apiVersion          = "5.95"
	saveVideoEndpoint   = "https://api.vk.com/method/video.save"
	deleteVideoEndpoint = "https://api.vk.com/method/video.delete"
	wallPostEndpoint    = "https://api.vk.com/method/wall.post"
)

type vkResponse struct {
	Response struct {
		UploadURL   string `json:"upload_url"`
		VideoID     int    `json:"video_id"`
		PostID      int    `json:"post_id"`
		OwnerID     int    `json:"owner_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AccessKey   string `json:"access_key"`
	} `json:"response,omitempty"`
	Error vkError `json:"error,omitempty"`
}

type vkError struct {
	ErrorCode     int    `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params"`
}

func saveVideo(link, accessToken, groupID string) (*vkResponse, error) {

	form := url.Values{}
	form.Add("access_token", accessToken)
	form.Add("group_id", groupID)
	form.Add("link", link)
	form.Add("v", apiVersion)

	reqBody := strings.NewReader(form.Encode())

	request, err := newRequest("POST", saveVideoEndpoint, reqBody)
	if err != nil {
		return nil, err
	}

	client, err := newClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var vkResp vkResponse
	json.Unmarshal(body, &vkResp)

	if vkResp.Error.ErrorCode != 0 {
		return nil, errors.New(vkResp.Error.ErrorMsg)
	}

	request, err = newRequest("GET", vkResp.Response.UploadURL, nil)
	if err != nil {
		return nil, err
	}

	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	return &vkResp, nil
}

func deleteVideo(groupID, videoID int, accessToken string) error {

	form := url.Values{}
	form.Add("access_token", accessToken)
	form.Add("owner_id", fmt.Sprintf("%d", groupID))
	form.Add("video_id", fmt.Sprintf("%d", videoID))
	form.Add("v", apiVersion)

	reqBody := strings.NewReader(form.Encode())

	request, err := newRequest("POST", deleteVideoEndpoint, reqBody)
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var vkResp vkResponse
	json.Unmarshal(body, &vkResp)

	if vkResp.Error.ErrorCode != 0 {
		return errors.New(vkResp.Error.ErrorMsg)
	}

	return nil
}

func addPost(ownerID, videoID int, entry *rssEntry, accessToken string) error {

	form := url.Values{}
	form.Add("v", apiVersion)
	form.Add("from_group", "1")
	form.Add("access_token", accessToken)
	form.Add("owner_id", fmt.Sprintf("%d", ownerID))
	form.Add("attachments", fmt.Sprintf("video%d_%d", ownerID, videoID))
	form.Add("message", fmt.Sprintf("%s\n\n%s", entry.Group.Title, entry.Group.Description))

	reqBody := strings.NewReader(form.Encode())
	request, err := newRequest("POST", wallPostEndpoint, reqBody)
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var vkResp vkResponse
	json.Unmarshal(body, &vkResp)
	if vkResp.Error.ErrorCode != 0 {
		return errors.New(vkResp.Error.ErrorMsg)
	}

	return nil
}
