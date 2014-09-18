'use strict'

function activeBookingsPresenter(element, options) {
  element = $(element);
  var model = options.model,
      tmpl  = options.tmpl.html();

  /* Listen to user events */
  element.on('click', '.btn-deactivate', function(e) {
    e.preventDefault();
    var bookID = $(this).parent().find('.book-id').attr('value');
    model.deactivateMachine(bookID);
  });

  /* Listen to model events */
  model.on('load:active-bookings', load);
  model.on('deactivateMachine:success', onDeacSuccess);
  model.on('deactivateMachine:fail', onDeacFail);
  model.on('logout', onLogout);

  /* Event listeners */

  // Loads active bookings page
  function load(bookings) {
    $.each(bookings, function(index, item) {
      var data = {
        machine_id: item.machine_id,
        machine_name: item.machine_name, 
        book_id: item.book_id
      };
      element.find('.bookings').append(riot.render(tmpl, data));
    });
    show();
  }

  // Handles the deactivateMachine:success event
  function onDeacSuccess(bookID) {
    // For now just hide, but later we should check
    // if if there are bookings left and remove just the
    // one that has been deactivated from the visible list.
    hide();
    // If there are no other bookings left (and in this case
    // there are none) load dashboard page
    model.load('dashboard');
  }

  // Handles the deactivateMachine:fail event
  function onDeacFail(message) {
    alert(message);
  }

  // Handle logout event
  function onLogout() {
    clear();
    hide();
  }

  // Shows active bookings page
  function show() {
    if (!element.hasClass('active')) {
      element.addClass('active');
    }
  }

  // Hides active bookings page
  function hide() {
    element.removeClass('active');
  }

  // clear bookings container
  function clear() {
    element.find('.bookings').empty();
  }
}