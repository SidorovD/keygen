package keygen_test

import (
	. "keygen"
	"testing"
)

var kg KeyGen

func setUp() {
	ks := NewStore()
	kg = New(ks)
}

func TestKeyGenerator_Gen(t *testing.T) {
	setUp()

	key, err := kg.Gen()
	if err != nil {
		t.Fatal(err)
	}

	if klen := len(key); klen != 4 {
		t.Errorf("key must be exactly 4 symbols, instead: %v", klen)
	}

	// Check a 2 generated keys are different
	key2, err := kg.Gen()
	if err != nil {
		t.Fatal(err)
	}

	if key == key2 {
		t.Error("key must diff, equals instead")
	}
}

func TestKeyGenerator_Submit(t *testing.T) {
	setUp()

	key, _ := kg.Gen()

	if st, _ := kg.Status(key); st != KeyStatus_Issued {
		t.Fatal("key not issued")
	}

	if err := kg.Submit(key); err != nil {
		t.Error(err)
	}

	if st, _ := kg.Status(key); st != KeyStatus_Submitted {
		t.Fatal("key not submitted")
	}

	// submit key shorter than a 4 symbols
	key = "sht"
	if err := kg.Submit(key); err != ErrKeyMustBeAKeyLenSymbols {
		t.Errorf("key is shorter than a 4 symbols. want %v, got %v", ErrKeyMustBeAKeyLenSymbols, err)
	}

	// submit key longer than a 4 symbols
	key = "longkey"
	if err := kg.Submit(key); err != ErrKeyMustBeAKeyLenSymbols {
		t.Errorf("key is longer than a 4 symbols. want %v, got %v", ErrKeyMustBeAKeyLenSymbols, err)
	}
}

func TestKeyGenerator_Submit_keyWasNotIssued(t *testing.T) {
	setUp()

	// not issued key
	key := "noti"

	st, err := kg.Status(key)
	if err != nil {
		t.Fatal(err)
	}

	if st != KeyStatus_NotIssued {
		t.Fatalf("want: %v, got: %v", KeyStatus_NotIssued, st)
	}

	if err := kg.Submit(key); err != ErrKeyWasNotIssued {
		t.Errorf("want: %v, got: %v", ErrKeyWasNotIssued, err)
	}
}

func TestKeyGenerator_Submit_keySubmitted(t *testing.T) {
	setUp()

	key, _ := kg.Gen()

	if st, _ := kg.Status(key); st != KeyStatus_Issued {
		t.Fatal("key not issued")
	}

	if err := kg.Submit(key); err != nil {
		t.Fatal(err)
	}

	if st, _ := kg.Status(key); st != KeyStatus_Submitted {
		t.Fatal("key not submitted")
	}

	if err := kg.Submit(key); err != ErrKeyWasAlreadySubmitted {
		t.Errorf("want: %v, got: %v", ErrKeyWasAlreadySubmitted, err)
	}
}

func TestKeyGen_FreeKeysCount(t *testing.T) {
	setUp()

	if n := kg.FreeKeysCount(); n != NOfCombinations {
		t.Fatalf("keys count must be %d, got: %d", NOfCombinations, n)
	}

	_, err := kg.Gen()
	if err != nil {
		t.Fatal(err)
	}

	if n := kg.FreeKeysCount(); n != NOfCombinations-1 {
		t.Errorf("keys count must be %d, got: %d", NOfCombinations-1, n)
	}
}
