"use strict";

(function ($) {
	var app = new App();
	routes({app: app});

	// Binds the Todo Presenter
	loginPresenter($("#login-page"), {model: app});
	dashboardPresenter($("#dashboard-page"), {model: app, productTemplate: $("#product-tmpl")});
	messagePresenter($("#message-page"), {model: app});

})(jQuery);
