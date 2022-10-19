package handler

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"

	"github.com/happyxhw/iself/model"
	"github.com/happyxhw/iself/pkg/query"
	"github.com/happyxhw/iself/pkg/util"
	"github.com/happyxhw/iself/service/user/handler/mocks"
	"github.com/happyxhw/iself/service/user/types"
)

const (
	mockURL = "https://mock.com"
	mockKey = "eShVkYp3s6v9y$B&E)H@McQfTjWnZq4t"
)

var mockUser = model.User{
	ID:       1,
	Name:     "mockName",
	Email:    "mock@email.com",
	Password: "mockPassword",
	Source:   "strava",
	SourceID: 1,
}

var mockOauth2Token = oauth2.Token{
	AccessToken:  "at",
	TokenType:    "code",
	RefreshToken: "rk",
}

var mockToken = util.NanoID(16)

func TestUser_InfoBySource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		ur: userRepo,
	}
	gomock.InOrder(
		userRepo.EXPECT().GetBySource(ctx, mockUser.Source, mockUser.SourceID, gomock.Any()).Return(&mockUser, nil),
	)

	_, err := h.Info(ctx, 0, mockUser.SourceID, mockUser.Source)

	require.NoError(t, err)
}

func TestUser_InfoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		ur: userRepo,
	}
	gomock.InOrder(
		userRepo.EXPECT().Get(ctx, mockUser.ID, gomock.Any()).Return(&mockUser, nil),
	)

	_, err := h.Info(ctx, mockUser.ID, 0, "")

	require.NoError(t, err)
}

func TestUser_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	cacher := mocks.NewMockCacher(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		cacher: cacher,
		mailer: mailer,
		ur:     userRepo,
	}
	email := mockUser.Email
	freqKey, activeKey := "active_freq:"+email, "active:"+email

	gomock.InOrder(
		userRepo.EXPECT().GetByEmail(ctx, email, query.Opt{Fields: []string{"id"}}).Return(nil, nil),
		userRepo.EXPECT().Create(ctx, gomock.Any()).Return(&mockUser, nil),
		cacher.EXPECT().SetNX(ctx, freqKey, nil, emailExpire).Return(true, nil),
		cacher.EXPECT().Set(ctx, activeKey, gomock.Any(), tokenExpire),
		mailer.EXPECT().Send(email, gomock.Any(), gomock.Any()),
	)

	var req = types.SignUpReq{
		Name:      mockUser.Name,
		Email:     mockUser.Email,
		Password:  mockUser.Password,
		ActiveURL: mockUser.AvatarURL,
	}

	err := h.SignUp(ctx, &req)

	require.NoError(t, err)
}

func TestUser_SignUpExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	cacher := mocks.NewMockCacher(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		cacher: cacher,
		mailer: mailer,
		ur:     userRepo,
	}
	email := mockUser.Email

	gomock.InOrder(
		userRepo.EXPECT().GetByEmail(ctx, email, query.Opt{Fields: []string{"id"}}).Return(&mockUser, nil),
	)

	var req = types.SignUpReq{
		Name:      mockUser.Name,
		Email:     mockUser.Email,
		Password:  mockUser.Password,
		ActiveURL: mockUser.AvatarURL,
	}

	err := h.SignUp(ctx, &req)

	require.Equal(t, err, ErrUserExists)
}

func TestUser_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		ur: userRepo,
	}
	email := mockUser.Email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(mockUser.Password), bcrypt.DefaultCost)
	req := types.SignInReq{
		Email:    email,
		Password: mockUser.Password,
	}
	mockUser.Password = string(hashedPassword)

	gomock.InOrder(
		userRepo.EXPECT().GetByEmail(ctx, email, query.Opt{}).Return(&mockUser, nil),
	)

	u, err := h.SignIn(ctx, &req)

	require.NoError(t, err)
	require.Equal(t, u.ID, mockUser.ID)
}

