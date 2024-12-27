package db

import (
	"context"
	"database/sql"
	"errors"
	repo2 "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
)

func TestRepository_GetFavorites(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, []*model.Favorite, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				rows := sqlmock.NewRows(
					[]string{
						"user_id", "item_id", "id", "title", "src", "alt", "price",
					},
				).
					AddRow(uid.String(), uuid.New().String(), uuid.New().String(), "Item 1", "src1", "alt1", 10.5).
					AddRow(uid.String(), uuid.New().String(), uuid.New().String(), "Item 2", "src2", "alt2", 15.5)

				mock.ExpectQuery(regexp.QuoteMeta(favGetQ)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res []*model.Favorite, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Len(t, res, 2)
				assert.Equal(t, "Item 1", res[0].Item.Title)
				assert.Equal(t, "Item 2", res[1].Item.Title)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(favGetQ)).
					WithArgs(uid).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res []*model.Favorite, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(favGetQ)).
					WithArgs(uid).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res []*model.Favorite, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				rows := sqlmock.NewRows(
					[]string{
						"user_id", "item_id", "id", "title", "src", "alt", "price",
					},
				).
					AddRow(uuid.Nil, uuid.New().String(), uuid.New().String(), "Item 1", "src1", "alt1", 10.5)

				mock.ExpectQuery(regexp.QuoteMeta(favGetQ)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res []*model.Favorite, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Scan")
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.GetFavorites(context.Background(), uid)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_AddToFavorites(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.Favorite, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(favAddQ)).
					WithArgs(uid, itemID).
					WillReturnRows(sqlmock.NewRows([]string{"item_id"}).AddRow(itemID.String()))
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, uid, res.UserID)
				assert.Equal(t, itemID, res.ItemID)
			},
		},
		{
			name: "AlreadyExists",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(favAddQ)).
					WithArgs(uid, itemID).
					WillReturnRows(sqlmock.NewRows([]string{"item_id"}).AddRow(uuid.Nil.String()))
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrAlreadyExists, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(favAddQ)).
					WithArgs(uid, itemID).
					WillReturnError(errors.New("violates foreign key constraint"))
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(favAddQ)).
					WithArgs(uid, itemID).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *model.Favorite, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.AddToFavorites(context.Background(), uid, itemID)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_RemoveFromFavorites(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()
	itemID := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(favDelQ)).
					WithArgs(uid, itemID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(favDelQ)).
					WithArgs(uid, itemID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
			},
		},
		{
			name: "ExecError",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(favDelQ)).
					WithArgs(uid, itemID).
					WillReturnError(errors.New("exec error"))
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "exec error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				err := repo.RemoveFromFavorites(context.Background(), uid, itemID)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}
