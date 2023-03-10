///usr/bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
  "os"
  "crypto/tls"
  "net/http"
)

func main() {
   http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
   resp, err := http.Get(os.Args[1])
   if err != nil {
    panic(err)
    os.Exit(-1)
   }

   if resp.StatusCode == 302 {
	  os.Exit(0)
   } else {
	  os.Exit(-2)
   }
}
