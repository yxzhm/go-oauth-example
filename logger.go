package main

import (
	"log"
	"os"
)

var (
	Trace   *log.Logger // Just about anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // Critical problem
)

func init() {

	Trace = log.New(os.Stdout,
		"TRACE: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
}
