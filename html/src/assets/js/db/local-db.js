"use strict";

function DB(key) {
	var store = window.localStorage;

  // Fill with test values
  var users = [
    {
      username: 'kris', 
      email: 'krisjanis.rijnieks@gmail.com', 
      password: 'asd123'
    }, 
    {
      username: 'alien', 
      email: 'alien@space.com',
      password: 'bada55'
    }
  ];

  var initialData = {
    users: users
  };

  store[key] = JSON.stringify(initialData);

	return {
		get: function() {
			return JSON.parse(store[key] || '{}')
		},

		put: function(data) {
			store[key] = JSON.stringify(data)
		}
	};
};
