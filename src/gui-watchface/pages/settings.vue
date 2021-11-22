<template>
  <v-container>
    <v-row align="center">
      <v-col cols="12" class="text-center">
        <v-btn @click="toggleFullscreen">Toggle fullscreen</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  name: 'Settings',

  methods: {
    toggleFullscreen() {
      // Based on https://developer.mozilla.org/en-US/docs/Web/API/Fullscreen_API
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen()
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen()
        }
      }
    },
  },
})
</script>
