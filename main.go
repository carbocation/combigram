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
		template.Must(template.New("error").Parse(tpl.all)).Execute(w, struct{ Error error }{Error: err})
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

func searchLatLong(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Searching lat long")
	ig, err := instagram.NewInstagram()
	if err != nil {
		log.Println(err)
		return
	}

	//Get the user's query (looks like /tags/iamatag%20soami%20metoo)
	vars := mux.Vars(r)
	lat, long := vars["lat"], vars["long"]

	out, err := ig.LocationSearch(lat, long)
	if err != nil {
		log.Println(err)
		template.Must(template.New("error").Parse(tpl.all)).Execute(w, struct{ Error error }{Error: err})
		return
	}

	var data = struct {
		Json  *[]instagram.InstagramData
		Query []string
	}{
		Json:  out,
		Query: []string{lat, long},
	}

	template.Must(template.New("searchLatLong").Parse(tpl.all)).Execute(w, data)
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
	By {{.User.Username}} on {{.Created}}
	<br />
	{{ if .Location }}
		In {{.Location.Name}}<br />
		Lat: {{.Location.Latitude}}<br />
		Long: {{.Location.Longitude}}<br />
		<br />
	{{end}}
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
{{end}}

{{define "error"}}
<h1>Error</h1>
{{.Error}}
<br />
This is usually because they blocked a query with adult language.
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
	r.HandleFunc("/lat/{lat:[0-9a-z -.]+}/long/{long:[0-9a-z -.]+}", searchLatLong).Name("searchLatLong")
	r.HandleFunc("/redirect/instagram", RedirectHandler)
	http.Handle("/", r)

	//Launch the server
	if err := http.ListenAndServe("localhost:9998", nil); err != nil {
		panic(err)
	}
}
