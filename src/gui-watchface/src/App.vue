<template>
  <v-app>
    <v-main>
      <router-view @toggle-fullscreen="toggleFullscreen" />
    </v-main>

    <v-snackbar v-model="showAlert" :color="alert.type" top absolute>
      {{ alert.message }}
    </v-snackbar>

    <v-bottom-navigation app grow color="primary" class="pt-2">
      <v-btn value="clock" to="/" exact>
        <span>CLOCK</span>
        <v-icon>mdi-clock-outline</v-icon>
      </v-btn>

      <v-btn value="alarms" to="/alarms" exact>
        <span>ALARMS</span>
        <v-icon>mdi-alarm</v-icon>
      </v-btn>

      <v-btn value="settings" to="/settings" exact>
        <span>SETTINGS</span>
        <v-icon>mdi-wrench-outline</v-icon>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script>
import mqtt from "mqtt";
import { mapState } from "vuex";

export default {
  name: "App",

  computed: {
    ...mapState(["alert"]),

    showAlert: {
      get() {
        return Boolean(this.alert.message);
      },
      set() {
        this.$store.commit("clearAlert");
      },
    },
  },

  methods: {
    toggleFullscreen() {
      // Based on https://developer.mozilla.org/en-US/docs/Web/API/Fullscreen_API
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen();
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        }
      }
    },

    generateClientID() {
      const clientSuffixNumber = Math.floor(
        Math.random() * Math.floor(2 ** 32)
      );
      const clientSuffix = clientSuffixNumber.toString(16).padStart(8, "0");
      return "gui-watchface-" + clientSuffix;
    },

    connectToMQTT() {
      // Settings
      const connectUrl = `mqtt://${window.location.hostname}:9001`;
      const clientId = this.generateClientID();
      const topics = ["sunrise_alarm/next_alarms_updated"];

      // Connect to MQTT
      try {
        this.client = mqtt.connect(connectUrl, { clientId });
      } catch (e) {
        const msg = `Unable to connect to MQTT: ${e.message}`;
        this.$store.commit("setAlert", { type: "error", message: msg });
      }
      this.client.on("connect", () => {
        const msg = `Connected to MQTT broker`;
        this.$store.commit("setAlert", { type: "success", message: msg });
      });
      this.client.on("error", (e) => {
        const msg = `Connection to MQTT broker failed: ${e.message}`;
        this.$store.commit("setAlert", { type: "error", message: msg });
      });
      this.client.on("message", (topic, message) => {
        this.$store.dispatch("handleMQTTMessage", {
          topic,
          payload: JSON.parse(message.toString()),
        });
      });

      // Subscribe to topics
      for (const topic of topics) {
        this.client.subscribe(topic, { qos: 1 }, (error, result) => {
          if (error) {
            const msg = `Failed to subscribe to topic "${topic}": ${error.message}`;
            this.$store.commit("setAlert", { type: "error", message: msg });
          } else {
            console.debug("Currently subscribed to MQTT topics", result);
          }
        });
      }
    },
  },

  mounted() {
    this.$nextTick(function () {
      this.$store.dispatch("getNextAlarms");
      this.connectToMQTT();
    });
  },
};
</script>
