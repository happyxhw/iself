package handler

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

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/pkg/strava"
	"git.happyxhw.cn/happyxhw/iself/service/strava/types"
	"git.happyxhw.cn/happyxhw/iself/service/user/handler"
)

var (
	// ErrStrava 获取 strava 数据错误
	ErrStrava = em.NewError(http.StatusServiceUnavailable, 60201, "strava types")
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
	tokenSrv *handler.Token

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
		tokenSrv:   &handler.Token{RDB: opt.RDB},
		oauth2Conf: opt.Oauth2Conf,
	}
}

// ActivityList get activity list
func (s *Strava) ActivityList(ctx context.Context, req *types.ActivityListReq, sourceID int64) (interface{}, error) {
	var results []*model.ActivityDetail
	fields := []string{
		"id", "name", "distance", "moving_time", "elapsed_time", "total_elevation_gain", "elev_high",
		"elev_low", "type", "start_date_local", "average_speed", "max_speed", "max_heartrate",
		"average_heartrate",
	}
	limit := req.Limit
	query := s.db.WithContext(ctx).Select(fields).Where("athlete_id = ?", sourceID).Limit(limit)
	if req.Type != types.All {
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
	query = s.db.WithContext(ctx).Model(&model.ActivityDetail{}).Select("MAX(id) AS max_id, MIN(id) AS min_id").
		Where("athlete_id = ?", sourceID)
	if req.Type != types.All {
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

// ActivityPageList get activity list
func (s *Strava) ActivityPageList(ctx context.Context, req *types.ActivityListPageReq, sourceID int64) (interface{}, error) {
	fields := []string{
		"id", "name", "distance", "moving_time", "elapsed_time", "total_elevation_gain", "elev_high",
		"elev_low", "type", "start_date_local", "average_speed", "max_speed", "max_heartrate",
		"average_heartrate", "summary_polyline", "calories",
	}
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	query := s.db.WithContext(ctx).Model(&model.ActivityDetail{}).Where("athlete_id = ?", sourceID)
	if req.Type != types.All {
		query = query.Where("type = ?", req.Type)
	}
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	if total == 0 {
		return echo.Map{
			"total": total,
			"list":  []*model.ActivityDetail{},
		}, nil
	}
	var ids []int64
	err = query.Order("id DESC").Limit(limit).Offset(offset).Pluck("id", &ids).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	var results []*model.ActivityDetail
	err = s.db.WithContext(ctx).Select(fields).Where("id IN (?)", ids).Order("id desc").Find(&results).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	for i, item := range results {
		results[i].AverageSpeed = transformVelocity(item.AverageSpeed, item.Type)
		results[i].MaxSpeed = transformVelocity(item.MaxSpeed, item.Type)
	}

	return echo.Map{
		"total": total,
		"list":  results,
	}, nil
}

func (s *Strava) Activity(ctx context.Context, id, sourceID int64) (echo.Map, error) {
	fields := []string{
		"id", "name", "distance", "moving_time", "elapsed_time", "total_elevation_gain", "elev_high",
		"elev_low", "type", "start_date_local", "average_speed", "max_speed", "max_heartrate", "average_heartrate",
		"created_at", "polyline", "calories",
	}
	var activity model.ActivityDetail
	err := s.db.WithContext(ctx).Select(fields).Where("id = ? AND athlete_id = ?", id, sourceID).Find(&activity).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	if activity.ID == 0 {
		return nil, em.ErrNotFound.Msg("activity not found")
	}

	var stream model.ActivityStream
	fields = []string{"id", "distance", "heartrate", "altitude", "velocity_smooth"}
	err = s.db.WithContext(ctx).Select(fields).Where("id = ?", id).Find(&stream).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}

	var streamSet strava.StreamSet
	_ = json.Unmarshal([]byte(stream.Distance), &streamSet.Distance)
	_ = json.Unmarshal([]byte(stream.VelocitySmooth), &streamSet.VelocitySmooth)
	_ = json.Unmarshal([]byte(stream.Heartrate), &streamSet.Heartrate)
	_ = json.Unmarshal([]byte(stream.Altitude), &streamSet.Altitude)

	if streamSet.VelocitySmooth != nil {
		for i, item := range streamSet.VelocitySmooth.Data {
			streamSet.VelocitySmooth.Data[i] = float32(transformVelocity(float64(item), activity.Type))
		}
	}
	activity.AverageSpeed = transformVelocity(activity.AverageSpeed, activity.Type)
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

func (s *Strava) SummaryStatsNow(ctx context.Context, req *types.StatsNowReq, sourceID int64) (*types.Stats, error) {
	now := time.Now()
	year, month, _ := now.Date()
	week := int(now.Weekday())
	monthStart := fmt.Sprintf("%d-%d-01", year, month)
	yearStart := fmt.Sprintf("%d-01-01", year)
	weekStart := now.AddDate(0, 0, -week+1).Format("2006-01-02")

	var tmp stats
	err := func() error {
		if err := s.db.WithContext(ctx).Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS week", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, weekStart).
			Find(&tmp).Error; err != nil {
			return err
		}
		if err := s.db.WithContext(ctx).Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS month", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, monthStart).
			Find(&tmp).Error; err != nil {
			return err
		}
		if err := s.db.WithContext(ctx).Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS year", req.Method, req.Field)).
			Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, yearStart).
			Find(&tmp).Error; err != nil {
			return err
		}
		if err := s.db.WithContext(ctx).Model(&model.ActivityDetail{}).Select(fmt.Sprintf("%s(%s) AS all", req.Method, req.Field)).
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
	err = s.db.WithContext(ctx).Select("freq, value").
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
	r := types.Stats{
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
		r.WeekProcess = fmt.Sprintf("%.0f", tmp.Week/goalMap["week"]*100)
	}
	if int(goalMap["month"]) == 0 {
		r.MonthGoal, r.MonthProcess = notExistsLabel, notExistsLabel
	} else {
		r.MonthProcess = fmt.Sprintf("%.0f", tmp.Month/goalMap["month"]*100)
	}
	if int(goalMap["year"]) == 0 {
		r.YearGoal, r.YearProcess = notExistsLabel, notExistsLabel
	} else {
		r.YearProcess = fmt.Sprintf("%.0f", tmp.Year/goalMap["year"]*100)
	}
	return &r, nil
}

const (
	limitWeek  = 12
	limitMonth = 12
	limitYear  = 12
)

// DateChart 以日期为横轴的统计数据：近一个月，近三个月，近半年，全年 by week || month || year
func (s *Strava) DateChart(ctx context.Context, req *types.DateChartReq, sourceID int64) (echo.Map, error) {
	valMap := make(map[string]float64)
	now := time.Now()
	startDate, start := findStartDate(req, now)

	rows, err := s.db.WithContext(ctx).Model(&model.ActivityDetail{}).
		Select(
			fmt.Sprintf("%s(%s) AS %s, date_trunc('%s', start_date_local) AS %s",
				req.Method, req.Field, req.Field, req.Freq, req.Freq),
		).
		Where("athlete_id = ? AND type = ? AND start_date_local >= ?", sourceID, req.Type, startDate).
		Group(req.Freq).Order(req.Freq).Rows()
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	if err := rows.Err(); err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var t float64
		var date time.Time
		_ = rows.Scan(&t, &date)
		valMap[date.Format("2006-01-02")] = t
	}

	date, value := makeChartData(req, valMap, start, now)
	if len(value) > req.Size {
		value = value[len(value)-req.Size:]
		date = date[len(date)-req.Size:]
	}

	max, min, avg, maxIndex, minIndex := findMarker(value)
	result := echo.Map{
		"value":     value,
		"time":      date,
		"max":       max,
		"min":       min,
		"avg":       avg,
		"max_index": maxIndex,
		"min_index": minIndex,
	}
	return result, nil
}

func findMarker(data []float64) (max, min, avg float64, maxIndex, minIndex int) {
	if len(data) == 0 {
		return
	}
	max, min, maxIndex, minIndex = data[0], -1, 0, -1
	var cnt int
	var sum float64
	for i, item := range data {
		// ignore zero value
		if item == 0 {
			continue
		}
		if item > max {
			max = item
			maxIndex = i
		}
		if item < min || min == -1 {
			min = item
			minIndex = i
		}
		sum += item
		cnt++
	}
	if minIndex == -1 || minIndex == maxIndex {
		minIndex = 0
		min = 0
	}
	if cnt != 0 {
		avg = sum / float64(cnt)
	}

	return max, min, avg, maxIndex, minIndex
}

func findStartDate(req *types.DateChartReq, now time.Time) (string, time.Time) {
	var dateStart string
	var start time.Time
	switch req.Freq {
	case types.Week:
		if req.Size > limitWeek {
			req.Size = limitWeek
		}
		// 向后找到起始日期（星期一）
		start = now.AddDate(0, 0, -7*req.Size)
		week := start.Weekday()
		if week != time.Monday {
			start = start.AddDate(0, 0, int(time.Monday)-int(week)+1)
		}
		year, month, day := start.AddDate(0, -req.Size, 0).Date()
		dateStart = fmt.Sprintf("%d-%d-%d", year, month, day)
	case types.Month:
		if req.Size > limitMonth {
			req.Size = limitMonth
		}
		start = now.AddDate(0, -req.Size, 0)
		year, month, _ := start.Date()
		dateStart = fmt.Sprintf("%d-%d-01", year, month)
	case types.Year:
		if req.Size > limitYear {
			req.Size = limitYear
		}
		start = now.AddDate(-req.Size, 0, 0)
		year, _, _ := start.Date()
		dateStart = fmt.Sprintf("%d-01-01", year)
	}
	start, _ = time.Parse("2006-01-02", dateStart)
	return dateStart, start
}

func makeChartData(req *types.DateChartReq, valMap map[string]float64, start, now time.Time) (date []string, value []float64) {
	var fraction float64 = 1
	if v, ok := fractionMap[req.Field]; ok {
		fraction = v
	}
	for now.After(start) {
		var key string
		switch req.Freq {
		case types.Week:
			key = start.Format("2006-01-02")
			date = append(date, start.Format("01-02"))
			start = start.AddDate(0, 0, 7)
		case types.Month:
			// do not use time.AddDate for month + 1
			key = start.Format("2006-01") + "-01"
			year, month, _ := start.Date()
			date = append(date, start.Format("01"))
			start = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		case types.Year:
			key = start.Format("2006") + "-01-01"
			date = append(date, start.Format("2006"))
			start = start.AddDate(1, 0, 0)
		}
		if v, ok := valMap[key]; ok {
			value = append(value, v/fraction)
		} else {
			value = append(value, 0.0)
		}
	}
	return date, value
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
func (s *Strava) Push(ctx context.Context, event *strava.SubscriptionEvent) error {
	ups, err := json.Marshal(event.Updates)
	if err != nil {
		return em.ErrBadRequest.Wrap(err)
	}

	var push model.PushEvent
	err = s.db.WithContext(ctx).Select("id, status").Where("object_id = ?", event.ObjectID).Find(&push).Error
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
		err = s.db.WithContext(ctx).Where("object_id = ?", event.ObjectID).Create(&newPush).Error
		if err != nil {
			return em.ErrDB.Wrap(err)
		}
	}
	if push.Status != 1 {
		emErr := s.push(ctx, event)
		if emErr != nil {
			log.Error("strava push", zap.Int64("object_id", event.ObjectID), log.Ctx(ctx))
			return emErr
		}
	}

	return nil
}

// push 处理推送
func (s *Strava) push(ctx context.Context, event *strava.SubscriptionEvent) error {
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

func (s *Strava) activityPush(ctx context.Context, event *strava.SubscriptionEvent) error {
	switch event.AspectType {
	case "create":
		return s.activityCreate(ctx, event)
	case "update":
		return em.ErrBadRequest.Msg("unknown aspect type")
	}
	return em.ErrBadRequest.Msg("unknown aspect type")
}

func (s *Strava) athletePush(ctx context.Context, event *strava.SubscriptionEvent) error {
	return nil
}

func (s *Strava) activityCreate(ctx context.Context, event *strava.SubscriptionEvent) error {
	token, err := s.tokenSrv.GetToken(ctx, "strava", event.OwnerID, s.oauth2Conf)
	if err != nil {
		return handler.ErrGetToken.Wrap(err)
	}
	stravaCli := strava.NewClient(s.oauth2Conf.Client(ctx, token))
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

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func formatActivityData(resp *strava.DetailedActivity, streamResp *strava.StreamSet, event *strava.SubscriptionEvent,
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
