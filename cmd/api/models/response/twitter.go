package response

type TwitterUserObject struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
	Verified        string `json:"verified"`
	WithHeld        string `json:"withheld"`
}

type TwitterUserResponse struct {
	Data TwitterUserObject `json:"data"`
}
