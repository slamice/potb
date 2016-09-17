function programTable(table_name, save_name, add_url, table_add, table_remove){
    var $TABLE = $('#' + table_name);
    var $save_button = $('#' + save_name);
    var $add_url = '/' + add_url;
    var $table_add = $('.' + table_add)
    var $table_remove = $('.' + table_remove)

    $table_add.click(function() {
      var $clone = $TABLE.find('tr.hide').clone(true).removeClass('hide table-line');
      $TABLE.find('table').append($clone);
    });

    $table_remove.click(function() {
      $(this).parents('tr').detach();
    });

    // A few jQuery helpers for exporting only
    jQuery.fn.pop = [].pop;
    jQuery.fn.shift = [].shift;

    $save_button.click(function() {
      var $rows = $TABLE.find('tr:not(:hidden)');
      var headers = [];
      var data = [];

      var $programDate = document.getElementById('program-date-field').value;

      if ($programDate != null) {
          $.ajax({
              url : "/addprogramdate",
              type: "POST",
              data: JSON.stringify({"ProgramDate": $programDate}),
              contentType: "application/json; charset=utf-8",
              dataType   : "json",
              success    : function(){
                  console.log("Pure jQuery Pure JS object");
              }
          });
          var $programDate = "";
      }


      // Get the headers (add special header logic here)
      $($rows.shift()).find('th:not(:empty)').each(function() {
        headers.push($(this).text().toLowerCase());
      });

      // Turn all existing rows into a loopable array
      $rows.each(function() {
        var $td = $(this).find('td');
        var h = {};

        // Use the headers from earlier to name our hash keys
        headers.forEach(function(header, i) {
          h[header] = $td.eq(i).text();
        });

        data.push(h);
      });

      $.ajax({
          url : $add_url,
          type: "POST",
          data: JSON.stringify(data),
          contentType: "application/json; charset=utf-8",
          dataType   : "json",
          success    : function(){
              console.log("Pure jQuery Pure JS object");
          }
      });

    });

}
