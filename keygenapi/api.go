package keygenapi

import (
	"errors"
	"fmt"
	"keygen"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type api struct {
	mux     *http.ServeMux
	service keygen.KeyGen
}

func New(kg keygen.KeyGen) http.Handler {
	m := http.NewServeMux()
	a := &api{service: kg, mux: m}

	m.HandleFunc("/keys/", a.keysHandler)
	m.HandleFunc("/key", a.genKeyHandler)
	m.HandleFunc("/count", a.countHandler)

	return a
}

func (a *api) keysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		a.submitKeyHandler(w, r)
		return
	} else if r.Method == "GET" {
		a.statusHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (a *api) submitKeyHandler(w http.ResponseWriter, r *http.Request) {
	key, err := cutKey(r.RequestURI, "/keys/")
	if err != nil {
		writeWithHeader(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	if err := a.service.Submit(key); err != nil {
		if err == keygen.ErrKeyMustBeAKeyLenSymbols || err == keygen.ErrKeyWasNotIssued || err == keygen.ErrKeyWasAlreadySubmitted {
			log.Print(err)

			writeWithHeader(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *api) statusHandler(w http.ResponseWriter, r *http.Request) {
	k, err := cutKey(r.RequestURI, "/keys/")
	if err != nil {
		writeWithHeader(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	st, err := a.service.Status(k)
	if err != nil {
		if err == keygen.ErrKeyMustBeAKeyLenSymbols {
			writeWithHeader(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if st == keygen.KeyStatus_NotIssued {
		w.Write([]byte("not issued"))
		return
	} else if st == keygen.KeyStatus_Issued {
		w.Write([]byte("issued"))
		return
	} else if st == keygen.KeyStatus_Submitted {
		w.Write([]byte("submitted"))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func cutKey(from, substr string) (string, error) {
	keys := strings.Split(from, substr)
	if len(keys) < 2 {
		msg := fmt.Sprintf("enter a key after \"%s\"", substr)
		return "", errors.New(msg)
	}

	return keys[1], nil
}

func (a *api) genKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	k, err := a.service.Gen()
	if err != nil {
		log.Print(err)

		st := http.StatusInternalServerError
		m := []byte(err.Error())
		writeWithHeader(w, st, m)

		return
	}

	writeWithHeader(w, http.StatusCreated, []byte(k))
}

func (a *api) countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	c := a.service.FreeKeysCount()
	cs := strconv.FormatUint(uint64(c), 10)
	w.Write([]byte(cs))
}

func writeWithHeader(w http.ResponseWriter, statusCode int, mes []byte) {
	w.WriteHeader(statusCode)
	w.Write(mes)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("from: %v, uri: %v", r.RemoteAddr, r.RequestURI)
	a.mux.ServeHTTP(w, r)
}
