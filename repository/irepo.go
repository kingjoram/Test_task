package repository

import "test/pkg/models"

type IDbRepo interface {
	InsertUrl(url models.Url) error
	GetId() (uint64, error)
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
	SaveUrl(short string, long string) error
}
