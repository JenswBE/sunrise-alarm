import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import cloneDeep from "lodash.clonedeep";

Vue.use(Vuex);
axios.defaults.baseURL = "http://localhost:8004"

export default new Vuex.Store({
  state: {
    alarms: {},
    alert: {},
  },

  getters: {
    getAlarm: (state) => (id) => cloneDeep(state.alarms[id]),

    sortedAlarms: (state) => (
      Object.values(state.alarms).sort(
        (a, b) => a.hour * 100 + a.minute - b.hour * 100 - b.minute
      )
    )
  },

  mutations: {
    setAlarms(state, data) {
      console.debug("mut setAlarms - Input", data)
      state.alarms = data.alarms.reduce((result, item) => { result[item.id] = item; return result }, {})
      console.debug("mut setAlarms - Output", state.alarms)
    },

    upsertAlarm(state, alarm) {
      console.debug("mut upsertAlarm - Input", alarm)
      Vue.set(state.alarms, alarm.id, alarm)
      console.debug("mut upsertAlarm - Output", state.alarms)
    },

    deleteAlarm(state, alarmID) {
      console.debug("mut deleteAlarm - Input", alarmID)
      Vue.delete(state.alarms, alarmID)
      console.debug("mut deleteAlarm - Output", state.alarms)
    },

    setAlert(state, alert) {
      state.alert = alert
    },

    clearAlert(state) {
      state.alert = {}
    },
  },

  actions: {
    async getAlarms(context) {
      const { data } = await axios.get(`/alarms`);
      context.commit("setAlarms", data);
    }
  },

  modules: {}
});
