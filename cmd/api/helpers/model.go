package helpers

type YoutubeAPIResponse struct {
	Video  interface{} `json:"video"`
	Author interface{} `json:"author"`
}
type DYoutubeItem struct {
	Kind           string          `json:"kind"`
	Etag           string          `json:"etag"`
	Id             string          `json:"id"`
	Snippet        DYoutubeSnippet `json:"snippet"`
	ContentDetails interface{}     `json:"contentDetails"`
	Statistics     interface{}     `json:"statistics"`
}

type DYoutubeAuthorItem struct {
	Kind           string          `json:"kind"`
	Etag           string          `json:"etag"`
	Id             string          `json:"id"`
	Snippet        DYoutubeSnippet `json:"snippet"`
	ContentDetails interface{}     `json:"contentDetails"`
	Statistics     interface{}     `json:"statistics"`
}

type DYoutubeSnippet struct {
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	ChannelId    string      `json:"channelId"`
	PublishedAt  string      `json:"publishedAt"`
	Thumbnails   interface{} `json:"thumbnails"`
	ChannelTitle string      `json:"channelTitle"`
	Tags         interface{} `json:"tags"`
}

type YoutubeVideoInformation struct {
	Kind     string         `json:"kind"`
	Etag     string         `json:"etag"`
	Items    []DYoutubeItem `json:"items"`
	PageInfo interface{}    `json:"pageInfo"`
}
type YoutubeAuthorInformation struct {
	Kind     string               `json:"kind"`
	Etag     string               `json:"etag"`
	Items    []DYoutubeAuthorItem `json:"items"`
	PageInfo interface{}          `json:"pageInfo"`
}
