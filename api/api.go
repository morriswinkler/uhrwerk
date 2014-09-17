// FabSmith - the Fab Lab Locksmith API layer.
package api

import (
  "strings"
  "net/url"
  "github.com/morriswinkler/uhrwerk/database"
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

// GetUser returns user data as JSON:
// {"status":"ok", "name":"name", etc...}
// Not implemented yet.
func (a *ApiCall) GetUser(sid string) string{
  return "{\"status\":\"not ok\"}"
}