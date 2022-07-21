package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const proxyAddr string = "localhost:9000"

var (
	counter            int    = 0
	firstInstanseHost  string = "http://localhost:8082"
	secondInstanseHost string = "http://localhost:8080"
)

func main() {
	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {

	if counter == 0 {
		if r.Method == "POST" {
			part := r.URL.Path
			content, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			ansBytes, err := http.Post(firstInstanseHost+part, r.Header.Get("Content-Type"), bytes.NewBuffer(content))
			if err != nil {
				log.Println(err)
			}
			log.Println(ansBytes)
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}

			w.Write(ans)
			defer r.Body.Close()

		} else if r.Method == "DELETE" {
			part := r.URL.Path
			content, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			ansBytes, err := http.NewRequest("DELETE", firstInstanseHost+part, bytes.NewBuffer(content))
			if err != nil {
				log.Println(err)
			}
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}
			w.Write(ans)
			log.Println(string(ans))
			defer r.Body.Close()

		} else if r.Method == "GET" {
			part := r.URL.Path
			ansBytes, err := http.Get(firstInstanseHost + part)
			if err != nil {
				log.Println(err)
			}
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}
			w.Write(ans)
			defer r.Body.Close()

		} else if r.Method == "PUT" {
			part := r.URL.Path
			content, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			ansBytes, err := http.NewRequest("PUT", firstInstanseHost+part, bytes.NewBuffer(content))
			if err != nil {
				log.Println(err)
			}
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}
			w.Write(ans)
			defer r.Body.Close()
		}
		counter++
	} else {
		if r.Method == "POST" {
			part := r.URL.Path
			content, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			ansBytes, err := http.Post(secondInstanseHost+part, r.Header.Get("Content-Type"), bytes.NewBuffer(content))
			if err != nil {
				log.Println(err)
			}

			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}

			w.Write(ans)
			defer r.Body.Close()

		} else if r.Method == "DELETE" {
			part := r.URL.Path
			content, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			ansBytes, err := http.NewRequest("DELETE", secondInstanseHost+part, bytes.NewBuffer(content))
			if err != nil {
				log.Println(err)
			}
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}
			w.Write(ans)
			defer r.Body.Close()

		} else if r.Method == "GET" {
			part := r.URL.Path
			ansBytes, err := http.Get(secondInstanseHost + part)
			if err != nil {
				log.Println(err)
			}
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}
			w.Write(ans)
			defer r.Body.Close()

		} else if r.Method == "PUT" {
			part := r.URL.Path
			content, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			ansBytes, err := http.NewRequest("PUT", secondInstanseHost+part, bytes.NewBuffer(content))
			if err != nil {
				log.Println(err)
			}
			ans, err := ioutil.ReadAll(ansBytes.Body)
			if err != nil {
				log.Println(err)
			}
			w.Write(ans)
			defer r.Body.Close()
		}
		counter--
	}
}
