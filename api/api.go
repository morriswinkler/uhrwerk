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
    return apiCall.Auth("username", "hashmash")
    break
  case "machines":
    // Return a list of machines
    return apiCall.GetMachines("tempsessionid")
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
func (a *ApiCall) Auth(username, passhash string) string{
  return "{\"status\":\"ok\", \"sid\":\"md5hash\"}";
}

// GetUser returns user data as JSON:
// {"status":"ok", "name":"name", etc...}
// Not implemented yet.
func (a *ApiCall) GetUser(sid string) string{
  return "{\"status\":\"not ok\"}"
}

// GetMachines returns a list of machines available to the user that is
// currently logged in. Pass logged user session id as arg.
func (a *ApiCall) GetMachines(sid string) string{
  return "{\"status\":\"ok\", \"machines\":[{\"name\":\"Machine Name\", \"other\":\"Other things\"}]}"
}