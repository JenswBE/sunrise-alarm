import Vue from "vue";
import VueRouter from "vue-router";
import Clock from "../views/Clock.vue";
import Alarms from "../views/Alarms.vue";

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
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
