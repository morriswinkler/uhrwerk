package main

import (
	"code.google.com/p/gcfg"
	"log"
	"os"
	"time"
)

func main() {
	var err error
	var cfg config
	var configFile string = "uhrwerk.ini"
	var db *Database

	err = gcfg.ReadFileInto(&cfg, configFile)
	if err != nil {
		log.Fatalln("Failed to open config file ", configFile, ":", err)
	}

	logFile, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file ", logFile.Name, ":", err)
	}

	Init(logFile, logFile, logFile, logFile)

	// nfc_c := make(chan string)
	// go authenticate_nfc(cfg.Nfc.Basedir, nfc_c)

	// for true {
	// 	fmt.Println(<-nfc_c)
	// }

	//DBInit("tcp(127.0.0.1:3306)", "root", "root", "test")
	//_, err = DBTest()
	db = new(Database)
	db.Init("tcp(127.0.0.1:3306)", "root", "root", "test")
	_, err = db.Test()
	if err != nil {
		ERROR.Printf("DBTest failed: %s", err)
		TRACE.Println("Exiting...")
		log.Printf("DBTest failed: %s\n", err);
		log.Println("Exiting...");
		os.Exit(1)
	} else {
		TRACE.Println("DBTest passed!");
		log.Println("DBTest passed!")
		httpdStart()
	}
	defer db.Close()

	for {
		INFO.Println("running")
		time.Sleep(time.Minute)
	}
}
