package types

import (
	"time"

	"github.com/jinzhu/copier"

	"github.com/happyxhw/iself/model"
)

// User schema
type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Source    string `json:"source"`
	SourceID  int64  `json:"source_id"`
	AvatarURL string `json:"avatar_url"`

	Role   int `json:"role,omitempty"`
	Status int `json:"status,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func NewUser(from *model.User) *User {
	var to User
	_ = copier.Copy(&to, from)

	return &to
}

func NewUserList(from []*model.User) []*User {
	list := make([]*User, 0, len(from))
	for _, item := range from {
		list = append(list, NewUser(item))
	}

	return list
}

type SignUpReq struct {
	Name      string `json:"name" validate:"gte=1,lte=64"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"gte=8,lte=16"`
	ActiveURL string `json:"active_url" validate:"required,url"`
}

type SignInReq struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"gte=8,lte=16"`
	RememberMe bool   `json:"remember_me"`
}

type ChangePasswordReq struct {
	Old string `json:"old" validate:"gte=8,lte=16"`
	New string `json:"new" validate:"gte=8,lte=16"`
}

type ResetPasswordReq struct {
	Password string `json:"password" validate:"gte=8,lte=16"`
	Token    string `json:"token" validate:"required"`
}

type ActiveReq struct {
	Token string `json:"token" validate:"required"`
}

type SendEmailReq struct {
	Email string `json:"email" validate:"required,email"`
	Type  string `json:"type" validate:"oneof=active reset"`
}

type SetOauth2StateReq struct {
	State string `json:"state" validate:"required"`
}

type Oauth2ExchangeReq struct {
	Code  string `query:"code" validate:"required"`
	State string `query:"state" validate:"required"`
}
