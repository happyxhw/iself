package repo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/happyxhw/iself/model"
	"github.com/happyxhw/iself/pkg/mymock"
	"github.com/happyxhw/iself/pkg/query"
)

var mockAct = model.StravaActivityDetail{
	ID: 1,
}

func TestStravaRepo_GetActivity(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "strava_activity_detail" WHERE (id = $1 AND athlete_id = $2) AND "strava_activity_detail"."deleted_at" = $3 LIMIT 1`
	mock.ExpectQuery(sql).
		WithArgs(mockAct.ID, mockUser.ID, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(mockAct.ID),
		)

	repo := NewStravaRepo(gdb)

	_, err := repo.GetDetailedActivity(context.TODO(), mockAct.ID, mockUser.ID, query.Opt{})

	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStravaRepo_GetStream(t *testing.T) {
	gdb, mock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "strava_activity_stream" WHERE id = $1 AND "strava_activity_stream"."deleted_at" = $2 LIMIT 1`
	mock.ExpectQuery(sql).
		WithArgs(mockAct.ID, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(mockAct.ID),
		)

	repo := NewStravaRepo(gdb)

	_, err := repo.GetStreamSet(context.TODO(), mockAct.ID, query.Opt{})

	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
