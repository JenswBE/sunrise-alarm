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
Vue.prototype.$emptyAlarm = {
  enabled: false,
  name: "",
  hour: 0,
  minute: 0,
  days: [],
  skip_next: false,
};
Vue.prototype.$scrollToTop = () => window.scrollTo(0, 0)

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
