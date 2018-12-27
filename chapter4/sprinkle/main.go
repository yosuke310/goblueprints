package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const otherWord = "*"

var transforms []string

func main() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadFile(filepath.Dir(exe) + "/transforms.txt")
	if err != nil {
		log.Fatal(err)
	}
	transforms = strings.Split(string(bytes), "\n")

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
