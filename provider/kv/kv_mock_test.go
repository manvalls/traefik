package kv

import (
	"errors"
	"strings"

	"github.com/docker/libkv/store"
	"github.com/manvalls/traefik/types"
)

type KvMock struct {
	Provider
}

func (provider *KvMock) loadConfig() *types.Configuration {
	return nil
}

// Override Get/List to return a error
type KvError struct {
	Get  error
	List error
}

// Extremely limited mock store so we can test initialization
type Mock struct {
	Error           KvError
	KVPairs         []*store.KVPair
	WatchTreeMethod func() <-chan []*store.KVPair
}

func (s *Mock) Put(key string, value []byte, opts *store.WriteOptions) error {
	return errors.New("Put not supported")
}

func (s *Mock) Get(key string, opts *store.ReadOptions) (*store.KVPair, error) {
	if err := s.Error.Get; err != nil {
		return nil, err
	}
	for _, kvPair := range s.KVPairs {
		if kvPair.Key == key {
			return kvPair, nil
		}
	}
	return nil, store.ErrKeyNotFound
}

func (s *Mock) Delete(key string) error {
	return errors.New("Delete not supported")
}

// Exists mock
func (s *Mock) Exists(key string, opts *store.ReadOptions) (bool, error) {
	if err := s.Error.Get; err != nil {
		return false, err
	}
	for _, kvPair := range s.KVPairs {
		if strings.HasPrefix(kvPair.Key, key) {
			return true, nil
		}
	}
	return false, store.ErrKeyNotFound
}

// Watch mock
func (s *Mock) Watch(key string, stopCh <-chan struct{}, opts *store.ReadOptions) (<-chan *store.KVPair, error) {
	return nil, errors.New("Watch not supported")
}

// WatchTree mock
func (s *Mock) WatchTree(prefix string, stopCh <-chan struct{}, opts *store.ReadOptions) (<-chan []*store.KVPair, error) {
	return s.WatchTreeMethod(), nil
}

// NewLock mock
func (s *Mock) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	return nil, errors.New("NewLock not supported")
}

// List mock
func (s *Mock) List(prefix string, opts *store.ReadOptions) ([]*store.KVPair, error) {
	if err := s.Error.List; err != nil {
		return nil, err
	}
	kv := []*store.KVPair{}
	for _, kvPair := range s.KVPairs {
		if strings.HasPrefix(kvPair.Key, prefix) && !strings.ContainsAny(strings.TrimPrefix(kvPair.Key, prefix), "/") {
			kv = append(kv, kvPair)
		}
	}
	return kv, nil
}

// DeleteTree mock
func (s *Mock) DeleteTree(prefix string) error {
	return errors.New("DeleteTree not supported")
}

// AtomicPut mock
func (s *Mock) AtomicPut(key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) {
	return false, nil, errors.New("AtomicPut not supported")
}

// AtomicDelete mock
func (s *Mock) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	return false, errors.New("AtomicDelete not supported")
}

// Close mock
func (s *Mock) Close() {}
