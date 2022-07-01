package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

type UserStore struct {
	Mu    sync.Mutex
	Store []user
}

func getBatch(n int64, pool int64) (res []user) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, pool)

	userStore := new(UserStore)
	var i int64

	for i = 0; i < n; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(j int64) {
			user := getOne(int64(j))
			userStore.Mu.Lock()
			userStore.Store = append(userStore.Store, user)
			userStore.Mu.Unlock()
			<-sem
			wg.Done()
		}(i)
	}

	wg.Wait()
	return userStore.Store
}
