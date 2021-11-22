<template>
  <v-container fluid fill-height>
    <v-row align="center">
      <v-col cols="12" class="text-center">
        <p
          class="text-h1 font-weight-bold"
          style="font-family: DejaVu Sans Mono, monospace !important"
        >
          {{ hour }}{{ sep }}{{ minute }}
        </p>
        <p class="text-h4">{{ date }}</p>
      </v-col>
      <v-col cols="12" class="text-center">
        <p class="text-subtitle-1">{{ nextAlarmText }}</p>
        <p class="text-subtitle-1" v-show="temperature">
          {{ temperature }}&#8451; - {{ humidity }}%
        </p>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import { DateTime } from 'luxon'

export default Vue.extend({
  name: 'Clock',

  data: () => ({
    hour: '00',
    minute: '00',
    sep: ':',
    date: '',
    timer: 0,
  }),

  computed: {
    ...mapState('alarms', ['nextAlarm']),
    ...mapState('general', ['temperature', 'humidity']),

    nextAlarmText() {
      // Check if set
      if (this.nextAlarm === '') {
        return ''
      }

      // Setup variables
      let day = ''
      const nextAlarmDate = DateTime.fromISO(this.nextAlarm)
      const now = DateTime.local()
      const tomorrow = DateTime.local().plus({ days: 1 })

      // Check if alarm is today
      if (nextAlarmDate.weekday == now.weekday) {
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
      const now = DateTime.now()
      this.hour = now.toFormat('HH')
      this.minute = now.toFormat('mm')
      this.sep = now.second % 2 == 0 ? ':' : ' '
      this.date = now.setLocale('en-UK').toLocaleString({
        weekday: 'long',
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      })
      this.timer = window.setTimeout(this.updateDateTime, 1000)
    },
  },

  mounted() {
    this.timer = window.setTimeout(this.updateDateTime, 1000)
    this.updateDateTime()
  },

  beforeDestroy() {
    window.clearTimeout(this.timer)
    this.updateDateTime()
  },
})
</script>
