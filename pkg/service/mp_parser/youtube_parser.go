package mp_parser

import (
	"fmt"
	"github.com/rs/zerolog"
	_ "google.golang.org/api/youtube/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func ParseYoutubeLink(youtubeLink string) (string, error) {
	pattern, _ := regexp.Compile("(?:https?:\\/{2})?(?:w{3}\\.)?youtu(?:be)?\\.(?:com|be)(?:\\/watch\\?v=|\\/)([^\\s&]+)")
	videoId := pattern.FindAllStringSubmatch(youtubeLink, -1)[0][1]
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%v&key=%v&part=snippet,contentDetails,statistics,status", videoId, os.Getenv("YOUTUBE_API_KEY")))
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
