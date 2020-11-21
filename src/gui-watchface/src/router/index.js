import Vue from "vue";
import VueRouter from "vue-router";
import Clock from "../views/Clock.vue";
import Alarms from "../views/Alarms.vue";
import Settings from "../views/Settings.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Clock",
    component: Clock,
  },
  {
    path: "/alarms",
    name: "Alarms",
    component: Alarms,
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
