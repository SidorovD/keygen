package keygen

import (
	"errors"
	"sync"
)

var (
	ErrKeyDoesNotExist = errors.New("key does not exist")
	ErrKeyAlreadyExist = errors.New("key already exists")
)

type KeyStore interface {
	Get(key string) (k *Key, err error)
	Add(k *Key) error
	Update(k *Key) error
	Save() error
}

// IMKS is a dummy store
type inMemoryKeyStore struct {
	mux   sync.RWMutex
	store map[string]*Key
}

func NewStore() KeyStore {
	s := make(map[string]*Key)
	return &inMemoryKeyStore{store: s}
}

func (s *inMemoryKeyStore) Get(key string) (k *Key, err error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if _, ok := s.store[key]; !ok {
		return nil, ErrKeyDoesNotExist
	}

	k = s.store[key]
	return
}

func (s *inMemoryKeyStore) Add(k *Key) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.store[k.Key()]; ok {
		return ErrKeyAlreadyExist
	}

	s.store[k.Key()] = k
	return nil
}

func (s *inMemoryKeyStore) Update(k *Key) error {
	// Do nothing, cause IMKS use a pointer to Key when it do Add
	// and if we change a Key state, it'll be already in memory
	return nil
}

func (s *inMemoryKeyStore) Save() error {
	// Do nothing, cause already in memory
	return nil
}
