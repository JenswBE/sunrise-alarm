<template>
  <v-container fluid fill-height>
    <v-row class="fill-height">
      <v-col align-self="end" cols="12" class="text-center">
        <p
          class="text-h1 font-weight-bold"
          style="font-family: DejaVu Sans Mono, monospace !important"
        >
          {{ hour }}{{ sep }}{{ minute }}
        </p>
        <p class="text-h4">{{ date }}</p>
      </v-col>
      <v-col align-self="end" cols="12" class="text-center">
        <p class="text-subtitle-1">Next alarm: Tomorrow at 00:00</p>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import helpers from "../helpers";

export default {
  name: "Clock",

  data: () => ({
    hour: "00",
    minute: "00",
    sep: ":",
    date: "",
  }),

  methods: {
    updateDateTime() {
      const now = new Date();
      this.hour = helpers.padTime(now.getHours());
      this.minute = helpers.padTime(now.getMinutes());
      this.sep = now.getSeconds() % 2 == 0 ? ":" : " ";
      this.date = now.toLocaleDateString("en-UK", {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
      });
      this.$options.timer = window.setTimeout(this.updateDateTime, 1000);
    },
  },

  mounted() {
    this.$options.timer = window.setTimeout(this.updateDateTime, 1000);
    this.updateDateTime();
  },

  beforeDestroy() {
    window.clearTimeout(this.$options.timer);
    this.updateDateTime();
  },
};
</script>
