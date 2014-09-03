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

        // Check password
        if (data.users[i].password === credentials.password) {

          // We should do some kind of session management here, but for now
          // it is ok just like this and that

          // Trigger login:success event and pass the user ID (for now i - the
          // iterator can be the user ID)
          self.trigger('login:success', i);

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

  self.loadProducts = function(){
    var products = [
      {name: 'Laser Cutter', status:'available', price:'1 €/h'},
      {name: '3D Printer', status:'available', price:'1 €/h'},
      {name: 'CNC Mill', status:'unavailable', price:'1 €/h'},
      {name: 'Hand Drill', status:'used', price:'1 €/h'},
      {name: 'Maker Bot', status:'Available', price:'1 €/h'}
    ];
    self.trigger('load:products', products);
  }
};
