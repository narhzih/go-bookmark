package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	_ "google.golang.org/api/youtube/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func ParseYoutubeLink(youtubeLink string) (YoutubeAPIResponse, error) {
	pattern, _ := regexp.Compile("(?:https?:\\/{2})?(?:w{3}\\.)?youtu(?:be)?\\.(?:com|be)(?:\\/watch\\?v=|\\/)([^\\s&]+)")
	videoId := pattern.FindAllStringSubmatch(youtubeLink, -1)[0][1]
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%v&key=%v&part=snippet,contentDetails,statistics,status", videoId, os.Getenv("YOUTUBE_API_KEY")))
	//resp, err := http.Get(fmt.Sprintf("https://www.youtube.com/oembed?url=%+v&format=json", youtubeLink))
	if err != nil {
		return YoutubeAPIResponse{}, err
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
		return YoutubeAPIResponse{}, err
	}
	var yResp YoutubeVideoInformation
	json.Unmarshal(respBody, &yResp)

	authorResp, err := http.Get(fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/channels?part=snippet,contentDetails,statistics&id=%v&key=%v", yResp.Items[0].Snippet.ChannelId, os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		return YoutubeAPIResponse{}, errors.New("author information unavailable")
	}
	authorRespBody, err := io.ReadAll(authorResp.Body)
	if err != nil {
		logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
		logger.Err(err).Msg("Error happened when reading stream")
		return YoutubeAPIResponse{}, err
	}
	var aResp YoutubeAuthorInformation
	json.Unmarshal(authorRespBody, &aResp)

	// build final response
	var youtubeApiResponse YoutubeAPIResponse
	youtubeApiResponse.Video = yResp
	youtubeApiResponse.Author = aResp
	return youtubeApiResponse, nil
}
