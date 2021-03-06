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
          badgeClass = 'badge-success';
          actions = '<button type="button" class="btn btn-sm btn-danger" id="stop-btn" data-item="' + item.Name + '">stop</button>';
          break
        case "inactive":
          badgeClass = 'badge-secondary';
          actions = '<button type="button" class="btn btn-sm btn-primary" id="start-btn" data-item="' + item.Name + '">start</button>';
          break
        case "failed":
          badgeClass = 'badge-danger';
          actions = '<button type="button" class="btn btn-sm btn-danger" id="start-btn" data-item="' + item.Name + '">start</button>';
          break
        default:
          badgeClass = 'badge-warning';
          actions = '<button type="button" class="btn btn-sm btn-primary" id="start-btn" data-item="' + item.Name + '">start</button>' +
                    '<button type="button" class="btn btn-sm btn-danger" id="stop-btn" data-item="' + item.Name + '">stop</button>';
          break
      };

      $("#units-table").append(
        '<tr class="unit-item">' +
        '<td data-toggle="tooltip" data-placement="bottom" title="' + item.Description + '">' + item.Name + '</td>' +
        '<td><span class="badge ' + badgeClass + '">' + item.ActiveState + ' / ' + item.SubState + '</span></td>' +
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
        alert("server reply: " + data.status + "/" + data.statusText + "\n" + data.responseText);
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
        alert("server reply: " + data.status + "/" + data.statusText + "\n" + data.responseText);
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

  $(function() {
    $('[data-toggle="tooltip"]').tooltip()
  })
});
