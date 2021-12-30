package api

import (
	"git.happyxhw.cn/happyxhw/iself/model"
)

// SignUpReq request
type SignUpReq struct {
	Name      string `json:"name" validate:"gte=1,lte=64"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"gte=8,lte=16"`
	ActiveURL string `json:"active_url" validate:"required,url"`
}

// SignInReq request
type SignInReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=8,lte=16"`
}

// Oauth2ExchangeReq oauth2 exchange req
type Oauth2ExchangeReq struct {
	// Source string `query:"source" validate:"oneof=github strava google"`
	Source string `query:"source"`
	Code   string `query:"code" validate:"required"`
	State  string `query:"state" validate:"required"`
}

// SetStateReq state req
type SetStateReq struct {
	State string `json:"state" validate:"required"`
	URL   string `json:"url" validate:"required"`
}

// ActiveReq active req
type ActiveReq struct {
	Token string `json:"token" validate:"required"`
}

// SendEmailReq send email req
type SendEmailReq struct {
	Email string `json:"email" validate:"required,email"`
	Type  string `json:"type" validate:"oneof=active reset"`
	URL   string `json:"url" validate:"required,url"`
}

// ChangePasswordReq change password req
type ChangePasswordReq struct {
	Old string `json:"old" validate:"gte=8,lte=16"`
	New string `json:"new" validate:"gte=8,lte=16"`
}

// ResetPasswordReq reset password req
type ResetPasswordReq struct {
	Password string `json:"password" validate:"gte=8,lte=16"`
	Token    string `json:"token" validate:"required"`
}

// Info for user
type Info struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// NewInfo return user info
func NewInfo(u *model.User) *Info {
	return &Info{
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}
