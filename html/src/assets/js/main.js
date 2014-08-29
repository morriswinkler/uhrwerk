var appViews = [
	$("#signin"),
	$("#dashboard"), 
	$("#message"),
	$("#deactivate")
];

var currentView = 0;

var hideAllViews = function() {
	for (var i=0; i<appViews.length; i++) {
		appViews[i].css("display", "none");
	}
};

var showView = function(view) {
	view.css("display", "block");
};

var nextView = function() {
	currentView++;
	if (currentView >= appViews.length) {
		currentView = 0;
	}
};

var centerVertically = function(arr) {
	for (var i=0; i<arr.length; i++) {
		var el = arr[i];
		var elHeight = el.height();
		var winHeight = $(window).height();
		var wspaceHeight = winHeight - elHeight;
		el.css("margin-top", wspaceHeight/2);
	}
};

$(document).ready( function() {
	centerVertically([
		$(".form-signin"), 
		$(".form-signout"), 
		$(".form-deactivate")
	]);
	hideAllViews();
	showView(appViews[currentView]);

	$(".advance").click( function(e){
		e.preventDefault();
		nextView();
		hideAllViews();
		showView(appViews[currentView]);
		centerVertically([
			$(".form-signin"), 
			$(".form-signout"), 
			$(".form-deactivate")
		]);
	});
});

$(window).resize( function() {
	centerVertically([
		$(".form-signin"), 
		$(".form-signout"), 
		$(".form-deactivate")
	]);
});