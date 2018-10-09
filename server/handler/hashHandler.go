/*
Package handler provides request handlers for http requests.
*/
package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ifIMust/encodeServer/server/model"
)

/*
A HashHandler handles requests to the '/hash' endpoint.
A POST request with a non-empty "password" field will be handled by computing the hash 5 seconds later.
The response to this request is the id of the stored hash.

A GET request to '/hash/N' where N is a stored hash ID will respond with the saved hash.
*/
type HashHandler struct {
	nextId        int
	nextIdMutex   sync.Mutex
	keyStore      map[int]string
	keyStoreMutex sync.Mutex
	run           atomic.Value
	delay         time.Duration
	stats         *model.Stats
	hasher        *Hasher
	waitGroup     *sync.WaitGroup
}

/*
NewHashHandler initializes and returns a new HashHandler.
It will log request statistics to stats.
When the handler begins processing delayed requests, it will add them to waitGroup so that
shutdown may be delayed until each request has completed.
*/
func NewHashHandler(stats *model.Stats, waitGroup *sync.WaitGroup) *HashHandler {
	h := new(HashHandler)
	h.keyStore = make(map[int]string)
	h.run.Store(true)
	h.delay = 5 * time.Second
	h.stats = stats
	h.hasher = NewHasher()
	h.waitGroup = waitGroup
	return h
}

/*
HandleRequest is an http request handler intended for use with http.ServeMux.
*/
func (h *HashHandler) HandleRequest(w http.ResponseWriter, request *http.Request) {
	if h.run.Load().(bool) {
		var err error = nil
		if request.Method == "POST" {
			err = h.handlePost(w, request)

		} else if request.Method == "GET" {
			err = h.handleGet(w, request)
		} else {
			err = errors.New(fmt.Sprintf("unsupported request type: %s", request.Method))
		}
		if err != nil {
			log.Println(err)
			http.NotFound(w, request)
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

}

/*
Shutdown disables further handling of requests by this handler.
*/
func (h *HashHandler) Shutdown() {
	h.run.Store(false)
}

/*
SetDelay modifies the hashing delay used by the handler. The default delay is 5s.
*/
func (h *HashHandler) SetDelay(t time.Duration) {
	h.delay = t
}

func (h *HashHandler) getHash(id int) string {
	h.keyStoreMutex.Lock()
	defer h.keyStoreMutex.Unlock()
	return h.keyStore[id]
}

func (h *HashHandler) getNextHashId() int {
	h.nextIdMutex.Lock()
	defer h.nextIdMutex.Unlock()
	h.nextId++
	return h.nextId
}

func (h *HashHandler) handlePost(w http.ResponseWriter, request *http.Request) error {
	startTime := time.Now()
	request.ParseForm()
	password := request.PostForm.Get("password")
	if password == "" {
		return errors.New("malformed request: empty 'password' field")
	}
	nextId := h.getNextHashId()
	h.waitGroup.Add(1)
	go h.delayedHash(nextId, password)
	io.WriteString(w, strconv.Itoa(nextId))
	processingTime := time.Now().Sub(startTime)
	h.stats.AddRequest(processingTime)
	return nil
}

func (h *HashHandler) handleGet(w http.ResponseWriter, request *http.Request) error {
	path := request.URL.Path
	elements := strings.Split(path, "/")
	numElements := len(elements)
	if numElements == 3 {
		reqId := elements[numElements-1]
		id, err := strconv.Atoi(reqId)
		if err != nil {
			return errors.New(fmt.Sprintf("malformed request: non-integer request id: %s", reqId))
		}
		hash := h.getHash(id)
		if hash == "" {
			return errors.New("failed hash lookup")
		}
		io.WriteString(w, hash)
	} else {
		return errors.New(fmt.Sprintf("malformed request: '%v'\n", path))
	}
	return nil
}

func (h *HashHandler) delayedHash(id int, pwd string) {
	time.Sleep(h.delay)
	h.processHash(id, pwd)
	h.waitGroup.Done()
}

func (h *HashHandler) processHash(id int, pwd string) {
	h.keyStoreMutex.Lock()
	defer h.keyStoreMutex.Unlock()
	h.keyStore[id] = h.hasher.Hash(pwd)
}
