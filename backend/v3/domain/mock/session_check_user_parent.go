package domainmock

import (
	gomock "go.uber.org/mock/gomock"
)

// func NewMockCheckUserParentWithExpectations(ctrl *gomock.Controller) *MockCheckUserParent {
// 	mock := NewMockCheckUserParent(ctrl)
// 	mock.EXPECT().SetUserConditionProvider(gomock.Any()).AnyTimes()
// 	return mock
// }

// func (m *MockCheckUserParent) ExpectFetchSession(session, err any) *MockCheckUserParent {
// 	m.EXPECT().FetchSession(session, err).1

func InitCheckUserParent(ctrl *gomock.Controller) *MockCheckUserParent {
	mock := NewMockCheckUserParent(ctrl)
	mock.EXPECT().SetUserConditionProvider(gomock.Any())
	return mock
}

func (m *MockCheckUserParent) AddExpectation(expectation func(recorder *MockCheckUserParentMockRecorder)) *MockCheckUserParent {
	expectation(m.EXPECT())
	return m
}
