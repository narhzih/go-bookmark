package models

type TwitterExpandedData struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Text           string `json:"text"`
	AuthorID       string `json:"author_id"`
}

type TwitterExpandedResponse struct {
	Data []TwitterExpandedData `json:"data"`
}
