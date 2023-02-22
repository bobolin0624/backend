package model

type PoliticianQuestionCreate struct {
	Category     string
	UserId       string
	Question     string
	PoliticianId int64
}

type PoliticianQuestion struct {
	Id           int64
	PoliticianId int64
	UserName     string

	Category string
	Question string
	Reply    string
	Likes    int64

	CreatedAt int64
	ReplyAt   int64
}
