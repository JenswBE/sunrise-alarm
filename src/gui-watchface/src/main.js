import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";

Vue.config.productionTip = false;
Vue.prototype.$days = [
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
  "Sunday",
]

Vue.filter('padTime', function (value) {
  if (!value) return "00";
  if (value < 10) return "0" + value.toString();
  return value;
})

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
