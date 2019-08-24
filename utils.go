package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

// Logger method for anything
func Logger(msg string) {
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
	log.Println(msg)
	f.Close()
	return
}

// GenerateGUID generates UUID/GUID
func GenerateGUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func getMdFile(name string) ([]byte, error) {
	dat, err := ioutil.ReadFile(Path + name + ".md")
	return dat, err
}

func updateMdFile(name string, payload []byte) error {
	err := ioutil.WriteFile(Path+name+".md", payload, 0644)
	return err
}

func generateImageName() string {
	t := time.Now()
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d-%X",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), b[0:3])
}

func updateHugo() error {
	cmd := exec.Command("hugo")
	log.Printf("Running hugo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	print(cmd.Stdout)
	err := cmd.Run()
	if err != nil {
		log.Printf("Hugo finished with error: %v", err)
		return err
	}
	return nil
}

func setHugoPrefix(name string) string {
	return fmt.Sprintf("---\ntitle: \"%s\"\ndate: %v\n---\n\n", name, time.Now().Format(time.RFC3339))
}
