package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	bolt "github.com/boltdb/bolt"
)

//Path contains route to md files.
const Path = "content/content/"

//Max uploading file size in bytes
const maxFileSize = 10 * 1024 * 1024 // 10mb

//DBase type used for storing BoltDB instance
type DBase struct {
	DB       *bolt.DB
	Settings map[string]string
}

//User - name/pass - used for login/signup
type User struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

//Token used for tokens array.
type Token struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

//Err used for error handling in http requests
type Err struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// MD struct used for file-related queries
type MD struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

var dB DBase
var auth DBase

var tokens map[string]string

// EnableCors Adds CORS to header
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// CORSHandler getting everythign work for CORS
func CORSHandler(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w, req)
}

// getting everythign work for CORS
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	updateHugo()
	tokens = make(map[string]string)
	auth = initAuthBase()
	setHugoPrefix("newline")
	router := mux.NewRouter()
	router.HandleFunc("/", getReadme).Methods("GET")
	router.HandleFunc("/signin", signInEndpoint).Methods("POST")
	router.HandleFunc("/signin", CORSHandler).Methods("OPTIONS")
	router.HandleFunc("/signup", signUpEndpoint).Methods("POST")
	router.HandleFunc("/signup", CORSHandler).Methods("OPTIONS")
	router.HandleFunc("/article", CORSHandler).Methods("OPTIONS")
	router.HandleFunc("/article/{name}", getArticle).Methods("GET")
	router.HandleFunc("/article", updateArticle).Methods("PUT")
	router.HandleFunc("/article", createArticle).Methods("POST")
	router.HandleFunc("/file", CORSHandler).Methods("OPTIONS")
	router.HandleFunc("/file", uploadFile).Methods("POST")
	router.HandleFunc("/section", createSection).Methods("POST")
	router.HandleFunc("/section", CORSHandler).Methods("OPTIONS")
	log.Fatal(http.ListenAndServe(":1337", router))
}
