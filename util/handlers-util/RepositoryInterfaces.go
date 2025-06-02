package handlersUtil

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
)

type PageRepository[E any, F any] interface {
	FindPage(repositoriesUtil.PageProps) (*entity.Page[E], error)
}

type RepositoryFindById[E any] interface {
	FindById(id string) (*E, error)
}

type RepositoryFindAll[E any] interface {
	FindAll(userId string) (*[]E, error)
}

type RepositoryAdd[E any] interface {
	Add(*E) (*E, error)
    Update(*E) error
}

type RepositoryDelete interface {
	Delete(id string) error
}
