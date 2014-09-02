"use strict";

function loginPresenter(element, options) {
  element = $(element);
  var model = options.model;

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

  function onLoginSuccess() {
    element.removeClass('active');
    model.load('dashboard');
  }

  function onLoginFail(message) {
    alert(message);
  }
}
