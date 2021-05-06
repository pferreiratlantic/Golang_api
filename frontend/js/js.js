var numOfOffsets = 0
var lineLimit = 50
var currentOffset = 0
var currentCountrySet = ""
$(document).ready(function(){
    $.ajax({
        url:"/statistics",
        type: "GET",
        success: function(userPerContryResults){
            $('#numOfUserPerCountry').empty();
            for (var i = 0; i < userPerContryResults.length; i++) {
                var line = '<tr><td>'+
                userPerContryResults[i].count+
                '</td><td>'+
                userPerContryResults[i].countryName+
                '</tr>'
                $('#numOfUserPerCountry').append(line);
            }
        }
    });
    $.ajax({
        url:"/countries",
        type: "GET",
        success: function(contriesResults){
            for (var i = 0; i < contriesResults.length; i++) {
                var line = '<option value="'+
                contriesResults[i].countryName+
                '">'+
                contriesResults[i].countryName+
                '</option>'
                $('#countrySelector').append(line);
            }
        }
    });
 });

$(document).ready(function(){
        $('#button2').on('click', function(e){
            currentOffset = 0
            currentCountrySet=$( "#countrySelector option:selected" ).text();
            e.preventDefault();
            $.ajax({
                url:"/statistics",
                type: "GET",
                data: {"country":currentCountrySet},
                success: function(numOfRows){
                    for(var k = 0 ; k < numOfRows.length; k++){
                        if ( numOfRows[k].countryName == currentCountrySet){
                            numOfOffsets = Math.trunc(parseInt(numOfRows[k].count) / lineLimit) +1;
                        }
                    }
                    var lineOffset = "On Page [ "+ (currentOffset+1)+" ] of "+numOfOffsets+" pages.";
                    $('#offsetDimension').empty();
                    $('#offsetDimension').append(lineOffset);
                }
            });

        $.ajax({
            url:"/users",
            type: "GET",
            data: {"country":currentCountrySet, "lines":lineLimit, "offset":currentOffset},
            success: function(results){
                $('#listUsersByCountry').empty();
                for (var i = 0; i < results.length; i++) {
                    var line = '<tr><td>'+
                        results[i].userId+
                        '</td><td>'+
                        results[i].userEmail+
                        '</td><td>'+
                        results[i].userPhone+
                        '</td><td>'+
                        results[i].userParcelWeight+
                        '</td></tr>'
                    $('#listUsersByCountry').append(line);
                }

            }
        });
    });
 });

$(document).ready(function(){
        $('#buttonFirst').on('click', function(e){
            currentOffset=0;
        $.ajax({
            url:"/users",
            type: "GET",
            data: {"country":currentCountrySet, "lines":lineLimit, "offset":currentOffset},
            success: function(results){
                $('#listUsersByCountry').empty();
                for (var i = 0; i < results.length; i++) {
                    var line = '<tr><td>'+
                        results[i].userId+
                        '</td><td>'+
                        results[i].userEmail+
                        '</td><td>'+
                        results[i].userPhone+
                        '</td><td>'+
                        results[i].userParcelWeight+
                        '</td></tr>'
                    $('#listUsersByCountry').append(line);
                }
                var lineOffset = "On Page [ "+ (currentOffset+1) +" ] of "+numOfOffsets+" pages.";
                $('#offsetDimension').empty();
                $('#offsetDimension').append(lineOffset);
            }
        });
    });
 });

$(document).ready(function(){
        $('#buttonNext').on('click', function(e){
            currentOffset=currentOffset+lineLimit;
            if(currentOffset > numOfOffsets){
                currentOffset = numOfOffsets;
            }
        $.ajax({
            url:"/users",
            type: "GET",
            data: {"country":currentCountrySet, "lines":lineLimit, "offset":currentOffset},
            success: function(results){
                $('#listUsersByCountry').empty();
                for (var i = 0; i < results.length; i++) {
                    var line = '<tr><td>'+
                        results[i].userId+
                        '</td><td>'+
                        results[i].userEmail+
                        '</td><td>'+
                        results[i].userPhone+
                        '</td><td>'+
                        results[i].userParcelWeight+
                        '</td></tr>'
                    $('#listUsersByCountry').append(line);
                }
                var lineOffset = "On Page [ "+(Math.trunc(currentOffset / lineLimit)+1) +" ] of "+numOfOffsets+" pages.";
                $('#offsetDimension').empty();
                $('#offsetDimension').append(lineOffset);
            }
        });
    });
 });

$(document).ready(function(){
        $('#buttonPrevious').on('click', function(e){
            currentOffset=currentOffset-lineLimit;
            if(currentOffset < 0){
                currentOffset = 0;
            }
        $.ajax({
            url:"/users",
            type: "GET",
            data: {"country":currentCountrySet, "lines":lineLimit, "offset":currentOffset},
            success: function(results){
                $('#listUsersByCountry').empty();
                for (var i = 0; i < results.length; i++) {
                    var line = '<tr><td>'+
                        results[i].userId+
                        '</td><td>'+
                        results[i].userEmail+
                        '</td><td>'+
                        results[i].userPhone+
                        '</td><td>'+
                        results[i].userParcelWeight+
                        '</td></tr>'
                    $('#listUsersByCountry').append(line);
                }
                var lineOffset = "On Page [ "+(Math.trunc(currentOffset / lineLimit)+1) +" ] of "+numOfOffsets+" pages.";
                $('#offsetDimension').empty();
                $('#offsetDimension').append(lineOffset);
            }
        });
    });
 });

$(document).ready(function(){
        $('#button1').on('click', function(e){
            var path=$("#filePath").val();
            e.preventDefault();
            console.log(path);
            alert("Starting loading process, this may take a while..");
            $.ajax({
                url:"/loader",
                type: "GET",
                data: {"path":path},
                success: function(results){
                    console.log(results);
                    
                    if (results == null || results < 1){
                        alert("No regs processed. Check url.");
                    }
                    else{
                        alert(results+" regs were successfully processed.");
                        $.ajax({
                            url:"/statistics",
                            type: "GET",
                            success: function(userPerContryResults){
                                $('#numOfUserPerCountry').empty();
                                for (var i = 0; i < userPerContryResults.length; i++) {
                                    var line = '<tr><td>'+
                                    userPerContryResults[i].count+
                                    '</td><td>'+
                                    userPerContryResults[i].countryName+
                                    '</tr>'
                                    $('#numOfUserPerCountry').append(line);
                                }
                            }
                        });
                    }
                }
            });
    });
 });