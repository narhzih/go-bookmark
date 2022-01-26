package mp_parser

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func ParseYoutubeLink(youtubeLink string) (string, error) {
	//timeout := time.Duration(10 * time.Second)
	//requestBody, err := json.Marshal(map[string]interface{}{})
	//client := http.Client{Timeout: timeout}
	//request, err := http.NewRequest("GET", fmt.Sprintf("https://www.youtube.com/oembed?url=%+v&format=json", youtubeLink), bytes.NewBuffer(requestBody))
	resp, err := http.Get(fmt.Sprintf("https://www.youtube.com/oembed?url=%+v&format=json", youtubeLink))
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
