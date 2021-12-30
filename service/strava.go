package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"git.happyxhw.cn/happyxhw/iself/api"
	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	sm "git.happyxhw.cn/happyxhw/iself/pkg/strava"
)

var (
	// ErrStrava 获取 strava 数据错误
	ErrStrava = em.NewError(http.StatusServiceUnavailable, 60201, "strava api")
)

const (
	notExistsLabel = "--"
)

type MaxMinID struct {
	MaxID int64 `gorm:"column:max_id"`
	MinID int64 `gorm:"column:min_id"`
}

// Strava service
type Strava struct {
	tokenSrv *Token

	db  *gorm.DB
	rdb *redis.Client

	oauth2Conf *oauth2.Config
}

// StravaOption config for strava
type StravaOption struct {
	DB     *gorm.DB
	RDB    *redis.Client
	AesKey string

	Oauth2Conf *oauth2.Config
}

// NewStrava return strava srv
func NewStrava(opt *StravaOption) *Strava {
	return &Strava{
		db:         opt.DB,
		rdb:        opt.RDB,
		tokenSrv:   &Token{rdb: opt.RDB},
		oauth2Conf: opt.Oauth2Conf,
	}
}

// ActivityList get activity list
func (s *Strava) ActivityList(_ context.Context, req *api.ActivityListReq, sourceID int64) (interface{}, *em.Error) {
	var results []*model.ActivityDetail
	fields := []string{
		"id", "name", "distance", "moving_time", "elapsed_time", "total_elevation_gain", "elev_high",
		"elev_low", "type", "start_date_local", "average_speed", "max_speed", "max_heartrate",
		"average_heartrate",
	}
	limit := req.Limit
	query := s.db.Select(fields).Where("athlete_id = ?", sourceID).Limit(limit)
	if req.Type != api.All {
		query = query.Where("type = ?", req.Type)
	}
	var reverse bool
	if req.After > 0 {
		query = query.Where("id < ?", req.After).Order("id DESC")
	} else if req.Before > 0 {
		reverse = true
		query = query.Where("id > ?", req.Before).Order("id ASC")
	} else {
		query = query.Order("id DESC")
	}
	err := query.Find(&results).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	if len(results) == 0 {
		return nil, nil
	}
	// 是否有下一页或上一页
	var tmp MaxMinID
	query = s.db.Model(&model.ActivityDetail{}).Select("MAX(id) AS max_id, MIN(id) AS min_id").
		Where("athlete_id = ?", sourceID)
	if req.Type != api.All {
		query = query.Where("type = ?", req.Type)
	}
	err = query.Find(&tmp).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	var after, before int64
	if reverse {
		sort.SliceStable(results, func(i, j int) bool {
			return true
		})
	}
	if len(results) > 0 {
		after = results[len(results)-1].ID
		if after == tmp.MinID {
			after = 0
		}
		before = results[0].ID
		if before == tmp.MaxID {
			before = 0
		}
	}

	return echo.Map{
		"after":  after,
		"before": before,
		"list":   results,
	}, nil
}

func (s *Strava) Activity(ctx context.Context, id, sourceID int64) (echo.Map, *em.Error) {
	fields := []string{
		"id", "name", "distance", "moving_time", "elapsed_time", "total_elevation_gain", "elev_high",
		"elev_low", "type", "start_date_local", "average_speed", "max_speed", "max_heartrate", "average_heartrate",
		"created_at", "polyline", "calories",
	}
	var activity model.ActivityDetail
	err := s.db.Select(fields).Where("id = ? AND athlete_id = ?", id, sourceID).Find(&activity).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	if activity.ID == 0 {
		return nil, em.ErrNotFound.Msg("activity not found")
	}

	var stream model.ActivityStream
	fields = []string{"id", "distance", "heartrate", "altitude", "velocity_smooth"}
	err = s.db.Select(fields).Where("id = ?", id).Find(&stream).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	var streamSet sm.StreamSet
	_ = json.Unmarshal([]byte(stream.Distance), &streamSet.Distance)
	_ = json.Unmarshal([]byte(stream.VelocitySmooth), &streamSet.VelocitySmooth)
	_ = json.Unmarshal([]byte(stream.Heartrate), &streamSet.Heartrate)
	_ = json.Unmarshal([]byte(stream.Altitude), &streamSet.Altitude)

	if streamSet.VelocitySmooth != nil {
		for i, item := range streamSet.VelocitySmooth.Data {
			streamSet.VelocitySmooth.Data[i] = float32(transformVelocity(float64(item), activity.Type))
		}
	}
	activity.AverageHeartrate = transformVelocity(activity.AverageHeartrate, activity.Type)
	res := echo.Map{
		"activity":        activity,
		"distance":        streamSet.Distance,
		"velocity_smooth": streamSet.VelocitySmooth,
		"heartrate":       streamSet.Heartrate,
		"altitude":        streamSet.Altitude,
	}
	return res, nil
}

