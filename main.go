package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/carbocation/go.instagram"
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
	return
}

func welcome(w http.ResponseWriter, r *http.Request) {
	return
}

func postSearchTags(w http.ResponseWriter, r *http.Request) {
	ig, err := instagram.NewInstagram()
	if err != nil {
		log.Println(err)
		return
	}
	
	//Get the user's query (looks like /tags/iamatag%20soami%20metoo)	
	vars := mux.Vars(r)
	tags := strings.Split(vars["tags"], ` `)

	out, err := ig.TagsMediaRecent(tags)
	if err != nil {
		log.Println(err)
		return
	}
	
	if out.Meta.Code != http.StatusOK {
		fmt.Fprintf(w, "Status %d. We received an error from Instagram.\n", out.Meta.Code)
		return
	}

	fmt.Fprintf(w, "Meta:\n%+v\n", out.Meta)
	fmt.Fprintf(w, "Pagination:\n%+v\n", out.Pagination)
	fmt.Fprintf(w, "Data:\n%+v\n", out.Data)
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
	r.HandleFunc("/", welcome)
	r.HandleFunc("/tags/{tags:[0-9a-z ]+}", postSearchTags)
	r.HandleFunc("/hello", hello)
	r.HandleFunc("/redirect/instagram", RedirectHandler)
	http.Handle("/", r)
	http.ListenAndServe(":9999", nil)
}
