package main

import (
	"log"
	"net/http"
	"tinder-matching-system/logic/matching"
)

func main() {
	ms := matching.NewMatchingSystem()

	// 注册HTTP处理器
	http.HandleFunc("/add_single_person_and_match", ms.AddSinglePersonAndMatchHandler)
	http.HandleFunc("/query_single_people", ms.QuerySinglePeople)
	http.HandleFunc("/RemoveSinglePerson", ms.RemoveSinglePerson)


	// 启动HTTP服务器
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
