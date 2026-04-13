package dbmock

func (m *MockConnection) AddExpectation(expectation func(recorder *MockConnectionMockRecorder)) *MockConnection {
	expectation(m.EXPECT())
	return m
}

func (m *MockPool) AddExpectation(expectation func(recorder *MockPoolMockRecorder)) *MockPool {
	expectation(m.EXPECT())
	return m
}

func (m *MockTransaction) AddExpectation(expectation func(recorder *MockTransactionMockRecorder)) *MockTransaction {
	expectation(m.EXPECT())
	return m
}
