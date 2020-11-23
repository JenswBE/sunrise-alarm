use chrono::prelude::*;

use sunrise_common::general::Alarm;

/// Calculates the next datetime the alarm should trigger
pub fn calculate_next_datetime(alarm: &Alarm, after: &DateTime<Local>) -> DateTime<Local> {
    // Create possible next alarm date
    let alarm_dt = after
        .with_hour(alarm.hour as u32)
        .unwrap()
        .with_minute(alarm.minute as u32)
        .unwrap()
        .with_second(0)
        .unwrap()
        .with_nanosecond(0)
        .unwrap();

    // Check if alarm is exactly at date
    let weekday_number = after.weekday().num_days_from_monday() as u8;
    if alarm.days.len() > 0 || alarm.days.contains(&weekday_number) {
        if *after > alarm_dt {
            // Alarm is at date and still to come
            return alarm_dt;
        }
    }

    // Next alarm is not on date => Calculate next weekday
    let next_day = calculate_next_day(alarm, after.weekday());

    // Shift currently calculated alarm time until the next occurence of this weekday
    let mut alarm_date = alarm_dt.date();
    loop {
        alarm_date = alarm_date.succ();
        if alarm_date.weekday() == next_day {
            break;
        }
    }
    return alarm_date.and_time(alarm_dt.time()).unwrap();
}

/// Calculates the next day the alarm should trigger
fn calculate_next_day(alarm: &Alarm, current_day: Weekday) -> Weekday {
    if alarm.days.len() == 0 {
        // No days set => next day is tomorrow
        return current_day.succ();
    }

    // Get next day => First day higher then current day. Else first day in list.
    let next_day = alarm
        .days
        .iter()
        .find(|&x| *x as u32 > current_day.num_days_from_monday());
    let next_day = next_day.unwrap_or(&alarm.days[0]);
    return weekday_from_u8(next_day).unwrap();
}

fn weekday_from_u8(day_number: &u8) -> Option<Weekday> {
    match day_number {
        0 => Some(Weekday::Mon),
        1 => Some(Weekday::Tue),
        2 => Some(Weekday::Wed),
        3 => Some(Weekday::Thu),
        4 => Some(Weekday::Fri),
        5 => Some(Weekday::Sat),
        6 => Some(Weekday::Sun),
        _ => None,
    }
}
