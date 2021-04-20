package main

import (
	"go-learn/week2/homework/controller"
	"net/http"
)

func main() {
	http.HandleFunc("/add", controller.AddPerson)
	http.HandleFunc("/del", controller.DelPerson)
	http.HandleFunc("/update", controller.UpdatePerson)
	http.HandleFunc("/fetch", controller.FetchPerson)
	http.ListenAndServe(":8081", nil)
}
