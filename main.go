package main

import (
	"net/http"
	"goJobSkills/controller"
	"log"
	"goJobSkills/job"
	"goJobSkills/client"
)

func main() {
	initialize()
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func initialize() {
	go job.Register()
	addHandler()
	client.Init()
}

func addHandler()  {
	http.HandleFunc("/", controller.HelloWorld)
	http.HandleFunc("/spider/jianshu", controller.JianShu)
}