package main

import (
	"net/http"

	"code.ysitd.cloud/component/exposer/internal/bootstrap"
)

func main() {
	go bootstrap.GetSyncer().Run()
	http.ListenAndServe(":8080", bootstrap.GetService())
}
