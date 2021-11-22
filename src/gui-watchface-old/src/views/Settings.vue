<template>
  <v-container>
    <v-row align="center">
      <v-col cols="12" class="text-center">
        <v-btn @click="$emit('toggle-fullscreen')">Toggle fullscreen</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "Settings",
};
</script>
