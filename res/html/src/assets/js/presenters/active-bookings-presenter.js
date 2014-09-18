'use strict'

function activeBookingsPresenter(element, options) {
  element = $(element);
  var model = options.model,
      tmpl  = options.tmpl.html(),
      timeout = null;

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
      // Calculate usage time for each booking
      var timeStart = item.time_start;
      var timeNow = item.time_now;
      var stampStart = getUnixTimestamp(timeStart);
      var stampNow = getUnixTimestamp(timeNow);
      var usageTimeSeconds = stampNow - stampStart;
      var usageTimeDisplay = secondsTimeSpanToHMS(usageTimeSeconds);

      var data = {
        machine_id: item.machine_id,
        machine_name: item.machine_name, 
        book_id: item.book_id,
        usage_time_seconds: usageTimeSeconds,
        usage_time_display: usageTimeDisplay
      };
      element.find('.bookings').append(riot.render(tmpl, data));
    });
    updateUsageTime();
    show();
  }

  function getUnixTimestamp(datetimeString) {
    var arr = datetimeString.split(" ");
    var d = arr[0].split("-");
    var t = arr[1].split(":");
    var date = new Date(d[0], d[1], d[2], t[0], t[1], t[2], 0);
    return Math.round(date.getTime() / 1000);
  }

  function secondsTimeSpanToHMS(s) {
    var h = Math.floor(s / 3600); // Get whole hours
    s -= h * 3600;
    var m = Math.floor(s / 60); // Get remaining minutes
    s -= m * 60;
    return (h > 0 ? h + ':' : '')
         + (m < 10 ? '0' + m : m) + ':'
         + (s < 10 ? '0' + s : s); //zero padding on minutes and seconds
  }

  function updateUsageTime() {
    //console.log('active book timeout running');
    element.find('.usage-time').each(function() {
      var currentSeconds = $(this).data('usage-seconds');
      currentSeconds++;
      $(this).data('usage-seconds', currentSeconds);
      $(this).html(secondsTimeSpanToHMS(currentSeconds));
    });
    timeout = setTimeout(updateUsageTime, 1000);
  }

  // Handles the deactivateMachine:success event
  function onDeacSuccess(bookID) {
    clearTimeout(timeout);
    timeout = null;
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
    clearTimeout(timeout);
    timeout = null;
    alert(message);
  }

  // Handle logout event
  function onLogout() {
    clearTimeout(timeout);
    timeout = null;
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