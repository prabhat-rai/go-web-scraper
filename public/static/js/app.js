var webScrapperApp= {
    tableId : '',
    tableIdArray : '',
    tableClassName : '.data-table-list',
    tableLimit : 10,
    sortCol : 0,
    sortColOrder : 'desc',
    dataTableLengthArray : [ 5, 10, 25, 50, 100 ],

    additionalDataTableFilters : {},

    /**
     * Check for one or more datatables in the current page and call loadDatatable for each instance
     */
    initializeDatatable : function () {
        var tableObj = $( webScrapperApp.tableClassName );

        if( tableObj.length > 1 )
        {
            $.each( tableObj, function () {
                var test = $( this );
                webScrapperApp.loadDatatable( test );
            } );
        }
        else
        {
            webScrapperApp.loadDatatable( tableObj );
        }

    },

    /**
     * Sets up the datatable by picking up the settings defined/ custom settings and load the data
     * @param tableObj
     */
    loadDatatable : function ( tableObj ) {
        var url = tableObj.attr( 'data-url' )
            ? tableObj.attr( 'data-url' ) : window.location.href;
        var tableLimit = tableObj.attr( 'data-limit' )
            ? tableObj.attr( 'data-limit' ) : webScrapperApp.tableLimit;
        var sortCol = tableObj.attr( 'data-sort-col' )
            ? tableObj.attr( 'data-sort-col' ) : webScrapperApp.sortCol;
        var sortColOrder = tableObj.attr( 'data-sort-order' )
            ? tableObj.attr( 'data-sort-order' ) : webScrapperApp.sortColOrder;
        var collapseColumns = tableObj.attr( 'data-collapse-data-columns' )
            ? tableObj.attr('data-collapse-data-columns') : false;
        var columnListToCollapse = collapseColumns !== false ? collapseColumns.split(',') : [];

        // Columns array to hold details of datatable columns
        var columns = [];

        tableObj. find( 'thead th' ).each( function () {
            var tHead = $(this);
            var visibleColAttr = tHead.attr( 'data-col-visible' );
            var sortableColAttr = tHead.attr( 'data-col-sortable' );
            var searchableColAttr = tHead.attr( 'data-col-searchable' );

            var colVisibility = '1';
            var colSortable = '1';
            var colSearchable = '1';

            if( utility.attrDefined( visibleColAttr , '0' ) )
            {
                colVisibility = '0';
            }

            if( utility.attrDefined( sortableColAttr , '0' ) )
            {
                colSortable = '0';
            }

            if( utility.attrDefined( searchableColAttr , '0' ) )
            {
                colSearchable = '0';
            }

            columns.push( {
                data: tHead.attr( 'data-dt-query' ),
                name : tHead.attr( 'data-dt-name' ),
                visible : ( colVisibility === '1' ),
                sortable : ( colSortable === '1' ),
                searchable : ( colSearchable === '1' )
            } );
        }).promise().done( function() {
            $( tableObj ).DataTable({
                processing: true,
                serverSide: true,
                ajax: {
                    "url" : url,
                    "data": function ( d ) {
                        $.each(webScrapperApp.additionalDataTableFilters, function (key, val) {
                            d[key] = val;
                        });
                    }
                },
                pagingType: "full_numbers",
                aoColumns: columns,
                searchDelay: 1000,
                lengthMenu: webScrapperApp.dataTableLengthArray,
                order: [ [ sortCol, ( sortColOrder ).toLowerCase() ] ],
                displayLength: tableLimit,
                columnDefs: [
                    {
                        targets: columnListToCollapse.map(Number),
                        createdCell: function(cell, cellData) {
                            if(collapseColumns === false || cellData.length < 100) {
                                return;
                            }

                            var $cell = $(cell);
                            $cell.contents().wrapAll("<div class='content'></div>");
                            var $content = $cell.find(".content");

                            $cell.append($("<button class='btn btn-info btn-icon-split btn-sm'>" +
                                    "<span class='text'>Read more</span>" +
                                    "<span class='icon text-white-50'>" +
                                        "<i class='fas fa-angle-down'></i>" +
                                    "</span>" +
                                "</button>"
                            ));

                            $btn = $cell.find("button");
                            $content.css({
                                "height": "50px",
                                "overflow": "hidden"
                            });
                            $cell.data("isLess", true);

                            $btn.click(function() {
                                var isLess = $cell.data("isLess");
                                $content.css("height", isLess ? "auto" : "50px");
                                $(this).find('.text')
                                    .text(isLess ? "Read less" : "Read more");
                                $(this).find('.icon')
                                    .html(isLess ? "<i class='fas fa-angle-up'></i>" : "<i class='fas fa-angle-down'></i>");
                                $cell.data("isLess", !isLess);
                            })
                        }
                    }
                ]
            });
        } );
    },

    applyReviewFilters: function () {
        let filters = {};
        $('.ratings-filter').each(function( index ) {
            if($( this ).val() !== '') {
                filters[$( this ).attr("data-filter-on")] = $( this ).val();
            }
        });

        this.additionalDataTableFilters = filters;

        $(this.tableClassName).DataTable().ajax.reload();
    }
};

$( document ).ready( function () {
    if( $( webScrapperApp.tableClassName ).length > 0 ) {
        webScrapperApp.initializeDatatable();
    }
});