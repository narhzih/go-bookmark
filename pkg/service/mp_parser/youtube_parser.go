package mp_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func ParseYoutubeLink(youtubeLink string) (string, error) {
	timeout := time.Duration(10 * time.Second)
	requestBody, err := json.Marshal(map[string]interface{}{})
	client := http.Client{Timeout: timeout}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://www.youtube.com/oembed?url=%v&format=json", youtubeLink), bytes.NewBuffer(requestBody))
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
