<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6" class="make-scrollable">
        <alarm-list
          v-model="selectedAlarm"
          @add-alarm="addAlarm"
          @toggle-alarm="toggleAlarm"
          @skip-alarm="skipAlarm"
        />
      </v-col>
      <v-col cols="6" class="make-scrollable">
        <alarm-edit :value="selectedAlarm" @delete-alarm="deleteAlarm" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import AlarmEdit from '../components/AlarmEdit.vue'
import AlarmList from '../components/AlarmList.vue'
import { EMPTY_ALARM, Alarm } from '../store/alarms'

export default Vue.extend({
  name: 'Alarms',
  components: { AlarmEdit, AlarmList },

  data: () => ({
    selectedAlarm: '',
  }),

  computed: {
    ...mapGetters({
      getAlarm: 'alarms/getAlarm',
    }),
  },

  methods: {
    async addAlarm() {
      this.$axios
        .post('/alarms', EMPTY_ALARM)
        .then(({ data }) => {
          this.$store.commit('alarms/UPSERT_ALARM', data)
          this.selectedAlarm = data.id
        })
        .catch((e: any) => {
          const msg = `Unable to add alarm: ${e.message}`
          this.$store.commit('general/SET_ALERT', {
            type: 'error',
            message: msg,
          })
        })
    },

    async toggleAlarm(alarmID: string) {
      let alarm = this.getAlarm(alarmID)
      alarm.enabled = !alarm.enabled
      alarm.skip_next = false
      await this.saveAlarm(alarm)
    },

    async skipAlarm(alarmID: string) {
      let alarm = this.getAlarm(alarmID)
      alarm.skip_next = !alarm.skip_next
      await this.saveAlarm(alarm)
    },

    async saveAlarm(alarm: Alarm) {
      this.$axios
        .put(`/alarms/${alarm.id}`, alarm)
        .then(() => this.$store.commit('alarms/UPSERT_ALARM', alarm))
        .catch((e: any) => {
          const msg = `Unable to save alarm: ${e.message}`
          this.$store.commit('general/SET_ALERT', {
            type: 'error',
            message: msg,
          })
        })
    },

    async deleteAlarm() {
      this.$axios
        .delete(`/alarms/${this.selectedAlarm}`)
        .then(() => {
          this.$store.commit('alarms/DELETE_ALARM', this.selectedAlarm)
          this.selectedAlarm = ''
        })
        .catch((e: any) => {
          const msg = `Unable to delete alarm: ${e.message}`
          this.$store.commit('general/SET_ALERT', {
            type: 'error',
            message: msg,
          })
        })
        .finally(() => this.$scrollToTop())
    },
  },

  mounted() {
    this.$store.dispatch('alarms/getAlarms')
  },
})
</script>

<style>
.make-scrollable {
  overflow-y: auto;
  max-height: 80vh;
}
</style>