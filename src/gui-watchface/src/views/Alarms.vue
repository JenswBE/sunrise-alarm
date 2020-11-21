<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6">
        <alarm-list
          v-model="selectedAlarm"
          @add-alarm="addAlarm"
          @toggle-alarm="toggleAlarm"
          @skip-alarm="skipAlarm"
        />
      </v-col>
      <v-col cols="6">
        <alarm-edit :value="selectedAlarm" @delete-alarm="deleteAlarm" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from "axios";
import AlarmEdit from "../components/AlarmEdit.vue";
import AlarmList from "../components/AlarmList.vue";

export default {
  name: "Alarms",
  components: { AlarmEdit, AlarmList },

  data: () => ({
    selectedAlarm: "",
  }),

  methods: {
    async addAlarm() {
      axios
        .post("/alarms", this.$emptyAlarm)
        .then(({ data }) => {
          this.$store.commit("upsertAlarm", data);
          this.selectedAlarm = data.id;
        })
        .catch((e) =>
          this.$store.commit("setAlert", { type: "error", message: e })
        );
    },

    async toggleAlarm(alarmID) {
      let alarm = this.$store.getters.getAlarm(alarmID);
      alarm.enabled = !alarm.enabled;
      alarm.skip_next = false;
      await this.saveAlarm(alarm);
    },

    async skipAlarm(alarmID) {
      let alarm = this.$store.getters.getAlarm(alarmID);
      alarm.skip_next = !alarm.skip_next;
      await this.saveAlarm(alarm);
    },

    async saveAlarm(alarm) {
      axios
        .put(`/alarms/${alarm.id}`, alarm)
        .then(() => this.$store.commit("upsertAlarm", alarm))
        .catch((e) =>
          this.$store.commit("setAlert", { type: "error", message: e })
        );
    },

    async deleteAlarm() {
      axios
        .delete(`/alarms/${this.selectedAlarm}`)
        .then(() => {
          this.$store.commit("deleteAlarm", this.selectedAlarm);
          this.selectedAlarm = "";
        })
        .catch((e) =>
          this.$store.commit("setAlert", { type: "error", message: e })
        );
    },
  },

  mounted() {
    this.$store.dispatch("getAlarms");
  },
};
</script>