type stats struct {
	Week  float64
	Month float64
	Year  float64
	All   float64
}

var unitMap = map[string]string{
	"distance":    "km",
	"moving_time": "s",
	"calories":    "Cal",
}

var fractionMap = map[string]float64{
	"distance": 1000,
}

func (s *Strava) SummaryStatsNow(ctx context.Context, req *api.StatsNowReq, sourceID int64) (*api.Stats, *em.Error) {
	now := time.Now()
	year, month, _ := now.Date()
	week := int(now.Weekday())
	monthStart := fmt.Sprintf("%d-%d-01", year, month)
	yearStart := fmt.Sprintf("%d-01-01", year)
	weekStart := now.AddDate(0, 0, -week+1).Format("2006-01-02")

	var tmp stats
	err := func() error {
		if err := s.db.Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS week", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, weekStart).
			Find(&tmp).Error; err != nil {
			return err
		}
		if err := s.db.Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS month", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, monthStart).
			Find(&tmp).Error; err != nil {
			return err
		}
		if err := s.db.Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS year", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, yearStart).
			Find(&tmp).Error; err != nil {
			return err
		}
		if err := s.db.Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS all", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ?", sourceID, req.Type).
			Find(&tmp).Error; err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	var goals []*model.Goal
	err = s.db.Select("freq, value").
		Where("athlete_id = ? AND type = ? AND field = ?", sourceID, req.Type, req.Field).
		Find(&goals).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	goalMap := make(map[string]float64, len(goals))
	for _, item := range goals {
		goalMap[item.Freq] = item.Value
	}
	var fraction float64 = 1
	if v, ok := fractionMap[req.Field]; ok {
		fraction = v
	}
	r := api.Stats{
		Type: req.Type,
		Unit: unitMap[req.Field],
		All:  fmt.Sprintf("%.0f", tmp.All/fraction),

		Week:     fmt.Sprintf("%.0f", tmp.Week/fraction),
		WeekGoal: fmt.Sprintf("%.0f", goalMap["week"]/fraction),

		Month:     fmt.Sprintf("%.0f", tmp.Month/fraction),
		MonthGoal: fmt.Sprintf("%.0f", goalMap["month"]/fraction),

		Year:     fmt.Sprintf("%.0f", tmp.Year/fraction),
		YearGoal: fmt.Sprintf("%.0f", goalMap["year"]/fraction),
	}
	if int(goalMap["week"]) == 0 {
		r.WeekGoal, r.WeekProcess = notExistsLabel, notExistsLabel
	} else {
		r.WeekProcess = fmt.Sprintf("%.0f", tmp.Week/goalMap["week"]/100)
	}
	if int(goalMap["month"]) == 0 {
		r.MonthGoal, r.MonthProcess = notExistsLabel, notExistsLabel
	} else {
		r.MonthProcess = fmt.Sprintf("%.0f", tmp.Month/goalMap["month"]/100)
	}
	if int(goalMap["year"]) == 0 {
		r.YearGoal, r.YearProcess = notExistsLabel, notExistsLabel
	} else {
		r.YearProcess = fmt.Sprintf("%.0f", tmp.Year/goalMap["year"]*100)
	}
	return &r, nil
}

func transformVelocity(vel float64, activityType string) float64 {
	if vel == 0 {
		return vel
	}
	switch activityType {
	case "run":
		// 转为配速 min/km
		t := 16.666666667 / vel
		minutes := float64(int(t))
		seconds := (t - minutes) * 60.0 / 100.0
		return minutes + seconds
	case "ride", "virtualride":
		// 转为速度 km/h
		return 3.6 * vel
	}
	return vel
}

