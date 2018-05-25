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

// unit item component
Vue.component('unit-item', {
  props: ['unit'],
  template: `<tr>
    <td v-bind:title="unit.Description">{{ unit.Name }}</td>
    <td><unit-badge v-bind:state="unit.ActiveState" v-bind:substate="unit.SubState"></unit-badge></td>
    <td><unit-buttons v-bind:name="unit.Name" v-bind:started="unit.ActiveState === 'active'"></unit-buttons></td>
  </tr>`,
});

// unit item badge
Vue.component('unit-badge', {
  props: ['state', 'substate'],
  template: `<span class="badge" v-bind:class="badgeClass">{{ state }} / {{ substate }}</span>`,

  computed: {
    badgeClass: function() {
      var self = this;

      switch(self.state) {
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
    }
  }
});

// unit item buttons
Vue.component('unit-buttons', {
  props: ['name', 'started'],
  template: `<div class="unit-buttons">
    <unit-button-start v-if="!started" v-bind:name="name"></unit-button-start>
    <unit-button-stop v-else v-bind:name="name"></unit-button-stop>
    <unit-button-journal v-bind:name="name"></unit-button-journal>
  </div>`,
});

Vue.component('unit-button-start', {
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

Vue.component('unit-button-stop', {
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

Vue.component('unit-button-journal', {
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
