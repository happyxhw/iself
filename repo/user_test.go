package repo

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"git.happyxhw.cn/happyxhw/iself/model"
	"git.happyxhw.cn/happyxhw/iself/pkg/mymock"
	"git.happyxhw.cn/happyxhw/iself/pkg/query"
	"git.happyxhw.cn/happyxhw/iself/pkg/util"
)

var mockUser = model.User{
	ID:       1,
	Name:     "mockName",
	Email:    "mockEmail",
	Password: "mockPassword",
	Source:   "strava",
	SourceID: 11,
}

// https://github.com/DATA-DOG/go-sqlmock/issues/118#issuecomment-450974462
// https://github.com/DATA-DOG/go-sqlmock/issues/118#issuecomment-614573409
func TestUser_Create(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `
		INSERT INTO "user" 
		("name","email","password","avatar_url","role","source","source_id","status","deleted_at","created_at","updated_at") 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "created_at","updated_at","id"
	`

	user := model.User{
		Name:      mockUser.Name,
		Email:     mockUser.Email,
		Password:  mockUser.Password,
		Status:    int(model.WaitActiveStatus),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mock.ExpectQuery(sql).
		WithArgs(mockUser.Name, mockUser.Email, mockUser.Password, "", 0, "", 0, 0, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"created_at", "updated_at", "id"}).
				AddRow(time.Now(), time.Now(), 1),
		)

	repo := NewUserRepo(gdb)

	u, err := repo.Create(context.TODO(), &user)

	require.NoError(t, err)

	require.Equal(t, u.ID, int64(1))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_Get(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE id = $1 AND "user"."deleted_at" = $2 LIMIT 1`
	mock.ExpectQuery(sql).
		WithArgs(mockUser.ID, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(mockUser.ID, mockUser.Name, mockUser.Email),
		)

	repo := NewUserRepo(gdb)

	u, err := repo.Get(context.TODO(), mockUser.ID, query.Opt{})

	require.NoError(t, err)

	require.Equal(t, u.Email, mockUser.Email)
	require.Equal(t, u.Name, mockUser.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetByEmail(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE email = $1 AND "user"."deleted_at" = $2 LIMIT 1`
	mock.ExpectQuery(sql).
		WithArgs(mockUser.Email, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(mockUser.ID, mockUser.Name, mockUser.Email),
		)

	repo := NewUserRepo(gdb)

	u, err := repo.GetByEmail(context.TODO(), mockUser.Email, query.Opt{})

	require.NoError(t, err)

	require.Equal(t, u.Email, mockUser.Email)
	require.Equal(t, u.Name, mockUser.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetBySource(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE (source = $1 AND source_id = $2) AND "user"."deleted_at" = $3 LIMIT 1`
	mock.ExpectQuery(sql).
		WithArgs(mockUser.Source, mockUser.SourceID, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(mockUser.ID, mockUser.Name, mockUser.Email),
		)

	repo := NewUserRepo(gdb)

	u, err := repo.GetBySource(context.TODO(), mockUser.Source, mockUser.SourceID, query.Opt{})

	require.NoError(t, err)

	require.Equal(t, u.Email, mockUser.Email)
	require.Equal(t, u.Name, mockUser.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_GetWithFields(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT "id","name","email" FROM "user" WHERE id = $1 AND "user"."deleted_at" = $2 LIMIT 1`
	mock.ExpectQuery(sql).
		WithArgs(mockUser.ID, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(mockUser.ID, mockUser.Name, mockUser.Email),
		)

	repo := NewUserRepo(gdb)

	u, err := repo.Get(context.TODO(), mockUser.ID, query.Opt{
		Fields: []string{"id", "name", "email"},
	})

	require.NoError(t, err)

	require.Equal(t, u.Email, mockUser.Email)
	require.Equal(t, u.Name, mockUser.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_QueryOnlyCount(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT count(*) FROM "user" WHERE "user"."name" = $1 AND "user"."deleted_at" = $2`
	mock.ExpectQuery(sql).
		WithArgs(mockUser.Name, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).
				AddRow(2),
		)
	repo := NewUserRepo(gdb)
	params := model.UserParam{
		Param: query.Param{
			OnlyCount: true,
		},
		Name: &mockUser.Name,
	}
	opt := query.Opt{}

	pr, list, err := repo.Query(context.TODO(), &params, opt)

	require.NoError(t, err)

	require.Equal(t, pr.Total, int64(2))
	require.Equal(t, len(list), 0)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_Query(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql1 := `SELECT count(*) FROM "user" WHERE "user"."name" = $1 AND "user"."status" = $2 AND "user"."deleted_at" = $3`
	mock.ExpectQuery(sql1).
		WithArgs(mockUser.Name, 1, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).
				AddRow(11),
		)
	sql2 := `
		SELECT "id" FROM "user"
		WHERE "user"."name" = $1 AND "user"."status" = $2 AND "user"."deleted_at" = $3
		ORDER BY id DESC LIMIT 10 OFFSET 10
	`
	mock.ExpectQuery(sql2).
		WithArgs(mockUser.Name, 1, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(mockUser.ID),
		)
	sql3 := `
		SELECT "id","name","email" FROM "user" 
		WHERE "user"."name" = $1 AND "user"."status" = $2 AND "user"."deleted_at" = $3 AND id IN ($4)
		ORDER BY id DESC
	`
	mock.ExpectQuery(sql3).
		WithArgs(mockUser.Name, 1, 0, sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(mockUser.ID, mockUser.Name, mockUser.Email),
		)
	repo := NewUserRepo(gdb)
	params := model.UserParam{
		Param: query.Param{
			Page:   2,
			Size:   10,
			SortBy: "-id",
		},
		Name:   &mockUser.Name,
		Status: util.Int(1),
	}
	opt := query.Opt{
		Fields: []string{"id", "name", "email"},
	}

	pr, list, err := repo.Query(context.TODO(), &params, opt)

	require.NoError(t, err)

	require.Equal(t, pr.Total, int64(11))
	require.Equal(t, len(list), 1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_Update(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `UPDATE "user" SET "password"=$1,"updated_at"=$2 WHERE id = $3 AND "user"."deleted_at" = $4`
	mock.ExpectExec(sql).
		WithArgs(mockUser.Password, sqlmock.AnyArg(), mockUser.ID, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewUserRepo(gdb)

	params := model.UserParam{
		Password: util.String(mockUser.Password),
	}

	rows, err := repo.Update(context.TODO(), mockUser.ID, &params)

	require.NoError(t, err)
	require.Equal(t, rows, int64(1))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_UpdateByEmail(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `UPDATE "user" SET "password"=$1,"updated_at"=$2 WHERE email = $3 AND "user"."deleted_at" = $4`
	mock.ExpectExec(sql).
		WithArgs(mockUser.Password, sqlmock.AnyArg(), mockUser.Email, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewUserRepo(gdb)

	params := model.UserParam{
		Password: util.String(mockUser.Password),
	}

	rows, err := repo.UpdateByEmail(context.TODO(), mockUser.Email, &params)

	require.NoError(t, err)
	require.Equal(t, rows, int64(1))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepo_Delete(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `UPDATE "user" SET "deleted_at"=$1 WHERE id = $2 AND "user"."deleted_at" = $3`
	mock.ExpectExec(sql).
		WithArgs(sqlmock.AnyArg(), mockUser.ID, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewUserRepo(gdb)

	rows, err := repo.Delete(context.TODO(), mockUser.ID)

	require.NoError(t, err)
	require.Equal(t, rows, int64(1))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
