package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type routerSwapper struct {
	mu   sync.Mutex
	root *mux.Router
}

func (rs *routerSwapper) Swap(newRouter *mux.Router) {
	rs.mu.Lock()
	rs.root = newRouter
	rs.mu.Unlock()
}

func (rs *routerSwapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs.mu.Lock()
	root := rs.root
	rs.mu.Unlock()
	root.ServeHTTP(w, r)
}

func main() {
	r := mux.NewRouter()
	rs := routerSwapper{}
	rs.Swap(r)
	//TODO Take a look at https://github.com/gorilla/mux/issues/82 to hot-swap the routes in the config file

	// parse config file here
	// if has grandchildren; then
	//	create subrouter
	// else
	//	create routes
	//

	// For initial testing, will use io writer on hardcoded paths:
	//   - css
	//     - styles.css
	//       - css/main.css
	//       - css/test.css

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	path := "/css/styles.css"
	arr := []string{"css/main.css", "css/test.css"}
	setupRoute(r, path, arr)
	path = "/html/index.html"
	arr = []string{"header.html", "index.html"}
	setupRoute(r, path, arr)
	var iface http.Handler = &rs
	http.ListenAndServe(":8080", iface)
}

func setupRoute(r *mux.Router, path string, concats []string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var readers []io.Reader
		var files []*os.File
		for _, x := range concats {
			file, err := os.Open(x)
			if err != nil {
				//TODO: Proper error handling
				panic(err)
			}
			readers = append(readers, file)
			files = append(files, file)
		}
		cat := io.MultiReader(readers...)
		if _, err := io.Copy(w, cat); err != nil {
			log.Fatal(err)
		}
		log.Println("Request for " + path + " served: " + strings.Join(concats, ", "))

		for _, x := range files {
			x.Close()
		}
	})

}
