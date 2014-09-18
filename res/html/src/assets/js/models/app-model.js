"use strict";

function App() {
  var self = riot.observable(this),
      db = DB('app'),
      data = db.get();

  // Start the application by passing the start page name
  self.start = function(startPage) {
    self.load(startPage);
  }

  // Load page by passing the page name
  self.load = function(page){
    page = $.trim(page);
    self.trigger('load:' + page);
  }

  // Attempt to login a specific user
  // credentials contains username and password fields
  self.login = function(credentials) {

    // send ajax authenticate call
    var creds = {
      username: credentials.username,
      password: md5(credentials.password)
    }
    $.post('http://localhost:8080/api/auth', creds).done(function(data) {
      var o = $.parseJSON(data);

      // Parse server response
      if (o.status == 'ok') {
        $.cookie('fabsmith', o.sessionID, { expires: 1, path: '/' });
        self.trigger('login:success', o.sessionID);

        // Check for active bookings
        var args = {sessionID:o.sessionID};
        $.ajax({
          url: 'http://localhost:8080/api/machines/activated',
          type: 'GET',
          data: args,
          success: function(data) {
            var o = $.parseJSON(data);
            if (o.status == 'ok') {
              if (o.bookings) {
                self.trigger('load:active-bookings', o.bookings);
              } else {
                self.load('dashboard');
              }
            } else if (o.status == "error") {
              self.load('dashboard');
            } else {
              self.load('dashboard');
            }
          }
        });
      } else if (o.status == 'error') {
        self.trigger('login:fail', {
          message: o.message
        });
      } else {
        self.trigger('login:fail', {
          message: 'Some error occured'
        });
      }
    });
  }

  self.logout = function() {
    // Clear session cookie
    $.removeCookie('fabsmith');

    // And on the server side
    $.ajax({
      url: 'http://localhost:8080/api/auth',
      type: 'DELETE',
      success: function(data) {
        var o = $.parseJSON(data);
        if (o.status == "ok") {
          self.trigger('logout');
        } else {
          alert("Some error occured");
          self.trigger('logout');
        }
      }
    });
  }

  self.loadProducts = function() {

    // Attempt to get stored cookie
    var sessionID = $.cookie('fabsmith');
    var args = {sessionID:sessionID};
    $.ajax({
      url: 'http://localhost:8080/api/machines',
      type: 'GET',
      data: args,
      success: function(data) {
        var o = $.parseJSON(data);
        if (o.status == "ok") {
          var products = [];
          for (var i = 0; i < o.machines.length; i++) {
            products[i] = {
              id: o.machines[i].machine_id,
              name: o.machines[i].machine_name,
              status: o.machines[i].available,
              price: o.machines[i].calc_by_time ? 
                o.machines[i].costs_per_min : 
                o.machines[i].costs_per_kwh
            }
          }
          self.trigger('load:products', products);
        } else if (o.status == "error") {
          console.log("Error: " + o.message);
          //self.trigger('logout');
        } else {
          console.log("Some error occured");
        }
      }
    });
  }

  self.activateMachine = function(machineID) {
    // Attempt to get stored cookie
    var sessionID = $.cookie('fabsmith');
    var args = {sessionID:sessionID, machineID:machineID};
    $.ajax({
      // TODO: here we should call something like
      // /api/bookings/create or just /api/bookings and use
      // POST to create a new booking
      url: 'http://localhost:8080/api/machines/activate',
      type: 'POST',
      data: args,
      success: function(data) {
        var o = $.parseJSON(data);
        if (o.status == "ok") {
          self.trigger('activateMachine:success');
          self.load('message');
        } else if (o.status == "error") {
          self.trigger('activateMachine:fail', o.message);
        } else {
          self.trigger('activateMachine:fail', 'Some error occured');
        }
      }
    });
  };

  self.deactivateMachine = function(bookID) {
    var sessionID = $.cookie('fabsmith');
    var args = {sessionID:sessionID, bookID:bookID};
    $.ajax({
      // TODO: reorganize API so that we have
      // /api/bookings/xy - by using DELETE request we delete it
      url: 'http://localhost:8080/api/machines/deactivate',
      type: 'POST',
      data: args,
      success: function(data) {
        var o = $.parseJSON(data);
        if (o.status == "ok") {
          self.trigger("deactivateMachine:success", bookID);
        } else if (o.status == "error") {
          self.trigger("deactivateMachine:fail", o.message);
        } else {
          self.trigger("deactivateMachine:fail", "Some error occured");
        }
      }
    });
  };
};
