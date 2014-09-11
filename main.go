package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"github.com/morriswinkler/uhrwerk/config"
	"github.com/morriswinkler/uhrwerk/database"
	"github.com/morriswinkler/uhrwerk/debug"
	"github.com/morriswinkler/uhrwerk/http"
)

const ConfigFile string = "config.ini"

// Main function
func main() {
	// These are the main building blocks of our Fab Lab Locksmith solution
	var cfg *config.Config
	var db *database.Database
	var web *http.Server

	// Read config file
	cfg = new(config.Config)
	cfg.Init(ConfigFile)

	// Define log file
	logFile, err := os.OpenFile(cfg.Log.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file ", logFile.Name, ":", err)
	}

	debug.Init(logFile, logFile, logFile, logFile)

	// nfc_c := make(chan string)
	// go authenticate_nfc(cfg.Nfc.Basedir, nfc_c)

	// for true {
	// 	fmt.Println(<-nfc_c)
	// Once a NFC ID is got, we check the db which user has it
	// Log in the user here, pass the session id to the client side
	// And this happens only with a single registered IP of the
	// main terminal device
	// }

	//DBInit("tcp(127.0.0.1:3306)", "root", "root", "test")
	//_, err = DBTest()
	db = new(database.Database)
	web = new(http.Server)

	// Extract database related values from the ini file
	host := fmt.Sprintf("tcp(%s:%s)", 
		cfg.Database.Host, 
		cfg.Database.Port)
	username := cfg.Database.Username
	password := cfg.Database.Password

	db.Init(host, username, password, "test")
	_, err = db.Test()
	if err != nil {
		debug.ERROR.Printf("DBTest failed: %s", err)
		log.Printf("DBTest failed: %s\n", err)
		os.Exit(1)
	} else {
		debug.INFO.Println("DBTest passed!");
		log.Println("DBTest passed!")

		// Init web server
		err = web.Init(cfg.Webserver.Host, cfg.Webserver.Port, cfg.Webserver.Dir)
		if err != nil {
			debug.ERROR.Printf("Could not start webserver: ", err)
			log.Printf("Could not start webserver: ", err)
			os.Exit(1)
		}
	}
	defer db.Close()

	for {
		debug.INFO.Println("running")
		time.Sleep(time.Minute)
	}
}
