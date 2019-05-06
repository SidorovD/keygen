package keygen

import (
	"errors"
	"log"
	"math/rand"
)

const (
	// 26+26+10
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	keyLen          = 4

	// Total number of all posible combinations
	// (26+26+10)^4
	NOfCombinations = uint(14776336)
)

var (
	ErrAllPossibleCombinationsWereIssued = errors.New("all combinations of keys were issued")
	ErrKeyMustBeAKeyLenSymbols           = errors.New("key must be keyLen symbols")
)

// KeyGen is a REST-API symbol key generator
// consisting of uppercase and lowercase letters of the Latin alphabet, as well as numbers.
type KeyGen interface {

	// Gen generates the uniq key
	//
	// returns ErrAllPosibleCombinationsWereIssued if all possible key combinations have been issued
	Gen() (key string, err error)

	// Submit accepts keys
	//
	// return ErrKeyMustBeAKeyLenSymbols, ErrKeyHasNotBeenIssued, ErrKeyHasBeenAlreadySubmitted
	Submit(key string) error

	// Status returns one of the statuses
	// - KeyStatus_NotIssued -- if the key was issued
	// - KeyStatus_Issued -- if the key was issued and not submitted
	// - KeyStatus_Submitted -- if the key was issued and submitted
	//
	// returns the ErrKeyMustBeAKeyLenSymbols if key shorter or longer than keyLen symbols
	Status(key string) (KeyStatus, error)

	// FreeKeysCount returns a count of keys that not issued
	FreeKeysCount() uint
}

type keyGen struct {
	store    KeyStore
	keysUsed uint
}

func New(s KeyStore) KeyGen {
	return &keyGen{store: s}
}

// Gen generates and returns a uniq key
//
// returns ErrAllPosibleCombinationsWereIssued if
// all possible combinations of keys were issued
func (g *keyGen) Gen() (string, error) {
	if NOfCombinations <= g.keysUsed {
		return "", ErrAllPossibleCombinationsWereIssued
	}

	k := g.gen()

	err := g.store.Add(k)
	if err == ErrKeyAlreadyExist {
		return g.Gen()
	} else if err != nil {
		log.Panic(err)
	}

	log.Printf("key %s generated", k.key)

	g.keysUsed++
	return k.key, nil
}

// gen does a random key generation
func (g *keyGen) gen() *Key {
	k := make([]byte, keyLen)
	for i := range k {

		// default seed for rand.Intn is 1, for production usage it's not ok,
		// but for example it's fine
		k[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	key, err := NewKey(string(k))
	if err != nil {
		log.Panic(err)
	}

	return key
}

// Submit checks the key if key was generated and not used,
// it return true, false instead
func (g *keyGen) Submit(key string) error {
	if len(key) != keyLen {
		return ErrKeyMustBeAKeyLenSymbols
	}

	k, err := g.store.Get(key)
	if err != nil {
		if err == ErrKeyDoesNotExist {
			return ErrKeyHasNotBeenIssued
		}

		log.Panic(err)
	}

	if err := k.submit(); err != nil {
		return err
	}

	if err := g.store.Update(k); err != nil {
		log.Panic(err)
	}

	return nil
}

// Status returns one of the statuses
// - KeyStatus_NotIssued -- if key doesn't issued
// - KeyStatus_Issued -- if key issued and not submitted
// - KeyStatus_Submitted -- if key issued and submitted
func (g *keyGen) Status(key string) (KeyStatus, error) {
	if len(key) != keyLen {
		return KeyStatus(-1), ErrKeyMustBeAKeyLenSymbols
	}

	k, err := g.store.Get(key)
	if err != nil {
		if err == ErrKeyDoesNotExist {
			return KeyStatus_NotIssued, nil
		}

		log.Panic(err)
	}

	return k.Status(), nil
}

// FreeKeysCount returns a count of keys that not issued
func (g *keyGen) FreeKeysCount() uint {
	return NOfCombinations - g.keysUsed
}
