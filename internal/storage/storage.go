package storage

import (
	"context"
	"errors"
)

var (
	UrlExist = errors.New("url exists")
	NotFound = errors.New("origin url not found")
)

type Orm interface {
	Save(originUrl string, alias string) (int64, error)
	GetByAlias(alias string) (*UrlModel, error)
	DeleteByAlias(alias string) error
}

type UrlModel struct {
	id     int
	origin string
	alias  string
}

func New(ctx context.Context, storageType string, config map[string]string) (Orm, error) {
	var storage Orm
	var err error

	switch storageType {
	case "sqlite":
		storage, err = NewSqlite(config)
	}

	if err != nil {
		return nil, err
	}

	return storage, nil
}
