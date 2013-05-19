package main

import (
	"fmt"
	"html/template"
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

func welcome(w http.ResponseWriter, r *http.Request) {
	template.Must(template.New("welcome").Parse(tpl.all)).Execute(w, struct{}{})
	
	return
}

func searchTags(w http.ResponseWriter, r *http.Request) {
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

	var data = struct {
		Json  *[]instagram.InstagramData
		Query []string
	}{
		Json:  out,
		Query: tags,
	}

	template.Must(template.New("searchTags").Parse(tpl.all)).Execute(w, data)
	/*
		fmt.Fprintf(w, "Meta:\n%+v\n", out.Meta)
		fmt.Fprintf(w, "Pagination:\n%+v\n", out.Pagination)
		fmt.Fprintf(w, "Data:\n%+v\n", out.Data)
	*/
}

var tpl = struct {
	all string
}{
	all: `{{define "searchTags"}}
<html>
<head>
<title>
{{range .Query}}{{.}}+{{end}}
</title>
</head>
<body>
<h1>You searched for {{range .Query}}#{{.}}+{{end}}</h1>
<br />
<br />
{{range .Json}}
	{{template "parseData" .}}
	<br />
	<br />
{{end}}
</body>
</html>
{{end}}

{{define "parseData"}}
<div>
	By {{.User.Username}}
	<br />
	Tags: 
	{{range .Tags}}#{{.}}, {{end}}
	<br />
	{{with .Images.LowResolution}}
		<img src="{{.URL}}" width={{.Width}} height={{.Height}}>
	{{end}}
	<br />
</div>
{{end}}

{{define "welcome"}}
<html>
<body>
<h1>Welcome</h1>
<a href="/tags/golang%20gopher">Try an example search for tags containing 'golang' and 'gopher'</a>
</body>
</html>
{{end}}`,
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
	r.HandleFunc("/", welcome).Name("welcome")
	r.HandleFunc("/tags/{tags:[0-9a-z ]+}", searchTags).Name("searchTags")
	r.HandleFunc("/redirect/instagram", RedirectHandler)
	http.Handle("/", r)
	http.ListenAndServe(":9998", nil)
}
