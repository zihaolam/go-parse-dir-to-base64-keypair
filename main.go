package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"sync"
)

type FileContentPair struct {
	name    string
	content *string
}

func main() {
	var (
		dirName string
	)

	flag.StringVar(&dirName, "dir", "d", "Directory to scan")
	flag.Parse()

	if dirName == "" {
		flag.PrintDefaults()
		return
	}

	contentCh := make(chan FileContentPair, 1)
	errorCh := make(chan error, 1)
	fmt.Println(dirName)
	paths, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		filepath := fp.Join(dirName, path.Name())
		go readAndConvertToBase64(filepath, contentCh, errorCh)
	}

	go func(contentCh chan FileContentPair, errorCh chan error, wg *sync.WaitGroup) {
		for {
			select {
			case content := <-contentCh:
				fmt.Println(content.name + "=" + *content.content)
				wg.Done()
			case err := <-errorCh:
				log.Fatal(err)
			}
		}
	}(contentCh, errorCh, &wg)

	wg.Wait()
}

func readAndConvertToBase64(filepath string, contentCh chan FileContentPair, errCh chan error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		errCh <- err
	}

	base64Encoded := base64.StdEncoding.EncodeToString(content)
	contentCh <- FileContentPair{
		content: &base64Encoded,
		name:    fp.Base(filepath),
	}
}
