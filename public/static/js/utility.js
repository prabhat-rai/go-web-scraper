var utility = {
    /**
     * Utility function to checik whether an attr is defined or not
     * @param attr
     * @param value
     * @returns {boolean}
     */
    attrDefined : function ( attr, value ) {
        return ( typeof attr !== typeof undefined && attr  !== false && attr === value );
    },

    /**
     * Shows AJAX loader
     */
    showAjaxLoader : function () {
        var loader = $("<div class='page-loader'></div>");
        $('body').prepend(loader, $('body:first-child'));
        $('.page-loader').show();
    },

    /**
     * Hides AJAX loader
     */
    hideAjaxLoader : function () {
        $('.page-loader').hide().remove();
    },

    /**
     * Show notifications depending on the classes provided
     * @param text
     * @param element_class
     * @param duration
     * @param parent_class
     */
    showNotification : function ( text, element_class, duration, parent_class ) {
        if( !element_class ) {
            element_class = 'text-info';
        }

        if( !parent_class ) {
            parent_class = 'custom';
        }

        if( !duration ) {
            duration = 10;
        }
        var element_html = '<strong class="' + element_class + '">' + text + '</strong>';
        $( '.notification-msg' )
            .finish()
            .attr( 'class', 'notification-msg collapse alert' )
            .addClass( parent_class )
            .html( element_html ).slideToggle()
            .delay( duration * 1000 )
            .slideToggle();
    }
};
