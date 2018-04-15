package util

// Just some auxilliary functions

import (
	"log"
)

// Checkerr is an antipettern and will be removed
func Checkerr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
