package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
)

func TestBatchExecutors(t *testing.T) {
	t.Parallel()
	t.Run("when commands are batched should run in expected order and return no error", func(t *testing.T) {
		t.Parallel()
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

		batcher := domain.BatchExecutors(mockCmd1, mockCmd2, mockCmd3)
		err := domain.Invoke(t.Context(), batcher)
		require.NoError(t, err)
		require.True(t, mockCtrl.Satisfied())
	})

	t.Run("when batched commands fail should return error from invoker", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		mockCmd1 := domainmock.NewMockCommander(mockCtrl)
		mockCmd2 := domainmock.NewMockCommander(mockCtrl)

		mockCmd1.EXPECT().String().Return("cmd1").AnyTimes()
		mockCmd2.EXPECT().String().Return("cmd2").AnyTimes()

		gomock.InOrder(
			mockCmd1.EXPECT().Validate(gomock.Any(), gomock.Any()).Times(1),
			mockCmd1.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(assert.AnError).Times(1),
		)

		batcher := domain.BatchExecutors(mockCmd1, mockCmd2)
		err := domain.Invoke(t.Context(), batcher)
		require.ErrorIs(t, err, assert.AnError)
	})
}
