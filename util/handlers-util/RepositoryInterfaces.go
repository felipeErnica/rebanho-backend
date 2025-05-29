package handlersUtil

import "github.com/felipeErnica/rebanho-backend/entity"

type PageRepository[E any, F any] interface {
	FindPage(sort string, order string, cursor string, filter F) (*entity.Page[E], error)
}

type RepositoryFindById[E any] interface {
	FindById(id string) (*E, error)
}

type RepositoryFindAll[E any] interface {
	FindAll() (*[]E, error)
}

type RepositoryAdd[E any] interface {
	Add(*E) (*E, error)
    Update(*E) error
}

type RepositoryDelete interface {
	Delete(id string) error
}
