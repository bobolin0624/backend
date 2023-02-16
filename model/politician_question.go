package model

type PoliticianQuestionCreate struct {
	UserId       string `json:"user_id"`
	PoliticianId int64  `json:"politician_id"`
	Type         string `json:"type"`
	Question     string `json:"question"`
}

type PoliticianQuestion struct {
	Id           int64
	PoliticianId int64
	UserName     string

	Type     string
	Question string
	Reply    string
	Likes    int64

	CreatedAt int64
	ReplyAt   int64
}
