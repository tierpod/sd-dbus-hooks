// const token = prompt("enter token");

const HTTP = axios.create({
  headers: {
    //'X-Token': token,
    'X-Token': '123',
  }
})

// header
Vue.component('navbar', {
  props: ['count'],
  template: `<nav class="navbar navbar-dark bg-dark">
    <a class="navbar-brand" href="#">sd-dbus-hooks ({{ count }})</a>
    <button class="btn btn-primary" type="button" v-on:click="$emit('get-units')">reload</button>
  </nav>`
});

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
