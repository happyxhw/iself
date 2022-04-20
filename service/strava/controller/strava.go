package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/strava"
	"git.happyxhw.cn/happyxhw/iself/service/strava/handler"
	"git.happyxhw.cn/happyxhw/iself/service/strava/types"
)

const (
	verifyToken = "4DpFAupeI7SGSgCCPCub" //nolint:gosec
)

const (
	defaultPageSize = 20
)

// Strava types
type Strava struct {
	srv     *handler.Strava
	goalSrv *handler.Goal
}

func NewStrava(srv *handler.Strava, goalSrv *handler.Goal) *Strava {
	return &Strava{
		srv:     srv,
		goalSrv: goalSrv,
	}
}

// Push Strava push event
func (s *Strava) Push(c echo.Context) error {
	var req strava.SubscriptionEvent
	if err := em.Bind(c, &req); err != nil {
		return em.ErrBadRequest.Wrap(err)
	}
	if emErr := s.srv.Push(em.GetCtx(c), &req); emErr != nil {
		return emErr
	}
	return em.OK(c, nil)
}

// VerifyPush 校验订阅时服务的可用性
func (s *Strava) VerifyPush(c echo.Context) error {
	mode := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")
	challenge := c.QueryParam("hub.challenge")

	if mode != "subscribe" || challenge == "" {
		return em.ErrBadRequest
	}
	if token != verifyToken {
		return em.ErrBadRequest.Msg("verify token")
	}

	return em.OK(c, echo.Map{
		"code":          http.StatusOK,
		"hub.challenge": challenge,
	})
}

// ActivityList 活动列表
func (s *Strava) ActivityList(c echo.Context) error {
	var req types.ActivityListPageReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	if req.PageSize > defaultPageSize || req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}
	sourceID, _ := c.Get("id").(int64)
	results, err := s.srv.ActivityPageList(em.GetCtx(c), &req, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, results)
}

func (s *Strava) Activity(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		return em.ErrParam
	}
	sourceID, _ := c.Get("id").(int64)
	result, err := s.srv.Activity(em.GetCtx(c), id, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, result)
}

func (s *Strava) SummaryStatsNow(c echo.Context) error {
	var req types.StatsNowReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	result, err := s.srv.SummaryStatsNow(em.GetCtx(c), &req, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, result)
}

func (s *Strava) DateChart(c echo.Context) error {
	var req types.DateChartReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	result, err := s.srv.DateChart(em.GetCtx(c), &req, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, result)
}

func (s *Strava) CreateGoal(c echo.Context) error {
	var req types.GoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	err := s.goalSrv.Create(em.GetCtx(c), &req)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

func (s *Strava) UpdateGoal(c echo.Context) error {
	var req types.GoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	err := s.goalSrv.UpdateValue(em.GetCtx(c), &req)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

func (s *Strava) DeleteGoal(c echo.Context) error {
	var req types.GoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	err := s.goalSrv.Delete(em.GetCtx(c), &req)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

func (s *Strava) QueryGoal(c echo.Context) error {
	var req types.QueryGoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	result, err := s.goalSrv.Query(em.GetCtx(c), sourceID, req.Type, req.Field)
	if err != nil {
		return err
	}
	return em.OK(c, echo.Map{
		"list": result,
	})
}
