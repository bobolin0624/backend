package model

type CandidateType string

const (
	CandidateTypeLyLocal CandidateType = "ly-local"
	CandidateTypeLyParty CandidateType = "ly-party"
	CandidateTypePres    CandidateType = "pres"
)

type Candidate struct {
	Type         CandidateType
	Term         int
	PoliticianId int
	Number       int
	Elected      bool

	PartyId int
	Area    string

	VicePresident bool
}

type CandidateLyRepr struct {
	Type         CandidateType `json:"type" binding:"required"`
	Term         int           `json:"term" binding:"required"`
	PoliticianId int           `json:"politicianId" binding:"required"`
	Number       int           `json:"number" binding:"required"`
	Elected      bool          `json:"elected"`
	PartyId      int           `json:"partyId"`
	Area         string        `json:"area"`
}

func (c *CandidateLyRepr) Model() *Candidate {
	return &Candidate{
		Type:         c.Type,
		Term:         c.Term,
		PoliticianId: c.PoliticianId,
		Number:       c.Number,
		Elected:      c.Elected,
		PartyId:      c.PartyId,
		Area:         c.Area,
	}
}

func (c *Candidate) ReprLy() *CandidateLyRepr {
	return &CandidateLyRepr{
		Type:         c.Type,
		Term:         c.Term,
		PoliticianId: c.PoliticianId,
		Number:       c.Number,
		Elected:      c.Elected,

		PartyId: c.PartyId,
		Area:    c.Area,
	}
}

type CandidatePresRepr struct {
	Type         string `json:"type"`
	Term         int    `json:"term"`
	PoliticianId int    `json:"politicianId"`
	Number       int    `json:"number"`
	Elected      bool   `json:"elected"`
	PartyId      int    `json:"partyId"`

	VicePresident bool `json:"vicePresident"`
}

func (c *Candidate) ReprPres() *CandidateLyRepr {
	return &CandidateLyRepr{
		Type:         c.Type,
		Term:         c.Term,
		PoliticianId: c.PoliticianId,
		Number:       c.Number,
		Elected:      c.Elected,

		PartyId: c.PartyId,
		Area:    c.Area,
	}
}
