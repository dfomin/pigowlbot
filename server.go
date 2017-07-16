package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := http.ListenAndServeTLS(":8444", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		fmt.Println("ListenAndServeTLS: ", err)
	}
}
