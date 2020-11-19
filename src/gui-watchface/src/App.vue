<template>
  <v-app>
    <v-main>
      <router-view />
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
    </v-bottom-navigation>
  </v-app>
</template>

<script>
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
};
</script>
