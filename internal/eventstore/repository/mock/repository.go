package mock

func NewMock(t *testing.T) *MockRepository {
	m := NewMockRepository(gomock.NewController(t))
	return nil
}
