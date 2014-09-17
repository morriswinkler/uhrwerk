// FabSmith - the Fab Lab Locksmith API layer.
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
  "github.com/morriswinkler/uhrwerk/database"
  "github.com/morriswinkler/uhrwerk/debug"
  "database/sql"
)

type Api struct {
  db *database.Database
}

// NewApi returns new Api instance. Pass a pointer to
// a Database instance as argument.
func NewApi(db *database.Database) *Api {
  return &Api{db: db}
}

// Executes a specific API call and returns response as JSON string.
// The string passed here has to be without the /api par at the beginning.
func (a *Api) Call(path, method string, vars *url.Values) string {
  
  // Route here
  // Split the URL for routing
  urlParts := strings.Split(path, "/")

  apiCall := NewApiCall(a, path, method, vars)

  switch urlParts[0] {
  case "auth":
    // Authentificate user
    return apiCall.Auth(method, vars)
    break
  case "machines":
    // Return a list of machines
    return apiCall.GetMachines(method, vars)
    break
  }

  return "{\"status\":\"error\", \"message\":\"No matching api call found\"}"
}

type ApiCall struct {
  api *Api
  path string
  method string
  vars *url.Values
}

// NewApiCall returns a new ApiCall struct instance that is being used to 
// execute Api calls.
func NewApiCall(api *Api, path, method string, vars *url.Values) *ApiCall {
  return &ApiCall{api: api, 
    path: path, 
    method: method, 
    vars: vars}
}

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

// GetUser returns user data as JSON:
// {"status":"ok", "name":"name", etc...}
// Not implemented yet.
func (a *ApiCall) GetUser(sid string) string{
  return "{\"status\":\"not ok\"}"
}

// GetMachines returns a list of machines available to the user that is
// currently logged in. Pass logged user session id as arg.
func (a *ApiCall) GetMachines(method string, vars *url.Values) string{

  if strings.ToLower(method) != "get" {
    return "{\"status\":\"error\",\"message\":\"Invalid request method\"}"
  }

  var db *sql.DB = a.api.db.GetHandle()

  // Get user ID
  sessionID := vars.Get("sessionID")
  var userID string
  err := db.QueryRow("SELECT user_id FROM sessions WHERE session_id=?", 
    sessionID).Scan(&userID)

  if err != nil {
    debug.ERROR.Printf("Could not get user ID for session %s: %s", 
      sessionID, 
      err)
    return "{\"status\":\"error\",\"message\":\"Session does not exist\"}"
  }

  // Get machine permissions for user id
  rows, err := db.Query("SELECT machine_id FROM permissions WHERE user_id=?", 
    userID)
  if err != nil {
    debug.ERROR.Printf("Could not get machine permissions for user ID %d: %s",
      userID, 
      err)
    return "{\"status\":\"error\",\"message\":\"Failed to get machine permissions\"}"
  }
  defer rows.Close()

  // Read machine IDs into an array
  var machineIDs []int
  var machineID int
  for rows.Next() {
    if err = rows.Scan(&machineID); err != nil {
      debug.ERROR.Printf("There was an error while getting machine permissions: %s", err)
      return "{\"status\":\"error\",\"message\":\"There was an error while getting machine permissions\"}"
    }
    // Make sure that we don't have duplicates
    dpl := false
    for _, v := range machineIDs {
      if v == machineID {
        dpl = true
        break
      }
    }
    if !dpl {
      // Append to array as no duplicates found
      machineIDs = append(machineIDs, machineID)
    }    
  } // rows.Close() is called automatically once rows.Next() becomes nil
  if err := rows.Err(); err != nil {
    rows.Close()
    debug.ERROR.Printf("There was an error with rows: %s", err)
    return "{\"status\":\"error\",\"message\":\"There was some error\"}"
  }

  // Notify if the user does not have any permissions
  if len(machineIDs) <= 0 {
    debug.INFO.Printf("User %d does not have the permission to use any of the machines", 
      userID)
    return "{\"status\":\"error\",\"message\":\"You do not have permissions to use any of the machines\"}"
  }

  // Get machines from database - construct query, prepare statement
  query := "SELECT machine_id, machine_name, machine_desc, available, calc_by_energy, calc_by_time, costs_per_kwh, costs_per_min FROM machines WHERE machine_id=?"
  stmt, err := db.Prepare(query)
  if err != nil {
    debug.ERROR.Printf("Error getting machines: %s", err)
    return "{\"status\":\"error\",\"message\":\"Could not get machines\"}"
  }
  response := "{\"status\":\"ok\",\"machines\":[\n"
  for i, v := range machineIDs {
    var machine_id int
    var machine_name string
    var machine_desc string
    var available bool
    var calc_by_energy bool
    var calc_by_time bool
    var costs_per_kwh float32
    var costs_per_min float32
    err := stmt.QueryRow(v).Scan(&machine_id, &machine_name, &machine_desc, 
      &available, &calc_by_energy, &calc_by_time, &costs_per_kwh, 
      &costs_per_min)
    if err != nil {
      debug.ERROR.Printf("Could not get machine data for ID %d: %s", v, err)
      return "{\"status\":\"error\",\"message\":\"Could not get machines\"}"
    }
    machineJSON := fmt.Sprintf("{\"machine_id\":%d, \"machine_name\":\"%s\", \"machine_desc\":\"%s\", \"available\":%t, \"calc_by_energy\":%t, \"calc_by_time\":%t, \"costs_per_kwh\":%f, \"costs_per_min\":%f}",
      machine_id, machine_name, machine_desc,
      available, calc_by_energy, calc_by_time,
      costs_per_kwh, costs_per_min)
    response = fmt.Sprintf("%s%s", response, machineJSON)
    if i < len(machineIDs) - 1 {
      response = fmt.Sprintf("%s%s", response, ",\n")
    }
  }
  response = fmt.Sprintf("%s%s", response, "\n]}")
  return response
}