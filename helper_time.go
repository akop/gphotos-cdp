package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// compiled year regex
var yearRegex = regexp.MustCompile(`\d{4}`)
var dayRegex = regexp.MustCompile(`\d{1,2}`)
var timeRegex = regexp.MustCompile(`(\d{1,2}):(\d\d)(?::\d\d)?.?([aApP][Mm])?$`)
var timeZoneRegex = regexp.MustCompile(`GMT([-+])?(\d{1,2})(?::(\d\d))?`)

func parseDate(dateStr, timeStr, tzStr string) (time.Time, error) {
	var year, month, day, hour, minute int
	yearStr := yearRegex.FindString(dateStr)
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
		dateStr = strings.Replace(dateStr, yearStr, "", 1)
	} else {
		year = time.Now().Year()
	}
	log.Trace().Msgf("parsed year: %d, dateStr: %s", year, dateStr)

	dayStr := dayRegex.FindString(dateStr)
	if dayStr != "" {
		day, _ = strconv.Atoi(dayStr)
	}
	dateStr = strings.Replace(dateStr, dayStr, "", 1)

	for i, v := range loc.ShortMonthNames {
		if strings.Contains(strings.ToUpper(dateStr), strings.ToUpper(v)) {
			month = i + 1
			break
		}
	}
	if month == 0 {
		return time.Time{}, fmt.Errorf("could not find month in string %s", dateStr)
	}
	log.Trace().Msgf("parsed month: %d, dateStr: %s", month, dateStr)

	if timeStr != "" {
		timeMatch := timeRegex.FindStringSubmatch(timeStr)
		if timeMatch == nil {
			return time.Time{}, fmt.Errorf("could not find time in string %s", timeStr)
		}
		hour, _ = strconv.Atoi(timeMatch[1])
		minute, _ = strconv.Atoi(timeMatch[2])
		if strings.EqualFold(timeMatch[3], "pm") && hour < 12 {
			hour += 12
		}
		if strings.EqualFold(timeMatch[3], "am") && hour == 12 {
			hour = 0
		}
	}

	// read time zone from timezoneStr in format GMT-05:00
	var timeZone *time.Location
	timeZoneStr := strings.Trim(tzStr, " ")
	if timeZoneStr != "" {
		timeZoneMatch := timeZoneRegex.FindStringSubmatch(timeZoneStr)
		if timeZoneMatch != nil {
			tzHour, _ := strconv.Atoi(timeZoneMatch[2])
			tzMinute, _ := strconv.Atoi(timeZoneMatch[3])
			offset := tzHour*60*60 + tzMinute*60
			if timeZoneMatch[1] == "-" {
				offset = -offset
			}
			timeZone = time.FixedZone("", offset)
		} else {
			return time.Time{}, fmt.Errorf("could not parse time zone in string %s", timeZoneStr)
		}
	} else {
		timeZone = time.Local
	}
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, timeZone), nil
}
