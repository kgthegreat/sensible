
$(function () {
    $('#myTab a:first').tab('show')
})
$( document ).ready(function() {
    console.log( "ready!" );
    var a = $('.tab-content .tweet-block')
    console.log(a)
    for (var i = 0; i < a.length; i++) {
        a[i].innerHTML = urlify(a[i].innerHTML)
    }
});

function urlify(text) {
    var urlRegex = /(https?:\/\/[^\s]+)/g;
    return text.replace(urlRegex, function(url) {
        return '<a href="' + url + '">' + url + '</a>';
    })
    // or alternatively
    // return text.replace(urlRegex, '<a href="$1">$1</a>')
}