// Push from strava event
func (s *Strava) Push(ctx context.Context, event *sm.SubscriptionEvent) *em.Error {
	ups, err := json.Marshal(event.Updates)
	if err != nil {
		return em.ErrBadRequest.Wrap(err)
	}

	var push model.PushEvent
	err = s.db.Select("id, status").Where("object_id = ?", event.ObjectID).Find(&push).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if push.ID == 0 {
		now := time.Now()
		newPush := model.PushEvent{
			OwnerID:    event.OwnerID,
			AspectType: event.AspectType,
			EventTime:  event.EventTime,
			ObjectID:   event.ObjectID,
			ObjectType: event.ObjectType,
			Updates:    string(ups),
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		err = s.db.Where("object_id = ?", event.ObjectID).Create(&newPush).Error
		if err != nil {
			return em.ErrDB.Wrap(err)
		}
	}
	if push.Status != 1 {
		emErr := s.push(ctx, event)
		if emErr != nil {
			log.Error("strava push", zap.Int64("object_id", event.ObjectID))
			return emErr
		}
	}

	return nil
}

// push 处理推送
func (s *Strava) push(ctx context.Context, event *sm.SubscriptionEvent) *em.Error {
	log.Info("strava push",
		zap.Int64("object_id", event.ObjectID), zap.Int64("owner_id", event.OwnerID),
		zap.String("aspect_type", event.AspectType), zap.String("object_type", event.ObjectType))
	switch event.ObjectType {
	case "athlete":
		return s.athletePush(ctx, event)
	case "activity":
		return s.activityPush(ctx, event)
	}
	return em.ErrBadRequest.Msg("unknown object type")
}

func (s *Strava) activityPush(ctx context.Context, event *sm.SubscriptionEvent) *em.Error {
	switch event.AspectType {
	case "create":
		return s.activityCreate(ctx, event)
	case "update":
		return em.ErrBadRequest.Msg("unknown aspect type")
	}
	return em.ErrBadRequest.Msg("unknown aspect type")
}

func (s *Strava) athletePush(ctx context.Context, event *sm.SubscriptionEvent) *em.Error {
	return nil
}

func (s *Strava) activityCreate(ctx context.Context, event *sm.SubscriptionEvent) *em.Error {
	token, err := s.tokenSrv.GetToken("strava", event.OwnerID, s.oauth2Conf)
	if err != nil {
		return ErrGetToken.Wrap(err)
	}
	stravaCli := sm.NewClient(s.oauth2Conf.Client(ctx, token))
	resp, body, err := stravaCli.Activity.Activity(ctx, event.ObjectID)
	if err != nil {
		return ErrStrava.Wrap(err)
	}
	raw := model.ActivityRaw{
		ID:   resp.Id,
		Data: string(body),
	}

	streamResp, err := stravaCli.Activity.ActivityStream(ctx, event.ObjectID)
	if err != nil {
		return ErrStrava.Wrap(err)
	}

	activity, stream := formatActivityData(resp, streamResp, event)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if dbErr := tx.Create(&raw).Error; dbErr != nil {
			return dbErr
		}
		if dbErr := tx.Create(&activity).Error; dbErr != nil {
			return dbErr
		}
		if dbErr := tx.Create(&stream).Error; dbErr != nil {
			return dbErr
		}
		if dbErr := tx.Model(&model.PushEvent{}).
			Where("object_id = ?", event.ObjectID).
			Update("status", 1).Error; dbErr != nil {
			return dbErr
		}
		return nil
	})
	if err != nil {
		return em.ErrDB.Wrap(err)
	}

	return nil
}

func formatActivityData(resp *sm.DetailedActivity, streamResp *sm.StreamSet, event *sm.SubscriptionEvent,
) (*model.ActivityDetail, *model.ActivityStream) {
	splitsMetric, _ := json.Marshal(resp.SplitsMetric)
	bestEfforts, _ := json.Marshal(resp.BestEfforts)
	ti, _ := json.Marshal(streamResp.Time)
	distance, _ := json.Marshal(streamResp.Distance)
	latlng, _ := json.Marshal(streamResp.Latlng)
	altitude, _ := json.Marshal(streamResp.Altitude)
	vel, _ := json.Marshal(streamResp.VelocitySmooth)
	heart, _ := json.Marshal(streamResp.Heartrate)
	cadence, _ := json.Marshal(streamResp.Cadence)
	watts, _ := json.Marshal(streamResp.Watts)
	temp, _ := json.Marshal(streamResp.Temp)
	moving, _ := json.Marshal(streamResp.Moving)
	grade, _ := json.Marshal(streamResp.GradeSmooth)

	activity := model.ActivityDetail{
		ID:                 resp.Id,
		AthleteID:          resp.Athlete.Id,
		Name:               resp.Name,
		Type:               strings.ToLower(resp.Type),
		Distance:           float64(resp.Distance),
		MovingTime:         resp.MovingTime,
		ElapsedTime:        resp.ElapsedTime,
		TotalElevationGain: float64(resp.TotalElevationGain),
		StartDateLocal:     resp.StartDateLocal,
		Polyline:           resp.Map.Polyline,
		SummaryPolyline:    resp.Map.SummaryPolyline,
		AverageSpeed:       float64(resp.AverageSpeed),
		MaxSpeed:           float64(resp.MaxSpeed),
		AverageHeartrate:   float64(resp.AverageHeartrate),
		MaxHeartrate:       float64(resp.MaxHeartrate),
		ElevHigh:           float64(resp.ElevHigh),
		ElevLow:            float64(resp.ElevLow),
		Calories:           float64(resp.Calories),
		SplitsMetric:       string(splitsMetric),
		BestEfforts:        string(bestEfforts),
		DeviceName:         resp.DeviceName,
	}

	stream := model.ActivityStream{
		ID:             event.ObjectID,
		Time:           string(ti),
		Distance:       string(distance),
		Latlng:         string(latlng),
		Altitude:       string(altitude),
		VelocitySmooth: string(vel),
		Heartrate:      string(heart),
		Cadence:        string(cadence),
		Watts:          string(watts),
		Temp:           string(temp),
		Moving:         string(moving),
		GradeSmooth:    string(grade),
	}

	return &activity, &stream
}
