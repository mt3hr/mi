package main

import (
	// _ "net/http/pprof"

	"github.com/mt3hr/mi/src/app/mi/mi/mi_server_cmd"
)

func main() {
	// go http.ListenAndServe(":8080", nil)
	mi_server_cmd.Execute()
}
