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

    // For now we add the user checking logic here, but it should be done in
    // some way on the backend side. We need to get some kind of session ID
    // back from the backend that would be stored in the local storage in the 
    // browser.
    var userFound = false;
    for (var i = 0; i < data.users.length; i++) {

      // It does not matter if email or username is used for login
      if (data.users[i].username === credentials.username || 
          data.users[i].email === credentials.username) {

        userFound = true;
        // This should be done on the server side thou
        var md5Password = md5(credentials.password);

        // Check password
        if (data.auth[i].password === md5Password) {

          // We should do some kind of session management here, but for now
          // it is ok just like this and that

          self.trigger('login:success', data.users[i].user_id);

          // Load dashboard, not sure if this should be here
          self.load('dashboard');
        } else {
          self.trigger('login:fail', {
            message: "Wrong password."
          });
        }

        // Break the for loop and look no further for users that match
        break;
      }
    }

    if (!userFound) {
      self.trigger('login:fail', {
        message: "User not found."
      });
    }
  }

  self.logout = function() {
    self.trigger('logout');
  }

  self.loadProducts = function() {
    var products = [];

    for (var i = 0; i < data.machines.length; i++) {
      products[i] = {
        name: data.machines[i].machine_name,
        status: data.machines[i].available,
        price: data.machines[i].calc_by_time ? 
               data.machines[i].costs_per_min : 
               data.machines[i].costs_per_kwh
      }
    }
    
    self.trigger('load:products', products);
  }
};
