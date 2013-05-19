package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/carbocation/go-instagram"
	"github.com/gorilla/mux"
)

func init() {
	instagram.Initialize(&instagram.Cfg{
		ClientID:     "dfd9ae5cf32b4cdd9afe0fd24500a86b",
		ClientSecret: "508160f256fb411290e770acaae31856",
		RedirectURL:  "http://localhost:9999/redirect/instagram", // "http://combigram.carbocation.com/ig-redirect",
	})
}

func searchTags(w http.ResponseWriter, r *http.Request) {

}

func postSearchTags(w http.ResponseWriter, r *http.Request) {
	ig, err := instagram.NewInstagram()
	if err != nil {
		log.Println(err)
		return
	}

	out := ig.TagsMediaRecent("golang")

	fmt.Fprintf(w, "%s", out)
}

func hello(w http.ResponseWriter, r *http.Request) {
	ig, err := instagram.NewInstagram()
	if err != nil {
		log.Println(err)
		return
	}

	url := ig.Authenticate("basic")
	log.Println(url)
	http.Redirect(w, r, url, http.StatusFound)

	return
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("inside redirect")
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(rbody))
	fmt.Fprintf(w, "%s", string(rbody))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", postSearchTags)
	r.HandleFunc("/hello", hello)
	r.HandleFunc("/redirect/instagram", RedirectHandler)
	http.Handle("/", r)
	http.ListenAndServe(":9999", nil)
}
