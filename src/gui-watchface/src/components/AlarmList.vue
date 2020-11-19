<template>
  <v-row>
    <v-col cols="12" class="text-center">
      <v-btn outlined block @click="$emit('add-alarm')">
        <v-icon>mdi-alarm-plus</v-icon>
      </v-btn>
    </v-col>
    <v-col cols="12">
      <v-list-item-group :value="value" @change="$emit('input', $event)">
        <v-list-item v-for="alarm in alarms" :key="alarm.id" :value="alarm.id">
          <template>
            <v-list-item-action>
              <v-btn icon @click.stop="">
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
              <v-btn icon @click.stop="">
                <v-icon>mdi-debug-step-over</v-icon>
              </v-btn>
            </v-list-item-action>
          </template>
        </v-list-item>
      </v-list-item-group>
    </v-col>
  </v-row>
</template>

<script>
export default {
  name: "AlarmList",

  props: {
    value: String,
    alarms: Object,
  },
};
</script>