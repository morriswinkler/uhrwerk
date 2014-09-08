'use strict';

function routes(models) {
	riot.route(function(hash) {
		models.app.trigger('load:message');
	});
}
