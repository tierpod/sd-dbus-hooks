$(document).ready(function(){
  function setHeader(xhr) {
    xhr.setRequestHeader("X-Token", "123");
  }

  $.ajax({
    url: "/units",
    type: "GET",
    dataType: "json",
    success: function(data) { unitsToTable(data); },
    beforeSend: setHeader
  });

  function unitsToTable(data) {
    $.each(data, function(i, item) {
      $("#units-table").append(
        '<tr>' +
        '<td id="units_' + item + '">' + item + '</td>' +
        '<td>' +
        '  <div class="btn-group" role="group" aria-label="unit-actions">' +
        '  <button type="button" class="btn btn btn-primary" id="start-btn" data-item="' + item + '">start</button>' +
        '  <button type="button" class="btn btn-danger" id="stop-btn" data-item="' + item + '">stop</button>' +
        '  </div>' +
        '  <button type="button" class="btn btn-info" id="journal-btn" data-item="' + item + '">journal</button>' +
        '</td>' +
        '</tr>')
    });
  };

  $(document).on("click", "#start-btn", function() {
    var self = this,
    value = $(self).data("item");
    $.ajax({
      url: "/unit/start/" + value,
      dataType: "json",
      type: "GET",
      success: function(data) { console.log(value + " started"); },
      beforeSend: setHeader
    });
  });

  $(document).on("click", "#stop-btn", function() {
    var self = this,
    value = $(self).data("item");
    $.ajax({
      url: "/unit/stop/" + value,
      dataType: "json",
      type: "GET",
      success: function(data) { console.log(value + " stopped"); },
      beforeSend: setHeader
    });
  });

  $(document).on("click", "#journal-btn", function() {
    var self = this,
    value = $(self).data("item");
    $.ajax({
      url: "/unit/journal/" + value,
      dataType: "text",
      type: "GET",
      success: function(data) {
        $("#journal-title").empty();
        $("#journal-title").append("journal for " + value);
        $("#journal-data").empty();
        $("#journal-data").append("<pre>" + data + "</pre>");
        $("#journal-modal").modal("handleUpdate");
        $("#journal-modal").modal("show");
      },
      beforeSend: setHeader
    });
  });
});