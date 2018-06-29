const HTTP = axios.create();

// header component
Vue.component('navbar', {
  props: ['count'],
  template: `
<nav class="navbar is-dark">
  <div class="navbar-brand">
    <a class="navbar-item">sd-dbus-hooks ({{ count }})</a>
  </div>
  <div class="navbar-menu">
    <div class="navbar-end">
      <div class="navbar-item">
        <a class="button is-success" v-on:click="$emit('get-units')">reload</a>
      </div>
    </div>
  </div>
</nav>`,
});

Vue.component('unit-start-button', {
  props: ['name'],
  template: `<a class="button is-primary" v-bind:name="name" v-on:click="start">start</a>`,

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
  template: `<a class="button is-danger" v-bind:name="name" v-on:click="stop">stop</a>`,

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
  template: `<a class="button is-info" v-bind:name="name" v-on:click="showModal">journal</a>`,

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
      app.showJournal = true;
    },
  }
});

// unit item component
Vue.component('unit-item', {
  props: ['unit'],
  template: `
<tr>
  <td>
    {{ title }} <span class="tag" v-bind:class="tagClass">{{ unit.ActiveState }} / {{ unit.SubState }}</span>
  </td>
  <td>
    <div class="field has-addons">
      <p class="control">
        <unit-start-button v-if="!isUnitActive" v-bind:name="unit.Name"></unit-start-button>
        <unit-stop-button v-else v-bind:name="unit.Name"></unit-stop-button>
      </p>
      <p class="control">
        <unit-journal-button v-bind:name="unit.Name"></unit-journal-button>
      </p>
    </div>
  </td>
</tr>`,

  computed: {
    tagClass: function() {
      switch(this.unit.ActiveState) {
        case "active":
          return 'is-success';
        case "inactive":
          return 'is-light';
        case "failed":
          return 'is-danger';
        default:
          return 'is-warning';
      };
    },

    isUnitActive: function() {
      return this.unit.ActiveState === 'active';
    },

    title: function() {
      if (this.unit.Description) {
        return `${ this.unit.Name } (${ this.unit.Description })`;
      } else {
        return this.unit.Name;
      };
    },
  }
});

// journal modal window
Vue.component('journal-modal', {
  props: ['item'],
  template: `
<div class="modal">
  <div class="modal-background"></div>
  <div class="modal-card">
    <header class="modal-card-head">
      <p class="modal-card-title">journal for {{ item.name }}</p>
      <button class="delete" aria-label="close" v-on:click="hide"></button>
    </header>
    <section class="modal-card-body">
      <pre>{{ item.data }}</pre>
    </section>
  </div>
</div>`,

  methods: {
    hide: function() { app.showJournal = false; }
  }
});

// enter x-token modal window
Vue.component('login-modal', {
  data: function() {
    return {
      token: "",
    }
  },

  template: `
  <div class="modal">
  <div class="modal-background"></div>
  <div class="modal-card">
    <section class="modal-card-body">
      <label class="label">Token</label>
      <div class="field has-addons">
        <div class="control is-expanded">
          <input class="input" type=text placeholder="enter token" v-model="token"></input>
        </div>
        <div class="control">
          <a class="button is-primary" v-on:click="login">login</a>
        </div>
      </div>
    </section>
  </div>
</div>
  `,

  methods: {
    hide: function() { app.showLogin = false; },
    login: function() {
      HTTP.defaults.headers.common['X-Token'] = this.token;
      this.hide();
      app.isLoggedIn = true;
      app.getUnits();
    },
  }
})

// units table
Vue.component('units-table', {
  props: ['units'],
  template: `
<table class="table is-fullwidth">
  <thead>
    <tr>
      <th>unit</th>
      <th>actions</th>
    </tr>
  </thead>
  <tbody>
    <tr is="unit-item" v-for="item in units" v-bind:unit="item"></tr>
  </tbody>
</table>`,
});

var app = new Vue({
  el: '#app',

  data: {
    units: [],
    showJournal: false,
    showLogin: true,
    isLoggedIn: false,
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
      if (!this.isLoggedIn) {
        return false;
      };

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
