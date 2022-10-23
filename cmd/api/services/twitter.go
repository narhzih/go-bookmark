package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mypipeapp/mypipeapi/cmd/api/models/response"
	"github.com/mypipeapp/mypipeapi/db/models"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (s Services) GetFullUserInformation(id string) (response.TwitterUserResponse, error) {
	var twitterUserResponse response.TwitterUserResponse
	reqUserFields := "user.fields=created_at,description,entities,id,location,name,pinned_tweet_id,profile_image_url,protected,url,username,verified,withheld"
	reqUrl := fmt.Sprintf("https://api.twitter.com/2/users/%v?%v", id, reqUserFields)
	twitterHttp, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return twitterUserResponse, err
	}
	twitterHttp.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("BEARER_TOKEN")))

	twitterResponse, err := http.DefaultClient.Do(twitterHttp)
	if err != nil {
		return twitterUserResponse, err
	}
	if twitterResponse.StatusCode != http.StatusOK {
		return twitterUserResponse, errors.New("request not successful")
	}

	respBody, err := io.ReadAll(twitterResponse.Body)
	json.Unmarshal(respBody, &twitterUserResponse)
	return twitterUserResponse, nil
}

func (s Services) ExpandTweet(tweetId string) (models.TwitterExpandedResponse, error) {
	var expandedResponse models.TwitterExpandedResponse
	reqFields := "tweet.fields=attachments,author_id,created_at,entities,conversation_id"
	reqExpansions := "expansions=attachments.media_keys,attachments.poll_ids,author_id,in_reply_to_user_id"
	reqMediaFields := "media.fields="
	reqUrl := fmt.Sprintf("https://api.twitter.com/2/tweets?ids=%v&%v&%v&%v", tweetId, reqFields, reqExpansions, reqMediaFields)
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
		return expandedResponse, errors.New("could not expand tweet")
	}

	twitterResBody, err := io.ReadAll(twitterRes.Body)
	if err != nil {
		return expandedResponse, err
	}
	json.Unmarshal(twitterResBody, &expandedResponse)
	return expandedResponse, nil
}

func (s Services) GetThreadByConversationID(conversationID, author string) (string, error) {
	reqFields := "tweet.fields=attachments,author_id,created_at,entities,conversation_id"
	reqExpansions := "expansions=attachments.media_keys,attachments.poll_ids,author_id,in_reply_to_user_id"
	//reqUserFields := ""
	reqQuery := fmt.Sprintf("from:%v conversation_id:%v", author, conversationID)
	reqUrl := fmt.Sprintf("https://api.twitter.com/2/tweets/search/recent?query=%v&max_results=50&%v&%v", url.QueryEscape(reqQuery), reqFields, reqExpansions)

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
	//s.Logger.Info().Msg(string(reqBody))
	return string(reqBody), nil

}
