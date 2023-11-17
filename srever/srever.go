package server

import (
	"log"
	"net/http"
	"path"
)

func Run() {
	http.HandleFunc("/", handerALL)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handerALL(w http.ResponseWriter, r *http.Request) {
	lenPath := len(r.URL.Path)
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./public/index.html")
	} else if lenPath > 8 && r.URL.Path[:7] == "/public" {
		fileOfpublic(w, r)
	} else if lenPath > len("/api")+2 && r.URL.Path[:len("/api")] == "/api" {
		 api(w , r )
	} else {
		http.NotFound(w, r)
	}

}

func fileOfpublic(w http.ResponseWriter, r *http.Request) {
	pathFile := "." + path.Clean(r.URL.Path)
	http.ServeFile(w, r, pathFile)
}
func api(w http.ResponseWriter, r *http.Request) {

}
