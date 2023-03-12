package model

type PoliticianPolicyCreate struct {
	PoliticianId int
	Category     string
	Content       string
}

type PoliticianPolicy struct {
	Id       int
	PoliticianId int
	Category string
	Content string
}

type PoliticianPolicyRepr struct {
	Category string `json:"category"`
	Content string `json:"content"`
}

func (pp *PoliticianPolicy) Repr() *PoliticianPolicyRepr {
	return &PoliticianPolicyRepr{
		Category: pp.Category,
		Content:  pp.Content,
	}
}
