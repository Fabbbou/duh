package repository

type MockDbRepositoryFactory struct{}

func NewMockDbRepositoryFactory() *MockDbRepositoryFactory {
	return &MockDbRepositoryFactory{}
}

func (factory *MockDbRepositoryFactory) NewDbRepository(repoPath string) (DbRepository, error) {
	return NewMockInmemoryDbRepository(), nil
}
