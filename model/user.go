package model

type User struct {
	Id        string
	Name      string
	AvatarURL string
	Email     string

	GoogleId string
}

type UserRepr struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	Email     string `json:"email"`
}

func (u *User) Repr() *UserRepr {
	return &UserRepr{
		Id:        u.Id,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		Email:     u.Email,
	}
}
