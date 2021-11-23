<template>
  <v-app dark>
    <v-snackbar
      v-model="showAlert"
      :color="alert.type"
      timeout="1500"
      top
      absolute
    >
      {{ alert.message }}
    </v-snackbar>

    <v-main>
      <Nuxt />
    </v-main>

    <v-bottom-navigation app grow color="primary">
      <v-btn to="/" exact>
        <span>CLOCK</span>
        <v-icon>mdi-clock-outline</v-icon>
      </v-btn>

      <v-btn to="/alarms" exact>
        <span>ALARMS</span>
        <v-icon>mdi-alarm</v-icon>
      </v-btn>

      <v-btn to="/settings" exact>
        <span>SETTINGS</span>
        <v-icon>mdi-wrench-outline</v-icon>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
  computed: {
    ...mapState('general', ['alert']),

    showAlert: {
      get(): boolean {
        return Boolean(this.alert.message)
      },
      set() {
        this.$store.commit('general/CLEAR_ALERT')
      },
    },
  },

  mounted() {
    this.$nextTick(function () {
      this.$store.dispatch('alarms/getNextAlarms')
    })
  },
})
</script>
