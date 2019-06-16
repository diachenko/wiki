package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	bolt "github.com/boltdb/bolt"
)

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

// Logger method for anything
func Logger(msg string, file string) {
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
	log.Println(msg + "\n")
	f.Close()
	return
}

// GenerateGUID generates UUID/GUID
func GenerateGUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func signUpEndpoint(w http.ResponseWriter, req *http.Request) {
	var user User
	var err Err
	json.NewDecoder(req.Body).Decode(&user)
	pass := md5.New()
	io.WriteString(pass, user.Pass)
	passHash := pass.Sum(nil)
	//is there user with same name?
	auth.DB.View(func(tx *bolt.Tx) error {
		bb := tx.Bucket([]byte("Users"))
		resp := bb.Get([]byte(user.Name))
		if resp != nil {
			err.Code = 500
			err.Text = "User already registered"
			str, _ := json.Marshal(err)
			http.Error(w, string(str), 500)
		}
		return nil
	})
	if err.Text != "" {
		return
	}
	auth.DB.Update(func(tx *bolt.Tx) error {
		users, _ := tx.CreateBucketIfNotExists([]byte("Users"))
		users.Put([]byte(user.Name), passHash)
		return nil
	})
	json.NewEncoder(w).Encode(user)
}

func getMdFile(name string) ([]byte, error) {
	dat, err := ioutil.ReadFile(name)
	return dat, err
}

func updateMdFile(name string, payload []byte) error {
	err := ioutil.WriteFile(name, payload, 0644)
	return err
}

//GET
func getFile(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	data, err := getMdFile(params["name"])
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
	json.NewEncoder(w).Encode(data)
	return
}

//PUT
func updateFile(w http.ResponseWriter, req *http.Request) {
	var md MD
	json.NewDecoder(req.Body).Decode(&md)
	err := updateMdFile(md.Name, []byte(md.Text))
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
	return
}

//POST
func createFile(w http.ResponseWriter, req *http.Request) {
	return
}

func signInEndpoint(w http.ResponseWriter, req *http.Request) {
	var user User
	var err Err
	var tt Token

	json.NewDecoder(req.Body).Decode(&user)
	pass := md5.New()
	io.WriteString(pass, user.Pass)
	passHash := pass.Sum(nil)

	auth.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		resp := b.Get([]byte(user.Name))
		if resp != nil {
			if bytes.Equal(resp, passHash) {
				t := make([]byte, 16)
				rand.Read(t)
				tt.Name = user.Name
				tt.Token = fmt.Sprintf("%X", t[0:16])
				tokens[tt.Token] = user.Name
			} else {
				err.Code = 500
				err.Text = "Wrong password"
				str, _ := json.Marshal(err)
				http.Error(w, string(str), 500)
				return nil
			}
		} else {
			err.Code = 500
			err.Text = "User not found"
			str, _ := json.Marshal(err)
			http.Error(w, string(str), 500)
		}
		return nil
	})
	if err.Text == "" {
		json.NewEncoder(w).Encode(tt)
	}
}

func initAuthBase() DBase {
	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Users"))
		return nil
	})
	return DBase{DB: db}
}

func main() {
	tokens = make(map[string]string)
	auth = initAuthBase()

	router := mux.NewRouter()
	router.HandleFunc("/signin", signInEndpoint).Methods("POST")
	router.HandleFunc("/signup", signUpEndpoint).Methods("POST")
	router.HandleFunc("/file/{name}", getFile).Methods("GET")
	router.HandleFunc("/file", updateFile).Methods("PUT")
	router.HandleFunc("/file", createFile).Methods("POST")
	log.Fatal(http.ListenAndServe(":1337", router))
}