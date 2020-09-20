package iohelper

import (
	"log"
	"os"
)

func CreateFile(path string) *os.File {
	file, err := os.Create(path)
	if ErrLog(err){
		return nil
	}
	return file
}

func ErrLog(err error) bool{
	if err != nil {
		log.Println(err)
		return true
	}else {
		return false
	}
}

func ErrFatal(err error) bool{
	if err != nil {
		log.Fatal(err)
		return true
	}else {
		return false
	}
}