package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/happyxhw/iself/model"
	"github.com/happyxhw/iself/pkg/cx"
	"github.com/happyxhw/iself/pkg/ex"
	"github.com/happyxhw/iself/pkg/strava"
	"github.com/happyxhw/iself/service/strava/handler"
	"github.com/happyxhw/iself/service/strava/types"
)

const (
	verifyToken = "4DpFAupeI7SGSgCCPCub" //nolint:gosec
)

type Strava struct {
	srv *handler.Strava
}

func NewStrava(srv *handler.Strava) *Strava {
	return &Strava{
		srv: srv,
	}
}

func (s *Strava) GetActivity(c echo.Context) error {
	activityID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if activityID == 0 {
		return ex.ErrParam.Msg("wrong activity id")
	}
	uc := ex.GetUser(c)
	result, err := s.srv.GetActivity(cx.NewTraceCx(c), activityID, uc.SourceID)
	if err != nil {
		return err
	}

	return ex.OK(c, result)
}

func (s *Strava) ListActivity(c echo.Context) error {
	var req types.ActivityQueryParam
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	if req.SortBy == "" {
		req.SortBy = "-id"
	}
	param := model.StravaActivityParam{
		Param: req.Param,
		Type:  req.ActivityType,
	}
	uc := ex.GetUser(c)
	results, err := s.srv.ListActivity(cx.NewTraceCx(c), uc.SourceID, &param)
	if err != nil {
		return err
	}

	return ex.OK(c, results)
}

func (s *Strava) GetProgressStats(c echo.Context) error {
	var req types.ProgressStatsReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	uc := ex.GetUser(c)
	result, err := s.srv.GetProgressStats(cx.NewTraceCx(c), uc.SourceID, &req)
	if err != nil {
		return err
	}

	return ex.OK(c, result)
}

func (s *Strava) GetAggStats(c echo.Context) error {
	var req types.AggStatsReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	uc := ex.GetUser(c)
	result, err := s.srv.GetAggStats(cx.NewTraceCx(c), uc.SourceID, &req)
	if err != nil {
		return err
	}

	return ex.OK(c, result)
}

func (s *Strava) CreateGoal(c echo.Context) error {
	var req types.CreateGoalReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	uc := ex.GetUser(c)
	err := s.srv.CreateGoal(cx.NewTraceCx(c), uc.SourceID, &req)
	if err != nil {
		return err
	}
	return ex.OK(c, nil)
}

func (s *Strava) QueryGoal(c echo.Context) error {
	var req types.QueryGoalReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	uc := ex.GetUser(c)
	result, err := s.srv.GetGoals(cx.NewTraceCx(c), uc.SourceID, req.Type, req.Field)
	if err != nil {
		return err
	}
	return ex.OK(c, echo.Map{"data": result})
}

func (s *Strava) UpdateGoal(c echo.Context) error {
	var req types.UpdateGoalReq
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	uc := ex.GetUser(c)
	err := s.srv.UpdateGoal(cx.NewTraceCx(c), uc.SourceID, &req)
	if err != nil {
		return err
	}
	return ex.OK(c, nil)
}

func (s *Strava) DeleteGoal(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		return ex.ErrParam.Msg("goal id")
	}
	uc := ex.GetUser(c)
	err := s.srv.DeleteGoal(cx.NewTraceCx(c), uc.SourceID, id)
	if err != nil {
		return err
	}
	return ex.OK(c, nil)
}

// Push Strava push event
func (s *Strava) Push(c echo.Context) error {
	var req strava.SubscriptionEvent
	if err := ex.Bind(c, &req); err != nil {
		return err
	}
	if emErr := s.srv.Push(cx.NewTraceCx(c), &req); emErr != nil {
		return emErr
	}
	return ex.OK(c, nil)
}

// VerifyPush 校验订阅时服务的可用性
func (s *Strava) VerifyPush(c echo.Context) error {
	mode := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")
	challenge := c.QueryParam("hub.challenge")

	if mode != "subscribe" || challenge == "" {
		return ex.ErrBadRequest
	}
	if token != verifyToken {
		return ex.ErrBadRequest.Msg("verify token")
	}

	return ex.OK(c, echo.Map{
		"code":          http.StatusOK,
		"hub.challenge": challenge,
	})
}
