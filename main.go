package main

import (
    "time"
    "log"
    "io"
    "encoding/json"
    "net/http"
    "bytes"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "strconv"
    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var (
    session *mgo.Session
    collection *mgo.Collection
    collectionT *mgo.Collection
)

type (

    // Tracker has information about identification.
    Tracker struct {
        ID  bson.ObjectId `bson:"_id,omitempty"`
        Name string `bson:"name" json:"name"`
        CreateAt time.Time `bson:"create_at" json:"create_at"`
    }

    // Tracked has information about UserAgent.
    Tracked struct {
        IDTracker bson.ObjectId `bson:"id_tracker,omitempty"`
        UserAgent string `bson:"user_agent"`
        ReadEmail bool `bson:"read_email"`
        CameToSite bool `bson:"came_to_site"`
        CreateAt time.Time `bson:"create_at"`
    }

    // ResponseJSON response before Insert or Error
    ResponseJSON struct {
        Message string `json:"message"`
    }
)

func insertTracked(trackID string, userAgent string){
    tracker := Tracker{}
    err := collection.FindId(bson.ObjectIdHex(trackID)).One(&tracker)

    if err != nil {
        log.Printf("Error: %s - TrackID: %s", err, trackID)
    } else {
        err = collectionT.Insert(&Tracked{tracker.ID, userAgent, true, false, time.Now()})

        if err != nil {
            log.Printf("Error: %s", err)
        } else {
            log.Println("Created successfully")
        }
    }
}

func index(response http.ResponseWriter, request *http.Request){
    io.WriteString(response, "Home Tracker")
}

func track(response http.ResponseWriter, request *http.Request){
    vars := mux.Vars(request)
    trackID := vars["track_id"]

    rgb := image.NewRGBA(image.Rect(0, 0, 1, 1))
    black := color.RGBA{0, 0, 0, 255}
    draw.Draw(rgb, rgb.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

    var img image.Image = rgb

    buffer := new(bytes.Buffer)
    if err := jpeg.Encode(buffer, img, nil); err != nil {
        log.Println("unable to encode image.")
    }

    response.Header().Set("Content-Type", "image/jpeg")
    response.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    if _, err := response.Write(buffer.Bytes()); err != nil {
        log.Println("unable to write image.")
    }

    go insertTracked(trackID, request.Header.Get("User-Agent"))
}

func addTrack(response http.ResponseWriter, request *http.Request){
    name := request.FormValue("name")

    // Create collection
    err := collection.Insert(&Tracker{Name: name, CreateAt: time.Now()})

    if err != nil {
        log.Fatal(err)
    }

    j, err := json.Marshal(ResponseJSON{"Created successfully"})
    if err != nil {
        panic(err)
    }

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusCreated)
    response.Write(j)
}


func main(){
    // Instance mux route
    route := mux.NewRouter()

    route.HandleFunc("/", index).Methods("GET")
    route.HandleFunc("/{track_id}.jpg", track).Methods("GET")
    route.HandleFunc("/add", addTrack).Methods("POST")

    // Connect MongoDb
    session, err := mgo.Dial("mongodb://localhost:27017")

	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    collection = session.DB("stalkers").C("tracker")
    collectionT = session.DB("stalkers").C("tracked")

    http.ListenAndServe(":5000", route)
}
