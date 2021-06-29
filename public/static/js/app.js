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
        var dateColumns = tableObj.attr( 'data-date-columns' )
            ? tableObj.attr('data-date-columns') : false;
        var columnListToCollapse = collapseColumns !== false ? collapseColumns.split(',') : [];
        var columnListForDate = dateColumns !== false ? dateColumns.split(',') : [];

        var statusColumns = tableObj.attr( 'data-status-columns' )
            ? tableObj.attr('data-status-columns') : false;
        var columnListForStatus = statusColumns !== false ? statusColumns.split(',') : [];

        var actionColumns = tableObj.attr( 'data-action-columns' )
            ? tableObj.attr('data-action-columns') : false;
        var columnListForAction = actionColumns !== false ? actionColumns.split(',') : [];
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
                    },
                    {
                        targets: columnListForStatus.map(Number),
                        render: function(cellData, type, row) {
                            let entity = $('.data-table-list').first().attr('data-entity');
                            let active = false;
                            let attributeId = row.id;

                            if(row.active) {
                                active = true;
                            }

                            if (cellData === true) {
                                return '<label><input id="active" name="active" type="checkbox" class="fas fa-check btn btn-success btn-circle" ' +
                                    'onchange="webScrapperApp.changeActiveStatus(\'' + entity + '\',\'' + attributeId + '\',' + active + ')" ' +
                                    'checked style="visibility:hidden;">' +
                                        '<span id="active-toggle" class="btn btn-success btn-circle"> ' +
                                        '<i id="active-toggle-icon" class="fas fa-check"></i>' +
                                    '</span></input></label>';
                            } else {
                                return '<label><input id="inactive" name="active" type="checkbox"  class="fas fa-check btn btn-danger btn-circle" ' +
                                    'onchange="webScrapperApp.changeActiveStatus(\'' + entity + '\',\'' + attributeId + '\',' + active + ')" ' +
                                    'style="visibility:hidden;" checked>' +
                                        '<span id="inactive-toggle" class="btn btn-danger btn-circle"> ' +
                                            '<i id="inactive-toggle-icon" class="fas fa-times"></i>' +
                                    '</span></input></label>';
                            }
                        }
                    },
                    {
                        targets: columnListForDate.map(Number),
                        render: function(cellData) {
                            var d = new Date(0);
                            d.setUTCSeconds(cellData);
                            return d.toLocaleDateString('en-GB', {
                                day : 'numeric',
                                month : 'numeric',
                                year : 'numeric',
                            });
                        }
                    },
                    {
                        
                        targets: columnListForAction.map(Number),
                            render: function(cellData, type, row) {
                                let attributeId = row.id;
                                let name = row.name;
                                let email = row.email;
                                let phone = row.phone;
                                let role = row.role;
                                return '<button onclick="webScrapperApp.openModel(\'' + attributeId + '\',\'' + name + '\',\'' + email + '\',\'' + phone + '\',\'' + role+'\'' + ')" data-toggle="modal" data-target="#userModal" class="btn btn-primary" type="button" id="edit" aria-haspopup="true" aria-expanded="false"> Edit </button>';
                            }
                    },
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
    },

    changeKeyGroupSubscription : function (keyGroupId, subscriptionStatus) {
        let subscriptionText = (subscriptionStatus === 1 ? "Subscribe" : "Unsubscribe");
        bootbox.confirm("Are you sure that you want to " + subscriptionText + "?", function (result) {
            if(result) {
                $.ajax({
                    url: "/ajax/keyword-groups/change-subscription",
                    dataType: 'json',
                    method : 'POST',
                    data : 'subscription=' + subscriptionStatus + "&id=" + keyGroupId,
                    success: function( response ) {
                        if ( response !== 0 ) {
                            $(webScrapperApp.tableClassName).DataTable().ajax.reload();
                            utility.showNotification(subscriptionText + 'd', 'text-success', 5, 'alert-info');
                        } else {
                            utility.showNotification( 'Something went wrong.', 'text-danger', 5, 'alert-warning' );
                        }
                    }
                });
            }
        });
    },

    fetchReviews : function (concept, fetchForPlatform) {
        $.ajax({
            url: '/ajax/reviews/fetch?concept=' + concept + '&platform=' + fetchForPlatform,
            dataType: 'json',
            method : 'GET',
            success: function( response ) {
                if ( response !== 0 ) {
                    let successMsg = 'We are fetching latest reviews of ' + concept + ' from ' + fetchForPlatform + ' platform';
                    utility.showNotification(successMsg, 'text-success', 5, 'alert-info');
                } else {
                    utility.showNotification( 'Something went wrong.', 'text-danger', 5, 'alert-warning' );
                }
            }
        });
    },

    loadAnalyticsData : function () {
        let noOfDays = $('#noOfDays').val();
        let dataForPlatformGraph = [];
        let dataForConceptGraph = [];
        let allConcepts = ["babyshop", "centrepoint", "homebox", "homecentre", "lifestyle", "max", "shoemart", "shukran", "splash"];
        let allConceptLabels = ["Babyshop", "Centrepoint", "Homebox", "Homecentre", "Lifestyle", "Max", "Shoemart", "Shukran", "Splash"];
        let conceptLineColors = [ '#2577B5', '#DD7CB5', '#FF1111', '#444444', '#111111' , '#222222', '#333333', '#666666', '#555555' ];
        let allPlatforms = ['android', 'ios'];
        let allPlatformLabels = ['Android', 'iOS'];
        let platformLineColors = [ '#2577B5', '#DD7CB5' ];
        $.ajax({
            url: "/ajax/analytics/counts?days=" + noOfDays,
            dataType: 'json',
            method : 'GET',
            success: function( response ) {
                let dataObj = {"day" : ""};

                $.each(allPlatforms, function (platformKey, platform) {
                    dataObj[platform] = 0;
                });

                $.each(response.platform, function (date, values) {
                    let currentDataObj = { ...dataObj };
                    currentDataObj["day"] = date;

                    $.each(values, function (platform, count) {
                        currentDataObj[platform] = count;
                    });

                    dataForPlatformGraph.push(currentDataObj);
                });

                webScrapperApp.loadChart('platformAreaChart', dataForPlatformGraph, allPlatforms,
                    allPlatformLabels, platformLineColors);

                let conceptDataObj = {"day" : ""};
                $.each(allConcepts, function (conceptKey, concept) {
                    conceptDataObj[concept] = 0;
                });

                $.each(response.concept, function (date, values) {
                    let currentDataObj = { ...conceptDataObj };
                    currentDataObj["day"] = date;

                    $.each(values, function (concept, count) {
                        currentDataObj[concept] = count;
                    });

                    dataForConceptGraph.push(currentDataObj);
                });

                webScrapperApp.loadChart('conceptAreaChart', dataForConceptGraph, allConcepts,
                    allConceptLabels, conceptLineColors);
            }
        });
    },

    loadChart : function (elementId, dataForGraph, yKeys, yLabels, lineColors) {
        $("#" + elementId).empty();
        // Line Chart statistics
        Morris.Line({
            element: elementId,
            data: dataForGraph,
            xkey: 'day',
            ykeys: yKeys,
            labels: yLabels,
            pointSize: 2,
            hideHover: 'auto',
            lineColors: lineColors,
        });
    },

    setupReviewPage : function () {
        let url = new URL(location.href);
        let searchParams = new URLSearchParams(url.search);
        let platform = searchParams.get('platform');

        if(platform) {
            $('#platformSelect').val(platform);
            webScrapperApp.applyReviewFilters();
        }
    },

    changeActiveStatus : function (type, name,status) {
        $.ajax({
            url: "/ajax/" + type + "/status?id="+name+"&active="+!status,
            dataType: 'json',
            method : 'POST',
            success: function( response ) {
                utility.showNotification("Updated active status!", 'text-success', 5, 'alert-info');
                if( $( webScrapperApp.tableClassName ).length > 0 ) {
                    $(webScrapperApp.tableClassName).DataTable().ajax.reload();
                } else {
                    window.location.reload();
                }
            }
        });

    },

    changeStatus : function () {

        if(document.getElementById('active').checked) {
            document.getElementById("toggle").className = "btn btn-success btn-circle";
            document.getElementById('active').value = true;
            document.getElementById("toggle-icon").className = "fas fa-check"

        } else {
            document.getElementById("toggle").className = "btn btn-danger btn-circle";
            document.getElementById('active').value = false;
            document.getElementById("toggle-icon").className = "fas fa-times"
        }
    },

    getKeywords: function () {
        $("#keyword-list").val($("#keywordGroupSelect").val())
    },
    openModel: function (id,name,email,phone,role,active){
        document.getElementById("userId").value = id;
        document.getElementById("userName").value = name;
        document.getElementById("userEmail").value = email;
        document.getElementById("userPhone").value = phone;
        if (role === "admin"){
            document.getElementById("admin").checked = true;
        }
        else{
            document.getElementById("user").checked = true;
        }
    }
};

