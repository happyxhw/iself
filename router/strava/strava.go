package strava

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	api2 "git.happyxhw.cn/happyxhw/iself/api"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	sm "git.happyxhw.cn/happyxhw/iself/pkg/strava"
	"git.happyxhw.cn/happyxhw/iself/service"
)

const (
	verifyToken = "4DpFAupeI7SGSgCCPCub" //nolint:gosec
)

const (
	defaultPageSize = 20
)

// strava api
type strava struct {
	srv     *service.Strava
	goalSrv *service.Goal
}

// Push strava push event
func (s *strava) Push(c echo.Context) error {
	var req sm.SubscriptionEvent
	if err := em.Bind(c, &req); err != nil {
		return em.ErrBadRequest.Wrap(err)
	}
	if emErr := s.srv.Push(c.Request().Context(), &req); emErr != nil {
		return emErr
	}
	return em.OK(c, nil)
}

// VerifyPush 校验订阅时服务的可用性
func (s *strava) VerifyPush(c echo.Context) error {
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
func (s *strava) ActivityList(c echo.Context) error {
	var req api2.ActivityListReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	if req.Limit > defaultPageSize || req.Limit == 0 {
		req.Limit = defaultPageSize
	}
	sourceID, _ := c.Get("id").(int64)
	results, err := s.srv.ActivityList(c.Request().Context(), &req, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, results)
}

func (s *strava) Activity(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		return em.ErrParam
	}
	sourceID, _ := c.Get("id").(int64)
	result, err := s.srv.Activity(c.Request().Context(), id, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, result)
}

func (s *strava) SummaryStatsNow(c echo.Context) error {
	var req api2.StatsNowReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	result, err := s.srv.SummaryStatsNow(c.Request().Context(), &req, sourceID)
	if err != nil {
		return err
	}

	return em.OK(c, result)
}

func (s *strava) CreateGoal(c echo.Context) error {
	var req api2.GoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	err := s.goalSrv.Create(&req)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

func (s *strava) UpdateGoal(c echo.Context) error {
	var req api2.GoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	err := s.goalSrv.UpdateValue(&req)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

func (s *strava) DeleteGoal(c echo.Context) error {
	var req api2.GoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	err := s.goalSrv.Delete(&req)
	if err != nil {
		return err
	}
	return em.OK(c, nil)
}

func (s *strava) QueryGoal(c echo.Context) error {
	var req api2.QueryGoalReq
	if err := em.Bind(c, &req); err != nil {
		return em.ErrParam.Wrap(err)
	}
	sourceID, _ := c.Get("id").(int64)
	req.SourceID = sourceID
	result, err := s.goalSrv.Query(sourceID, req.Type, req.Field)
	if err != nil {
		return err
	}
	return em.OK(c, echo.Map{
		"list": result,
	})
}
