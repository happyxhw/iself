package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/ex"
	"git.happyxhw.cn/happyxhw/iself/pkg/log"
	"git.happyxhw.cn/happyxhw/iself/pkg/oauth2x"
	"git.happyxhw.cn/happyxhw/iself/pkg/query"
	"git.happyxhw.cn/happyxhw/iself/pkg/strava"
	"git.happyxhw.cn/happyxhw/iself/pkg/trans"
	"git.happyxhw.cn/happyxhw/iself/pkg/util"
	"git.happyxhw.cn/happyxhw/iself/repo"
	"git.happyxhw.cn/happyxhw/iself/service/strava/types"
)

type Strava struct {
	sr        *repo.StravaRepo
	tr        *repo.TokenRepo
	transRepo *trans.Trans

	auth oauth2x.Oauth2x
}

func NewStrava(sr *repo.StravaRepo, tr *repo.TokenRepo, transRepo *trans.Trans, auth oauth2x.Oauth2x) *Strava {
	return &Strava{
		sr:        sr,
		tr:        tr,
		auth:      auth,
		transRepo: transRepo,
	}
}

func (s *Strava) GetActivity(ctx context.Context, activityID, athleteID int64) (*types.Activity, error) {
	detailed, err := s.sr.GetDetailedActivity(ctx, activityID, athleteID, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	if detailed == nil {
		return nil, ex.ErrNotFound.Msg("activity not found")
	}
	streamSet, err := s.sr.GetStreamSet(ctx, detailed.ID, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	return &types.Activity{
		DetailedActivity: types.NewDetailedActivity(detailed),
		StreamSet:        types.NewStreamSet(streamSet),
	}, nil
}

func (s *Strava) ListActivity(ctx context.Context, athleteID int64, req *model.StravaActivityParam) (*types.ActivityQueryResult, error) {
	p, r, err := s.sr.QueryDetailedActivity(ctx, athleteID, req, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	for i, item := range r {
		r[i].AverageSpeed = transformVelocity(item.AverageSpeed, item.Type)
		r[i].MaxSpeed = transformVelocity(item.MaxSpeed, item.Type)
	}

	return &types.ActivityQueryResult{PageResult: p, Data: types.NewDetailedActivityList(r)}, nil
}

func (s *Strava) GetProgressStats(ctx context.Context, athleteID int64,
	req *types.ProgressStatsReq) (*types.ActivityProgressStats, error) {
	now := time.Now()
	year, month, _ := now.Date()
	week := int(now.Weekday())
	monthStart := fmt.Sprintf("%d-%d-01", year, month)
	yearStart := fmt.Sprintf("%d-01-01", year)
	weekStart := now.AddDate(0, 0, -week+1).Format("2006-01-02")
	var weekVal, monthVal, yearVal, allVal float64
	err := func() error {
		var dbErr error
		if weekVal, dbErr = s.sr.GetActivityProgressStats(ctx, athleteID, req.Type, req.Method, req.Field, weekStart); dbErr != nil {
			return dbErr
		}
		if monthVal, dbErr = s.sr.GetActivityProgressStats(ctx, athleteID, req.Type, req.Method, req.Field, monthStart); dbErr != nil {
			return dbErr
		}
		if yearVal, dbErr = s.sr.GetActivityProgressStats(ctx, athleteID, req.Type, req.Method, req.Field, yearStart); dbErr != nil {
			return dbErr
		}
		if allVal, dbErr = s.sr.GetActivityProgressStats(ctx, athleteID, req.Type, req.Method, req.Field, ""); dbErr != nil {
			return dbErr
		}
		return nil
	}()
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	g := model.StravaGoal{
		AthleteID: athleteID,
		Type:      req.Type,
		Field:     req.Field,
	}
	goals, err := s.sr.GetAllGoal(ctx, &g, query.Fields("freq, value"))
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}
	goalMap := make(map[string]float64, len(goals))
	for _, item := range goals {
		goalMap[item.Freq] = item.Value
	}
	var fraction float64 = 1
	if v, ok := fractionMap[req.Field]; ok {
		fraction = v
	}
	r := types.ActivityProgressStats{
		Type: req.Type,
		Unit: unitMap[req.Field],
		All:  fmt.Sprintf("%.0f", allVal/fraction),

		Week:     fmt.Sprintf("%.0f", weekVal/fraction),
		WeekGoal: fmt.Sprintf("%.0f", goalMap["week"]/fraction),

		Month:     fmt.Sprintf("%.0f", monthVal/fraction),
		MonthGoal: fmt.Sprintf("%.0f", goalMap["month"]/fraction),

		Year:     fmt.Sprintf("%.0f", yearVal/fraction),
		YearGoal: fmt.Sprintf("%.0f", goalMap["year"]/fraction),
	}
	if int(goalMap["week"]) == 0 {
		r.WeekGoal, r.WeekProcess = notExistsLabel, notExistsLabel
	} else {
		r.WeekProcess = fmt.Sprintf("%.0f", weekVal/goalMap["week"]*100)
	}
	if int(goalMap["month"]) == 0 {
		r.MonthGoal, r.MonthProcess = notExistsLabel, notExistsLabel
	} else {
		r.MonthProcess = fmt.Sprintf("%.0f", monthVal/goalMap["month"]*100)
	}
	if int(goalMap["year"]) == 0 {
		r.YearGoal, r.YearProcess = notExistsLabel, notExistsLabel
	} else {
		r.YearProcess = fmt.Sprintf("%.0f", yearVal/goalMap["year"]*100)
	}

	return &r, nil
}

// GetAggStats 以日期为横轴的统计数据：近一个月，近三个月，近半年，全年 by week || month || year
func (s *Strava) GetAggStats(ctx context.Context, athleteID int64, req *types.AggStatsReq) (*types.ActivityAggStats, error) {
	now := time.Now()
	startDate, start := findStartDate(req, now)

	val, err := s.sr.GetActivityAggStats(ctx, athleteID, req.Type, req.Method, req.Field, startDate, req.Freq)
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}

	date, value := makeChartData(req, val, start, now)
	if len(value) > req.Size {
		value = value[len(value)-req.Size:]
		date = date[len(date)-req.Size:]
	}

	max, min, avg, maxIndex, minIndex := findMarker(value)

	r := types.ActivityAggStats{
		Value:    value,
		Time:     date,
		Max:      max,
		Min:      min,
		Avg:      avg,
		MaxIndex: maxIndex,
		MinIndex: minIndex,
	}

	return &r, nil
}

func (s *Strava) Push(ctx context.Context, event *strava.SubscriptionEvent) error {
	ups, err := json.Marshal(event.Updates)
	if err != nil {
		return ex.ErrBadRequest.Wrap(err)
	}

	data, err := s.sr.GetPushEvent(ctx, event.ObjectID, query.Fields("id", "status"))
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	if data == nil {
		newPush := model.StravaPushEvent{
			OwnerID:    event.OwnerID,
			AspectType: event.AspectType,
			EventTime:  event.EventTime,
			ObjectID:   event.ObjectID,
			ObjectType: event.ObjectType,
			Updates:    string(ups),
		}
		if txErr := s.sr.CreatePushEvent(ctx, &newPush); txErr != nil {
			return ex.ErrDB.Wrap(err)
		}
	}

	if data == nil || data.Status != int(model.EventProcessedStatus) {
		err = s.push(ctx, event)
		if err != nil {
			log.Error("strava push", zap.Int64("object_id", event.ObjectID), log.Ctx(ctx))
			return err
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
	return ex.ErrBadRequest.Msg("unknown object type")
}

func (s *Strava) activityPush(ctx context.Context, event *strava.SubscriptionEvent) error {
	switch event.AspectType {
	case "create":
		return s.activityCreate(ctx, event)
	case "update":
		return ex.ErrBadRequest.Msg("unknown aspect type")
	}
	return ex.ErrBadRequest.Msg("unknown aspect type")
}

func (s *Strava) athletePush(ctx context.Context, event *strava.SubscriptionEvent) error {
	return nil
}

func (s *Strava) activityCreate(ctx context.Context, event *strava.SubscriptionEvent) error {
	token, err := s.tr.GetToken(ctx, "strava", event.OwnerID, s.auth.Refresh)
	if err != nil {
		return ErrStravaToken.Wrap(err)
	}
	stravaCli := strava.NewClient(s.auth.Client(ctx, token))
	activityData, body, err := stravaCli.Activity.Activity(ctx, event.ObjectID)
	if err != nil {
		return ErrStravaAPI.Wrap(err)
	}
	activityRawData := model.StravaActivityRaw{
		ID:   event.ObjectID,
		Data: string(body),
	}

	streamSet, err := stravaCli.Activity.ActivityStream(ctx, event.ObjectID)
	if err != nil {
		return ErrStravaAPI.Wrap(err)
	}

	detailedActivityData := model.StravaActivityDetail{
		ID:                 event.ObjectID,
		AthleteID:          activityData.Athlete.ID,
		Name:               activityData.Name,
		Type:               activityData.Type,
		Distance:           activityData.Distance,
		MovingTime:         activityData.MovingTime,
		ElapsedTime:        activityData.ElapsedTime,
		TotalElevationGain: activityData.TotalElevationGain,
		StartDateLocal:     activityData.StartDateLocal,

		AverageSpeed:     activityData.AverageSpeed,
		MaxSpeed:         activityData.MaxSpeed,
		AverageHeartrate: activityData.AverageHeartrate,
		MaxHeartrate:     activityData.MaxHeartrate,
		ElevHigh:         activityData.ElevHigh,
		ElevLow:          activityData.ElevLow,
		Calories:         activityData.Calories,
		DeviceName:       activityData.DeviceName,
		SplitsMetricJSON: activityData.SplitsMetric,
		BestEffortsJSON:  activityData.BestEfforts,
	}
	if activityData.Map != nil {
		detailedActivityData.SummaryPolyline = activityData.Map.SummaryPolyline
		detailedActivityData.Polyline = activityData.Map.Polyline
	}

	streamData := model.StravaActivityStream{
		ID:                   detailedActivityData.ID,
		TimeStream:           streamSet.Time,
		DistanceStream:       streamSet.Distance,
		LatlngStream:         streamSet.Latlng,
		AltitudeStream:       streamSet.Altitude,
		VelocitySmoothStream: streamSet.VelocitySmooth,
		HeartrateStream:      streamSet.Heartrate,
		CadenceStream:        streamSet.Cadence,
		WattsStream:          streamSet.Watts,
		TempStream:           streamSet.Temp,
		MovingStream:         streamSet.Moving,
		GradeSmoothStream:    streamSet.GradeSmooth,
	}

	err = s.transRepo.Exec(ctx, func(ctx context.Context) error {
		if txErr := s.sr.CreateActivityRaw(ctx, &activityRawData); txErr != nil {
			return txErr
		}
		if txErr := s.sr.CreateDetailedActivity(ctx, &detailedActivityData); txErr != nil {
			return txErr
		}
		if txErr := s.sr.CreateStreamSet(ctx, &streamData); txErr != nil {
			return txErr
		}
		updateParams := &model.StravaPushEventParam{Status: util.Int(int(model.EventProcessedStatus))}
		if _, txErr := s.sr.UpdatePushEvent(ctx, event.ObjectID, updateParams); txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}

	return nil
}

func (s *Strava) CreateGoal(ctx context.Context, athleteID int64, req *types.CreateGoalReq) error {
	param := model.StravaGoal{
		AthleteID: athleteID,
		Type:      req.Type,
		Field:     req.Field,
		Freq:      req.Field,
	}
	g, err := s.sr.GetGoal(ctx, &param, query.Fields("id"))
	if err != nil {
		return err
	}
	if g != nil {
		return ex.ErrConflict.Msg("goal exists")
	}

	g = &model.StravaGoal{
		Type:      req.Type,
		Field:     req.Field,
		Freq:      req.Freq,
		Value:     req.Value,
		AthleteID: athleteID,
	}
	err = s.sr.CreateGoal(ctx, g)
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	return nil
}

func (s *Strava) UpdateGoal(ctx context.Context, athleteID int64, req *types.UpdateGoalReq) error {
	g, err := s.sr.GetGoalByID(ctx, athleteID, req.ID, query.Fields("id"))
	if err != nil {
		return err
	}
	if g == nil {
		return nil
	}
	_, err = s.sr.UpdateGoal(ctx, athleteID, req.ID, &model.StravaGoalParam{Value: &req.Value})
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}
	return nil
}

func (s *Strava) DeleteGoal(ctx context.Context, athleteID, goalID int64) error {
	g, err := s.sr.GetGoalByID(ctx, athleteID, goalID, query.Fields("id"))
	if err != nil {
		return err
	}
	if g == nil {
		return nil
	}
	_, err = s.sr.DeleteGoal(ctx, athleteID, goalID)
	if err != nil {
		return ex.ErrDB.Wrap(err)
	}

	return nil
}

func (s *Strava) GetGoals(ctx context.Context, athleteID int64, activityType, field string) ([]*types.Goal, error) {
	param := model.StravaGoal{
		AthleteID: athleteID,
		Type:      activityType,
		Field:     field,
	}
	gs, err := s.sr.GetAllGoal(ctx, &param, query.Opt{})
	if err != nil {
		return nil, ex.ErrDB.Wrap(err)
	}

	return types.NewGoals(gs), nil
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

func findStartDate(req *types.AggStatsReq, now time.Time) (string, time.Time) {
	var dateStart string
	var start time.Time
	switch req.Freq {
	case Week:
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
	case Month:
		if req.Size > limitMonth {
			req.Size = limitMonth
		}
		start = now.AddDate(0, -req.Size, 0)
		year, month, _ := start.Date()
		dateStart = fmt.Sprintf("%d-%d-01", year, month)
	case Year:
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

func makeChartData(req *types.AggStatsReq, valMap map[string]float64, start, now time.Time) (date []string, value []float64) {
	var fraction float64 = 1
	if v, ok := fractionMap[req.Field]; ok {
		fraction = v
	}
	for now.After(start) {
		var key string
		switch req.Freq {
		case Week:
			key = start.Format("2006-01-02")
			date = append(date, start.Format("01-02"))
			start = start.AddDate(0, 0, 7)
		case Month:
			// do not use time.AddDate for month + 1
			key = start.Format("2006-01") + "-01"
			year, month, _ := start.Date()
			date = append(date, start.Format("01"))
			start = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		case Year:
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
