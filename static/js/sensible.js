
$(function () {
  $('#myTab a:first').tab('show')

})
$( document ).ready(function() {
  console.log( "ready!" );
    $('[data-toggle="tooltip"]').tooltip()
    var a = $('.tab-content .tweet-block')
    console.log(a)
    for (var i = 0; i < a.length; i++) {
        a[i].innerHTML = urlify(a[i].innerHTML)
    }

  var selection;
  $('.tweet-content').mouseup(function() {
    selection = getSelected();

    if (selection && (selection = new String(selection).replace(/^\s+|\s+$/g,''))) {
      openModal(selection);
    }
  });

  $(".categorise").click(function(e){
    e.preventDefault()
    console.log("Clicked on category")
    var buttonEl = $(this)
    var category = buttonEl.attr('type')
    console.log(selection)
    console.log(category)
    $.ajax({
      url: '/categorise',
      data: JSON.stringify({
        "phrase": selection,
        "category": category
      }),
      type: 'POST',
      success: function(res) {
        console.log("success1")
        $('#myModal').modal('hide');
      },
      error: function(error) {
        $('#myModal').modal('hide');
        console.log(error);
      }
    })
  });

  callTwitterInteraction(".tweet-retweet", "data-tweet-id", "/retweet")
  callTwitterInteraction(".tweet-fav", "data-tweet-id", "/fav")

  
  $('.nav-tabs').scrollingTabs({
    bootstrapVersion: 4,
    cssClassLeftArrow: 'arrow left i',
    cssClassRightArrow: 'arrow right i',
    enableSwiping: true,
    scrollToTabEdge: true,
    disableScrollArrowsOnFullyScrolled: true
    //    leftArrowContent: ['']

    
  });
});

function callTwitterInteraction(clickedEl, dataEl, remoteApi) {
  $(clickedEl).click(function(e){
    e.preventDefault()
    console.log("Clicked on " + clickedEl)
    var linkEl = $(this)
//    console.log(linkEl[0])
    console.log(linkEl)
    var tweetId = linkEl.attr('data-tweet-id')
    console.log(tweetId)
    $.ajax({
      url: remoteApi,
      data: JSON.stringify({
        "id_str": tweetId
      }),
      type: 'POST',
      success: function(res) {
        console.log("success " + remoteApi)
      },
      error: function(error) {
        console.log(error);
      }
    })
  });
}

function openModal(selection) {

  $('#myModal').on('show.bs.modal', function (event) {
    var modalEl = $(this)
    modalEl.find('.modal-title').text(selection)
    
  }).modal({
    keyboard: true
  })
}




function urlify(text) {
    var urlRegex = /(https?:\/\/[^\s]+)/g;
    return text.replace(urlRegex, function(url) {
        return '<a href="' + url + '">' + url + '</a>';
    })
    // or alternatively
    // return text.replace(urlRegex, '<a href="$1">$1</a>')
}

function getSelected() {
    if (window.getSelection) {
        return window.getSelection();
    }
    else if (document.getSelection) {
        return document.getSelection();
    }
    else {
        var selection = document.selection && document.selection.createRange();
        if (selection.text) {
            return selection.text;
        }
        return false;
    }
    return false;
}

