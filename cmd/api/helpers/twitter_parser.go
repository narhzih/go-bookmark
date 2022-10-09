package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func ParseTwitterLink(twitterLink string) (string, error) {
	chatId := getChatId(twitterLink)
	timeout := time.Duration(10 * time.Second)
	requestBody, err := json.Marshal(map[string]interface{}{})
	client := http.Client{Timeout: timeout}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitter.com/1.1/statuses/show.json?id=%v&tweet_mode=extended", chatId), bytes.NewBuffer(requestBody))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("BEARER_TOKEN")))
	resp, err := client.Do(request)
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
		return "", err
	}
	responseToString := string(respBody)
	return responseToString, nil
}

func getChatId(twitterLink string) interface{} {
	linkSlice := strings.Split(twitterLink, "/")
	chatId := linkSlice[len(linkSlice)-1]
	return chatId
}
