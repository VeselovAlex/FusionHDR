package main

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

func getId() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(rand.Int63(), 16)
}

func deleteCache(id string) {
	path := "input/cache-" + id
	os.RemoveAll(path)
}

func cache(files []*multipart.FileHeader) (string, error) {
	id := getId()
	for !createDirIfNotExists("input/cache-" + id) {
		id = getId()
	}

	lst, err := os.Create(fmt.Sprintf("input/cache-%s/list.txt", id))
	if err != nil {
		return "", err
	}
	defer lst.Close()
	for i, f := range files {
		ext := f.Filename[strings.LastIndex(f.Filename, "."):]
		got, err := f.Open()
		if err != nil {
			return "", err
		}
		defer got.Close()

		cpy, err := os.Create(fmt.Sprintf("input/cache-%s/file-%d.%s", id, i, ext))
		if err != nil {
			return "", err
		}
		defer cpy.Close()
		_, err = io.Copy(cpy, got)
		fmt.Fprintln(lst, cpy.Name())
	}
	return id, nil
}

func createDirIfNotExists(path string) bool {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		os.Mkdir(path, 0777)
		return true
	}
	return false
}
