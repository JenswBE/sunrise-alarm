let clock = document.getElementById("sunrise-clock");
let date = document.getElementById("sunrise-date");
let nextAlarm = document.getElementById("sunrise-next-alarm");

function updateClockAndDate() {
  // Update clock
  const now = luxon.DateTime.local();
  const format = now.second % 2 == 0 ? "HH:mm" : "HH mm";
  clock.textContent = now.toFormat(format);

  // Update date
  date.textContent = now.setLocale("en-UK").toLocaleString({
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

function updateNextAlarm() {
  fetch("/api/next-alarm-to-ring")
    .then((response) => {
      if (response.status !== 200) {
        throw new Error(response.status);
      }
      return response.json();
    })
    .then((data) => {
      nextAlarm.textContent = generateNextAlarmText(data.alarm_time);
    })
    .catch((error) => {
      if (error.message === "204") {
        nextAlarm.textContent = "";
      } else {
        nextAlarm.textContent = `Unable to fetch next alarm: ${error.message}`;
      }
    });
}

function generateNextAlarmText(time) {
  // Setup variables
  let day = "";
  const now = luxon.DateTime.local();
  const nextAlarmDate = luxon.DateTime.fromISO(time);
  const tomorrow = now.plus({ days: 1 });

  // Check if alarm is today
  if (nextAlarmDate.weekday == now.weekday) {
    day = "Today";
  } else if (nextAlarmDate.weekday == tomorrow.weekday) {
    day = "Tomorrow";
  } else {
    day = nextAlarmDate.toFormat("cccc");
  }

  // Format and set text
  const alarm_time = nextAlarmDate.toFormat("HH:mm");
  return `Next alarm: ${day} at ${alarm_time}`;
}

updateClockAndDate();
setInterval(updateClockAndDate, 1000);

updateNextAlarm();
setInterval(updateNextAlarm, 5000);
