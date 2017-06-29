package controller

import (
	"net/http"
	"fmt"
	"strings"
	"GoJobSkills/log"
)

var logger = log.GetLogger()

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	logger.Println(r.Form)  
	logger.Println("path", r.URL.Path)
	logger.Println("scheme", r.URL.Scheme)
	logger.Println(r.Form["url_long"])
	for k, v := range r.Form {
		logger.Println("key:", k)
		logger.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Go!") 

}