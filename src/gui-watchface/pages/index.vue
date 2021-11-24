<template>
  <v-container fluid fill-height>
    <v-row align="center">
      <v-col cols="12" class="text-center">
        <p
          class="text-h1 font-weight-bold"
          style="font-family: DejaVu Sans Mono, monospace !important"
        >
          {{ currentTime }}
        </p>
        <p class="text-h4">{{ currentDate }}</p>
      </v-col>
      <v-col cols="12" class="text-center">
        <p class="text-subtitle-1">{{ nextAlarmText }}</p>
      </v-col>
      <v-col
        cols="weatherCols"
        class="text-center"
        v-if="weatherInside.temperature"
      >
        <p class="text-subtitle-1">
          <v-icon class="mr-2">mdi-home</v-icon>
          {{ weatherInside.temperature.toFixed(1) }}&#8451; -
          {{ weatherInside.humidity.toFixed(0) }}%
        </p>
      </v-col>
      <v-col
        cols="weatherCols"
        class="text-center"
        v-if="weatherOutside.temperature"
      >
        <p class="text-subtitle-1">
          <v-icon class="mr-2">mdi-sun-thermometer</v-icon>
          {{ weatherOutside.temperature.toFixed(1) }}&#8451; -
          {{ weatherOutside.humidity.toFixed(0) }}%
        </p>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import { DateTime, Duration } from 'luxon'
import { UpdateWeatherOutsidePayload } from '~/store/general'

export default Vue.extend({
  name: 'Clock',

  data: () => ({
    now: DateTime.local(),
    timer: 0,
    timerWeather: 1,
  }),

  computed: {
    ...mapState('alarms', ['nextAlarm']),
    ...mapState('general', ['weatherInside', 'weatherOutside']),

    currentTime(): string {
      const format = this.now.second % 2 == 0 ? 'HH:mm' : 'HH mm'
      return this.now.toFormat(format)
    },

    currentDate(): string {
      return this.now.setLocale('en-UK').toLocaleString({
        weekday: 'long',
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      })
    },

    weatherCols(): number {
      return this.weatherInside.temperature && this.weatherOutside.temperature
        ? 6
        : 12
    },

    nextAlarmText() {
      // Check if set
      if (this.nextAlarm === '') {
        return ''
      }

      // Setup variables
      let day = ''
      const nextAlarmDate = DateTime.fromISO(this.nextAlarm)
      const tomorrow = this.now.plus({ days: 1 })

      // Check if alarm is today
      if (nextAlarmDate.weekday == this.now.weekday) {
        day = 'Today'
      } else if (nextAlarmDate.weekday == tomorrow.weekday) {
        day = 'Tomorrow'
      } else {
        day = nextAlarmDate.toFormat('cccc')
      }

      // Format and set text
      const alarm_time = nextAlarmDate.toFormat('HH:mm')
      return `Next alarm: ${day} at ${alarm_time}`
    },
  },

  methods: {
    updateDateTime() {
      this.now = DateTime.now()
      this.timer = window.setTimeout(this.updateDateTime, 1000)
    },

    updateWeather() {
      const payload: UpdateWeatherOutsidePayload = {
        cityID: this.$config.openWeather.cityID,
        apiKey: this.$config.openWeather.apiKey,
      }
      this.$store.dispatch('general/updateWeatherOutside', payload)
      const timeout = Duration.fromObject({ minutes: 10 }).toMillis()
      this.timer = window.setTimeout(this.updateWeather, timeout)
    },
  },

  mounted() {
    this.updateDateTime()
    this.updateWeather()
  },

  beforeDestroy() {
    window.clearTimeout(this.timer)
    window.clearTimeout(this.timerWeather)
  },
})
</script>
