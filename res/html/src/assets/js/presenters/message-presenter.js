'use strict';

function messagePresenter(element, options) {
  element = $(element);
  var model = options.model,
      timeout = null;

  /* Listen to user events */
  element.on('submit', 'form', function(e) {
    e.preventDefault();
    model.logout();
  });

  /* Listen to model events */
  model.on('load:message', load);
  model.on('logout', hide);

  /* Event handlers */
  function load() {
    if (!element.hasClass('active')) {
      element.addClass('active');
    }

    // Activate countDown
    countdown(10);
  }

  function hide() {

    // Deactivate timeout.
    if (timeout != null) {
      window.clearTimeout(timeout);
      timeout = null;
    }

    element.removeClass('active');
  }

  /* Private methods */
  function countdown(timeLeft) {
    if (timeLeft <= 0) {
      model.logout();
    } else {
      element.find('.countdown').html(timeLeft);
      timeout = setTimeout(function() {
        countdown(--timeLeft);
      }, 1000);
    }
  }
}
