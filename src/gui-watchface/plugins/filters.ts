import Vue from 'vue'
import { Alarm } from '../store/alarms'

declare module 'vue/types/vue' {
  // Vue components
  interface Vue {
    $padTime(value: string): string
    $formatTime(alarm: Alarm): string
    $scrollToTop(): void
  }
}

function padTime(value: number): string {
  if (!value) return '00'
  if (value < 10) return '0' + value.toString()
  return value.toString()
}
Vue.filter('padTime', padTime)
Vue.prototype.$padTime = padTime

function formatTime(alarm: Alarm): string {
  if (!alarm) return '00:00'
  const { hour, minute } = alarm
  return `${padTime(hour)}:${padTime(minute)}`
}
Vue.filter('formatTime', formatTime)
Vue.prototype.$formatTime = formatTime

function scrollToTop() {
  window.scrollTo(0, 0)
}
Vue.prototype.$scrollToTop = scrollToTop
