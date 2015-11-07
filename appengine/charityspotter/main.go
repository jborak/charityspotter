package hello

import (
    "fmt"
    "net/http"
		_"time"
)

// var (
// 	counter int
// )

// func num(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, fmt.Sprintf("Num: %d", counter));
// }

func init() {
	  // http.HandleFunc("/num/", num)
    http.HandleFunc("/api/hello/", helloWorld)

    // // Run the polling mechanism in the background.
    // go func() {
    //   for {
    //     counter++
    //     time.Sleep(time.Duration(5) * time.Second)
    //   }
    // }()
}

/////
// Data
////

// Image stores information about images created.
type Image struct {
  ID      int      `json:"uid"`
  URL     string   `json:"url"`
  Created int64    `json:"created"`
  Tags    []string `json:"tags"`
}

// Poll stores information about the last time we polled, so we can
// skip ones we've already processed.
type Poll struct {
  Created int64
}

////
// Handlers
////
const (
  ImageIndex string = "ImageIndex"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}



