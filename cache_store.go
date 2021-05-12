package main

import (
	"errors"
	"sync"
	"time"
)

var (
	NotFoundErr           = errors.New("not found on cache")
	CannotWriteToCacheErr = errors.New("not sure even this works")
)

type CacheItem struct {
	CreatedAt time.Time
	Content   []byte
}

type CacheStore struct {
	mu sync.Mutex

	store map[string]*CacheItem
}

func NewCacheStore() *CacheStore {
	return &CacheStore{
		mu:    sync.Mutex{},
		store: make(map[string]*CacheItem),
	}
}

func (s *CacheStore) StartLock() {
	s.mu.Lock()
}

func (s *CacheStore) EndLock() {
	s.mu.Unlock()
}

func (s *CacheStore) Get(item string) (*CacheItem, error) {
	val, ok := s.store[item]
	if !ok {
		return nil, NotFoundErr
	}

	return val, nil
}

func (s *CacheStore) Set(key string, item *CacheItem) error {
	return func() error {
		s.store[key] = item

		if a := recover(); a != nil {
			return CannotWriteToCacheErr
		}

		return nil
	}()
}
