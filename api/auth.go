package api

import (
  "io"
  "fmt"
  "time"
  "strconv"
  "strings"
  "net/url"
  "math/rand"
  "crypto/md5"
  "github.com/morriswinkler/uhrwerk/debug"
  "database/sql"
)

// Authentificates a user and returns JSON:
// {"status":"ok", "sid":"md5hash"} on success
// {"status":"error", "message":"why"} on failure
// One has to provide username and hashed password as args.
func (a *ApiCall) Auth(method string, vars *url.Values) string {

  var db *sql.DB = a.api.db.GetHandle()

  if strings.ToLower(method) == "delete" {
    
    // Destroy session
    _, err := db.Exec("DELETE FROM sessions WHERE session_id=?",
      vars.Get("sessionID"))

    if err != nil {
      debug.ERROR.Printf("Failed to delete session ID %s: %s", 
        vars.Get("sessionID"),
        err)
    }

    debug.INFO.Printf("A session was destroyed")

    // And exit here
    return "{\"status\":\"ok\"}"
  }

  // TODO: check method and allow only POST method to pass through

  var username, password string
  username = vars.Get("username")
  password = vars.Get("password") // this will be a md5 string in our case

  if username == "" {
    debug.ERROR.Printf("Username not set")
    return "{\"status\":\"error\", \"message\":\"Username not set\"}"
  }

  if password == "" {
    debug.ERROR.Printf("No password")
    return "{\"status\":\"error\", \"message\":\"No password\"}"
  }

  // Get user ID from db
  
  var user_id int
  err := db.QueryRow("SELECT user_id FROM users WHERE username=? OR email=?", 
    username, username).Scan(&user_id)
  
  switch {
  case err == sql.ErrNoRows:
    debug.ERROR.Printf("User %s not found", username)
    return "{\"status\":\"error\", \"message\": \"No user found\"}"
  case err != nil:
    debug.ERROR.Printf("Could not get user_id: %s", err)
    return "{\"status\":\"error\", \"message\":\"There was an error\"}"
  }

  // Check if passwords match
  var user_pass string
  err = db.QueryRow("SELECT password FROM auth WHERE user_id=?", 
    user_id).Scan(&user_pass)

  switch {
  case err == sql.ErrNoRows:
    debug.ERROR.Printf("Could not get password for user %s with user_id %d", 
      username, user_id)
    return "{\"status\":\"error\", \"message\":\"Could not get password\"}"
  case err != nil:
    debug.ERROR.Printf("Could not get password: %s", err)
    return "{\"status\":\"error\", \"message\":\"There was an error\"}"
  }

  if (password == user_pass) {

    // Success - create session ID
    var tstamp string = strconv.FormatInt(time.Now().Unix(), 10)
    rand.Seed(time.Now().UnixNano())
    var random string = strconv.FormatInt(rand.Int63n(time.Now().UnixNano()), 10)
    var combi string = fmt.Sprintf("%s-%s-%s", 
      username, 
      tstamp, 
      random)
    h := md5.New()
    io.WriteString(h, combi)
    var sessionID string = fmt.Sprintf("%x", h.Sum(nil))

    // Store session ID
    //var expiryTimestamp int64 = time.Now().Unix() + 3600 // 1 hour from now
    _, err = db.Exec("INSERT INTO sessions VALUES (?, ?, ?, ?)", 
      nil, 
      user_id, 
      sessionID,
      time.Now().Format("2006-01-02 15:04:05"))

    if err != nil {
      debug.ERROR.Printf("Failed to store session for user: %s: %s", 
        username, 
        err)
      return "{\"status\":\"error\", \"message\":\"There was an error\"}"
    }

    // Create response
    var response string = fmt.Sprintf("{\"status\":\"ok\", \"sessionID\":\"%s\"}", 
      sessionID)
    debug.INFO.Printf("User %s successfully authenticated", username)
    return response;
  } else {
    debug.ERROR.Printf("User %s failed to authenticate", username)
    return "{\"status\":\"error\", \"message\":\"Failed to authenticate\"}"
  }

  // TODO: think if we really need to get content from two tables 
  // (users and auth) in order to authenticate
}