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
  template: `<button type="button" class="btn btn-sm btn-danger" v-bind:name="name" v-on:click="start">start</button>`,

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
  template: `<button type="button" class="btn btn-sm btn-info" v-bind:name="name" v-on:click="journal">journal</button>`,

  methods: {
    journal: function() {
      var self = this;
      console.log("get journal " + this.name);

      HTTP.get('/unit/journal/' + this.name)
      .then(function(responce) {
        console.log(responce.data);
      })
      .catch(function(error) {
        alert(error);
      });
    }
  }
});

var app = new Vue({
  el: '#app',

  data: {
    units: [],
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
