package db

import (
	"context"
	"database/sql"
	"errors"
	repo2 "github.com/JMURv/par-pro/products/internal/repo"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
	"time"
)

func TestRepository_ItemAttrSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	query := "test"
	page := 1
	size := 10
	expectedCount := int64(2)
	expectedTotalPages := int((expectedCount + int64(size) - 1) / int64(size))

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *md.PaginatedItemAttrData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountAttrsQ)).
					WithArgs("%" + query + "%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"id", "name", "value"}).
					AddRow(1, "Attr 1", "Value 1").
					AddRow(2, "Attr 2", "Value 2")

				mock.ExpectQuery(regexp.QuoteMeta(itemSearchAttrQ)).
					WithArgs("%"+query+"%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemAttrData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Attr 1", res.Data[0].Name)
				assert.Equal(t, "Attr 2", res.Data[1].Name)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountAttrsQ)).
					WithArgs("%" + query + "%").
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemAttrData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "SearchError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountAttrsQ)).
					WithArgs("%" + query + "%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(itemSearchAttrQ)).
					WithArgs("%"+query+"%", (page-1)*size, size).
					WillReturnError(errors.New("search error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemAttrData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "search error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountAttrsQ)).
					WithArgs("%" + query + "%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"id", "name", "value"}).
					AddRow(1, "Attr 1", "Value 1").
					AddRow(nil, "Attr 2", "Value 2")

				mock.ExpectQuery(regexp.QuoteMeta(itemSearchAttrQ)).
					WithArgs("%"+query+"%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemAttrData, err error) {
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
				res, err := repo.ItemAttrSearch(context.Background(), query, size, page)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ItemSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	query := "test item"
	page := 1
	size := 10
	expectedCount := int64(2)
	expectedTotalPages := int((expectedCount + int64(size) - 1) / int64(size))
	countQ := "SELECT COUNT(*) FROM item WHERE title ILIKE ? AND title ILIKE ?"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *md.PaginatedItemsData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(countQ)).
					WithArgs("%test%", "%item%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "article", "description", "price", "src", "alt",
						"in_stock", "is_hit", "is_rec", "parent_id", "created_at", "updated_at", "categories",
					},
				).
					AddRow(
						uuid.New().String(),
						"Test Item 1",
						"Article 1",
						"Description 1",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						"{Category 1|slug1}",
					).
					AddRow(
						uuid.New().String(),
						"Test Item 2",
						"Article 2",
						"Description 2",
						15.5,
						"src2",
						"alt2",
						false,
						true,
						true,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						"{Category 2|slug2}",
					)

				mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
					WithArgs("%test%", "%item%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Test Item 1", res.Data[0].Title)
				assert.Equal(t, "Test Item 2", res.Data[1].Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(countQ)).
					WithArgs("%test%", "%item%").
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(countQ)).
					WithArgs("%test%", "%item%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery("SELECT").
					WithArgs("%test%", "%item%", (page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(countQ)).
					WithArgs("%test%", "%item%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "article", "description", "price", "src", "alt",
						"in_stock", "is_hit", "is_rec", "parent_id", "created_at", "updated_at", "categories",
					},
				).
					AddRow(
						1,
						"Test Item 1",
						"Article 1",
						"Description 1",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						0,
						"2024-12-25 18:00:00",
						"2024-12-25 18:00:00",
						pq.Array([]string{"Category 1|slug1"}),
					).
					AddRow(
						nil,
						"Test Item 2",
						"Article 2",
						"Description 2",
						15.5,
						"src2",
						"alt2",
						false,
						true,
						true,
						0,
						"2024-12-25 18:00:00",
						"2024-12-25 18:00:00",
						pq.Array([]string{"Category 2|slug2"}),
					)

				mock.ExpectQuery("SELECT").
					WithArgs("%test%", "%item%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
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
				res, err := repo.ItemSearch(context.Background(), query, page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ListItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	page := 1
	size := 10
	expectedCount := int64(2)
	expectedTotalPages := int((expectedCount + int64(size) - 1) / int64(size))

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *md.PaginatedItemsData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "article", "description", "price", "src", "alt",
						"in_stock", "is_hit", "is_rec", "parent_id", "created_at", "updated_at", "categories",
					},
				).
					AddRow(
						uuid.New().String(),
						"Test Item 1",
						"Article 1",
						"Description 1",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						"{Category 1|slug1}",
					).
					AddRow(
						uuid.New().String(),
						"Test Item 2",
						"Article 2",
						"Description 2",
						15.5,
						"src2",
						"alt2",
						false,
						true,
						true,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						"{Category 2|slug2}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemListQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Test Item 1", res.Data[0].Title)
				assert.Equal(t, "Test Item 2", res.Data[1].Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountQ)).
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(itemListQ)).
					WithArgs((page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "article", "description", "price", "src", "alt",
						"in_stock", "is_hit", "is_rec", "parent_id", "created_at", "updated_at", "categories",
					},
				).
					AddRow(
						nil,
						"Test Item 1",
						"Article 1",
						"Description 1",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						pq.Array([]string{"Category 1|slug1"}),
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemListQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
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
				res, err := repo.ListItems(context.Background(), page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_GetItemByUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *md.Item, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				rows := sqlmock.NewRows(
					[]string{
						"id",
						"title",
						"article",
						"description",
						"price",
						"src",
						"alt",
						"in_stock",
						"is_hit",
						"is_rec",
						"parent_id",
						"created_at",
						"updated_at",
						"media",
						"attrs",
						"categories",
					},
				).
					AddRow(
						uid.String(),
						"Test Item",
						"Article",
						"Description",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						"{src1|alt1}",
						"{attr1|value1}",
						"{Category 1|slug1}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemGetByIDQ)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.Item, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, "Test Item", res.Title)
				assert.Len(t, res.Media, 1)
				assert.Len(t, res.Attributes, 1)
				assert.Len(t, res.Categories, 1)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemGetByIDQ)).
					WithArgs(uid).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res *md.Item, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemGetByIDQ)).
					WithArgs(uid).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *md.Item, err error) {
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
						"id",
						"title",
						"article",
						"description",
						"price",
						"src",
						"alt",
						"in_stock",
						"is_hit",
						"is_rec",
						"parent_id",
						"created_at",
						"updated_at",
						"media",
						"attrs",
						"categories",
					},
				).
					AddRow(
						nil,
						"Test Item",
						"Article",
						"Description",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						time.Now(),
						time.Now(),
						pq.Array([]string{"src1|alt1"}),
						pq.Array([]string{"attr1|value1"}),
						pq.Array([]string{"Category 1|slug1"}),
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemGetByIDQ)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.Item, err error) {
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
				res, err := repo.GetItemByUUID(context.Background(), uid)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_CreateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}

	item := &md.Item{
		Title:        "Test Item",
		Article:      "Article 1",
		Description:  "Description 1",
		Price:        10.5,
		Src:          "src1",
		Alt:          "alt1",
		InStock:      true,
		IsHit:        false,
		IsRec:        false,
		ParentItemID: uuid.Nil,
		Media: []md.ItemMedia{
			{Src: "src1", Alt: "alt1"},
		},
		Attributes: []md.ItemAttribute{
			{Name: "attr1", Value: "value1"},
		},
		Categories: []md.Category{
			{Slug: "slug1"},
		},
	}

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, uuid.UUID, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(itemCreateQ)).
					WithArgs(
						item.Title,
						item.Article,
						item.Description,
						item.Price,
						item.Src,
						item.Alt,
						item.InStock,
						item.IsHit,
						item.IsRec,
						item.ParentItemID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))

				mock.ExpectExec(regexp.QuoteMeta(itemMediaCreateQ)).
					WithArgs(sqlmock.AnyArg(), item.Media[0].Src, item.Media[0].Alt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(itemAttrCreateQ)).
					WithArgs(sqlmock.AnyArg(), item.Attributes[0].Name, item.Attributes[0].Value).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(itemCategoryCreateQ)).
					WithArgs(sqlmock.AnyArg(), item.Categories[0].Slug).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, id uuid.UUID, err error) {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, id)
			},
		},
		{
			name: "BeginError",
			mockExpect: func() {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			expectedResp: func(t *testing.T, id uuid.UUID, err error) {
				assert.Error(t, err)
				assert.Equal(t, "begin error", err.Error())
				assert.Equal(t, uuid.Nil, id)
			},
		},
		{
			name: "CreateError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(itemCreateQ)).
					WithArgs(
						item.Title,
						item.Article,
						item.Description,
						item.Price,
						item.Src,
						item.Alt,
						item.InStock,
						item.IsHit,
						item.IsRec,
						item.ParentItemID,
					).
					WillReturnError(errors.New("create error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, id uuid.UUID, err error) {
				assert.Error(t, err)
				assert.Equal(t, "create error", err.Error())
				assert.Equal(t, uuid.Nil, id)
			},
		},
		{
			name: "CommitError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(itemCreateQ)).
					WithArgs(
						item.Title,
						item.Article,
						item.Description,
						item.Price,
						item.Src,
						item.Alt,
						item.InStock,
						item.IsHit,
						item.IsRec,
						item.ParentItemID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))

				mock.ExpectExec(regexp.QuoteMeta(itemMediaCreateQ)).
					WithArgs(sqlmock.AnyArg(), item.Media[0].Src, item.Media[0].Alt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(itemAttrCreateQ)).
					WithArgs(sqlmock.AnyArg(), item.Attributes[0].Name, item.Attributes[0].Value).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(itemCategoryCreateQ)).
					WithArgs(sqlmock.AnyArg(), item.Categories[0].Slug).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, id uuid.UUID, err error) {
				assert.Error(t, err)
				assert.Equal(t, "commit error", err.Error())
				assert.Equal(t, uuid.Nil, id)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				id, err := repo.CreateItem(context.Background(), item)
				tt.expectedResp(t, id, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

//func TestRepository_UpdateItem(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := Repository{conn: db}
//	uid := uuid.New()
//
//	item := &md.Item{
//		Title:        "Updated Item",
//		Article:      "Updated Article",
//		Description:  "Updated Description",
//		Price:        20.0,
//		Src:          "updated-src",
//		Alt:          "updated-alt",
//		InStock:      true,
//		IsHit:        false,
//		IsRec:        false,
//		ParentItemID: uuid.Nil,
//		Media: []md.ItemMedia{
//			{Src: "updated-src1", Alt: "updated-alt1"},
//		},
//		Attributes: []md.ItemAttribute{
//			{Name: "updated-attr1", Value: "updated-value1"},
//		},
//		Categories: []md.Category{
//			{Slug: "updated-slug1"},
//		},
//		RelatedProducts: []md.RelatedProduct{
//			{RelatedItemID: uuid.New()},
//		},
//		Variants: []md.Item{
//			{
//				Title:       "Variant Item",
//				Article:     "Variant Article",
//				Description: "Variant Description",
//				Price:       15.0,
//				Src:         "variant-src",
//				Alt:         "variant-alt",
//				InStock:     false,
//				IsHit:       true,
//				IsRec:       true,
//			},
//		},
//	}
//
//	tests := []struct {
//		name         string
//		mockExpect   func()
//		expectedResp func(*testing.T, error)
//	}{
//		{
//			name: "Success",
//			mockExpect: func() {
//				mock.ExpectBegin()
//
//				mock.ExpectExec(regexp.QuoteMeta(itemUpdateQ)).
//					WithArgs(
//						item.Title,
//						item.Article,
//						item.Description,
//						item.Price,
//						item.Src,
//						item.Alt,
//						item.InStock,
//						item.IsHit,
//						item.IsRec,
//						item.ParentItemID,
//						uid,
//					).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				// Expect media update
//				mock.ExpectQuery(regexp.QuoteMeta(itemMediaList)).
//					WithArgs(uid).
//					WillReturnRows(
//						sqlmock.NewRows([]string{"id", "src", "alt"}).
//							AddRow(uint64(1), "existing-src1", "existing-alt1"),
//					)
//
//				mock.ExpectExec(regexp.QuoteMeta(itemMediaCreateQ)).
//					WithArgs(uid, item.Media[0].Src, item.Media[0].Alt).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				// Expect attribute update
//				mock.ExpectQuery(regexp.QuoteMeta(itemAttrList)).
//					WithArgs(uid).
//					WillReturnRows(
//						sqlmock.NewRows([]string{"id", "name", "value"}).
//							AddRow(uint64(1), "existing-attr1", "existing-value1"),
//					)
//
//				mock.ExpectExec(regexp.QuoteMeta(itemAttrCreateQ)).
//					WithArgs(uid, item.Attributes[0].Name, item.Attributes[0].Value).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				// Expect category update
//				mock.ExpectQuery(regexp.QuoteMeta(itemCategoriesListQ)).
//					WithArgs(uid).
//					WillReturnRows(
//						sqlmock.NewRows([]string{"item_id", "category_slug"}).
//							AddRow(uid, "existing-slug1"),
//					)
//
//				mock.ExpectExec(regexp.QuoteMeta(itemCategoryCreateQ)).
//					WithArgs(uid, item.Categories[0].Slug).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				// Expect related product update
//				mock.ExpectQuery(regexp.QuoteMeta(itemRelatedProductList)).
//					WithArgs(uid).
//					WillReturnRows(
//						sqlmock.NewRows([]string{"item_id", "related_item_id"}).
//							AddRow(uid, "existing-related-id-1"),
//					)
//
//				mock.ExpectExec(regexp.QuoteMeta(itemRelatedProductCreateQ)).
//					WithArgs(uid, item.RelatedProducts[0].RelatedItemID).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				// Expect variant update
//				mock.ExpectQuery("SELECT id, title, article, description, price, src, alt, in_stock, is_hit, is_rec FROM item WHERE parent_id = $1").
//					WithArgs(uid).
//					WillReturnRows(
//						sqlmock.NewRows([]string{"id"}).
//							AddRow("existing-var-id"),
//					)
//
//				mock.ExpectExec(regexp.QuoteMeta(itemUpdateQ)).
//					WithArgs(
//						item.Variants[0].Title,
//						item.Variants[0].Article,
//						item.Variants[0].Description,
//						item.Variants[0].Price,
//						item.Variants[0].Src,
//						item.Variants[0].Alt,
//						item.Variants[0].InStock,
//						item.Variants[0].IsHit,
//						item.Variants[0].IsRec,
//						uid,
//						"existing-var-id",
//					).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				mock.ExpectCommit()
//			},
//			expectedResp: func(t *testing.T, err error) {
//				require.NoError(t, err)
//			},
//		},
//		{
//			name: "BeginError",
//			mockExpect: func() {
//				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
//			},
//			expectedResp: func(t *testing.T, err error) {
//				assert.Error(t, err)
//				assert.Equal(t, "begin error", err.Error())
//			},
//		},
//		{
//			name: "UpdateError",
//			mockExpect: func() {
//				mock.ExpectBegin()
//
//				mock.ExpectExec(regexp.QuoteMeta(itemUpdateQ)).
//					WithArgs(
//						item.Title,
//						item.Article,
//						item.Description,
//						item.Price,
//						item.Src,
//						item.Alt,
//						item.InStock,
//						item.IsHit,
//						item.IsRec,
//						item.ParentItemID,
//						uid,
//					).
//					WillReturnError(errors.New("update error"))
//
//				mock.ExpectRollback()
//			},
//			expectedResp: func(t *testing.T, err error) {
//				assert.Error(t, err)
//				assert.Equal(t, "update error", err.Error())
//			},
//		},
//		{
//			name: "CommitError",
//			mockExpect: func() {
//				mock.ExpectBegin()
//
//				mock.ExpectExec(regexp.QuoteMeta(itemUpdateQ)).
//					WithArgs(
//						item.Title,
//						item.Article,
//						item.Description,
//						item.Price,
//						item.Src,
//						item.Alt,
//						item.InStock,
//						item.IsHit,
//						item.IsRec,
//						item.ParentItemID,
//						uid,
//					).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//
//				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
//			},
//			expectedResp: func(t *testing.T, err error) {
//				assert.Error(t, err)
//				assert.Equal(t, "commit error", err.Error())
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(
//			tt.name, func(t *testing.T) {
//				tt.mockExpect()
//				err := repo.UpdateItem(context.Background(), uid, item)
//				tt.expectedResp(t, err)
//				err = mock.ExpectationsWereMet()
//				assert.NoError(t, err)
//			},
//		)
//	}
//}

func TestRepository_DeleteItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(itemDeleteQ)).
					WithArgs(uid).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResp: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(itemDeleteQ)).
					WithArgs(uid).
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
				mock.ExpectExec(regexp.QuoteMeta(itemDeleteQ)).
					WithArgs(uid).
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
				err := repo.DeleteItem(context.Background(), uid)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ListItemVariants(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, []*md.Item, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "article", "description", "price", "src", "alt",
						"in_stock", "is_hit", "is_rec", "parent_id", "media", "attrs",
					},
				).
					AddRow(
						uid.String(),
						"Variant Item 1",
						"Article 1",
						"Description 1",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						"{src1|alt1}",
						"{attr1|value1}",
					).
					AddRow(
						uid.String(),
						"Variant Item 2",
						"Article 2",
						"Description 2",
						15.5,
						"src2",
						"alt2",
						false,
						true,
						true,
						uuid.New().String(),
						"{src2|alt2}",
						"{attr2|value2}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemListVarsQ)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res []*md.Item, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Len(t, res, 2)
				assert.Equal(t, "Variant Item 1", res[0].Title)
				assert.Equal(t, "Variant Item 2", res[1].Title)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemListVarsQ)).
					WithArgs(uid).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res []*md.Item, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemListVarsQ)).
					WithArgs(uid).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res []*md.Item, err error) {
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
						"id", "title", "article", "description", "price", "src", "alt",
						"in_stock", "is_hit", "is_rec", "parent_id", "media", "attrs",
					},
				).
					AddRow(
						nil,
						"Variant Item 1",
						"Article 1",
						"Description 1",
						10.5,
						"src1",
						"alt1",
						true,
						false,
						false,
						uuid.New().String(),
						pq.Array([]string{"src1|alt1"}),
						pq.Array([]string{"attr1|value1"}),
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemListVarsQ)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res []*md.Item, err error) {
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
				res, err := repo.ListItemVariants(context.Background(), uid)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ListRelatedItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	uid := uuid.New()

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, []*md.RelatedProduct, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "price", "src", "alt",
					},
				).
					AddRow(
						uid.String(), "Related Item 1", 10.5, "src1", "alt1",
					).
					AddRow(
						uid.String(), "Related Item 2", 15.5, "src2", "alt2",
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemListRelated)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res []*md.RelatedProduct, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Len(t, res, 2)
				assert.Equal(t, "Related Item 1", res[0].RelatedItem.Title)
				assert.Equal(t, "Related Item 2", res[1].RelatedItem.Title)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemListRelated)).
					WithArgs(uid).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res []*md.RelatedProduct, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemListRelated)).
					WithArgs(uid).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res []*md.RelatedProduct, err error) {
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
						"id", "title", "price", "src", "alt",
					},
				).
					AddRow(
						uuid.Nil, "Related Item 1", 10.5, "src1", "alt1",
					)

				mock.ExpectQuery(regexp.QuoteMeta(itemListRelated)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res []*md.RelatedProduct, err error) {
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
				res, err := repo.ListRelatedItems(context.Background(), uid)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_HitItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	page := 1
	size := 10
	expectedCount := int64(2)
	expectedTotalPages := int((expectedCount + int64(size) - 1) / int64(size))

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *md.PaginatedItemsData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountHitQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "price", "src", "alt", "is_hit", "is_rec",
					},
				).
					AddRow(uuid.New().String(), "Hit Item 1", 10.5, "src1", "alt1", true, false).
					AddRow(uuid.New().String(), "Hit Item 2", 15.5, "src2", "alt2", true, true)

				mock.ExpectQuery(regexp.QuoteMeta(itemListHitQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Hit Item 1", res.Data[0].Title)
				assert.Equal(t, "Hit Item 2", res.Data[1].Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountHitQ)).
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountHitQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(itemListHitQ)).
					WithArgs((page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountHitQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "price", "src", "alt", "is_hit", "is_rec",
					},
				).
					AddRow(uuid.Nil, "Hit Item 1", 10.5, "src1", "alt1", true, false)

				mock.ExpectQuery(regexp.QuoteMeta(itemListHitQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
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
				res, err := repo.HitItems(context.Background(), page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_RecItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	page := 1
	size := 10
	expectedCount := int64(2)
	expectedTotalPages := int((expectedCount + int64(size) - 1) / int64(size))

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *md.PaginatedItemsData, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				// Setup count query expectation
				mock.ExpectQuery(regexp.QuoteMeta(itemCountRecQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				// Setup item list query expectation
				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "price", "src", "alt", "is_hit", "is_rec",
					},
				).
					AddRow(uuid.New().String(), "Rec Item 1", 10.5, "src1", "alt1", true, false).
					AddRow(uuid.New().String(), "Rec Item 2", 15.5, "src2", "alt2", true, true)

				mock.ExpectQuery(regexp.QuoteMeta(itemListRecQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, expectedCount, res.Count)
				assert.Len(t, res.Data, 2)
				assert.Equal(t, expectedTotalPages, res.TotalPages)
				assert.Equal(t, "Rec Item 1", res.Data[0].Title)
				assert.Equal(t, "Rec Item 2", res.Data[1].Title)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountRecQ)).
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountRecQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(itemListRecQ)).
					WithArgs((page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ScanError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(itemCountRecQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "price", "src", "alt", "is_hit", "is_rec",
					},
				).
					AddRow(uuid.Nil, "Rec Item 1", 10.5, "src1", "alt1", true, false)

				mock.ExpectQuery(regexp.QuoteMeta(itemListRecQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *md.PaginatedItemsData, err error) {
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
				res, err := repo.RecItems(context.Background(), page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}
