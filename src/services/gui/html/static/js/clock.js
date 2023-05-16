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

function updateNextRing() {
  fetch("/api/next-ring-time")
    .then((response) => {
      if (response.status !== 200) {
        throw new Error(response.status);
      }
      return response.json();
    })
    .then((data) => {
      nextAlarm.textContent = generateNextRingTimeText(data.ring_time);
    })
    .catch((error) => {
      if (error.message === "204") {
        nextAlarm.textContent = "";
      } else {
        nextAlarm.textContent = `Unable to fetch next ring time: ${error.message}`;
      }
    });
}

function generateNextRingTimeText(time) {
  // Setup variables
  let day = "";
  const now = luxon.DateTime.local();
  const nextRingDate = luxon.DateTime.fromISO(time);
  const tomorrow = now.plus({ days: 1 });
  const nextWeek = now.plus({ days: 6 }); // 6 days instead of 7 to prevent time-based edge cases.

  // Check if ring is today
  if (nextRingDate.toISODate() == now.toISODate()) {
    day = "Today";
  } else if (nextRingDate.toISODate() == tomorrow.toISODate()) {
    day = "Tomorrow";
  } else if (nextRingDate < nextWeek) {
    day = nextRingDate.toFormat("cccc");
  } else {
    daysDiff = Math.round(nextRingDate.diff(now, "days").as("days"));
    day = nextRingDate.toFormat("cccc") + ` (in ${daysDiff} days)`;
  }

  // Format and set text
  const ring_time = nextRingDate.toFormat("HH:mm");
  return `Next ring time: ${day} at ${ring_time}`;
}

updateClockAndDate();
setInterval(updateClockAndDate, 1000);

updateNextRing();
setInterval(updateNextRing, 5000);
