package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
