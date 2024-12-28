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
	"time"
)

var (
	page               = 1
	size               = 10
	expectedCount      = int64(2)
	expectedTotalPages = int((expectedCount + int64(size) - 1) / int64(size))
)

func TestRepository_PromotionSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	query := "test promo"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.PaginatedPromosData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoSearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"slug", "title", "description", "src", "alt", "lasts_to",
					},
				).
					AddRow("promo-1", "Test Promo 1", "Description 1", "src1", "alt1", time.Now()).
					AddRow("promo-2", "Test Promo 2", "Description 2", "src2", "alt2", time.Now())

				mock.ExpectQuery(regexp.QuoteMeta(promoSearchQ)).
					WithArgs("%"+query+"%", "%"+query+"%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Test Promo 1", res.Data[0].Title)
				assert.Equal(t, "Test Promo 2", res.Data[1].Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoSearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoSearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(promoSearchQ)).
					WithArgs("%"+query+"%", "%"+query+"%", (page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoSearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"slug", "title", "description", "src", "alt", "lasts_to",
					},
				).
					AddRow(nil, "Test Promo 1", "Description 1", "src1", "alt1", "2024-12-31")

				mock.ExpectQuery(regexp.QuoteMeta(promoSearchQ)).
					WithArgs("%"+query+"%", "%"+query+"%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
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
				res, err := repo.PromotionSearch(context.Background(), query, page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ListPromotionItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	slug := "test-promo"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.PaginatedPromoItemsData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountItemsQ)).
					WithArgs(slug).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"discount", "promotion_slug", "item_id", "title", "price", "src", "alt",
					},
				).
					AddRow(10, "test-promo", uuid.New().String(), "Item 1", 10.5, "src1", "alt1").
					AddRow(20, "test-promo", uuid.New().String(), "Item 2", 15.5, "src2", "alt2")

				mock.ExpectQuery(regexp.QuoteMeta(promoItemsQ)).
					WithArgs(slug, (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromoItemsData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Item 1", res.Data[0].Item.Title)
				assert.Equal(t, "Item 2", res.Data[1].Item.Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountItemsQ)).
					WithArgs(slug).
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromoItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountItemsQ)).
					WithArgs(slug).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(promoItemsQ)).
					WithArgs(slug, (page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromoItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountItemsQ)).
					WithArgs(slug).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"discount", "promotion_slug", "item_id", "title", "price", "src", "alt",
					},
				).
					AddRow(nil, "test-promo", uuid.New().String(), "Item 1", 10.5, "src1", "alt1")

				mock.ExpectQuery(regexp.QuoteMeta(promoItemsQ)).
					WithArgs(slug, (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromoItemsData, err error) {
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
				res, err := repo.ListPromotionItems(context.Background(), slug, page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_GetPromotion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	slug := "test-promo"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.Promotion, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoGetQ)).
					WithArgs(slug).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"slug", "title", "description", "src", "alt", "lasts_to",
							},
						).AddRow(
							"test-promo",
							"Test Promotion",
							"Description of test promotion",
							"src1",
							"alt1",
							time.Now(),
						),
					)
			},
			expectedResp: func(t *testing.T, res *model.Promotion, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, "Test Promotion", res.Title)
				assert.Equal(t, "Description of test promotion", res.Description)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoGetQ)).
					WithArgs(slug).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res *model.Promotion, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoGetQ)).
					WithArgs(slug).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *model.Promotion, err error) {
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
						"slug", "title", "description", "src", "alt", "lasts_to",
					},
				).AddRow(nil, "Test Promotion", "Description of test promotion", "src1", "alt1", "2024-12-31")

				mock.ExpectQuery(regexp.QuoteMeta(promoGetQ)).
					WithArgs(slug).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.Promotion, err error) {
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
				res, err := repo.GetPromotion(context.Background(), slug)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ListPromotions(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.PaginatedPromosData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"slug", "title", "description", "src", "alt", "lasts_to",
					},
				).
					AddRow("promo-1", "Promo 1", "Description 1", "src1", "alt1", time.Now()).
					AddRow("promo-2", "Promo 2", "Description 2", "src2", "alt2", time.Now())

				mock.ExpectQuery(regexp.QuoteMeta(promoListQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Promo 1", res.Data[0].Title)
				assert.Equal(t, "Promo 2", res.Data[1].Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountQ)).
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(promoListQ)).
					WithArgs((page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(promoCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"slug", "title", "description", "src", "alt", "lasts_to",
					},
				).
					AddRow(nil, "Promo 1", "Description 1", "src1", "alt1", "2024-12-31")

				mock.ExpectQuery(regexp.QuoteMeta(promoListQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.PaginatedPromosData, err error) {
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
				res, err := repo.ListPromotions(context.Background(), page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_CreatePromotion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}

	promotion := &model.Promotion{
		Slug:        "promo-test",
		Title:       "Promotion Test",
		Description: "Description of promotion test",
		Src:         "src-test",
		Alt:         "alt-test",
		LastsTo:     time.Now(),
		PromotionItems: []*model.PromotionItem{
			{Discount: 10, ItemID: uuid.New()},
			{Discount: 20, ItemID: uuid.New()},
		},
	}

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, string, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(promoCreateQ)).
					WithArgs(
						promotion.Slug,
						promotion.Title,
						promotion.Description,
						promotion.Src,
						promotion.Alt,
						promotion.LastsTo,
					).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(promotion.Slug))

				for _, item := range promotion.PromotionItems {
					mock.ExpectExec(regexp.QuoteMeta(promoItemCreateQ)).
						WithArgs(item.Discount, promotion.Slug, item.ItemID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				require.NoError(t, err)
				assert.Equal(t, promotion.Slug, slug)
			},
		},
		{
			name: "BeginError",
			mockExpect: func() {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "begin error", err.Error())
				assert.Equal(t, "", slug)
			},
		},
		{
			name: "CreateError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(promoCreateQ)).
					WithArgs(
						promotion.Slug,
						promotion.Title,
						promotion.Description,
						promotion.Src,
						promotion.Alt,
						promotion.LastsTo,
					).
					WillReturnError(errors.New("create error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "create error", err.Error())
				assert.Equal(t, "", slug)
			},
		},
		{
			name: "CommitError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(promoCreateQ)).
					WithArgs(
						promotion.Slug,
						promotion.Title,
						promotion.Description,
						promotion.Src,
						promotion.Alt,
						promotion.LastsTo,
					).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(promotion.Slug))

				for _, item := range promotion.PromotionItems {
					mock.ExpectExec(regexp.QuoteMeta(promoItemCreateQ)).
						WithArgs(item.Discount, promotion.Slug, item.ItemID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "commit error", err.Error())
				assert.Equal(t, "", slug)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				slug, err := repo.CreatePromotion(context.Background(), promotion)
				tt.expectedResp(t, slug, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_UpdatePromotion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}

	promotion := &model.Promotion{
		Slug:        "promo-test",
		Title:       "Updated Promotion Test",
		Description: "Updated description of promotion test",
		Src:         "updated-src-test",
		Alt:         "updated-alt-test",
		LastsTo:     time.Now(),
		PromotionItems: []*model.PromotionItem{
			{Discount: 15, ItemID: uuid.New()},
			{Discount: 25, ItemID: uuid.New()},
		},
	}

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(promoUpdateQ)).
					WithArgs(
						promotion.Title,
						promotion.Description,
						promotion.Src,
						promotion.Alt,
						promotion.LastsTo,
						promotion.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows([]string{"discount", "promotion_slug", "item_id"}).
					AddRow(
						promotion.PromotionItems[0].Discount,
						promotion.Slug,
						promotion.PromotionItems[0].ItemID.String(),
					).
					AddRow(
						promotion.PromotionItems[1].Discount,
						promotion.Slug,
						promotion.PromotionItems[1].ItemID.String(),
					)

				mock.ExpectQuery(regexp.QuoteMeta(promoItemListQ)).
					WithArgs(promotion.Slug).
					WillReturnRows(rows)

				for _, item := range promotion.PromotionItems {
					mock.ExpectExec(regexp.QuoteMeta(promoItemUpdateQ)).
						WithArgs(item.Discount, promotion.Slug, item.ItemID, promotion.Slug, item.ItemID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "BeginError",
			mockExpect: func() {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "begin error", err.Error())
			},
		},
		{
			name: "UpdateError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(promoUpdateQ)).
					WithArgs(
						promotion.Title,
						promotion.Description,
						promotion.Src,
						promotion.Alt,
						promotion.LastsTo,
						promotion.Slug,
					).
					WillReturnError(errors.New("update error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "update error", err.Error())
			},
		},
		{
			name: "CommitError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(promoUpdateQ)).
					WithArgs(
						promotion.Title,
						promotion.Description,
						promotion.Src,
						promotion.Alt,
						promotion.LastsTo,
						promotion.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows([]string{"discount", "promotion_slug", "item_id"}).
					AddRow(
						promotion.PromotionItems[0].Discount,
						promotion.Slug,
						promotion.PromotionItems[0].ItemID.String(),
					).
					AddRow(
						promotion.PromotionItems[1].Discount,
						promotion.Slug,
						promotion.PromotionItems[1].ItemID.String(),
					)

				mock.ExpectQuery(regexp.QuoteMeta(promoItemListQ)).
					WithArgs(promotion.Slug).
					WillReturnRows(rows)

				for _, item := range promotion.PromotionItems {
					mock.ExpectExec(regexp.QuoteMeta(promoItemUpdateQ)).
						WithArgs(item.Discount, promotion.Slug, item.ItemID, promotion.Slug, item.ItemID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "commit error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				err := repo.UpdatePromotion(context.Background(), promotion.Slug, promotion)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_DeletePromotion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	slug := "promo-test"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(promoDeleteQ)).
					WithArgs(slug).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(promoDeleteQ)).
					WithArgs(slug).
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
				mock.ExpectExec(regexp.QuoteMeta(promoDeleteQ)).
					WithArgs(slug).
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
				err := repo.DeletePromotion(context.Background(), slug)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}
