$(function() {
  var $lightbox = $('<div class="lightbox hide"></div>');
  $("body").append($lightbox);

  $(".lightbox-available").on("click", function(event) {
    event.preventDefault();

    var origin = $(this).attr("href");
    $lightbox.html('<img src="' + origin + '">');
    $lightbox.show();
  });

  $lightbox.on("click", function() {
    $(this).hide();
  });
})