func TestUser_SignInByOauth2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	userRepo := mocks.NewMockUserRepo(ctrl)
	tokenRepo := mocks.NewMockTokenRepo(ctrl)
	oauth2Provider := mocks.NewMockOauth2x(ctrl)

	h := User{
		ur: userRepo,
		tr: tokenRepo,
	}

	mockCode := "mockCode"
	email := fmt.Sprintf("%d@%s", mockUser.SourceID, mockUser.Source)

	gomock.InOrder(
		oauth2Provider.EXPECT().Exchange(ctx, mockCode).Return(&mockOauth2Token, nil),
		oauth2Provider.EXPECT().GetUser(ctx, &mockOauth2Token).Return(&mockUser, nil),
		tokenRepo.EXPECT().SaveToken(ctx, &mockOauth2Token, mockUser.Source, mockUser.SourceID),

		userRepo.EXPECT().GetByEmail(ctx, email, query.Opt{}).Return(nil, nil),
		userRepo.EXPECT().GetBySource(ctx, mockUser.Source, mockUser.SourceID, query.Opt{}).Return(nil, nil),
		userRepo.EXPECT().Create(ctx, gomock.Any()).Return(&mockUser, nil),
	)

	u, err := h.SignInByOauth2(ctx, mockUser.Source, mockCode, oauth2Provider)

	require.NoError(t, err)
	require.Equal(t, u.ID, mockUser.ID)
}

func TestUser_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		ur: userRepo,
	}
	req := types.ChangePasswordReq{
		Old: mockUser.Password,
		New: "456",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(mockUser.Password), bcrypt.DefaultCost)
	mockUser.Password = string(hashedPassword)
	gomock.InOrder(
		userRepo.EXPECT().Get(ctx, mockUser.ID, query.Fields("id", "password")).Return(&mockUser, nil),
		userRepo.EXPECT().Update(ctx, mockUser.ID, gomock.Any()).Return(int64(1), nil),
	)

	err := h.ChangePassword(ctx, mockUser.ID, &req)

	require.NoError(t, err)
}

func TestUser_ResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	cacher := mocks.NewMockCacher(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		cacher: cacher,
		mailer: mailer,
		ur:     userRepo,
	}
	email := mockUser.Email
	token, err := h.encryptToken(email, mockToken)
	if err != nil {
		t.Fatal(err)
	}
	token, _ = url.QueryUnescape(token)
	key := fmt.Sprintf("reset:%s", email)
	req := types.ResetPasswordReq{
		Password: "456",
		Token:    token,
	}
	gomock.InOrder(
		cacher.EXPECT().GetString(ctx, key).Return(mockToken, nil),
		userRepo.EXPECT().UpdateByEmail(ctx, mockUser.Email, gomock.Any()).Return(int64(1), nil),
		cacher.EXPECT().Del(ctx, key).Return(int64(1), nil),
	)

	err = h.ResetPassword(ctx, &req)

	require.NoError(t, err)
}

func TestUser_Activate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	cacher := mocks.NewMockCacher(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		cacher: cacher,
		mailer: mailer,
		ur:     userRepo,
	}
	email := mockUser.Email

	token, err := h.encryptToken(email, mockToken)
	if err != nil {
		t.Fatal(err)
	}
	token, _ = url.QueryUnescape(token)
	key := fmt.Sprintf("active:%s", email)
	params := model.UserParam{
		Status: util.Int(int(model.ActivatedStatus)),
	}
	gomock.InOrder(
		cacher.EXPECT().GetString(ctx, key).Return(mockToken, nil),
		userRepo.EXPECT().UpdateByEmail(ctx, mockUser.Email, &params).Return(int64(1), nil),
		cacher.EXPECT().Del(ctx, key).Return(int64(1), nil),
	)

	err = h.Active(ctx, token)

	require.NoError(t, err)
}

func TestUser_SendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()

	cacher := mocks.NewMockCacher(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	h := User{
		aesKey: []byte(mockKey),

		cacher: cacher,
		mailer: mailer,
		ur:     userRepo,
	}
	email := mockUser.Email
	freqKey, activeKey := "active_freq:"+email, "active:"+email

	gomock.InOrder(
		userRepo.EXPECT().GetByEmail(ctx, email, query.Opt{Fields: []string{"id"}}).Return(&model.User{}, nil),
		cacher.EXPECT().SetNX(ctx, freqKey, nil, emailExpire).Return(true, nil),
		cacher.EXPECT().Set(ctx, activeKey, gomock.Any(), tokenExpire),
		mailer.EXPECT().Send(email, gomock.Any(), gomock.Any()),
	)

	err := h.SendEmail(ctx, email, ActiveEmail, mockURL)

	require.NoError(t, err)
}
