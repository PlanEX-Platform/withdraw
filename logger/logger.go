package logger

import (
	"log"
	"flag"
	"os"
	"strconv"
	"time"
)

var (
	Log *log.Logger
)

func init() {
	// set location of log file
	//var logpath = "/root/logs/sluise-" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"
	var logpath = "./withdraw-" + strconv.FormatInt(time.Now().Unix(), 10) + ".log"

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)
}
