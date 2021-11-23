import { Plugin, Context } from '@nuxt/types'
import mqtt from 'mqtt'

const topics = [
  'sunrise_alarm/next_alarms_updated',
  'sunrise_alarm/temp_humid_updated',
]

declare module 'vue/types/vue' {
  // this.$mqtt inside Vue components
  interface Vue {
    $mqtt: mqtt.MqttClient
  }
}

declare module '@nuxt/types' {
  // nuxtContext.app.$mqtt inside asyncData, fetch, plugins, middleware, nuxtServerInit
  interface NuxtAppOptions {
    $mqtt: mqtt.MqttClient
  }
  // nuxtContext.$mqtt
  interface Context {
    $mqtt: mqtt.MqttClient
  }
}

declare module 'vuex/types/index' {
  // this.$mqtt inside Vuex stores
  interface Store<S> {
    $mqtt: mqtt.MqttClient
  }
}

const mqttPlugin: Plugin = (context, inject) => {
  let client = connectToMQTT(context)
  inject('mqtt', client)
}

export default mqttPlugin

function connectToMQTT(context: Context): mqtt.MqttClient {
  // Settings
  const connectUrl = `ws://${window.location.hostname}:9001/mqtt`
  const clientId = generateClientID()

  // Connect to MQTT
  let client = {} as mqtt.MqttClient
  try {
    client = mqtt.connect(connectUrl, { clientId })
  } catch (e: any) {
    const msg = `Unable to connect to MQTT: ${e.message}`
    context.store.commit('general/SET_ALERT', { type: 'error', message: msg })
  }
  client.on('connect', () => {
    const msg = `Connected to MQTT broker`
    context.store.commit('general/SET_ALERT', {
      type: 'success',
      message: msg,
    })
  })
  client.on('error', (e: any) => {
    const msg = `Connection to MQTT broker failed: ${e.message}`
    context.store.commit('general/SET_ALERT', { type: 'error', message: msg })
  })
  client.on('message', (topic, message) => {
    context.store.dispatch('general/handleMQTTMessage', {
      topic,
      payload: JSON.parse(message.toString()),
    })
  })

  // Subscribe to topics
  for (const topic of topics) {
    client.subscribe(topic, { qos: 1 }, (error, result) => {
      if (error) {
        const msg = `Failed to subscribe to topic "${topic}": ${error.message}`
        context.store.commit('general/SET_ALERT', {
          type: 'error',
          message: msg,
        })
      } else {
        console.debug('Currently subscribed to MQTT topics', result)
      }
    })
  }

  // Return client
  return client
}

function generateClientID(): string {
  const clientSuffixNumber = Math.floor(Math.random() * Math.floor(2 ** 32))
  const clientSuffix = clientSuffixNumber.toString(16).padStart(8, '0')
  return 'gui-watchface-' + clientSuffix
}
