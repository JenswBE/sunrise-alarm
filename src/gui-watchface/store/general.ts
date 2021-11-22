import { GetterTree, MutationTree, ActionTree } from 'vuex'
import cloneDeep from 'lodash.clonedeep'
import Vue from 'vue'

export type Alert = {
  type: AlertType
  message: string
}

export enum AlertType {
  // Based on https://vuetifyjs.com/en/api/v-alert/#props-type
  Success = 'success',
  Info = 'info',
  Warning = 'warning',
  Error = 'error',
}

export const state = () => ({
    alert: {} as Alert,
    temperature: "",
    humidity: "",
})

export type RootState = ReturnType<typeof state>

export const mutations: MutationTree<RootState> = {
  SET_ALERT(state, alert) {
    console.debug("mut setAlert - Input", alert);
    state.alert = alert;
  },

  CLEAR_ALERT(state) {
    console.debug("mut clearAlert");
    state.alert = {} as Alert;
  },

  SET_TEMP_HUMID(state, tempHumid) {
    console.debug("mut setTempHumid - Input", tempHumid);
    state.temperature = tempHumid.temperature;
    state.humidity = tempHumid.humidity;
  },
}

export const actions: ActionTree<RootState, RootState> = {
  handleMQTTMessage(context, { topic, payload }) {
    console.debug("action handleMQTTMessage - Input", { topic, payload });
    switch (topic) {
      case "sunrise_alarm/next_alarms_updated":
        context.commit("setNextAlarms", payload);
        break;
      case "sunrise_alarm/temp_humid_updated":
        context.commit("setTempHumid", payload);
        break;
    }
  },
}