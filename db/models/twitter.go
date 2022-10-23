package models

type TwitterExpandedData struct {
	ID               string      `json:"id"`
	InReplyToUserId  string      `json:"in_reply_to_user_id"`
	Entities         interface{} `json:"entities"`
	ConversationID   string      `json:"conversation_id"`
	Text             string      `json:"text"`
	AuthorID         string      `json:"author_id"`
	CreatedAt        string      `json:"created_at"`
	Attachments      interface{} `json:"attachments"`
	ExtendedEntities interface{} `json:"extended_entities"`
}

type TwitterExpandedResponse struct {
	Data []TwitterExpandedData `json:"data"`
}
