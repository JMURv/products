package db

import (
	"context"
	"database/sql"
	"errors"
	repo2 "github.com/JMURv/par-pro/products/internal/repo"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
)

func TestRepository_CategorySearch(t *testing.T) {
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
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categorySearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug"}).AddRow(
					"test-slug",
					"Test title",
					"test-src",
					"test-alt",
					"test-parent",
				).AddRow(
					"test-slug-2",
					"Test title",
					"test-src",
					"test-alt",
					"test-parent",
				)

				mock.ExpectQuery(regexp.QuoteMeta(categorySearchQ)).
					WithArgs("%"+query+"%", "%"+query+"%", (page-1)*size, size).
					WillReturnRows(rows)

			},
			expectedResp: func(t *testing.T, res any, err error) {
				resp, ok := res.(*model.PaginatedCategoryData)
				require.True(t, ok)

				assert.Equal(t, expectedCount, resp.Count)
				assert.Len(t, resp.Data, 2)
				assert.Equal(t, expectedTotalPages, resp.TotalPages)
			},
		},
		{
			name: "Count error",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categorySearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnError(errors.New("count error"))

			},
			expectedResp: func(t *testing.T, res any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
			},
		},
		{
			name: "ErrInternal",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categorySearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(categorySearchQ)).
					WithArgs("%"+query+"%", "%"+query+"%", (page-1)*size, size).
					WillReturnError(errors.New("find error"))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "find error", err.Error())
			},
		},
		{
			name: "Empty",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categorySearchCountQ)).
					WithArgs("%"+query+"%", "%"+query+"%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery(regexp.QuoteMeta(categorySearchQ)).
					WithArgs("%"+query+"%", "%"+query+"%", (page-1)*size, size).
					WillReturnRows(sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug"}))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				resp, ok := res.(*model.PaginatedCategoryData)
				require.True(t, ok)
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				assert.Equal(t, int64(0), resp.Count)
				assert.Len(t, resp.Data, 0)
				assert.Equal(t, 0, resp.TotalPages)
				assert.False(t, resp.HasNextPage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.CategorySearch(context.Background(), query, page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)

			},
		)
	}
}

