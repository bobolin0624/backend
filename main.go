package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, Taiwan Voting Guide!")
    })

    http.ListenAndServe(":8080", nil)
}
