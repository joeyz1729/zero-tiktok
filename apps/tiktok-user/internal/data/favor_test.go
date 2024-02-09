package data

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestRepo_UpdateFavoriteRelation(t *testing.T) {
	repo := NewRepo(dsn, addr)
	num := 30
	var wg sync.WaitGroup
	wg.Add(2 * num)
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			err := repo.AddFavoriteRelation(int64(i%5+1), int64((i+1)%5+1))
			assert.Nil(t, err)
		}(i)
	}
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			err := repo.DelFavoriteRelation(int64(i%5+1), int64((i+1)%5+1))
			assert.Nil(t, err)
		}(i)
	}
}
