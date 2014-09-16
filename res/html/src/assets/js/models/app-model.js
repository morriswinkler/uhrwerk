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
      console.log(o);

      // Parse server response
      if (o.status == 'ok') {
        $.cookie('fabsmith', o.sessionID, { expires: 1, path: '/' });
        self.trigger('login:success', o.sessionID);
        self.load('dashboard');
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
    /*
    $.ajax({
      url: 'http://localhost:8080/api/auth',
      type: 'DELETE',
      success: function(data) {
        var o = $.parseJSON(data);
      }
    });
    */

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
