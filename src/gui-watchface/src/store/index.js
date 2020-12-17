import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import cloneDeep from "lodash.clonedeep";

Vue.use(Vuex);
axios.defaults.baseURL = `http://${window.location.hostname}:8004`;

export default new Vuex.Store({
  state: {
    alarms: {},
    alert: {},
    nextAlarm: "",
  },

  getters: {
    getAlarm: (state) => (id) => cloneDeep(state.alarms[id]),

    sortedAlarms: (state) =>
      Object.values(state.alarms).sort(
        (a, b) => a.hour * 100 + a.minute - b.hour * 100 - b.minute
      ),
  },

  mutations: {
    setAlarms(state, data) {
      console.debug("mut setAlarms - Input", data);
      state.alarms = data.reduce((result, item) => {
        result[item.id] = item;
        return result;
      }, {});
      console.debug("mut setAlarms - Output alarms", state.alarms);
    },

    setNextAlarms(state, data) {
      console.debug("mut setNextAlarms - Input", data);
      if (
        data["ring"] == undefined ||
        data["ring"]["alarm_datetime"] == undefined
      ) {
        state.nextAlarm = "";
      } else {
        state.nextAlarm = data["ring"]["alarm_datetime"];
      }
      console.debug("mut setNextAlarms - Output nextAlarm", state.nextAlarm);
    },

    upsertAlarm(state, alarm) {
      console.debug("mut upsertAlarm - Input", alarm);
      Vue.set(state.alarms, alarm.id, cloneDeep(alarm));
      console.debug("mut upsertAlarm - Output alarms", state.alarms);
    },

    deleteAlarm(state, alarmID) {
      console.debug("mut deleteAlarm - Input", alarmID);
      Vue.delete(state.alarms, alarmID);
      console.debug("mut deleteAlarm - Output alarms", state.alarms);
    },

    setAlert(state, alert) {
      state.alert = alert;
    },

    clearAlert(state) {
      state.alert = {};
    },
  },

  actions: {
    handleMQTTMessage(context, { topic, payload }) {
      console.debug("action handleMQTTMessage - Input", { topic, payload });
      switch (topic) {
        case "sunrise_alarm/next_alarms_updated":
          context.commit("setNextAlarms", payload);
          break;
      }
    },

    async getAlarms(context) {
      axios
        .get(`/alarms`)
        .then(({ data }) => {
          context.commit("setAlarms", data);
        })
        .catch((e) => {
          const msg = `Unable to fetch alarms: ${e.message}`;
          context.commit("setAlert", { type: "error", message: msg });
        });
    },

    async getNextAlarms(context) {
      axios
        .get(`/alarms/next`)
        .then(({ data }) => {
          context.commit("setNextAlarms", data);
        })
        .catch((e) => {
          const msg = `Unable to fetch next alarms: ${e.message}`;
          context.commit("setAlert", { type: "error", message: msg });
        });
    },
  },

  modules: {},
});
