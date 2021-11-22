<template>
  <v-row>
    <v-col cols="12" class="text-center">
      <v-btn outlined block @click="$emit('add-alarm')">
        <v-icon>mdi-alarm-plus</v-icon>
      </v-btn>
    </v-col>
    <v-col cols="12">
      <v-list-item-group :value="value" @change="$emit('input', $event)">
        <v-list-item
          v-for="alarm in sortedAlarms"
          :key="alarm.id"
          :value="alarm.id"
        >
          <v-list-item-action>
            <v-btn
              icon
              @click.stop="$emit('toggle-alarm', alarm.id)"
              :class="getTextColor(alarm)"
            >
              <v-icon v-if="alarm.enabled">mdi-alarm</v-icon>
              <v-icon v-else>mdi-alarm-off</v-icon>
            </v-btn>
          </v-list-item-action>

          <v-list-item-content>
            <v-list-item-title :class="getTextColor(alarm)">
              {{ formatTime(alarm) }}
            </v-list-item-title>
            <v-list-item-subtitle :class="getTextColor(alarm)">
              {{ formatDays(alarm) }}
            </v-list-item-subtitle>
          </v-list-item-content>

          <v-list-item-action>
            <v-btn
              icon
              @click.stop="$emit('skip-alarm', alarm.id)"
              :class="getTextColor(alarm)"
            >
              <v-icon>mdi-debug-step-over</v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list-item-group>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { DAYS, Alarm } from '../store/alarms'

export default Vue.extend({
  name: 'AlarmList',

  props: {
    value: String,
  },

  computed: {
    ...mapGetters('alarms', ['sortedAlarms']),
  },

  methods: {
    formatTime(alarm: Alarm): string {
      const time = this.$formatTime(alarm)
      return alarm.name ? `${time} - ${alarm.name}` : time
    },

    getTextColor(alarm: Alarm) {
      if (alarm.skip_next) return 'red--text'
      if (!alarm.enabled) return 'grey--text'
      return ''
    },

    formatDays(alarm: Alarm) {
      return DAYS.map((day, i) => (alarm.days.includes(i) ? day[0] : '_')).join(
        ' '
      )
    },
  },
})
</script>