func TestRepository_ListCategories(t *testing.T) {
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
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categoryCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug", "children"}).
					AddRow(
						"test-slug",
						"Test title",
						"test-src",
						"test-alt",
						"test-parent",
						"{child-slug|child-title|child-src|child-alt|test-slug}",
					).
					AddRow(
						"test-slug-2",
						"Test title 2",
						"test-src-2",
						"test-alt-2",
						"test-parent-2",
						"{child-slug-2|child-title-2|child-src-2|child-alt-2|test-slug-2}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(categoryListQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				resp, ok := res.(*model.PaginatedCategoryData)
				require.True(t, ok)
				assert.NoError(t, err)
				assert.Equal(t, expectedCount, resp.Count)
				assert.Len(t, resp.Data, 2)
				assert.Equal(t, expectedTotalPages, resp.TotalPages)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categoryCountQ)).
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categoryCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(categoryListQ)).
					WithArgs((page-1)*size, size).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
			},
		},
		{
			name: "EmptyResult",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categoryCountQ)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery(regexp.QuoteMeta(categoryListQ)).
					WithArgs((page-1)*size, size).
					WillReturnRows(sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug", "children"}))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				resp, ok := res.(*model.PaginatedCategoryData)
				require.True(t, ok)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), resp.Count)
				assert.Len(t, resp.Data, 0)
				assert.Equal(t, 0, resp.TotalPages)
				assert.False(t, resp.HasNextPage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.ListCategories(context.Background(), page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_GetCategoryBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	slug := "test-slug"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, *model.Category, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug", "children", "filters"}).
					AddRow(
						"test-slug", "Test title", "test-src", "test-alt", "test-parent",
						"{child-slug|child-title|child-src|child-alt|test-parent}",
						"{123|filter-name|filter-values|filter-type|1.23|1.23}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(categoryGetQ)).
					WithArgs(slug).
					WillReturnRows(rows)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM item JOIN item_category ON item_category.item_id = item.id WHERE item_category.category_slug = $1`)).
					WithArgs("child-slug").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
			},
			expectedResp: func(t *testing.T, res *model.Category, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "test-slug", res.Slug)
				assert.Equal(t, "Test title", res.Title)
				assert.Equal(t, 10, res.Children[0].ProductQuantity)
			},
		},
		{
			name: "NotFound",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categoryGetQ)).
					WithArgs(slug).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResp: func(t *testing.T, res *model.Category, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(categoryGetQ)).
					WithArgs(slug).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res *model.Category, err error) {
				assert.Error(t, err)
				assert.Equal(t, "query error", err.Error())
				assert.Nil(t, res)
			},
		},
		{
			name: "ChildrenScanError",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug", "children", "filters"}).
					AddRow(
						"test-slug", "Test title", "test-src", "test-alt", "test-parent", "invalid-children",
						"{123|filter-name|filter-values|filter-type|1.23|1.23}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(categoryGetQ)).
					WithArgs(slug).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.Category, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unable to parse array")
				assert.Nil(t, res)
			},
		},
		{
			name: "FiltersScanError",
			mockExpect: func() {
				rows := sqlmock.NewRows([]string{"slug", "title", "src", "alt", "parent_slug", "children", "filters"}).
					AddRow(
						"test-slug", "Test title", "test-src", "test-alt", "test-parent",
						"{child-slug|child-title|child-src|child-alt|test-parent", "invalid-filters}",
					)

				mock.ExpectQuery(regexp.QuoteMeta(categoryGetQ)).
					WithArgs(slug).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res *model.Category, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unable to parse array")
				assert.Nil(t, res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.GetCategoryBySlug(context.Background(), slug)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_CreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	category := &model.Category{
		Slug:       "test-slug",
		Title:      "Test title",
		Src:        "test-src",
		Alt:        "test-alt",
		ParentSlug: "test-parent",
		Filters: []model.Filter{
			{
				Name:       "filter-name",
				Values:     []string{"value1", "value2"},
				FilterType: "filter-type",
				MinValue:   1.23,
				MaxValue:   1.23,
			},
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

				mock.ExpectQuery(regexp.QuoteMeta(categoryCreateQ)).
					WithArgs(category.Slug, category.Title, category.Src, category.Alt, category.ParentSlug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(category.Slug))

				mock.ExpectQuery(regexp.QuoteMeta(filterCreateQ)).
					WithArgs(
						category.Filters[0].Name,
						pq.Array(category.Filters[0].Values),
						category.Filters[0].FilterType,
						category.Filters[0].MinValue,
						category.Filters[0].MaxValue,
						category.Slug,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("filter-id"))

				mock.ExpectCommit()
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				require.NoError(t, err)
				assert.Equal(t, category.Slug, slug)
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
				assert.Empty(t, slug)
			},
		},
		{
			name: "CreateError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(categoryCreateQ)).
					WithArgs(category.Slug, category.Title, category.Src, category.Alt, category.ParentSlug).
					WillReturnError(errors.New("create error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "create error", err.Error())
				assert.Empty(t, slug)
			},
		},
		{
			name: "UniqueConstraintError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(categoryCreateQ)).
					WithArgs(category.Slug, category.Title, category.Src, category.Alt, category.ParentSlug).
					WillReturnError(errors.New("unique constraint violation"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrAlreadyExists, err)
				assert.Empty(t, slug)
			},
		},
		{
			name: "CreateFiltersError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(categoryCreateQ)).
					WithArgs(category.Slug, category.Title, category.Src, category.Alt, category.ParentSlug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(category.Slug))

				mock.ExpectQuery(regexp.QuoteMeta(filterCreateQ)).
					WithArgs(
						category.Filters[0].Name,
						pq.Array(category.Filters[0].Values),
						category.Filters[0].FilterType,
						category.Filters[0].MinValue,
						category.Filters[0].MaxValue,
						category.Slug,
					).
					WillReturnError(errors.New("create filter error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "create filter error", err.Error())
				assert.Empty(t, slug)
			},
		},
		{
			name: "CommitError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(categoryCreateQ)).
					WithArgs(category.Slug, category.Title, category.Src, category.Alt, category.ParentSlug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(category.Slug))

				mock.ExpectQuery(regexp.QuoteMeta(filterCreateQ)).
					WithArgs(
						category.Filters[0].Name,
						pq.Array(category.Filters[0].Values),
						category.Filters[0].FilterType,
						category.Filters[0].MinValue,
						category.Filters[0].MaxValue,
						category.Slug,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("filter-id"))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResp: func(t *testing.T, slug string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "commit error", err.Error())
				assert.Empty(t, slug)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				slug, err := repo.CreateCategory(context.Background(), category)
				tt.expectedResp(t, slug, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_UpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	category := &model.Category{
		Slug:       "test-slug",
		Title:      "Updated title",
		Src:        "updated-src",
		Alt:        "updated-alt",
		ParentSlug: "updated-parent",
		Filters: []model.Filter{
			{
				ID:         1,
				Name:       "updated-filter-name",
				Values:     []string{"updated-value1", "updated-value2"},
				FilterType: "updated-filter-type",
				MinValue:   1.23,
				MaxValue:   1.23,
			},
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

				mock.ExpectExec(regexp.QuoteMeta(categoryUpdateQ)).
					WithArgs(category.Title, category.Src, category.Alt, category.ParentSlug, category.Slug).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta(filterListQ)).
					WithArgs(category.Slug).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"name",
								"values",
								"filter_type",
								"min_value",
								"max_value",
								"category_slug",
							},
						).
							AddRow(
								1,
								"filter-name",
								pq.StringArray{"value1", "value2"},
								"filter-type",
								1.23,
								1.23,
								category.Slug,
							),
					)

				mock.ExpectExec(regexp.QuoteMeta(filterUpdateQ)).
					WithArgs(
						"updated-filter-name",
						pq.Array([]string{"updated-value1", "updated-value2"}),
						"updated-filter-type",
						1.23,
						1.23,
						category.Slug,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

				mock.ExpectExec(regexp.QuoteMeta(categoryUpdateQ)).
					WithArgs(category.Title, category.Src, category.Alt, category.ParentSlug, category.Slug).
					WillReturnError(errors.New("update error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "update error", err.Error())
			},
		},
		{
			name: "NotFoundError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(categoryUpdateQ)).
					WithArgs(category.Title, category.Src, category.Alt, category.ParentSlug, category.Slug).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, repo2.ErrNotFound, err)
			},
		},
		{
			name: "UpdateFiltersError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(categoryUpdateQ)).
					WithArgs(category.Title, category.Src, category.Alt, category.ParentSlug, category.Slug).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta(filterListQ)).
					WithArgs(category.Slug).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"name",
								"values",
								"filter_type",
								"min_value",
								"max_value",
								"category_slug",
							},
						).
							AddRow(
								1,
								"filter-name",
								pq.StringArray{"value1", "value2"},
								"filter-type",
								1.23,
								1.23,
								category.Slug,
							),
					)

				mock.ExpectExec(regexp.QuoteMeta(filterUpdateQ)).
					WithArgs(
						"updated-filter-name",
						pq.Array([]string{"updated-value1", "updated-value2"}),
						"updated-filter-type",
						1.23,
						1.23,
						category.Slug,
						1,
					).
					WillReturnError(errors.New("update filters error"))

				mock.ExpectRollback()
			},
			expectedResp: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "update filters error", err.Error())
			},
		},
		{
			name: "CommitError",
			mockExpect: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(categoryUpdateQ)).
					WithArgs(category.Title, category.Src, category.Alt, category.ParentSlug, category.Slug).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta(filterListQ)).
					WithArgs(category.Slug).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"name",
								"values",
								"filter_type",
								"min_value",
								"max_value",
								"category_slug",
							},
						).
							AddRow(
								1,
								"filter-name",
								pq.StringArray{"value1", "value2"},
								"filter-type",
								1.23,
								1.23,
								category.Slug,
							),
					)

				mock.ExpectExec(regexp.QuoteMeta(filterUpdateQ)).
					WithArgs(
						"updated-filter-name",
						pq.Array([]string{"updated-value1", "updated-value2"}),
						"updated-filter-type",
						1.23,
						1.23,
						category.Slug,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

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
				err := repo.UpdateCategory(context.Background(), category.Slug, category)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_DeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	slug := "test-slug"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectExec(regexp.QuoteMeta(categoryDeleteQ)).
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
				mock.ExpectExec(regexp.QuoteMeta(categoryDeleteQ)).
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
				mock.ExpectExec(regexp.QuoteMeta(categoryDeleteQ)).
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
				err := repo.DeleteCategory(context.Background(), slug)
				tt.expectedResp(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_CategoryFiltersSearch(t *testing.T) {
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
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(filterCountQ)).
					WithArgs("%" + query + "%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				rows := sqlmock.NewRows(
					[]string{
						"id",
						"name",
						"values",
						"filter_type",
						"min_value",
						"max_value",
						"category_slug",
					},
				).
					AddRow(
						1,
						"filter-name",
						pq.StringArray{"value1", "value2"},
						"filter-type",
						1.23,
						1.23,
						"test-category",
					).
					AddRow(
						2,
						"filter-name-2",
						pq.StringArray{"value3", "value4"},
						"filter-type-2",
						1.23,
						1.23,
						"test-category-2",
					)

				mock.ExpectQuery(regexp.QuoteMeta(filterSearchQ)).
					WithArgs("%"+query+"%", (page-1)*size, size).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				resp, ok := res.(*model.PaginatedFilterData)
				require.True(t, ok)

				assert.Equal(t, expectedCount, resp.Count)
				assert.Len(t, resp.Data, 2)
				assert.Equal(t, expectedTotalPages, resp.TotalPages)
			},
		},
		{
			name: "CountError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(filterCountQ)).
					WithArgs("%" + query + "%").
					WillReturnError(errors.New("count error"))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "count error", err.Error())
			},
		},
		{
			name: "SearchError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(filterCountQ)).
					WithArgs("%" + query + "%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

				mock.ExpectQuery(regexp.QuoteMeta(filterSearchQ)).
					WithArgs("%"+query+"%", (page-1)*size, size).
					WillReturnError(errors.New("search error"))
			},
			expectedResp: func(t *testing.T, res any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "search error", err.Error())
			},
		},
		{
			name: "EmptyResult",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(filterCountQ)).
					WithArgs("%" + query + "%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery(regexp.QuoteMeta(filterSearchQ)).
					WithArgs("%"+query+"%", (page-1)*size, size).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{
								"id",
								"name",
								"values",
								"filter_type",
								"min_value",
								"max_value",
								"category_slug",
							},
						),
					)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				resp, ok := res.(*model.PaginatedFilterData)
				require.True(t, ok)
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				assert.Equal(t, int64(0), resp.Count)
				assert.Len(t, resp.Data, 0)
				assert.Equal(t, 0, resp.TotalPages)
				assert.False(t, resp.HasNextPage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.mockExpect()
				res, err := repo.CategoryFiltersSearch(context.Background(), query, page, size)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}

func TestRepository_ListCategoryFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	slug := "test-slug"

	tests := []struct {
		name         string
		mockExpect   func()
		expectedResp func(*testing.T, any, error)
	}{
		{
			name: "Success",
			mockExpect: func() {
				rows := sqlmock.NewRows(
					[]string{
						"id",
						"name",
						"values",
						"filter_type",
						"min_value",
						"max_value",
						"category_slug",
					},
				).
					AddRow(
						1, "filter-name", pq.StringArray{"value1", "value2"}, "filter-type", 1.23, 1.23, slug,
					).
					AddRow(
						2,
						"filter-name-2",
						pq.StringArray{"value3", "value4"},
						"filter-type-2",
						1.23,
						1.23,
						slug,
					)

				mock.ExpectQuery(regexp.QuoteMeta(filterListQ)).
					WithArgs(slug).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res any, err error) {
				require.NoError(t, err)
				assert.NotNil(t, res)
				resp, ok := res.([]*model.Filter)
				require.True(t, ok)
				assert.Len(t, resp, 2)
				assert.Equal(t, "filter-name", resp[0].Name)
				assert.Equal(t, "filter-name-2", resp[1].Name)
			},
		},
		{
			name: "QueryError",
			mockExpect: func() {
				mock.ExpectQuery(regexp.QuoteMeta(filterListQ)).
					WithArgs(slug).
					WillReturnError(errors.New("query error"))
			},
			expectedResp: func(t *testing.T, res any, err error) {
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
						"name",
						"values",
						"filter_type",
						"min_value",
						"max_value",
						"category_slug",
					},
				).
					AddRow(1, "filter-name", pq.StringArray{"value1", "value2"}, "filter-type", "min", "max", slug).
					AddRow(
						nil,
						"filter-name-2",
						pq.StringArray{"value3", "value4"},
						"filter-type-2",
						"min-2",
						"max-2",
						slug,
					)

				mock.ExpectQuery(regexp.QuoteMeta(filterListQ)).
					WithArgs(slug).
					WillReturnRows(rows)
			},
			expectedResp: func(t *testing.T, res any, err error) {
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
				res, err := repo.ListCategoryFilters(context.Background(), slug)
				tt.expectedResp(t, res, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		)
	}
}
