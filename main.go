package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

    mux.HandleFunc("GET /comment", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w,"Comentários")
    })

    err:= http.ListenAndServe("localhost:8080", mux)

    if err != nil {
        fmt.Println(err.Error())
    }
}
