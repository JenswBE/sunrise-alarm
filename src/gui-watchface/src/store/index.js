import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);
axios.defaults.baseURL = "http://localhost:8004"

export default new Vuex.Store({
  state: {
    alarms: [],
  },
  mutations: {
    setAlarms(state, data) {
      console.log(data)
      state.alarms = data.alarms;
    }
  },
  actions: {
    async getAlarms(context) {
      const { data } = await axios.get(`/alarms`);
      context.commit("setAlarms", data);
    }
  },
  modules: {}
});
