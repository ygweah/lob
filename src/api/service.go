package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type AddressService struct {
	addr              string
	addressManager    AddressManager
	server            *http.Server
	wg                sync.WaitGroup
	errListenAndServe atomic.Value
}

func NewAddressService(addr string, m AddressManager) *AddressService {
	return &AddressService{
		addr:           addr,
		addressManager: m,
	}
}

func (s *AddressService) Start() error {
	sm := http.NewServeMux()
	sm.Handle("/address", s)

	s.server = &http.Server{
		Handler: sm,
		Addr:    s.addr,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != nil {
			s.errListenAndServe.Store(err)
		}
	}()

	// sometime is needed on linux to allow server to start
	time.Sleep(time.Second)

	return s.getErrListenAndServe()
}

func (s *AddressService) Wait() error {
	s.wg.Wait()
	if err := s.getErrListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *AddressService) Stop() error {
	return s.server.Shutdown(context.TODO())
}

func (s *AddressService) getErrListenAndServe() error {
	v := s.errListenAndServe.Load()
	if v != nil {
		return v.(error)
	}
	return nil
}

func (s *AddressService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("http request: %+v", r)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	q := r.URL.Query().Get("query")
	log.Printf("http request query: %v", q)

	from := r.URL.Query().Get("from")
	log.Printf("http request from: %v", from)

	fromIndex, err := strconv.Atoi(from)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	to := r.URL.Query().Get("to")
	log.Printf("http request to: %v", to)

	toIndex, err := strconv.Atoi(to)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	adrs := s.addressManager.Find(q, fromIndex, toIndex)
	b, err := toJson(adrs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("http response with content: %v", string(b))
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func toJson(adrs []*Address) ([]byte, error) {
	b, err := json.Marshal(adrs)
	if err != nil {
		return nil, err
	}
	return b, nil
}
