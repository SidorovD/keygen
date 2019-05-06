package keygen

import (
	"errors"
)

type KeyStatus int

const (
	KeyStatus_NotIssued = iota
	KeyStatus_Issued
	KeyStatus_Submitted
)

var (
	ErrKeyWasNotIssued        = errors.New("Key.Submit: key has not been issued")
	ErrKeyWasAlreadySubmitted = errors.New("Key.Submit: key already submitted")
)

type Key struct {
	key    string
	status KeyStatus
}

func NewKey(key string) (*Key, error) {
	if len(key) != keyLen {
		return nil, ErrKeyMustBeAKeyLenSymbols
	}

	return &Key{key: key, status: KeyStatus_Issued}, nil
}

func (k *Key) Key() string {
	return k.key
}

func (k *Key) Status() KeyStatus {
	return k.status
}

func (k *Key) submit() error {
	// !k.isSubmitted() hack
	if !k.isIssued() && !k.isSubmitted() {
		return ErrKeyWasNotIssued
	}

	if k.isSubmitted() {
		return ErrKeyWasAlreadySubmitted
	}

	k.status = KeyStatus_Submitted
	return nil
}

func (k *Key) isIssued() bool {
	return k.status == KeyStatus_Issued
}

func (k *Key) isSubmitted() bool {
	return k.status == KeyStatus_Submitted
}
