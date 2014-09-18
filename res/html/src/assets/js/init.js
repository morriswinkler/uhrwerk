/*
So this is an attempt to use riot.js for everything, init.js does the 
initialization of the main application, the App model so to say. This script
has to be placed at the end of the load queue as it has to be executed after
all of the parts of the application has been loaded.
*/

'use strict';

(function ($) {

	// Create main application model
	var app = new App();

	// Routes??? Not sure how to use them yet
	// Not using before not sure how routing works
	//routes({app: app});

	// Bind Active Booking presenter to the App model
	activeBookingsPresenter( $('#active-bookings-page'), {
		model: app,
		tmpl: $('#active-booking-tmpl')
	});

	// Bind Login presenter to the App model
	loginPresenter( $("#login-page"), { 
		model: app,
		alertTmpl: $('#login-alert-tmpl')
	});
	
	// Bind Dashboard presenter to the App model
	dashboardPresenter($('#dashboard-page'), {
		model: app, 
		tmpl: $('#product-tmpl')
	});
	
	// Bind Message presenter to the App model
	messagePresenter($('#message-page'), {
		model: app
	});

	// Bind status bar presenter to the App model
	statusBarPresenter($('#status-bar'), {
		model: app,
		tmpl: $('#status-bar-tmpl')
	});

	// Start the application with a specific page: login, dashboard, message
	app.start('login');

})(jQuery);
