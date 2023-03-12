///usr/bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
  "os"
  "crypto/tls"
  "net/http"
  "strconv"
  "fmt"
)

func main() {

   http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

   expectedCode, err := strconv.Atoi(os.Args[1])

   if err != nil {
	panic(err)
	os.Exit(-1)
   }

   resp, err := http.Get(os.Args[2])

   if err != nil {
	panic(err)
	os.Exit(-2)
   }

   if resp.StatusCode == expectedCode {
	fmt.Printf("[OK]\n")
	os.Exit(0)
   } else {
	fmt.Printf("[KO] Got status code %v and it should be %v.\n", resp.StatusCode, expectedCode)
	os.Exit(-3)
   }

}

