"use strict";

// Maybe this part should handle also the db communication logic.

function DB(key) {
	var store = window.localStorage;

  // Users table
  var users = [
    {
      user_id: 1,
      first_name: 'Krisjanis',
      last_name: 'Rijnieks',
      username: 'kris', 
      email: 'krisjanis.rijnieks@gmail.com'
      // There are more fields, but it's not needed in the testing phase so
      // I'm going to skip them.
    }, 
    {
      user_id: 2,
      first_name: 'Kruger',
      last_name: 'Ultimus',
      username: 'kruxy', 
      email: 'kru@xy.io',
    }
  ];

  // Authentification table
  var auth = [
    {
      user_id: 1,
      password: 'bfd59291e825b5f2bbf1eb76569f8fe7' // md5 of asd123
    },
    {
      user_id: 2,
      password: 'bfd59291e825b5f2bbf1eb76569f8fe7' // md5 of asd123
    }
  ];

  // Machines table
  var machines = [
    {
      machine_id: 1,
      machine_name: 'i3Berlin 3D Printer',
      machine_desc: 'The tools you make. Your tools, your make.',
      available: true,
      unavail_msg: '',
      unavail_till: 0,
      calc_by_energy: false,
      calc_by_time: true,
      costs_per_kwh: 0.0,
      costs_per_min: 0.1
    },
    {
      machine_id: 2,
      machine_name: 'MakerBot 3D Printer',
      machine_desc: 'NYC 3D printer 4 real and 4 life.',
      available: false,
      unavail_msg: 'We need some parts from the future to make it work again.',
      unavail_till: 4105119600,
      calc_by_energy: false,
      calc_by_time: true,
      costs_per_kwh: 0.0,
      costs_per_min: 0.2
    },
    {
      machine_id: 3,
      machine_name: 'Zing Laser Cutter',
      machine_desc: 'Cuts wood, plastic, paper. Fast.',
      available: true,
      unavail_msg: '',
      unavail_till: 0,
      calc_by_energy: false,
      calc_by_time: true,
      costs_per_kwh: 0.0,
      costs_per_min: 1.0
    },
    {
      machine_id: 4,
      machine_name: 'CNC Router',
      machine_desc: 'Cuts steel, plutanium, uranium. Drill on steroids.',
      available: true,
      unavail_msg: '',
      unavail_till: 0,
      calc_by_energy: false,
      calc_by_time: true,
      costs_per_kwh: 0.0,
      costs_per_min: 1.0
    },
    {
      machine_id: 5,
      machine_name: 'Hand Drill',
      machine_desc: 'A man is a man if he does not know how to handle one.',
      available: true,
      unavail_msg: '',
      unavail_till: 0,
      calc_by_energy: false,
      calc_by_time: true,
      costs_per_kwh: 0.0,
      costs_per_min: 0.2
    }
  ];

  // Put all tables in a single and fake database
  var initialData = {
    users: users,
    auth: auth,
    machines: machines
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
