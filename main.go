package main

import (
	"net/http"
	"GoJobSkills/controller"
	"log"
	"GoJobSkills/job"
)

func main() {
	go job.Register()
	http.HandleFunc("/", controller.HelloWorld)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}