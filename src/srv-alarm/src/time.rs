use chrono::prelude::*;

use crate::models::AlarmConfig;
use sunrise_common::alarm::{Alarm, NextAction, NextAlarm};

/// Calculates the next alarm for the alarm
pub fn calculate_next_alarm(alarm: &Alarm, config: &AlarmConfig) -> NextAlarm {
    calculate_next_alarm_testable(alarm, &Local::now(), config)
}

/// Calculates the next alarm for the alarm (testable version)
fn calculate_next_alarm_testable(
    alarm: &Alarm,
    now: &DateTime<Local>,
    config: &AlarmConfig,
) -> NextAlarm {
    // Check enabled
    if !alarm.enabled {
        return NextAlarm::default();
    }

    // Calculate next alarm
    let next_alarm = calculate_next_datetime(alarm, now);

    // Check if we need to skip
    if !alarm.skip_next {
        // Shouldn't skip next alarm
        NextAlarm {
            id: alarm.id.clone(),
            alarm_datetime: Some(next_alarm),
            next_action: NextAction::Ring,
            next_action_datetime: Some(next_alarm - config.light_duration),
        }
    } else {
        // Skip next alarm
        NextAlarm {
            id: alarm.id.clone(),
            alarm_datetime: Some(calculate_next_datetime(alarm, &next_alarm)),
            next_action: NextAction::Skip,
            next_action_datetime: Some(next_alarm),
        }
    }
}

/// Calculates the next datetime the alarm should trigger
fn calculate_next_datetime(alarm: &Alarm, after: &DateTime<Local>) -> DateTime<Local> {
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
    if alarm.days.len() == 0 || alarm.days.contains(&weekday_number) {
        if *after < alarm_dt {
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

#[cfg(test)]
mod tests {
    use chrono::Duration;
    use uuid::Uuid;

    use super::*;

    fn alarm(hour: u8, minute: u8) -> Alarm {
        Alarm {
            id: uuid(),
            enabled: true,
            hour,
            minute,
            ..Default::default()
        }
    }

    fn parse_date(date: &str) -> DateTime<Local> {
        Local.datetime_from_str(date, "%Y-%m-%d %H:%M").unwrap()
    }

    fn now() -> DateTime<Local> {
        parse_date("2020-09-19 14:30")
    }

    fn uuid() -> Uuid {
        Uuid::parse_str("51d2e380-3611-4857-8659-98f787858e98").unwrap()
    }

    fn assert_alarm(alarm: Alarm, alarm_dt: &str, action: NextAction, action_dt: &str) {
        let config = AlarmConfig {
            light_duration: Duration::minutes(7),
            ..Default::default()
        };

        let next_alarm = calculate_next_alarm_testable(&alarm, &now(), &config);
        assert_eq!(uuid(), next_alarm.id);
        assert_eq!(action, next_alarm.next_action);
        assert_eq!(Some(parse_date(alarm_dt)), next_alarm.alarm_datetime);
        assert_eq!(Some(parse_date(action_dt)), next_alarm.next_action_datetime);
    }

    #[test]
    fn test_still_today() {
        let alarm = alarm(16, 0);
        assert_alarm(
            alarm,
            "2020-09-19 16:00",
            NextAction::Ring,
            "2020-09-19 15:53",
        )
    }

    #[test]
    fn test_for_tomorrow() {
        let alarm = alarm(14, 0);
        assert_alarm(
            alarm,
            "2020-09-20 14:00",
            NextAction::Ring,
            "2020-09-20 13:53",
        )
    }

    #[test]
    fn test_for_today_but_skipped() {
        let mut alarm = alarm(16, 0);
        alarm.skip_next = true;
        assert_alarm(
            alarm,
            "2020-09-20 16:00",
            NextAction::Skip,
            "2020-09-19 16:00",
        )
    }

    #[test]
    fn test_for_tomorrow_but_skipped() {
        let mut alarm = alarm(14, 0);
        alarm.skip_next = true;
        assert_alarm(
            alarm,
            "2020-09-21 14:00",
            NextAction::Skip,
            "2020-09-20 14:00",
        )
    }

    #[test]
    fn test_repeated_and_for_today() {
        let mut alarm = alarm(16, 0);
        alarm.days = vec![5];
        assert_alarm(
            alarm,
            "2020-09-19 16:00",
            NextAction::Ring,
            "2020-09-19 15:53",
        )
    }

    #[test]
    fn test_repeated_but_not_for_today() {
        let mut alarm = alarm(14, 0);
        alarm.days = vec![5];
        assert_alarm(
            alarm,
            "2020-09-26 14:00",
            NextAction::Ring,
            "2020-09-26 13:53",
        )
    }

    #[test]
    fn test_repeated_for_monday_time_still_to_come() {
        let mut alarm = alarm(16, 0);
        alarm.days = vec![0, 1, 2, 3, 4];
        assert_alarm(
            alarm,
            "2020-09-21 16:00",
            NextAction::Ring,
            "2020-09-21 15:53",
        )
    }

    #[test]
    fn test_repeated_for_monday_time_past() {
        let mut alarm = alarm(14, 0);
        alarm.days = vec![0, 1, 2, 3, 4];
        assert_alarm(
            alarm,
            "2020-09-21 14:00",
            NextAction::Ring,
            "2020-09-21 13:53",
        )
    }

    #[test]
    fn test_repeated_for_today_but_skipped() {
        let mut alarm = alarm(16, 0);
        alarm.days = vec![5];
        alarm.skip_next = true;
        assert_alarm(
            alarm,
            "2020-09-26 16:00",
            NextAction::Skip,
            "2020-09-19 16:00",
        )
    }

    #[test]
    fn test_repeated_but_not_for_today_and_skipped() {
        let mut alarm = alarm(14, 0);
        alarm.days = vec![5];
        alarm.skip_next = true;
        assert_alarm(
            alarm,
            "2020-10-03 14:00",
            NextAction::Skip,
            "2020-09-26 14:00",
        )
    }

    #[test]
    fn test_repeated_for_monday_but_skipped_time_still_to_come() {
        let mut alarm = alarm(16, 0);
        alarm.days = vec![0, 1, 2, 3, 4];
        alarm.skip_next = true;
        assert_alarm(
            alarm,
            "2020-09-22 16:00",
            NextAction::Skip,
            "2020-09-21 16:00",
        )
    }

    #[test]
    fn test_repeated_for_monday_but_skipped_time_past() {
        let mut alarm = alarm(14, 0);
        alarm.days = vec![0, 1, 2, 3, 4];
        alarm.skip_next = true;
        assert_alarm(
            alarm,
            "2020-09-22 14:00",
            NextAction::Skip,
            "2020-09-21 14:00",
        )
    }
}
