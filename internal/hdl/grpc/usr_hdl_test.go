package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/par-pro/products/api/pb"
	m2 "github.com/JMURv/par-pro/products/internal/controller/mocks"
	"github.com/JMURv/par-pro/products/internal/handler/grpc/mocks"
	utils "github.com/JMURv/par-pro/products/pkg/utils/http"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestUserSearch(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockCtrl(ctrlMock)
	auth := m2.NewMockAuth(ctrlMock)
	h := New(auth, mockCtrl)

	ctx := context.Background()
	page := uint64(1)
	size := uint64(10)
	query := "test"

	expectedData := &utils.PaginatedData{}

	// Case 1: Invalid request (query, page, or size is empty/zero)
	t.Run("Invalid request", func(t *testing.T) {
		req := &pb.UserSearchReq{Query: "", Page: 0, Size: 0}
		res, err := h.UserSearch(ctx, req)

		assert.Nil(t, res)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	// Case 2: Controller returns user data successfully
	t.Run("Success", func(t *testing.T) {
		req := &pb.UserSearchReq{Query: query, Page: page, Size: size}

		mockCtrl.EXPECT().UserSearch(gomock.Any(), query, int(page), int(size)).Return(expectedData, nil).Times(1)

		res, err := h.UserSearch(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int64(expectedData.TotalPages), res.TotalPages)
	})

	// Case 3: Controller returns an internal error
	t.Run("Controller error", func(t *testing.T) {
		req := &pb.UserSearchReq{Query: query, Page: page, Size: size}

		mockCtrl.EXPECT().UserSearch(gomock.Any(), query, int(page), int(size)).Return(nil, errors.New("internal error")).Times(1)

		res, err := h.UserSearch(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, codes.Internal, status.Code(err))
	})
}
