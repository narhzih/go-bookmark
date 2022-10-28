package services

import (
	"fmt"
	"github.com/mypipeapp/mypipeapi/db/actions/postgres"
	"regexp"
	"strings"
)

func (s Services) PipeExists(pipeId, userId int64) (bool, error) {
	pipe, err := s.Repositories.Pipe.GetPipe(pipeId, userId)
	if err != nil {
		return false, err
	}
	if pipe.UserID != userId {
		return false, postgres.ErrNoRecord
	}

	return true, nil
}

func (s Services) UserOwnsPipe(pipeId, userId int64) (bool, error) {
	pipe, err := s.Repositories.Pipe.GetPipe(pipeId, userId)
	if err != nil {
		return false, err
	}

	if pipe.UserID != userId {
		return false, fmt.Errorf("pipe does not belong to this user")
	}

	return true, nil
}

func (s Services) GetPlatformFromLink(link string) (string, error) {
	linkSplit := strings.Split(link, "://")[1]
	r, _ := regexp.Compile("^((?:https?:)?\\/\\/)?((?:www|m)\\.)?((?:youtube(-nocookie)?\\.com|youtu.be))(\\/(?:[\\w\\-]+\\?v=|embed\\/|v\\/)?)([\\w\\-]+)(\\S+)?$")
	if strings.HasPrefix(linkSplit, "twitter") || strings.HasPrefix(linkSplit, "www.twitter") {
		return "twitter", nil
	} else if r.MatchString(linkSplit) {
		return "youtube", nil
	} else {
		return "others", nil
	}
}
