import { GetterTree, MutationTree, ActionTree } from 'vuex'
import cloneDeep from 'lodash.clonedeep'
import Vue from 'vue'

export const DAYS = [
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
    "Sunday",
];

export const EMPTY_ALARM = {
    enabled: false,
    name: "",
    hour: 0,
    minute: 0,
    days: [],
    skip_next: false,
};

export type Alarm = {
  id?: string
  name: string
  enabled: boolean
  hour: number
  minute: number
  skip_next: boolean
  days: number[]
}

export type AlarmsMap = { [id: string]: Alarm }

export const state = () => ({
    alarms: {} as AlarmsMap,
    nextAlarm: "",
})

export type RootState = ReturnType<typeof state>

export const getters: GetterTree<RootState, RootState> = {
  getAlarm(state) {
    return function (id: string) {
      return cloneDeep(state.alarms[id])
    }
  },

  sortedAlarms(state): Alarm[] {
    return Object.values(state.alarms).sort(sortAlarms)
  },
}

export const mutations: MutationTree<RootState> = {
  SET_ALARMS(state, data) {
    console.debug("mut setAlarms - Input", data);
    state.alarms = data.reduce((result: AlarmsMap, item: Alarm) => {
      result[item.id as string] = item;
      return result;
    }, {});
    console.debug("mut setAlarms - Output alarms", state.alarms);
  },

  SET_NEXT_ALARMS(state, data) {
    console.debug("mut setNextAlarms - Input", data);
    if (
      data["ring"] == undefined ||
      data["ring"]["alarm_datetime"] == undefined
    ) {
      state.nextAlarm = "";
    } else {
      state.nextAlarm = data["ring"]["alarm_datetime"];
    }
    console.debug("mut setNextAlarms - Output nextAlarm", state.nextAlarm);
  },

  UPSERT_ALARM(state, alarm) {
    console.debug("mut upsertAlarm - Input", alarm);
    Vue.set(state.alarms, alarm.id, cloneDeep(alarm));
    console.debug("mut upsertAlarm - Output alarms", state.alarms);
  },

  DELETE_ALARM(state, alarmID) {
    console.debug("mut deleteAlarm - Input", alarmID);
    Vue.delete(state.alarms, alarmID);
    console.debug("mut deleteAlarm - Output alarms", state.alarms);
  },
}

export const actions: ActionTree<RootState, RootState> = {
  async getAlarms(context) {
    this.$axios
      .get(`/alarms`)
      .then(({ data }) => {
        context.commit("SET_ALARMS", data);
      })
      .catch((e: any) => {
        const msg = `Unable to fetch alarms: ${e.message}`;
        context.commit("SET_ALERT", { type: "error", message: msg });
      });
  },

  async getNextAlarms(context) {
    this.$axios
      .get(`/alarms/next`)
      .then(({ data }) => {
        context.commit("SET_NEXT_ALARMS", data);
      })
      .catch((e: any) => {
        const msg = `Unable to fetch next alarms: ${e.message}`;
        context.commit("SET_ALERT", { type: "error", message: msg });
      });
  },
}

function sortAlarms(a: Alarm, b: Alarm): number {
  // Sort by time
  const byTime = a.hour * 100 + a.minute - b.hour * 100 - b.minute;
  if (byTime !== 0) return byTime;

  // Sort by enabled
  const byEnabled = a.enabled !== b.enabled;
  if (byEnabled) return a.enabled ? -1 : 1;

  // Sort by skip next
  const bySkipNext = a.skip_next !== b.skip_next;
  if (bySkipNext) return a.skip_next ? 1 : -1;

  // Sort by repeated
  const byRepeated = (a.days.length === 0) !== (b.days.length === 0);
  if (byRepeated) return a.days.length === 0 ? -1 : 1;

  // Sort by repeat days
  return a.days[0] - b.days[0];
}