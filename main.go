package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

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

	r.HandleFunc("/css/styles.css", func(w http.ResponseWriter, r *http.Request) {
		arr := [2]string{"css/main.css", "css/test.css"}
		var readers []io.Reader
		var files []*os.File
		for _, x := range arr {
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

		for _, x := range files {
			x.Close()
		}
	})
	http.ListenAndServe(":8080", r)
}
