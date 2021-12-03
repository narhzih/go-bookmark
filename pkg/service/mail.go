package service

import "gitlab.com/trencetech/mypipe-api/db/model"

func mailTo() {}

func (s Service) SendWelcomeMail(user model.User)       {}
func (s Service) SendVerificationMail(user model.User)  {}
func (s Service) SendPasswordResetMail(user model.User) {}
