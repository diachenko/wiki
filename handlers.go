package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/h2non/filetype"
)

func getReadme(w http.ResponseWriter, req *http.Request) {
	file, _ := ioutil.ReadFile("README.md")
	output := markdown.ToHTML(file, nil, nil)
	w.Write(output)
}

//GET
func getArticle(w http.ResponseWriter, req *http.Request) {
	EnableCors(&w)
	params := mux.Vars(req)
	data, err := getMdFile(params["name"])
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
	var md MD
	tmp := strings.SplitAfterN(string(data), "\n---", 2)[1]

	md.Name = params["name"]
	md.Text = strings.TrimSpace(tmp)
	json.NewEncoder(w).Encode(md)
	return
}

//PUT
func updateArticle(w http.ResponseWriter, req *http.Request) {
	EnableCors(&w)
	var md MD
	err := json.NewDecoder(req.Body).Decode(&md)
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 400)
		return
	}
	err = updateMdFile(md.Name, []byte(md.Text))
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
	updateHugo()
	return
}

//POST
func createArticle(w http.ResponseWriter, req *http.Request) {
	EnableCors(&w)
	var md MD

	err := json.NewDecoder(req.Body).Decode(&md)
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 400)
		return
	}

	if md.Name == "" {
		http.Error(w, "Article name should be filled", 400)
		return
	}

	_, err = os.Create("content/" + md.Name + ".md")

	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
	md.Text = setHugoPrefix(md.Name) + md.Text

	updateMdFile(md.Name, []byte(md.Text))
	updateHugo()
	return
}

//File upload (POST)
func uploadFile(w http.ResponseWriter, req *http.Request) {
	Logger("Upload file started")
	// limiting the file size
	req.ParseMultipartForm(maxFileSize)
	file, handler, err := req.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "Error receiving the file", 500)
		Logger("Error receiving the file")
		return
	}

	Logger(fmt.Sprintf("File: %+v", handler.Filename))
	Logger(fmt.Sprintf("Size: %+v", handler.Size))
	fmt.Println(file)
	// try to read file and check is it an image
	tempfile, _ := ioutil.ReadAll(file)
	fmt.Println(!filetype.IsImage(tempfile))
	if !filetype.IsImage(tempfile) {
		Logger("Not an image")
		http.Error(w, "Not an image", 400)

		return
	}
	//defining type to save file with correct extension (just in case)
	kind, _ := filetype.Match(tempfile)
	filename := generateImageName() + "." + kind.Extension

	f, err := os.Create("img/" + filename)
	if err != nil {
		Logger("Unable to create file: " + filename)
		http.Error(w, "Unable to create file", 500)
		return
	}
	defer f.Close()

	_, err = f.Write(tempfile)
	if err != nil {
		Logger("Unable to write file: " + filename)
		http.Error(w, "Unable to create file", 500)
		return
	}

	Logger("File " + filename + " successfully uploaded")
	json.NewEncoder(w).Encode(filename)
}

func createSection(w http.ResponseWriter, req *http.Request) {
	var md MD
	err := json.NewDecoder(req.Body).Decode(&md)
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
	err = os.Mkdir(md.Name, 0666)
	if err != nil {
		str, _ := json.Marshal(err)
		http.Error(w, string(str), 500)
		return
	}
}
