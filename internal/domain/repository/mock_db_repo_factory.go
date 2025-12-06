package repository

type MockDbRepositoryFactory struct {
	InmemoryCreated []MockInmemoryDbRepository
}

func NewMockDbRepositoryFactory() *MockDbRepositoryFactory {
	return &MockDbRepositoryFactory{
		InmemoryCreated: []MockInmemoryDbRepository{},
	}
}

func (factory *MockDbRepositoryFactory) NewDbRepository(repoPath string) (DbRepository, error) {
	factory.InmemoryCreated = append(factory.InmemoryCreated, *NewMockInmemoryDbRepository())
	return &factory.InmemoryCreated[len(factory.InmemoryCreated)-1], nil
}
