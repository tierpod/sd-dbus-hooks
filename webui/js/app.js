// const token = prompt("enter token");

const HTTP = axios.create({
  headers: {
    //'X-Token': token,
    'X-Token': '123',
  }
})

// header component
Vue.component('navbar', {
  props: ['count'],
  template: `<nav class="navbar navbar-dark bg-dark">
    <a class="navbar-brand" href="#">sd-dbus-hooks ({{ count }})</a>
    <button class="btn btn-primary" type="button" v-on:click="$emit('get-units')">reload</button>
  </nav>`,
});

Vue.component('unit-start-button', {
  props: ['name'],
  template: `<button type="button" class="btn btn-sm btn-primary" v-bind:name="name" v-on:click="start">start</button>`,

  methods: {
    start: function() {
      var self = this;
      console.log("start " + this.name);

      HTTP.get('/unit/start/' + this.name)
      .then(function(responce) {
        app.getUnits();
      })
      .catch(function(error) {
        alert(error);
      });
    }
  }
});

Vue.component('unit-stop-button', {
  props: ['name'],
  template: `<button type="button" class="btn btn-sm btn-danger" v-bind:name="name" v-on:click="stop">stop</button>`,

  methods: {
    stop: function() {
      var self = this;
      console.log("stop " + this.name);

      HTTP.get('/unit/stop/' + this.name)
      .then(function(responce) {
        app.getUnits();
      })
      .catch(function(error) {
        alert(error);
      });
    }
  }
});

Vue.component('unit-journal-button', {
  props: ['name'],
  template: `<button type="button" class="btn btn-sm btn-info" v-bind:name="name" v-on:click="showModal">journal</button>`,

  methods: {
    getData: function() {
      var self = this;
      console.log("get journal " + this.name);

      HTTP.get('/unit/journal/' + this.name)
      .then(function(responce) {
        app.journalItem = {
          name: self.name,
          data: responce.data,
        }
      })
      .catch(function(error) {
        alert(error);
      });
    },

    showModal: function() {
      var self = this;
      self.getData();
      $("#journal-modal").modal("show");
    },
  }
});

// unit item component
Vue.component('unit-item', {
  props: ['unit'],
  template: `
<tr>
  <td v-bind:title="unit.Description">{{ unit.Name }}</td>
  <td><span class="badge" v-bind:class="badgeClass">{{ unit.ActiveState }} / {{ unit.SubState }}</span></td>
  <td>
    <unit-start-button v-if="!isUnitActive" v-bind:name="unit.Name"></unit-start-button>
    <unit-stop-button v-else v-bind:name="unit.Name"></unit-stop-button>
    <unit-journal-button v-bind:name="unit.Name"></unit-journal-button>
  </td>
</tr>`,

  computed: {
    badgeClass: function() {
      switch(this.unit.ActiveState) {
        case "active":
          return 'badge-success';
          break
        case "inactive":
          return 'badge-secondary';
          break
        case "failed":
          return 'badge-danger';
          break
        default:
          return 'badge-warning';
          break
      };
    },

    isUnitActive: function() {
      return this.unit.ActiveState === 'active';
    }
  }
});

// journal modal window
Vue.component('journal-modal', {
  props: ['item'],
  template: `
  <div class="modal fade" id="journal-modal" role="dialog" aria-labelledby="JournalLabel" aria-hidden="true">
    <div class="modal-dialog modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">journal for: {{ item.name }}</h5>
            <button type="button" class="close" data-dismiss="modal" aria-label="Close" v-on:click="hide">
              <span aria-hidden="true">&times;</span>
            </button>
        </div>
        <div class="modal-body">
          <pre>
{{ item.data }}
          </pre>
        </div>
      </div>
    </div>
  </div>`,

  methods: {
    hide: function() { $("#journal-modal").modal("hide") }
  }
});

// units table
Vue.component('units-table', {
  props: ['units'],
  template: `
<div class="container">
  <div class="row">
      <div class="col">
          <table class="table table-striped" v-if="units.length">
              <thead>
                  <tr>
                  <th scope="col">unit name</th>
                  <th scope="col">status</th>
                  <th scope="col">actions</th>
                  </tr>
              </thead>
              <tbody>
                  <tr is="unit-item" v-for="item in units" v-bind:unit="item"></tr>
              </tbody>
          </table>
      </div>
  </div>
</div>`,
})

var app = new Vue({
  el: '#app',

  data: {
    units: [],
    journalItem: {
      name: "",
      data: "",
    },
  },

  created: function() {
    this.getUnits();
  },

  computed: {
    count: function() {
      return this.units.length;
    },
  },

  methods: {
    getUnits: function() {
      var self = this;

      HTTP.get('/unit/status/')
      .then(function(responce) {
        self.units = responce.data;
      })
      .catch(function(error) {
        alert(error);
      });
    }
  }
});
