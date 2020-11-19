<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6">
        <alarm-list
          :alarms="alarms"
          v-model="selectedAlarm"
          @add-alarm="addAlarm"
        />
      </v-col>
      <v-col cols="6">
        <alarm-edit
          :alarm="alarms[selectedAlarm]"
          @delete-alarm="deleteAlarm"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from "axios";
import { mapState } from "vuex";
import AlarmEdit from "../components/AlarmEdit.vue";
import AlarmList from "../components/AlarmList.vue";

export default {
  name: "Alarms",
  components: { AlarmEdit, AlarmList },

  data: () => ({
    selectedAlarm: "",
  }),

  computed: {
    ...mapState(["alarms"]),
  },

  methods: {
    async addAlarm() {
      const { data } = await axios.post("/alarms", {
        enabled: false,
        name: "",
        hour: 0,
        minute: 0,
        days: [],
        skip_next: false,
      });
      this.$store.commit("upsertAlarm", data);
      this.selectedAlarm = data.id;
    },

    async deleteAlarm() {
      await axios.delete(`/alarms/${this.selectedAlarm}`);
      this.$store.commit("deleteAlarm", this.selectedAlarm);
      this.selectedAlarm = "";
    },
  },

  mounted() {
    this.$store.dispatch("getAlarms");
  },
};
</script>