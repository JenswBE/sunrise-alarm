<template>
  <v-app>
    <v-main>
      <router-view />
    </v-main>

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
export default {
  name: "App",
};
</script>
