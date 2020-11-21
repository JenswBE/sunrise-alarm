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
              :class="alarm | getTextColor"
            >
              <v-icon v-if="alarm.enabled">mdi-alarm</v-icon>
              <v-icon v-else>mdi-alarm-off</v-icon>
            </v-btn>
          </v-list-item-action>

          <v-list-item-content>
            <v-list-item-title :class="alarm | getTextColor">
              {{ alarm | formatTime }}
            </v-list-item-title>
            <v-list-item-subtitle :class="alarm | getTextColor">
              _ _ _ _ _ _ _
            </v-list-item-subtitle>
          </v-list-item-content>

          <v-list-item-action>
            <v-btn
              icon
              @click.stop="$emit('skip-alarm', alarm.id)"
              :class="alarm | getTextColor"
            >
              <v-icon>mdi-debug-step-over</v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list-item-group>
    </v-col>
  </v-row>
</template>

<script>
import { mapGetters } from "vuex";
import helpers from "../helpers";

export default {
  name: "AlarmList",

  props: {
    value: String,
  },

  computed: {
    ...mapGetters(["sortedAlarms"]),
  },

  filters: {
    formatTime: (alarm) => {
      const time = helpers.formatTime(alarm);
      return alarm.name ? `${time} - ${alarm.name}` : time;
    },

    getTextColor: (alarm) => {
      if (alarm.skip_next) return "red--text";
      if (!alarm.enabled) return "grey--text";
      return "";
    },
  },
};
</script>