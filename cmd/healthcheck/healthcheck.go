///usr/bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
  "os"
  "net/http"
)

func main() {
   resp, err := http.Head(os.Args[1])
   if err != nil {
    panic(err)
    os.Exit(-1)
   }

   if resp.StatusCode == 200 {
	  os.Exit(0)
   } else {
	  os.Exit(-2)
   }
}
