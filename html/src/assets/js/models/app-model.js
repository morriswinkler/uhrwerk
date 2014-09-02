"use strict";

function App() {
  var self = riot.observable(this),
    db = DB("app");

  self.load = function(page){
    page = $.trim(page);
    self.trigger('load:' + page);
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

  self.login = function(data) {
    var username = data.username;
    var password = data.password;

    console.log('username: ' + username);
    console.log('password: ' + password);

    // TODO: connect to database to sort this out
    if (username === 'asd@asd') {
      if (password === 'asd') {
        // TODO: create new user instance
        self.trigger('login:success', data);
      } else {
        self.trigger('login:fail', 'Wrong password.');
      }
    } else {
      self.trigger('login:fail', 'User does not exist.');
    }
  }

  self.logout = function() {
    self.trigger('logout');
  }
};
