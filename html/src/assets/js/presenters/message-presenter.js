"use strict";

function messagePresenter(element, options) {
  element = $(element);
  var model = options.model;

  /* Listen to user events */
  element.on('submit', 'form', function(e) {
    e.preventDefault();
    element.removeClass('active');
    model.logout();
  });

  /* Listen to model events */
  model.on("load:message", load);

  function load() {
    if (!element.hasClass('active')) {
      element.addClass('active');
    }
  }
}
