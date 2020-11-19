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

      <v-text-field label="Name" :disabled="disabled"></v-text-field>
    </v-col>
    <v-col cols="6" class="pt-0">
      <v-list-item-group v-model="repeat" multiple>
        <v-list-item
          v-for="day in $days"
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
      <v-btn outlined block :disabled="disabled">
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
export default {
  name: "AlarmEdit",

  props: {
    alarm: Object,
  },

  data: ({ $root }) => ({
    timePicker: false,
    confirmDelete: false,
    repeat: [],
    dayPickers: {
      None: [],
      Week: $root.$days.slice(0, 5),
      Weekend: $root.$days.slice(5, 7),
      All: $root.$days,
    },
  }),

  computed: {
    alarmTime: {
      get() {
        if (!this.alarm) return "00:00";
        const { hour, minute } = this.alarm;
        const padTime = this.$options.filters.padTime;
        return `${padTime(hour)}:${padTime(minute)}`;
      },
      set(value) {
        const [hour, minute] = value.split(":");
        this.alarm.hour = Number(hour);
        this.alarm.minute = Number(minute);
      },
    },

    disabled() {
      return !this.alarm;
    },
  },

  methods: {
    pickDays(days) {
      this.repeat = days;
    },
  },
};
</script>