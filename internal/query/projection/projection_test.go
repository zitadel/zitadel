package projection

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestStart(t *testing.T) {
	duplicateName := gofakeit.Name()
	tests := []struct {
		name        string
		projections func(t *testing.T) []projection
		err         error
	}{
		{
			name: "happy path",
			projections: func(t *testing.T) []projection {
				ctrl := gomock.NewController(t)
				projections := make([]projection, 5)

				for i := range 5 {
					mock := NewMockprojection(ctrl)
					mock.EXPECT().Start(gomock.Any())
					mock.EXPECT().String().Return(gofakeit.Name())
					projections[i] = mock
				}

				return projections
			},
		},
		{
			name: "same projection used twice error",
			projections: func(t *testing.T) []projection {
				projections := make([]projection, 5)

				ctrl := gomock.NewController(t)
				mock := NewMockprojection(ctrl)
				mock.EXPECT().String().Return(duplicateName)
				mock.EXPECT().Start(gomock.Any())
				projections[0] = mock

				for i := 1; i < 4; i++ {
					mock := NewMockprojection(ctrl)
					mock.EXPECT().String().Return(gofakeit.Name())
					mock.EXPECT().Start(gomock.Any())
					projections[i] = mock
				}

				mock = NewMockprojection(ctrl)
				mock.EXPECT().String().Return(duplicateName)
				projections[4] = mock

				return projections
			},
			err: fmt.Errorf("projection for %s already added", duplicateName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projections = tt.projections(t)
			err := Start(t.Context())
			require.Equal(t, tt.err, err)
		})
	}
}
