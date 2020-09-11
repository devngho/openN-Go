package iohelper

import (
	"log"
	"os"
)

func CreateFolder(name string, perm os.FileMode){
	err := os.Mkdir(name, perm)
	ErrLog(err)
}

func CreateFile(path string) *os.File {
	file, err := os.Create(path)
	if ErrLog(err){
		return nil
	}
	return file
}

func ErrLog(err error) bool{
	if err != nil {
		log.Fatal(err)
		return true
	}else {
		return false
	}
}