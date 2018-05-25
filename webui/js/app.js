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
    <td><unit-buttons v-bind:name="unit.Name"></unit-buttons></td>
  </tr>`,
  //template: `<li>{{ unit.Name }}</li>`,
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

Vue.component('unit-buttons', {
  props: ['name'],
  template: `<button type="button" class="btn btn-sm btn-danger" v-bind:name="name">stop</button>`,
})

var app = new Vue({
  el: '#app',

  data: {
    message: 'Hello Vue!',
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
        console.log(responce.data);
        self.units = responce.data;
      })
      .catch(function(error) {
        alert(error);
      });
    }
  }
});
