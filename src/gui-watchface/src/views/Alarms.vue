<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6" max-height="100%" class="overflow-y-auto">
        <v-list-item-group v-model="settings">
          <v-list-item v-for="alarm in alarms" :key="alarm.id">
            <template>
              <v-list-item-action>
                <v-btn icon>
                  <v-icon v-if="alarm.enabled">mdi-alarm-off</v-icon>
                  <v-icon v-else>mdi-alarm</v-icon>
                </v-btn>
              </v-list-item-action>

              <v-list-item-content>
                <v-list-item-title>
                  {{ alarm.hour | padTime }}:{{ alarm.minute | padTime }}
                </v-list-item-title>
                <v-list-item-subtitle>_ _ _ _ _ _ _</v-list-item-subtitle>
              </v-list-item-content>

              <v-list-item-action>
                <v-btn icon>
                  <v-icon>mdi-debug-step-over</v-icon>
                </v-btn>
              </v-list-item-action>
            </template>
          </v-list-item>
        </v-list-item-group>
      </v-col>
      <v-col cols="6">
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
                <v-btn
                  text
                  color="primary"
                  @click="$refs.dialog.save(alarmTime)"
                >
                  OK
                </v-btn>
              </v-time-picker>
            </v-dialog>

            <v-text-field label="Name"></v-text-field>
          </v-col>
          <v-col cols="6" class="pt-0">
            <v-list-item-group v-model="repeat" multiple>
              <v-list-item v-for="day in days" :key="day" :value="day">
                <template v-slot:default="{ active }">
                  <v-list-item-content>
                    <v-list-item-title v-text="day"></v-list-item-title>
                  </v-list-item-content>

                  <v-list-item-action>
                    <v-checkbox
                      :input-value="active"
                      color="primary"
                    ></v-checkbox>
                  </v-list-item-action>
                </template>
              </v-list-item>
            </v-list-item-group>
          </v-col>
          <v-col cols="6">
            <v-row v-for="(days, name) in dayPickers" :key="name">
              <v-col>
                <v-btn outlined v-text="name" @click="pickDays(days)"></v-btn>
              </v-col>
            </v-row>
          </v-col>
          <v-col cols="12" class="text-center">
            <v-btn outlined>
              <v-icon left>mdi-alarm-off</v-icon>
              Delete alarm
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "Home",

  data: () => ({
    settings: [],
    alarmTime: "07:00",
    timePicker: false,
    days: [
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday",
      "Sunday",
    ],
    repeat: [],
    dayPickers: {
      None: [],
      Week: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"],
      Weekend: ["Saturday", "Sunday"],
      All: [
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
        "Sunday",
      ],
    },
  }),

  computed: {
    ...mapState(["alarms"]),
  },

  methods: {
    pickDays(days) {
      this.repeat = days;
    },
  },

  filters: {
    padTime(value) {
      if (!value) return "00";
      if (value < 10) return "0" + value.toString();
      return value;
    },
  },

  mounted() {
    this.$store.dispatch("getAlarms");
  },
};
</script>