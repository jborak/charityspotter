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
  // Setup URI functions
  http.HandleFunc("/api/hello/", hello)
  http.HandleFunc("/api/index/", updateAndIndex)
  http.HandleFunc("/api/search/", searchIndex)
  http.HandleFunc("/api/debug/", debugIndex)
}

/////
// Data
////

// Image stores information about images created.
type Image struct {
  ID      string   `json:"uid"`
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
  // ID      string    `json:""`
  Data    string    `json:"data"`
  URL     string    `json:"url"`
  Created time.Time `json:"created"`
}

func (this *ImageDoc) String() (string) {
  return fmt.Sprintf("Data: %s, URL: %s", this.Data, this.URL)
}

////
// API
////
func hello(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Hello, world!")
}

func debugIndex(w http.ResponseWriter, r *http.Request) {
  context := appengine.NewContext(r)

  // Fetch document index.
  index, _ := search.Open(ImageIndex)
  num := 0
  for t := index.List(context, nil); ; {
    doc := &ImageDoc{}
    _, err := t.Next(doc)
    if err == search.Done {
      break
    }
    if err != nil {
      fmt.Fprintf(w, err.Error())
      break
    }
    fmt.Fprintf(w, doc.String())
    num++
  }

  fmt.Fprintf(w, "%d documents in index", num)

  updated := getLastUpdated(context)
  fmt.Fprintf(w, "%d last updated", updated)
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
  } else {
    // Hack for now.
    setLastUpdate(context, time.Now().Unix()) 
  }
}

type Query struct {
  Terms string `json:"terms"`
}

func searchIndex(w http.ResponseWriter, r *http.Request) {
  context := appengine.NewContext(r)

  // Open document index.
  index, err := search.Open(ImageIndex)
  if err != nil {
    context.Errorf("%s", err.Error())
    // fmt.Fprintf(w, "error: %s", err.Error())
    return
  }

  // Get search terms.
  decoder := json.NewDecoder(r.Body)
  query := &Query{}
  if err := decoder.Decode(query); err != nil {
    context.Errorf("%s", err.Error())
    // fmt.Fprintf(w, "error: %s", err.Error())
    return
  }
  defer r.Body.Close()

  //fmt.Fprintf(w, "%s", query.Terms)


  // Execute the query using the terms from the request and saves the images
  images := make([]*ImageDoc, 0, 50)
  for t := index.Search(context, query.Terms, nil); ; {
    doc := &ImageDoc{}
    _, err := t.Next(doc)
    if err == search.Done {
      break
    }
    if err != nil {
      break
    }
    images = append(images, doc)
    // fmt.Fprintf(w, "Image: %s, URL: %s<br>", id, doc.URL)
  }

  // Encode images into json and send them out.
  encoder := json.NewEncoder(w)
  if err := encoder.Encode(images); err != nil {
    context.Errorf("%s", err.Error())
    fmt.Fprintf(w, "%s", err.Error())
  }
}

////
// Firebase funcs
////
const (
  // Firebase URLs
  FirebaseURL string = "https://charitysandbox.firebaseio.com"
  ItemsPath   string = "/items.json"  // GET
  UpdatePath  string = "/settings/update.json" // PUT

  // AppEngine constants
  ImageIndex string = "ImageIndex2"
)

func getLatestImages(context appengine.Context) (images []*Image) {
  updated := getLastUpdated(context) + 1

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
  resp, err := client.Get(FirebaseURL + UpdatePath)
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
    Method: "PUT",
    URL: updateURL,
    Body: data,
  }
  
  client := urlfetch.Client(context)
  if _, err := client.Do(request); err != nil {
    log.Print(err.Error())
  }
}
