<template>
  <v-row align="center">
    <v-col cols="12">
      <v-dialog
        ref="dialog"
        v-model="timePicker"
        :return-value.sync="alarmTime"
        persistent
        width="500px"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="alarmTime"
            label="Time"
            prepend-icon="mdi-clock-outline"
            readonly
            :disabled="disabled"
            v-bind="attrs"
            v-on="on"
          ></v-text-field>
        </template>
        <v-time-picker
          v-if="timePicker"
          v-model="alarmTime"
          landscape
          format="24hr"
          header-color="primary"
          full-width
        >
          <v-spacer></v-spacer>
          <v-btn text color="primary" @click="timePicker = false">
            Cancel
          </v-btn>
          <v-btn text color="primary" @click="$refs.dialog.save(alarmTime)">
            OK
          </v-btn>
        </v-time-picker>
      </v-dialog>

      <v-text-field
        label="Name"
        :disabled="disabled"
        v-model="alarm.name"
      ></v-text-field>
    </v-col>
    <v-col cols="6" class="pt-0">
      <v-list-item-group v-model="alarmDays" multiple>
        <v-list-item
          v-for="day in days"
          :key="day"
          :value="day"
          :disabled="disabled"
        >
          <template v-slot:default="{ active }">
            <v-list-item-content>
              <v-list-item-title v-text="day"></v-list-item-title>
            </v-list-item-content>

            <v-list-item-action>
              <v-checkbox
                :input-value="active"
                color="primary"
                :disabled="disabled"
              ></v-checkbox>
            </v-list-item-action>
          </template>
        </v-list-item>
      </v-list-item-group>
    </v-col>
    <v-col cols="6">
      <v-row v-for="(days, name) in dayPickers" :key="name">
        <v-col>
          <v-btn
            outlined
            v-text="name"
            @click="pickDays(days)"
            :disabled="disabled"
          ></v-btn>
        </v-col>
      </v-row>
    </v-col>
    <v-col cols="6" class="text-center">
      <v-btn outlined block :disabled="disabled" @click="saveAlarm(alarm)">
        <v-icon left>mdi-alarm-check</v-icon>
        Save alarm
      </v-btn>
    </v-col>
    <v-col cols="6" class="text-center">
      <v-dialog v-model="confirmDelete" max-width="290">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            v-bind="attrs"
            v-on="on"
            outlined
            block
            color="error"
            :disabled="disabled"
          >
            <v-icon left>mdi-alarm-off</v-icon>
            Delete alarm
          </v-btn>
        </template>
        <v-card>
          <v-card-title> Are you sure? </v-card-title>
          <v-card-actions>
            <v-btn outlined @click="confirmDelete = false"> Cancel </v-btn>
            <v-spacer></v-spacer>
            <v-btn
              color="red darken-1"
              outlined
              @click="
                confirmDelete = false;
                $emit('delete-alarm');
              "
            >
              Confirm
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import axios from "axios";
import helpers from "../helpers";
import { DAYS, EMPTY_ALARM } from "../constants";

export default {
  name: "AlarmEdit",

  props: {
    value: String,
  },

  watch: {
    value: {
      immediate: true,
      handler(alarmID) {
        if (alarmID) {
          this.alarm = this.$store.getters.getAlarm(alarmID);
        } else {
          this.alarm = EMPTY_ALARM;
        }
      },
    },
  },

  data: () => ({
    alarm: {},
    timePicker: false,
    confirmDelete: false,
    days: DAYS,
    dayPickers: {
      None: [],
      Week: DAYS.slice(0, 5),
      Weekend: DAYS.slice(5, 7),
      All: DAYS,
    },
  }),

  computed: {
    disabled() {
      return !this.value;
    },

    alarmTime: {
      get() {
        return helpers.formatTime(this.alarm);
      },
      set(value) {
        const [hour, minute] = value.split(":");
        this.alarm.hour = Number(hour);
        this.alarm.minute = Number(minute);
      },
    },

    alarmDays: {
      get() {
        return this.alarm.days.map((dayIndex) => DAYS[dayIndex]);
      },
      set(value) {
        this.alarm.days = value.map((day) => DAYS.indexOf(day));
      },
    },
  },

  methods: {
    pickDays(days) {
      this.alarmDays = days;
    },

    async saveAlarm(alarm) {
      axios
        .put(`/alarms/${alarm.id}`, alarm)
        .then(() => {
          this.$store.commit("upsertAlarm", alarm);
          this.$store.commit("setAlert", {
            type: "success",
            message: "Alarm successfully saved",
          });
          helpers.scrollToTop();
        })
        .catch((e) => {
          const msg = `Unable to save alarm: ${e.message}`;
          this.$store.commit("setAlert", { type: "error", message: msg });
          helpers.scrollToTop();
        });
    },
  },
};
</script>
