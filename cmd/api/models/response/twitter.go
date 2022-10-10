package response

type TwitterUserObject struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type TwitterUserResponse struct {
	Data TwitterUserObject `json:"data"`
}
