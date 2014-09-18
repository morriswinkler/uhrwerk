package api

import (
  "strings"
  "net/url"
  "database/sql"
  "fmt"
  "time"
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

// ActivateMachine activates a machine. Pass sessionID and machineID as 
// function arguments. Available through /api/machines/activate
// and accepts POST method variables
func (a *ApiCall) ActivateMachine(method string, vars *url.Values) string {
  /*
  // Dissallow any other request method other than POST
  if strings.ToLower(method) == "post" {
    debug.ERROR.Printf("Request method is not POST")
    return "{\"status\":\"error\",\"message\":\"Failed to activate\"}"
  }
  */

  sessionID := vars.Get("sessionID")
  machineID := vars.Get("machineID")
  db := a.api.db.GetHandle()

  // Get user ID from session
  userID, err := a.GetUserID(sessionID)
  if err != nil {
    debug.ERROR.Printf("Could not get user ID: %s", err)
    return "{\"status\":\"error\", \"message\":\"Failed to activate machine\"}"
  }

  // Check if machine with the given ID exists
  var machineName string
  query := "SELECT machine_name FROM machines WHERE machine_id=?"
  err = db.QueryRow(query, machineID).Scan(&machineName)
  if err != nil {
    debug.ERROR.Printf("Machine with ID %d does not exist: %s", machineID, err)
    return "{\"status\":\"error\", \"message\":\"Failed to activate machine\"}"
  }
  if machineName == "" {
    debug.ERROR.Printf("Machine with ID %d does not exist", machineID)
    return "{\"status\":\"error\", \"message\":\"Failed to activate machine\"}"
  }

  // Check if there is no active booking in the database already
  var bookExists bool = a.HasActiveBookings(userID)
  if bookExists {
    debug.ERROR.Printf("There is an active booking for user %d already", userID)
    return "{\"status\":\"error\", \"message\":\"There is an active booking already\"}"
  }

  // Create new booking for machine and activate it
  query = "INSERT INTO bookings VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)" // 14 vals
  _, err = db.Exec(query,
    nil, nil, userID, machineID, true, 
    time.Now().Format("2006-01-02 15:04:05"),
    nil, 0, 0, 0, 0, "", false, false)
  if err != nil {
    debug.ERROR.Printf("Could not add new booking: %s", err)
    return "{\"status\":\"error\", \"message\":\"Failed to activate machine\"}"
  }

  // Success, share it with others
  return "{\"status\":\"ok\"}"
}

// Checks if an user has at least one active booking
func (a *ApiCall) HasActiveBookings(userID int) bool {
  db := a.api.db.GetHandle()
  var bookExists bool
  query := "SELECT EXISTS (SELECT book_id FROM bookings WHERE user_id=? AND active=1)"
  err := db.QueryRow(query, userID).Scan(&bookExists)
  if err != nil {
    debug.ERROR.Printf("Could not exec query: %s", err)
    return false
  }
  return bookExists
}

func (a *ApiCall) GetActivatedMachines(method string, vars *url.Values) string {
  sessionID := vars.Get("sessionID")
  db := a.api.db.GetHandle()

  // Get user ID from session
  userID, err := a.GetUserID(sessionID)
  if err != nil {
    debug.ERROR.Printf("Could not get user ID: %s", err)
    return "{\"status\":\"error\", \"message\":\"Could not get user ID\"}"
  }

  // Check if at least one booking exists
  var bookExists bool = a.HasActiveBookings(userID)
  
  // Get list of active bookings if exist
  if bookExists {
    query := "SELECT b.book_id, b.machine_id, b.time_start, m.machine_name FROM bookings b, machines m WHERE b.user_id=? AND b.active=1 AND b.machine_id=m.machine_id"
    rows, err := db.Query(query, userID)
    if err != nil {
      debug.ERROR.Printf("Could not select bookings for user %d: %s",
        userID, err)
      return "{\"status\":\"error\",\"message\":\"Failed to get bookings\"}"
    }
    defer rows.Close()

    response := "{\"status\":\"ok\", \"bookings\":["
    for rows.Next() {
      var book_id int
      var machine_id int
      var time_start string
      var machine_name string
      if err = rows.Scan(&book_id, &machine_id, &time_start, &machine_name); err != nil {
        debug.ERROR.Printf("Could not get booking record: %s", err)
        return "{\"status\":\"error\",\"message\":\"Could not get active bookings\"}"
      }
      bookingString := fmt.Sprintf("{\"book_id\":%d, \"machine_id\":%d, \"time_start\":\"%s\", \"machine_name\":\"%s\", \"time_now\":\"%s\"},", 
        book_id, 
        machine_id, 
        time_start,
        machine_name, 
        time.Now().Format("2006-01-02 15:04:05")) 
        // Time now is for reference as the time on the client side can be 
        // f**d up.
      response = fmt.Sprintf("%s%s", response, bookingString)
    }

    if err := rows.Err(); err != nil {
      rows.Close()
      debug.ERROR.Printf("There was an error with rows: %s", err)
      return "{\"status\":\"error\",\"message\":\"There was some error\"}"
    }

    // Remove last comma from response string, add tail and return
    response = strings.TrimSuffix(response, ",")
    response = fmt.Sprintf("%s%s", response, "]}")
    return response
  }

  return "{\"status\":\"ok\"}"
}

// This should be called DeleteBooking or similar. 
// Removes active booking that has been created by user.
func (a *ApiCall) DeactivateMachine(method string, vars *url.Values) string {
  sessionID := vars.Get("sessionID")
  bookID := vars.Get("bookID")
  db := a.api.db.GetHandle()

  // Get user ID from session
  userID, err := a.GetUserID(sessionID)
  if err != nil {
    debug.ERROR.Printf("Could not get user ID: %s", err)
    return "{\"status\":\"error\", \"message\":\"Could not get user ID\"}"
  }

  query := "DELETE FROM bookings WHERE book_id=? AND user_id=?"
  _, err = db.Exec(query, bookID, userID)
  if err != nil {
    debug.ERROR.Printf("There was an error while deleting booking %d user %d: %s",
      bookID, userID, err)
    return "{\"status\":\"error\", \"message\":\"Error while deleting booking\"}"
  }

  return "{\"status\":\"ok\"}"
}