// FabSmith - the Fab Lab Locksmith API layer.
// The following link seems a better way to do this: 
// https://github.com/astaxie/build-web-application-with-golang/blob/master/en/eBook/08.3.md
package api

import (
  "strings"
  "net/url"
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
  case "machines":
    // Return a list of machines
    if len(urlParts) > 1 {
      switch urlParts[1] {
      case "activate":
        return apiCall.ActivateMachine(method, vars)
      case "activated":
        return apiCall.GetActivatedMachines(method, vars)
      case "deactivate":
        return apiCall.DeactivateMachine(method, vars)
      default:
        return "{\"status\":\"error\", \"message\":\"No matching api call found\"}"
      }
    }
    return apiCall.GetMachines(method, vars)
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

// GetUser returns user data as JSON:
// {"status":"ok", "name":"name", etc...}
// Not implemented yet.
func (a *ApiCall) GetUser(sid string) string{
  return "{\"status\":\"not ok\"}"
}

// Gets user ID by using existing session ID as argument
func (a *ApiCall) GetUserID(sessionID string) (int, error) {
  var db *sql.DB = a.api.db.GetHandle()

  var userID int
  err := db.QueryRow("SELECT user_id FROM sessions WHERE session_id=?", 
    sessionID).Scan(&userID)

  if err != nil {
    debug.ERROR.Printf("Could not get user ID for session %s: %s", 
      sessionID, 
      err)
    return 0, err
  }

  return userID, nil
}