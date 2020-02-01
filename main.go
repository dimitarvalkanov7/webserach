package main

import (
	"flag"
	"sort"
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/dimitarvalkanov7/websearch/webpages"
)

const (
	basePath string = "/home/leron/workspace/src/github.com/dimitarvalkanov7/websearch"
)

var (
	searchData []webpages.PageContent
	addr       = flag.String("addr", ":1234", "http service address")
)

func init() {
	searchData = webpages.GetData()
}

type Context struct {
	CurrentSearch string
	Urls          []string
	ShortUrls     []string
}

func search(w http.ResponseWriter, r *http.Request) {
	searchInput := r.URL.Query()["search"]
	searchOutput := Context{
		CurrentSearch: "",
		Urls:          make([]string, 0),
		ShortUrls:     make([]string, 0),
	}

	if len(searchInput) == 0 {
		tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "search.html")))
		err := tmpl.Execute(w, searchOutput)
		if err != nil {
			log.Printf("Error executing template: %v\n", err)
		}
		return
	}

	searchOutput.CurrentSearch = searchInput[0]
	searchWords := strings.Split(searchInput[0], " ")
	weightToUrl := make(map[string]int)

	for _, v := range searchWords {
		for _, wordsToUrl := range searchData {
			for word, weight := range wordsToUrl.Words {
				if word == v {
					weightToUrl[wordsToUrl.URL] += weight
				}
			}
		}
	}

	type kv struct {
		Key   string
		Value int
	}

	var weitghtToUrlSorted []kv
	for k, v := range weightToUrl {
		weitghtToUrlSorted = append(weitghtToUrlSorted, kv{k, v})
	}
	sort.Slice(weitghtToUrlSorted, func(i, j int) bool {
		return weitghtToUrlSorted[i].Value > weitghtToUrlSorted[j].Value
	})

	index := 0
	for _, url := range weitghtToUrlSorted {
		searchOutput.Urls = append(searchOutput.Urls, url.Key)
		searchOutput.ShortUrls = append(searchOutput.ShortUrls, string(url.Key[:20]))
		index++
		if index > 10 {
			break
		}
	}

	tmpl := template.Must(template.ParseFiles(path.Join(basePath, "templates", "search.html")))
	err := tmpl.Execute(w, searchOutput)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
	}
	return
}

func main() {
	http.HandleFunc("/search/", search)
	http.HandleFunc("/search", search)

	log.Println("Starting the HTTP server ...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
