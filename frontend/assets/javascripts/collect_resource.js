var ready;
ready = function() {

  $('#collect').magnificPopup({
    type: 'inline',
    preloader: false,
    focus: '#username',
    modal: true
  });

  $(document).on('click', '.popup-modal-dismiss', function (e) {
    e.preventDefault();
    $.magnificPopup.close();
  });


  $(document).on('click', '#collect', function () {
    $('#collect-button').hide();
    $('#collected-image').hide();

    $('#clock').countdown(Date.now()+10000)
    .on('update.countdown', function(event) {
      var format = 'Wait for %S seconds...';

      $(this).html(event.strftime(format));
    })
    .on('finish.countdown', function(event) {
      $(this).html('Now collect it!');
      $('#collect-status').text('Done!');
      $('#cancel-button').hide();
      $('#collect-button').show();
      $('#collect-image').hide();
      $('#collected-image').show();
    });
  });

};

$(document).ready(ready);
$(document).on('page:load', ready);