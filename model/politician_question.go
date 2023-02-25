package model

import (
	"time"
)

type PoliticianQuestionCreate struct {
	Category     string
	UserId       string
	Question     string
	PoliticianId int64
}

type PoliticianQuestion struct {
	Id       int64
	Category string

	UserName string
	Question string
	AskedAt  time.Time

	PoliticianId int64
	Reply        *string
	RepliedAt    *time.Time

	Likes int64
}

type PoliticianQuestionRepr struct {
	Category string `json:"category"`

	UserName string `json:"userName"`
	Question string `json:"question"`
	AskedAt  int64  `json:"askedAt"`

	Replied   bool   `json:"replied"`
	Reply     string `json:"reply"`
	RepliedAt int64  `json:"repliedAt"`

	Likes int64 `json:"likes"`
}

func (pq *PoliticianQuestion) Repr() *PoliticianQuestionRepr {
	reply := ""
	repliedAt := int64(0)
	replied := false
	if pq.Reply != nil && pq.RepliedAt != nil {
		reply = *pq.Reply
		repliedAt = pq.RepliedAt.Unix()
		replied = true
	}
	return &PoliticianQuestionRepr{
		Category: pq.Category,
		UserName: pq.UserName,
		Question: pq.Question,
		AskedAt:  pq.AskedAt.Unix(),

		Replied:   replied,
		Reply:     reply,
		RepliedAt: repliedAt,

		Likes: pq.Likes,
	}
}
