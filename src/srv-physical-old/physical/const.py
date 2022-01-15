"""This module contains constants"""

from datetime import timedelta
from pathlib import Path

# ========================================
# =            ALARM SETTINGS            =
# ========================================
SOUND_INCREASE_DURATION = timedelta(minutes=2)

# ========================================
# =         LOCALIZATION SETTINGS        =
# ========================================
WEEKDAYS = ("Monday", "Tuesday", "Wednesday",
            "Thursday", "Friday", "Saturday", "Sunday")

LOCAL_TIMEZONE = 'Europe/Brussels'

# ========================================
# =        FILE LOCATION SETTINGS        =
# ========================================
DIR_DATA = 'data'
DIR_ALARM_SOUNDS = 'alarm_sounds'
PATH_DEFAULT_SOUND = Path(DIR_DATA) / DIR_ALARM_SOUNDS / 'default.mp3'
DIR_CONFIG = 'config'
PATH_CONFIG_ALARMS = Path(DIR_DATA) / DIR_CONFIG / 'alarms.yml'
