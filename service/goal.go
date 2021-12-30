package service

import (
	"time"

	"gorm.io/gorm"

	"git.happyxhw.cn/happyxhw/iself/api"
	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/em"
)

type Goal struct {
	db *gorm.DB
}

func NewGoal(db *gorm.DB) *Goal {
	return &Goal{
		db: db,
	}
}

func (g *Goal) Create(req *api.GoalReq) *em.Error {
	var m model.Goal
	err := g.db.Select("id").
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
	err = g.db.Create(&m).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	return nil
}

func (g *Goal) UpdateValue(req *api.GoalReq) *em.Error {
	var m model.Goal
	err := g.db.Select("id").
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
	err = g.db.Model(&model.Goal{}).
		Where(
			"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
		).Update("value", req.Value).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	return nil
}

func (g *Goal) Delete(req *api.GoalReq) *em.Error {
	var m model.Goal
	err := g.db.Select("id").
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
	err = g.db.Where(
		"athlete_id = ? AND type = ? AND field = ? AND freq = ?", req.SourceID, req.Type, req.Field, req.Freq,
	).Delete(&model.Goal{}).Error
	if err != nil {
		return em.ErrDB.Wrap(err)
	}
	return nil
}

func (g *Goal) Query(sourceID int64, activityType, field string) ([]*api.QueryGoalResp, *em.Error) {
	var goals []*model.Goal
	err := g.db.Select("type, field, freq, value").
		Where("athlete_id = ? AND type = ? AND field = ?", sourceID, activityType, field).
		Find(&goals).Error
	if err != nil {
		return nil, em.ErrDB.Wrap(err)
	}
	results := make([]*api.QueryGoalResp, 0, len(goals))
	for _, item := range goals {
		results = append(results, &api.QueryGoalResp{
			Type:  item.Type,
			Field: item.Field,
			Freq:  item.Freq,
			Value: item.Value,
		})
	}
	return results, nil
}
