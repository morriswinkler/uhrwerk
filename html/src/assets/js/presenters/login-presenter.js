"use strict";

function loginPresenter(element, options) {
  element = $(element);
  var model = options.model,
      alertTmpl = options.alertTmpl.html();

  /* Listen to user events */
  element.on('submit', 'form', function(e) {
    e.preventDefault();
    model.login({
      username: this.username.value, 
      password: this.password.value
    });
  });

  /* Listen to model events */
  model.on('load:login logout', load);
  model.on('login:success', onLoginSuccess);
  model.on('login:fail', onLoginFail);

  /* Event handlers */
  function load() {
    if (!element.hasClass('active')) {
      element.addClass('active');
    }
  }

  function onLoginSuccess(userID) {
    element.find('.alerts').empty();
    alert('User ID: ' + userID);
    element.removeClass('active');
  }

  function onLoginFail(response) {
    // Clear the alert placeholder from previous rubish. The response object
    // contains a message variable that is rendered into alertTmpl accordingly.
    element.find('.alerts').empty().append(
      riot.render(alertTmpl, response
    ));
  }
}