$( document ).ready( function () {
    if( $( webScrapperApp.tableClassName ).length > 0 ) {
        webScrapperApp.initializeDatatable();
    }

    if( $('.select2-dropdown').length > 0 ) {
        $('.select2-dropdown').each(function () {
            let idKey = $(this).attr('data-id-key') || 'id';
            let textKey = $(this).attr('data-text-key') || 'text';

            $(this).select2({
                ajax: {
                    url: $(this).attr('data-url'),
                    dataType: 'json',
                    processResults: function (response) {
                        let result = [{"id" : "", "text" : "-- No Selection --"}];

                        $.each(response.data, function (key, val) {
                            result.push({"id" : val[idKey], "text" : val[textKey]});
                        });

                        // Transforms the top-level key of the response object from 'items' to 'results'
                        return {
                            results: result
                        };
                    },
                    data: function (params) {
                        var query = {
                            "search[value]" : params.term,
                            "active": 1,
                            "length" : 5,
                        };

                        // Query parameters will be ?search=[term]&type=public
                        return query;
                    },

                }
            });
        })
    }

    let specificPage = $('#specific_page_method');
    if( specificPage.length > 0 ) {
        let page = specificPage.attr('data-page');

        switch (page) {
            case "analytics":
                webScrapperApp.loadAnalyticsData();
                break;
            case "reviews":
                webScrapperApp.setupReviewPage();
                break;
            default:
                console.log("Unknown Page");
        }
    }
});