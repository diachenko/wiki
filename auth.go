package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	bolt "github.com/boltdb/bolt"
)

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
	db, err := bolt.Open("db/users.db", 0600, nil)
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Users"))
		return nil
	})
	return DBase{DB: db}
}
