package memory

import (
	"log/slog"
	"sync"
	"test/pkg/models"
)

type RepoInMemory struct {
	mu     sync.RWMutex
	memory map[string]string
}

func GetMemoryRepo(lg *slog.Logger) (*RepoInMemory, error) {
	lg.Info("creating memory repo")
	return &RepoInMemory{
		memory: make(map[string]string),
	}, nil
}

func (repo *RepoInMemory) InsertUrl(url models.Url) error {
	repo.mu.RLock()
	repo.memory[url.Short] = url.Long
	repo.mu.Unlock()

	return nil
}

func (repo *RepoInMemory) GetId() (uint64, error) {
	id := uint64(1)

	repo.mu.Lock()
	for range repo.memory {
		id++
	}
	repo.mu.Unlock()

	return id, nil
}

func (repo *RepoInMemory) GetShort(long string) (string, error) {
	short := ""

	repo.mu.Lock()
	for key, value := range repo.memory {
		if value == long {
			short = key
			break
		}
	}
	repo.mu.Unlock()

	return short, nil
}

func (repo *RepoInMemory) GetLong(short string) (string, error) {
	long := ""

	repo.mu.Lock()
	long = repo.memory[short]
	repo.mu.Unlock()

	return long, nil
}

func (repo *RepoInMemory) SaveUrl(short string, long string) error {

	return nil
}
