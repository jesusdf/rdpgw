///usr/bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	expectedCode, err := strconv.Atoi(os.Args[1])

	if err != nil {
		panic(err)
	}

	url := os.Args[2]
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == expectedCode {
		fmt.Printf("[OK]\n")
		os.Exit(0)
	} else {
		fmt.Printf("[KO] Got status code %v when accessing %v, and it should be %v.\n\n%v", resp.StatusCode, url, expectedCode, resp)
		os.Exit(-1)
	}

}
