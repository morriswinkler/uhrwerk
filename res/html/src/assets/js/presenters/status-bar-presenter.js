'use strict';

function statusBarPresenter(element, options) {
  element = $(element);
  var model = options.model,
      tmpl = options.tmpl.html();

  /* Listen to user events */
  element.on('click', '.btn-log-out', function(e) {
    e.preventDefault();
    model.logout();
  });

  /* Listen to model events */
  model.on('login:success load:dashboard load:message', show);
  model.on('logout', hide);

  /* Event handlers */
  function show() {
    element.empty();
    var data = {title: 'Dashboard'};
    element.append(riot.render(tmpl, data));
  }

  function hide() {
    element.empty();
  }
}