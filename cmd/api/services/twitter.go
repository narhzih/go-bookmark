package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/models"
	"io"
	"net/http"
	"os"
)

func (s Services) ExpandTweet(tweetId string) (models.TwitterExpandedResponse, error) {
	var expandedResponse models.TwitterExpandedResponse
	reqParams := "tweet.fields=author_id,conversation_id"
	reqUrl := fmt.Sprintf("https://api.twitter.com/2/tweets?ids=%v&%v", tweetId, reqParams)
	twitterReq, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return expandedResponse, err
	}
	twitterReq.Header.Set("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))
	twitterRes, err := http.DefaultClient.Do(twitterReq)
	if err != nil {
		return expandedResponse, err
	}

	if twitterRes.StatusCode != http.StatusOK {
		return expandedResponse, errors.New("an error occurred while connecting to twitter api")
	}

	twitterResBody, err := io.ReadAll(twitterRes.Body)
	if err != nil {
		return expandedResponse, err
	}

	json.Unmarshal(twitterResBody, &expandedResponse)
	return expandedResponse, nil
}

func (s Services) GetThreadByConversationID(conversationID string) (string, error) {
	reqParams := "tweet.fields=author_id"
	reqUrl := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=conversation_id:%v&%v", conversationID, reqParams)

	// build request
	twitterReq, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return "", err
	}
	twitterReq.Header.Set("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))
	twitterRes, err := http.DefaultClient.Do(twitterReq)
	if err != nil {
		return "", err
	}

	if twitterRes.StatusCode != http.StatusOK {
		return "", errors.New("http error while connecting to twitter")
	}
	reqBody, _ := io.ReadAll(twitterRes.Body)
	return string(reqBody), nil

}
