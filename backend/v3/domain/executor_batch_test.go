package domain_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
)

func TestBatchExecutors(t *testing.T) {
	// Given
	mockCtrl := gomock.NewController(t)
	mockCmd1 := domainmock.NewMockCommander(mockCtrl)
	mockCmd2 := domainmock.NewMockCommander(mockCtrl)
	mockCmd3 := domainmock.NewMockCommander(mockCtrl)

	mockCmd1.EXPECT().String().Return("cmd1").AnyTimes()
	mockCmd2.EXPECT().String().Return("cmd2").AnyTimes()
	mockCmd3.EXPECT().String().Return("cmd3").AnyTimes()

	gomock.InOrder(
		mockCmd1.EXPECT().Validate(gomock.Any(), gomock.Any()).Times(1),
		mockCmd1.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(1),
		mockCmd1.EXPECT().Events(gomock.Any(), gomock.Any()).Times(1),
		mockCmd2.EXPECT().Validate(gomock.Any(), gomock.Any()).Times(1),
		mockCmd2.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(1),
		mockCmd2.EXPECT().Events(gomock.Any(), gomock.Any()).Times(1),
		mockCmd3.EXPECT().Validate(gomock.Any(), gomock.Any()).Times(1),
		mockCmd3.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(1),
		mockCmd3.EXPECT().Events(gomock.Any(), gomock.Any()).Times(1),
	)

	batcher := domain.BatchExecutor(mockCmd1, mockCmd2, mockCmd3)

	// Test
	err := domain.Invoke(context.Background(), batcher)

	// Verify
	require.NoError(t, err)
	require.True(t, mockCtrl.Satisfied())
}
