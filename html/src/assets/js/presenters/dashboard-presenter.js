"use strict";

function dashboardPresenter(element, options) {
  element = $(element);
  var model = options.model,
      tmpl = options.tmpl.html();

  /* Listen to user events */
  element.on('click', '.btn-activate', function(e) {
    e.preventDefault();
    element.removeClass('active');
    model.load('message');
  });

  /* Listen to model events */
  model.on('load:dashboard', load);
  model.on('load:products', listProducts);
  model.on('logout', hide);

  function load() {
    // Activate page
    if (!element.hasClass('active')) {
      element.addClass('active');
    }

    // Fill with products
    model.loadProducts();
  }

  function listProducts(products) {
    element.find('.products').empty();
    $.each(products, function(index, item) {
      var data = {
        name: item.name, 
        status: (item.status.charAt(0).toUpperCase() + item.status.slice(1)),
        statusClass: item.status.toLowerCase(), 
        price: item.price,
        disabled: item.status.toLowerCase() !== 'available' ? ' disabled' : ''
      };
      element.find('.products').append(riot.render(tmpl, data));
    });
  }

  function hide() {
    element.removeClass('active');
  }
}
