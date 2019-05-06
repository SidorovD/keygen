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
	ErrKeyHasNotBeenIssued        = errors.New("Key.Submit: key has not been issued")
	ErrKeyHasBeenAlreadySubmitted = errors.New("Key.Submit: key already submitted")
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
	if !k.isIssued() {
		return ErrKeyHasNotBeenIssued
	}

	if k.isSubmitted() {
		return ErrKeyHasBeenAlreadySubmitted
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
