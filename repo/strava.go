package repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/query"
	"git.happyxhw.cn/happyxhw/iself/pkg/trans"
)

type StravaRepo struct {
	db *gorm.DB
}

func NewStravaRepo(db *gorm.DB) *StravaRepo {
	return &StravaRepo{
		db: db,
	}
}

func (sr *StravaRepo) CreateDetailedActivity(ctx context.Context, e *model.StravaActivityDetail) error {
	err := trans.DB(ctx, sr.db.WithContext(ctx)).Create(e).Error

	return err
}

func (sr *StravaRepo) GetDetailedActivity(ctx context.Context, activityID, athleteID int64,
	opt query.Opt) (*model.StravaActivityDetail, error) {
	var r model.StravaActivityDetail
	tx := trans.DB(ctx, sr.db.WithContext(ctx))
	tx = tx.Where("id = ? AND athlete_id = ?", activityID, athleteID)
	if err := query.Take(tx, opt, &r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &r, nil
}

func (sr *StravaRepo) QueryDetailedActivity(ctx context.Context, athleteID int64,
	params *model.StravaActivityParam, opt query.Opt) (*query.PagingResult, []*model.StravaActivityDetail, error) {
	var list []*model.StravaActivityDetail
	db := trans.DB(ctx, sr.db.WithContext(ctx)).Model(&model.StravaActivityDetail{}).Where("athlete_id = ?", athleteID)

	if params.Type != nil {
		db = db.Where("type = ?", *params.Type)
	}
	if params.Filter != "" {
		db = db.Where("name ILIKE ?", "%"+params.Filter+"%")
	}

	if len(opt.Fields) > 0 {
		db = db.Select(opt.Fields)
	}
	if params.SortBy != "" {
		if sortBy := query.ParseOrder(params.SortBy, activitySortFn); sortBy != "" {
			db = db.Order(sortBy)
		}
	}

	pr, err := query.WrapPageQuery(db, params.Param, &list)
	if err != nil {
		return nil, nil, err
	}

	return pr, list, nil
}

func (sr *StravaRepo) GetActivityProgressStats(ctx context.Context, athleteID int64,
	activityType, method, field, start string) (float64, error) {
	result := map[string]interface{}{}
	tx := sr.db.WithContext(ctx).Model(&model.StravaActivityDetail{}).
		Select(fmt.Sprintf("%s(%s) AS value", method, field)).
		Where("athlete_id = ? AND type = ?", athleteID, activityType)
	if start != "" {
		tx = tx.Where("start_date_local >= ?", start)
	}
	err := tx.Take(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	r, ok := result["value"].(float64)
	if ok {
		return r, nil
	}
	r2, _ := result["value"].(int64)

	return float64(r2), nil
}

func (sr *StravaRepo) GetActivityAggStats(ctx context.Context, athleteID int64,
	activityType, method, field, start, freq string) (map[string]float64, error) {
	valMap := make(map[string]float64)
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).Model(&model.StravaActivityDetail{})
	rows, err := tx.
		Select(
			fmt.Sprintf("%s(%s) AS %s, date_trunc('%s', start_date_local) AS %s",
				method, field, field, freq, freq),
		).
		Where("athlete_id = ? AND type = ? AND start_date_local >= ?", athleteID, activityType, start).
		Group(freq).Order(freq).Rows()
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var t float64
		var date time.Time
		_ = rows.Scan(&t, &date)
		valMap[date.Format("2006-01-02")] = t
	}

	return valMap, nil
}

func (sr *StravaRepo) CreateStreamSet(ctx context.Context, e *model.StravaActivityStream) error {
	err := trans.DB(ctx, sr.db.WithContext(ctx)).Create(e).Error

	return err
}

func (sr *StravaRepo) GetStreamSet(ctx context.Context, activityID int64, opt query.Opt) (*model.StravaActivityStream, error) {
	var r model.StravaActivityStream
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).Where("id = ?", activityID)
	if err := query.Take(tx, opt, &r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &r, nil
}

func (sr *StravaRepo) CreateActivityRaw(ctx context.Context, m *model.StravaActivityRaw) error {
	err := trans.DB(ctx, sr.db.WithContext(ctx)).Create(m).Error

	return err
}

func (sr *StravaRepo) CreatePushEvent(ctx context.Context, e *model.StravaPushEvent) error {
	err := trans.DB(ctx, sr.db.WithContext(ctx)).Create(e).Error

	return err
}

func (sr *StravaRepo) GetPushEvent(ctx context.Context, athleteID int64, opt query.Opt) (*model.StravaPushEvent, error) {
	var r model.StravaPushEvent
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).Where("object_id = ?", athleteID)
	if err := query.Take(tx, opt, &r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &r, nil
}

func (sr *StravaRepo) UpdatePushEvent(ctx context.Context, activityID int64, params *model.StravaPushEventParam) (int64, error) {
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).Table((&model.StravaPushEvent{}).TableName())
	r := tx.Where("object_id = ?", activityID).Updates(params)

	return r.RowsAffected, r.Error
}

func (sr *StravaRepo) CreateGoal(ctx context.Context, g *model.StravaGoal) error {
	err := trans.DB(ctx, sr.db.WithContext(ctx)).Create(g).Error

	return err
}

func (sr *StravaRepo) GetGoal(ctx context.Context, g *model.StravaGoal, opt query.Opt) (*model.StravaGoal, error) {
	var r model.StravaGoal
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).
		Where("athlete_id = ?", g.AthleteID).
		Where("type = ? AND field = ? AND freq = ?", g.Type, g.Field, g.Freq)
	if err := query.Take(tx, opt, &r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (sr *StravaRepo) GetGoalByID(ctx context.Context, athleteID, goalID int64, opt query.Opt) (*model.StravaGoal, error) {
	var r model.StravaGoal
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).
		Where("athlete_id = ? AND id = ?", athleteID, goalID)
	if err := query.Take(tx, opt, &r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (sr *StravaRepo) GetAllGoal(ctx context.Context, g *model.StravaGoal, opt query.Opt) ([]*model.StravaGoal, error) {
	var r []*model.StravaGoal
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).
		Where("athlete_id = ?", g.AthleteID).
		Where("type = ? AND field = ?", g.Type, g.Field)
	err := tx.Find(&r).Error

	return r, err
}

func (sr *StravaRepo) UpdateGoal(ctx context.Context, athleteID, goalID int64, params *model.StravaGoalParam) (int64, error) {
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).Table((&model.StravaGoal{}).TableName())
	r := tx.Where("athlete_id = ? AND id = ?", athleteID, goalID).Updates(params)

	return r.RowsAffected, r.Error
}

func (sr *StravaRepo) DeleteGoal(ctx context.Context, athleteID, goalID int64) (int64, error) {
	tx := trans.DB(ctx, sr.db.WithContext(ctx)).Table((&model.StravaGoal{}).TableName())
	r := tx.Where("athlete_id = ? AND id = ?", athleteID, goalID).Delete(&model.StravaGoal{})

	return r.RowsAffected, r.Error
}

func activitySortFn(key string) string {
	k := map[string]bool{
		"id":               true,
		"start_date_local": true,
	}
	if k[key] {
		return key
	}
	return ""
}
