package repository

import "test/pkg/models"

//go:generate mockgen -source=irepo.go -destination=repo_mock.go -package=repository
type IDbRepo interface {
	InsertUrl(url models.Url) error
	GetId() (uint64, error)
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
}
