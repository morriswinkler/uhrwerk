// Main config package for Fab Lab Locksmith
package config

import (
  "log"
  "code.google.com/p/gcfg"
)

// Main configuration struct
type Config struct {
  Log struct{
    LogFile string
  }

  Nfc struct {
    Basedir string
  }

  Database struct {
    Host, Port, Username, Password, DBName string
  }

  Webserver struct {
    Host, Port, Dir string
  }
}

// Loads and parses config ini file
func (c *Config) Init(configFile string) {
  err := gcfg.ReadFileInto(c, configFile)
  if err != nil {
    log.Fatalln("Failed to open config file ", configFile, ":", err)
  }
}