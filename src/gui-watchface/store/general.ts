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

export type Weather = {
  temperature: number
  humidity: number
}

export const state = () => ({
  alert: {} as Alert,
  weatherInside: {} as Weather,
  weatherOutside: {} as Weather,
})

export type RootState = ReturnType<typeof state>

export const mutations: MutationTree<RootState> = {
  SET_ALERT(state, alert) {
    console.debug('mut SET_ALERT - Input', alert)
    state.alert = alert
  },

  CLEAR_ALERT(state) {
    console.debug('mut CLEAR_ALERT')
    state.alert = {} as Alert
  },

  SET_WEATHER_INSIDE(state, weather: Weather) {
    console.debug('mut SET_WEATHER_INSIDE - Input', weather)
    state.weatherInside = weather
  },

  SET_WEATHER_OUTSIDE(state, weather: Weather) {
    console.debug('mut SET_WEATHER - Input', weather)
    state.weatherOutside = weather
  },
}

export type UpdateWeatherOutsidePayload = {
  cityID: string
  apiKey: string
}

export const actions: ActionTree<RootState, RootState> = {
  handleMQTTMessage(context, { topic, payload }) {
    console.debug('action handleMQTTMessage - Input', { topic, payload })
    switch (topic) {
      case 'sunrise_alarm/next_alarms_updated':
        context.commit('alarms/SET_NEXT_ALARMS', payload, { root: true })
        break
      case 'sunrise_alarm/temp_humid_updated':
        const weather: Weather = {
          temperature: payload.temperature,
          humidity: payload.humidity,
        }
        context.commit('SET_WEATHER_INSIDE', weather)
        break
    }
  },

  async updateWeatherOutside(context, payload: UpdateWeatherOutsidePayload) {
    const url = `/openweather/data/2.5/weather?id=${payload.cityID}&appid=${payload.apiKey}&units=metric`
    this.$axios
      .get(url, { baseURL: '' })
      .then(({ data }) => {
        const weather: Weather = {
          temperature: data.main.temp,
          humidity: data.main.humidity,
        }
        context.commit('SET_WEATHER_OUTSIDE', weather)
      })
      .catch((e: any) => {
        const msg = `Unable to fetch weather: ${e.message}`
        context.commit(
          'general/SET_ALERT',
          { type: 'error', message: msg },
          { root: true }
        )
      })
  },
}
