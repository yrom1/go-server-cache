package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type Cache struct {
	mu    sync.Mutex
	items map[string]string
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.items[key]
	return value, ok
}

func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func handleCache(w http.ResponseWriter, r *http.Request, cache *Cache) {
	switch r.Method {
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if data["key"] == "" || data["value"] == "" {
			http.Error(w, "key or value missing", http.StatusBadRequest)
			return
		}

		cache.Set(data["key"], data["value"])
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		key := r.URL.Path[len("/cache/"):]
		if key == "" {
			keys := cache.Keys()
			resp, err := json.Marshal(keys)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(resp)
			return
		}

		value, ok := cache.Get(key)
		if !ok {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}

		resp, err := json.Marshal(map[string]string{key: value})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	cache := NewCache()

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		handleCache(w, r, cache)
	})

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

