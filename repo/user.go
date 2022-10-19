package repo

import (
	"context"

	"gorm.io/gorm"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/query"
	"git.happyxhw.cn/happyxhw/iself/pkg/trans"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) Create(ctx context.Context, u *model.User) (*model.User, error) {
	err := trans.DB(ctx, ur.db.WithContext(ctx)).Create(u).Error
	return u, err
}

func (ur *UserRepo) Get(ctx context.Context, id int64, opt query.Opt) (*model.User, error) {
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Where("id = ?", id)
	return ur.get(tx, opt)
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string, opt query.Opt) (*model.User, error) {
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Where("email = ?", email)
	return ur.get(tx, opt)
}

func (ur *UserRepo) GetBySource(ctx context.Context, source string, sourceID int64, opt query.Opt) (*model.User, error) {
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Where("source = ? AND source_id = ?", source, sourceID)
	return ur.get(tx, opt)
}

func (ur *UserRepo) get(tx *gorm.DB, opt query.Opt) (*model.User, error) {
	var r model.User
	if err := query.Take(tx, opt, &r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &r, nil
}

func (ur *UserRepo) Query(ctx context.Context, params *model.UserParam, opt query.Opt) (*query.PagingResult, []*model.User, error) {
	var list []*model.User
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Model(&model.User{}).Where(params)
	if len(opt.Fields) > 0 {
		tx = tx.Select(opt.Fields)
	}
	if params.SortBy != "" {
		if sortBy := query.ParseOrder(params.SortBy, userSortFn); sortBy != "" {
			tx = tx.Order(sortBy)
		}
	}

	pr, err := query.WrapPageQuery(tx, params.Param, &list)
	if err != nil {
		return nil, nil, err
	}

	return pr, list, nil
}

func (ur *UserRepo) Update(ctx context.Context, id int64, params *model.UserParam) (int64, error) {
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Table((&model.User{}).TableName())
	r := tx.Where("id = ?", id).Updates(params)
	return r.RowsAffected, r.Error
}

func (ur *UserRepo) UpdateByEmail(ctx context.Context, email string, params *model.UserParam) (int64, error) {
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Table((&model.User{}).TableName())
	r := tx.Where("email = ?", email).Updates(params)
	return r.RowsAffected, r.Error
}

func (ur *UserRepo) Delete(ctx context.Context, id int64) (int64, error) {
	tx := trans.DB(ctx, ur.db.WithContext(ctx)).Model(&model.User{})
	r := tx.Where("id = ?", id).Delete(&model.User{})
	return r.RowsAffected, r.Error
}

func userSortFn(key string) string {
	k := map[string]bool{
		"id": true,
	}
	if k[key] {
		return key
	}
	return ""
}
