package handlersUtil

type RepositoryFindById[E any] interface {
	FindById(id string) (*E, error)
}

type RepositoryFindAll[E any] interface {
	FindAll(userId string) (*[]E, error)
}

type RepositoryAdd[E any] interface {
	Add(*E) (*E, error)
}

type RepositoryUpdate[E any] interface {
	Update(*E) error
}

type RepositoryDelete interface {
	Delete(id string) error
}
