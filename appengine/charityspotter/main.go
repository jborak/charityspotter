package charityspotter

import (
  "bytes"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "net/url"
  "strings"
	"time"

  "appengine"
  "appengine/search"
  "appengine/urlfetch"
)

func init() {
  // Initialize globals
  firebase = &http.Client{}

  // Setup URI functions
  http.HandleFunc("/api/hello/", helloWorld)
  http.HandleFunc("/api/index/", updateAndIndex)
  http.HandleFunc("/api/debug/", debugIndex)
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
type LastUpdate struct {
  Updated int64 `json:"updated"`
}

type ImageDoc struct {
  // ID      string
  Data    string
  URL     string
  Created time.Time
}

////
// Handlers
////
func helloWorld(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello, world!")
}

////
// Firebase funcs
////
const (
  // Firebase URLs
  FirebaseURL string = "https://charitysandbox.firebaseio.com"
  ItemsPath   string = "/items.json"  // GET
  UpdatePath  string = "/update.json" // PUT

  // AppEngine constants
  ImageIndex string = "ImageIndex"
)

var (
  firebase *http.Client
)

func debugIndex(w http.ResponseWriter, r *http.Request) {
  context := appengine.NewContext(r)

  // Fetch document index.
  index, _ := search.Open(ImageIndex)
  num := 0
  for t := index.List(context, nil); ; {
    var doc ImageDoc
    _, err := t.Next(&doc)
    if err == search.Done {
      break
    } else if err != nil {
      break
    }
    num++
  }

  fmt.Fprintf(w, "%d documents in index", num)
}

func updateAndIndex(w http.ResponseWriter, r *http.Request) {
  context := appengine.NewContext(r)

  images := getLatestImages(context)

  // Fetch document index.
  index, err := search.Open(ImageIndex)
  if err != nil {
    context.Errorf("%s", err.Error())
    fmt.Fprint(w, err.Error())
    return
  }

  // Create image docs from images so we can insert them into the index.
  for _, image := range images {
    imageDoc := &ImageDoc{
      // ID: fmt.Sprintf("%d", image.ID),
      URL: image.URL,
      Created: time.Unix(image.Created, 0),
      Data: strings.Join(image.Tags, " "),
    }
    _, err := index.Put(context, fmt.Sprintf("%d", image.ID), imageDoc) 
    if err != nil {
      fmt.Fprintf(w, err.Error())
    }
  }

  // Debug
  fmt.Fprintf(w, "Added %d new images to index\n", len(images))

  // Set new last updated time for next iteration.
  if len(images) > 0 {
    setLastUpdate(context, images[len(images) - 1].Created) 
  }
}

func getLatestImages(context appengine.Context) (images []*Image) {
  updated := getLastUpdated(context)

  client := urlfetch.Client(context)
  resp, err := client.Get(fmt.Sprintf("%s%s?orderBy=\"created\"&startAt=%d&", FirebaseURL, ItemsPath, updated))
  if err != nil {
    context.Errorf("%s", err.Error())
    return images
  }

  results := make(map[string]*json.RawMessage)
  decoder := json.NewDecoder(resp.Body)
  err = decoder.Decode(&results)
  if err != nil {
    log.Print(err.Error())
    return
  }
  defer resp.Body.Close()

  images = make([]*Image, 0, 50)
  for key, message := range results {
    image := &Image{}
    data, _ := message.MarshalJSON()
    if err := json.Unmarshal(data, image); err == nil {
      images = append(images, image)
      // Debug
      context.Debugf("\nKey: %s\nUID: %d\nURL: %s\nCreated: %s\nNum tags: %d", 
        key, image.ID, image.URL, time.Unix(image.Created, 0).String(), len(image.Tags),
      )
    } else {
      log.Print(err.Error())
    }
  }

  return images
}

func getLastUpdated(context appengine.Context) (int64) {
  client := urlfetch.Client(context)
  resp, err := client.Get(FirebaseURL + ItemsPath)
  if err != nil {
    log.Print(err.Error())
    return 0
  }

  decoder := json.NewDecoder(resp.Body)
  defer resp.Body.Close()

  updated := &LastUpdate{}
  if err := decoder.Decode(updated); err == nil {
    return updated.Updated
  }
  return 0
}

func setLastUpdate(context appengine.Context, updated int64) {
  
  data := &ClosingBuffer{bytes.NewBufferString("")}
  encoder := json.NewEncoder(data)
  if err := encoder.Encode(&LastUpdate{Updated: updated}); err != nil {
    log.Print(err.Error())
    return
  }

  updateURL, _ := url.Parse(FirebaseURL + UpdatePath)
  request := &http.Request{
    Method: "Put",
    URL: updateURL,
    Body: data,
  }
  
  client := urlfetch.Client(context)
  if _, err := client.Do(request); err != nil {
    log.Print(err.Error())
  }
}
