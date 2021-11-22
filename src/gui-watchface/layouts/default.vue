<template>
  <v-app dark>
    <v-snackbar
      v-model="showAlert"
      :color="alert.type"
      timeout="1500"
      top
      absolute
    >
      {{ alert.message }}
    </v-snackbar>

    <v-main>
      <Nuxt />
    </v-main>

    <v-bottom-navigation app grow color="primary">
      <v-btn to="/" exact>
        <span>CLOCK</span>
        <v-icon>mdi-clock-outline</v-icon>
      </v-btn>

      <v-btn to="/alarms" exact>
        <span>ALARMS</span>
        <v-icon>mdi-alarm</v-icon>
      </v-btn>

      <v-btn to="/settings" exact>
        <span>SETTINGS</span>
        <v-icon>mdi-wrench-outline</v-icon>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import mqtt from 'mqtt'

export default Vue.extend({
  data() {
    return {
      appel: 'lol',
    }
  },

  computed: {
    ...mapState('general', ['alert']),

    showAlert: {
      get() {
        return Boolean(this.alert.message)
      },
      set() {
        this.$store.commit('general/CLEAR_ALERT')
      },
    },
  },

  methods: {
    generateClientID(): string {
      const clientSuffixNumber = Math.floor(Math.random() * Math.floor(2 ** 32))
      const clientSuffix = clientSuffixNumber.toString(16).padStart(8, '0')
      return 'gui-watchface-' + clientSuffix
    },

    // connectToMQTT() {
    //   // Settings
    //   const connectUrl = `mqtt://${window.location.hostname}:9001`
    //   const clientId = this.generateClientID()
    //   const topics = [
    //     'sunrise_alarm/next_alarms_updated',
    //     'sunrise_alarm/temp_humid_updated',
    //   ]

    //   // Connect to MQTT
    //   try {
    //     this.client = mqtt.connect(connectUrl, { clientId })
    //   } catch (e) {
    //     const msg = `Unable to connect to MQTT: ${e.message}`
    //     this.$store.commit('setAlert', { type: 'error', message: msg })
    //   }
    //   this.client.on('connect', () => {
    //     const msg = `Connected to MQTT broker`
    //     this.$store.commit('setAlert', { type: 'success', message: msg })
    //   })
    //   this.client.on('error', (e: any) => {
    //     const msg = `Connection to MQTT broker failed: ${e.message}`
    //     this.$store.commit('setAlert', { type: 'error', message: msg })
    //   })
    //   this.client.on('message', (topic, message) => {
    //     this.$store.dispatch('handleMQTTMessage', {
    //       topic,
    //       payload: JSON.parse(message.toString()),
    //     })
    //   })

    //   // Subscribe to topics
    //   for (const topic of topics) {
    //     this.client.subscribe(topic, { qos: 1 }, (error, result) => {
    //       if (error) {
    //         const msg = `Failed to subscribe to topic "${topic}": ${error.message}`
    //         this.$store.commit('setAlert', { type: 'error', message: msg })
    //       } else {
    //         console.debug('Currently subscribed to MQTT topics', result)
    //       }
    //     })
    //   }
    // },
  },

  // mounted() {
  //   this.$nextTick(function () {
  //     this.$store.dispatch('alarms/getNextAlarms')
  //     this.connectToMQTT()
  //   })
  // },
})
</script>
