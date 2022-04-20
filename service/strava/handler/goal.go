package handler

import (
	"context"
	"time"

	"gorm.io/gorm"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
	"git.happyxhw.cn/happyxhw/iself/service/strava/types"
)

type Goal struct {
	db *gorm.DB
}

func NewGoal(db *gorm.DB) *Goal {
	return &Goal{
		db: db,
	}
}

func (g *Goal) Create(ctx context.Context, req *types.GoalReq) error {
	var m model.Goal
	err := g.db.WithContext(ctx).Select("id").
		Where(
			"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
		).
		Find(&m).Limit(1).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if m.ID > 0 {
		return em.ErrConflict.Msg("goal exists")
	}
	now := time.Now()
	m = model.Goal{
		AthleteID: req.SourceID,
		Type:      req.Type,
		Field:     req.Field,
		Freq:      req.Freq,
		Value:     req.Value,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = g.db.WithContext(ctx).Create(&m).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	return nil
}

func (g *Goal) UpdateValue(ctx context.Context, req *types.GoalReq) error {
	var m model.Goal
	err := g.db.WithContext(ctx).Select("id").
		Where(
			"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
		).
		Find(&m).Limit(1).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if m.ID == 0 {
		return em.ErrNotFound.Msg("goal not found")
	}
	err = g.db.WithContext(ctx).Model(&model.Goal{}).
		Where(
			"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
		).Update("value", req.Value).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	return nil
}

func (g *Goal) Delete(ctx context.Context, req *types.GoalReq) error {
	var m model.Goal
	err := g.db.WithContext(ctx).Select("id").
		Where(
			"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
		).
		Find(&m).Limit(1).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	if m.ID == 0 {
		return em.ErrNotFound.Msg("goal not found")
	}
	err = g.db.WithContext(ctx).Where(
		"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
	).Delete(&model.Goal{}).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	return nil
}

func (g *Goal) Query(ctx context.Context, sourceID int64, activityType, field string) ([]*types.QueryGoalResp, error) {
	var goals []*model.Goal
	err := g.db.WithContext(ctx).Select("type, field, freq, value").
		Where("athlete_id = ? AND type = ? AND field = ?", sourceID, activityType, field).
		Find(&goals).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	results := make([]*types.QueryGoalResp, 0, len(goals))
	for _, item := range goals {
		results = append(results, &types.QueryGoalResp{
			Type:  item.Type,
			Field: item.Field,
			Freq:  item.Freq,
			Value: item.Value,
		})
	}
	return results, nil
}
