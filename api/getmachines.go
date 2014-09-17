package api

import (
  "strings"
  "net/url"
  "database/sql"
  "fmt"
  "github.com/morriswinkler/uhrwerk/debug"
)

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