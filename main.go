package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const configDir string = "config"

var appconfig map[string]string

func main() {
	appconfig = make(map[string]string)
	load()

	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)
	periodicReload(ticker, done)

	r := mux.NewRouter()
	r.HandleFunc("/config/{key}", getConfig).Methods("GET")
	http.ListenAndServe(":8080", r)

	ticker.Stop()
	done <- true
}

func periodicReload(ticker *time.Ticker, done <-chan bool) {
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				load()
			}
		}
	}()
}

func load() {
	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		fmt.Println("cannot read dir "+configDir, err)
		return
	}
	for _, file := range files {
		key := file.Name()
		filename := configDir + "/" + file.Name()

		value, _ := ioutil.ReadFile(filename)
		if string(value) == "" {
			fmt.Println("Unable to read config value from", configDir+file.Name())
			continue
		}

		appconfig[key] = string(value)
		fmt.Println(appconfig)
	}
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	fmt.Println("getConfig!")
	key := mux.Vars(req)["key"]
	value, err := read(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not found %s", key), http.StatusNotFound)
		return
	}

	w.Write([]byte(value))
}

func read(key string) (string, error) {
	value, exist := appconfig[key]
	fmt.Println("ronen")
	fmt.Println(value)
	if !exist {
		return "", errors.New("could not found " + key)
	}

	return value, nil
}
