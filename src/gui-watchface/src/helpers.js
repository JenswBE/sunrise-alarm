export default {
  padTime,
  formatTime,
  scrollToTop
};

function padTime(value) {
  if (!value) return "00";
  if (value < 10) return "0" + value.toString();
  return value;
}

function formatTime(alarm) {
  if (!alarm) return "00:00";
  const { hour, minute } = alarm;
  return `${padTime(hour)}:${padTime(minute)}`;
}

function scrollToTop() {
  window.scrollTo(0, 0);
}
