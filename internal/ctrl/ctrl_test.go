package ctrl

import (
	"context"
	"errors"
	"github.com/JMURv/par-pro/products/internal/controller/mocks"
	repo "github.com/JMURv/par-pro/products/internal/repository"
	md "github.com/JMURv/par-pro/products/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestIsUserExist(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	authRepo := mocks.NewMockAuth(ctrlMock)
	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)
	mockSMTP := mocks.NewMockSMTPRepo(ctrlMock)

	ctrl := New(authRepo, mockRepo, mockCache, mockSMTP)

	ctx := context.Background()
	email := "test@example.com"

	// Test case 1: User exists
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Return(&md.User{}, nil).Times(1)

	isExist, err := ctrl.IsUserExist(ctx, email)
	assert.Nil(t, err)
	assert.True(t, isExist)

	// Test case 2: User does not exist (ErrNotFound)
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Return(nil, repo.ErrNotFound).Times(1)

	isExist, err = ctrl.IsUserExist(ctx, email)
	assert.Nil(t, err)
	assert.False(t, isExist)

	// Test case 3: Repo error (other than ErrNotFound)
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), email).Return(nil, errors.New("some repo error")).Times(1)

	isExist, err = ctrl.IsUserExist(ctx, email)
	assert.NotNil(t, err)
	assert.True(t, isExist)
}
