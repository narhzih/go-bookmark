package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type MpParser struct {
	InstagramParser string
	TwitterParser   string
	YoutubeParser   string
}

func ParseLink(link string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{})
	resp, err := http.Post(
		fmt.Sprintf("https://graph.facebook.com/v12.0/?scrape=true&id=%v&access_token=%v", url.QueryEscape(link), os.Getenv("FACEBOOK_ACCESS_TOKEN")),
		"application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
		logger.Err(err).Msg("Error happened when reading stream")
		return "", err
	}
	responseToString := string(respBody)
	return responseToString, nil
}
