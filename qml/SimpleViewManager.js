// holds the currently displayed view
var currentView;

// holds the current view history
var viewHistory = new Array();

function showView(nextView) {
    if(nextView.visible) {
		var isPreviousView = false;

		for(var i = 0; i < viewHistory.length; i++) {
	    	if(viewHistory[i] == nextView) {
				isPreviousView = true;
				viewHistory.splice(i);
				break;
	    	}
		}

		if(nextView != currentView) {
		    /*currentView.anchors.fill = nextView.anchors.fill = null;
		    currentView.width = nextView.width = currentView.width;
		    currentView.height = nextView.height = currentView.height;*/
		    nextView.y = viewManager.height * (isPreviousView ? -1 : 1);
	
		    if(currentViewAnimation.running)
				currentViewAnimation.complete();
	
		    if(nextViewAnimation.running)
				nextViewAnimation.complete();
	
		    currentViewAnimation.target = currentView;
		    currentViewAnimation.to = viewManager.height * (isPreviousView ? 1 : -1);
		    currentViewAnimation.running = true;
	
		    nextViewAnimation.target = nextView;
		    nextViewAnimation.to = - 0;
		    nextViewAnimation.running = true;
	
		    if(currentView && !isPreviousView) {
				viewHistory.push(currentView);
		    }
		    currentView = nextView;
		}
    }
}

function connectViewEvents(view)
{
    view.visibleChanged.connect( function() {
		SVM.showView(view);
    });
}
