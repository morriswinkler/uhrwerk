package main

import (
	"code.google.com/p/gcfg"
	"log"
	"os"
)

func main() {
	var err error
	var cfg config
	var configFile string = "uhrwerk.ini"

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

	go httpdStart()

	for true {
		INFO.Println("running")

	}
}
