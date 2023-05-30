package main

import (
	"encoding/json"
	"fmt"
  "math/rand"
  "net/http"
  "strings"
   b64 "encoding/base64"
   "github.com/google/uuid"
   "time"
   "unsafe"
   "strconv"
)

type version struct {
	Version     int
  Subversion string
  Codename   string
}

type siteLibrary struct {
	Uuid     string
  Name string
  NbElements int
}

const (
    letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 possibilities
    letterIdxBits = 6                    // 6 bits to represent 64 possibilities / indexes
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

type JsonElement struct {
      Uuid string 
      HadAttachment bool
      Attachment string
      Values map[string]string
} 


// see https://stackoverflow.com/questions/61930633/dynamically-encode-json-keys-in-go
// for generating json with 2 differents formats

// per https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func generate_random_word(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return *(*string)(unsafe.Pointer(&b))
}


func generate_random_bin_as_base64() string {
  // first we generate bunch of numbers
  data:="DATA"+generate_random_word(53+rand.Intn(22))
  return b64.StdEncoding.EncodeToString([]byte(data))
}


func generate_random_json_dict() map[string]string {
  
  // maybe location, value, serialNumber,clientName,flag
  var keys = []string{"location","value","serialNumber","clientName","flag","D53","F46"}
  nbKeys := rand.Intn(5) + 1
  
  // always id: string
  result:=map[string]string{"id": strconv.Itoa(rand.Intn(65535)+rand.Intn(250)+33)}
  
  var k int
  for i := 0; i < nbKeys; i++ {
    k=rand.Intn(len(keys))
    if   _,keyExists:=result[keys[k]]; keyExists {
      continue
    }
    result[keys[k]]=generate_random_word(rand.Intn(67) + rand.Intn(4)+2)
  }
  
  return result
}

func generate_one_element() JsonElement {
  id := uuid.New()
  var r JsonElement
  r.Uuid = id.String()
  if (rand.Intn(100)>67) {
    r.HadAttachment = true
    r.Attachment = generate_random_bin_as_base64()
  } else {
    r.Values = generate_random_json_dict()
  }
  
  return r
  
}

func generate_json(nbElements int, counter int) []byte {
  
  // TODO create an array of nBElements
	//data := map[string]interface{}{
  var data []JsonElement
  for i:= 0; i < nbElements; i++ {
    data = append(data, generate_one_element())
	}

	jsonData, err := json.Marshal(data)
  
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return nil
	}
  size := len(jsonData)
  where_to_cut := rand.Intn(size)
  subset := string(jsonData)[0:where_to_cut]


	//fmt.Printf("json data: %s\n", string(jsonData))
  //fmt.Printf(subset)
  if(counter%3==0) { return []byte(subset)}
  
  return jsonData
  
}

func VersionHandler() http.HandlerFunc{
  var v version;
  v.Version=rand.Intn(7)+2
  v.Subversion="47b2"
  v.Codename="Tumba Lobos"
  jsonVersion, err := json.Marshal(v)
  
	if err != nil {
		fmt.Printf("could not marshal version: %s\n", err)
		jsonVersion= nil
	}
  
  
  return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                w.Header().Set("Content-Type", "application/json")
                w.Write(jsonVersion)
       default:
                w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprintf(w, "I can't do that.")
        }
  }
}


func LibraryHandler(libs []siteLibrary) http.HandlerFunc{
  
  jsonData, err := json.Marshal(libs)
  
	if err != nil {
		fmt.Printf("could not marshal libraries: %s\n", err)
		jsonData= nil
	}
  
  
  return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                w.Header().Set("Content-Type", "application/json")
                w.Write(jsonData)

        default:
                w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprintf(w, "I can't do that.")
        }
  }
}

func OneLibraryHandler(libs []siteLibrary) http.HandlerFunc{
  counter:=0
  // create libraries map
  mapLibs := map[string]siteLibrary{}
  for _, element := range libs {
    mapLibs[element.Uuid] = element
  }
  return func(w http.ResponseWriter, r *http.Request) {
      var oneLib siteLibrary
      ok:=false
      if(len(r.URL.Path)>len("/lib/")) {
        parts := strings.SplitN(r.URL.Path[len("/lib/"):len(r.URL.Path)],"/",2)
        //fmt.Println(parts[0])
        oneLib,ok =mapLibs[parts[0]]
      } 

      switch r.Method {
        case "GET":
                if rand.Intn(10) >7 {
                  http.Error(w, "An Error occured on the server. Fugu meditating...", http.StatusInternalServerError)
                } else if ok {
                  w.Header().Set("Content-Type", "application/json")
                  
                  w.Write(generate_json(oneLib.NbElements,counter))
                  counter++
                } else {
                  // 404
                  http.Error(w, "404 not found.", http.StatusNotFound)
                }
        default:
                w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprintf(w, "I can't do that.")
        }
  }
}

func generateLibrary(nbLibrary int) []siteLibrary {
  var results []siteLibrary
  for i:= 0; i < nbLibrary; i++ {
    var library siteLibrary
    library.Uuid=uuid.New().String()
    library.NbElements = rand.Intn(7) + 2
    if rand.Intn(50)>40 {
      library.Name = "This one is a long one"
      library.NbElements+=23
    }
    results = append(results, library)
	}
  return results
}

func main() {
        libsConfig := generateLibrary(rand.Intn(3) + 3)
        http.HandleFunc("/libs/", LibraryHandler(libsConfig))
        
        http.HandleFunc("/lib/", OneLibraryHandler(libsConfig))
        
        http.HandleFunc("/version/", VersionHandler())
        
        http.HandleFunc("/healthz", k8sHealthHandler())

        fmt.Println("⇨ http serving on 8080")
        http.ListenAndServe(":8080", nil)
}