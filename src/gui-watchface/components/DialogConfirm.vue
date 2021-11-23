<template>
  <v-dialog :value="value" @input="$emit('input', false)" max-width="500px">
    <v-card class="pa-5">
      <v-card-title class="headline"> Are you sure? </v-card-title>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn x-large text @click="emitAndClose('cancel')"> Cancel </v-btn>
        <v-btn
          x-large
          :color="confirmColor"
          text
          @click="emitAndClose('confirm')"
        >
          Confirm
        </v-btn>
        <v-spacer></v-spacer>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue, { PropOptions } from 'vue'

export default Vue.extend({
  name: 'DialogConfirm',

  props: {
    value: {
      type: Boolean,
      required: true,
    } as PropOptions<boolean>,
    confirmColor: {
      type: String,
      default: 'primary',
    } as PropOptions<string>,
  },

  methods: {
    emitAndClose(event: string) {
      this.$emit(event)
      this.$emit('input', false)
    },
  },
})
</script>
