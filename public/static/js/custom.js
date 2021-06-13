// definition of insertAt function to insert string at a given position
String.prototype.insertAt = function( index, string ) {
    return this.substr( 0, index ) + string + this.substr( index );
};

// Definition of ucFirst function to upper case the first letter of a string
String.prototype.ucFirst = function() {
    return this.charAt( 0 ).toUpperCase() + this.slice( 1 );
};

// String prototype for removing HTML string
String.prototype.stripHtml = function(){
    var rex = /(<([^>]+)>)|(&lt;([^>]+)&gt;)|(&amp;)/ig;
    return this.replace(rex , "");

};

// Close notification on clicking the notification
$( '.close-notif,.alert' ).click( function () {
    $( this ).hide();
} );

// Bind AJAX functions before AJAX call start
$( document ).ajaxStart( function(){
    utility.showAjaxLoader();
});

// Bind AJAX functions after AJAX call ends
$( document ).ajaxComplete( function(){
    utility.hideAjaxLoader();
});

// Generic Error handler
$( document ).ajaxError( function( event, jqXhr ){
    var errorMessage = 'The request could not be completed';

    if ( jqXhr.status === 401 ) {
        location.reload();
    } else {
        if ( jqXhr.statusText && jqXhr.statusText.length > 0 ) {
            errorMessage = jqXhr.statusText;
        }

        utility.showNotification( errorMessage, 'text-danger', 5, 'alert-warning' );
    }
});
