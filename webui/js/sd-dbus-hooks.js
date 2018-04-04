$(document).ready(function(){
  function setHeader(xhr) {
    xhr.setRequestHeader("X-Token", token);
  }

  function reloadTable() {
    $(".unit-item").remove();
    getUnits();
  }

  function getUnits() {
    $.ajax({
      url: "/unit/status/",
      type: "GET",
      dataType: "json",
      success: function(data) { 
        unitsToTable(data);
        $("#navbar-title").html("Systemd dbus hooks (" + data.length + ")");
       },
      beforeSend: setHeader
    });
  };

  if (typeof token == "undefined" || !token) {
    var token = prompt("enter token");
  }

  getUnits();

  function unitsToTable(data) {
    $.each(data, function(i, item) {
      switch(item.ActiveState) {
        case "active":
          badge = '<span class="badge badge-success">' + item.ActiveState + '</span>';
          actions = '<button type="button" class="btn btn-sm btn-danger" id="stop-btn" data-item="' + item.Name + '">stop</button>'
          break
        case "inactive":
          badge = '<span class="badge badge-secondary">' + item.ActiveState + '</span>';
          actions = '<button type="button" class="btn btn-sm btn-danger" id="start-btn" data-item="' + item.Name + '">start</button>'
          break
        case "failed":
          badge = '<span class="badge badge-danger">' + item.ActiveState + '</span>';
          actions = '<button type="button" class="btn btn-sm btn-danger" id="s-btn" data-item="' + item.Name + '">stop</button>'
          break
        default:
          badge = '<span class="badge badge-warning">' + item.ActiveState + '</span>';
          actions = '<button type="button" class="btn btn-sm btn-primary" id="start-btn" data-item="' + item.Name + '">start</button>' +
                    '<button type="button" class="btn btn-sm btn-danger" id="stop-btn" data-item="' + item.Name + '">stop</button>'
          break
      };

      $("#units-table").append(
        '<tr class="unit-item">' +
        '<td>' + item.Name + '</td>' +
        '<td>' + badge + '</td>' +
        '<td>' +
        '  <div class="btn-group" role="group" aria-label="unit-actions">' +
        actions +
        '  </div>' +
        '  <button type="button" class="btn btn-sm btn-info" id="journal-btn" data-item="' + item.Name + '">journal</button>' +
        '</td>' +
        '</tr>')
    });
  };

  $(document).on("click", "#start-btn", function() {
    var self = this,
    value = $(self).data("item");
    $.ajax({
      url: "/unit/start/" + value,
      dataType: "text",
      type: "GET",
      success: function(data) {
        console.log(value + " started");
        reloadTable();
      },
      error: function(data) {
        alert("server reply: " + data.status + "/" + data.statusText);
      },
      beforeSend: setHeader
    });
  });

  $(document).on("click", "#stop-btn", function() {
    var self = this,
    value = $(self).data("item");
    $.ajax({
      url: "/unit/stop/" + value,
      dataType: "text",
      type: "GET",
      success: function(data) {
        console.log(value + " stopped");
        reloadTable();
      },
      error: function(data) {
        alert("server reply: " + data.status + "/" + data.statusText);
      },
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
        $("#journal-title").html("journal for " + value);
        $("#journal-data").html("<pre>" + data + "</pre>");
        $("#journal-modal").modal("handleUpdate");
        $("#journal-modal").modal("show");
      },
      beforeSend: setHeader
    });
  });

  $(document).on("click", "#reload-btn", function() {
    reloadTable();
  });
});
