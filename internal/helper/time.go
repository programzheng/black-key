package helper

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const Iso8601 = "2006-01-02"
const Yyyymmddhhmmss = "2006/01/02 15:04:05"
const Rfc2822 = "Mon Jan 02 15:04:05 -0700 2006"

func GetTimeByTimeString(ts string) (*time.Time, error) {
	dt := fmt.Sprintf("%s %s", GetNowDateTimeByFormat("2006-01-02"), ts)
	pdtl, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.Now().Local().Location())
	if err != nil {
		return nil, err
	}
	return &pdtl, nil
}

func GetWeekDays() []string {
	return []string{
		time.Sunday.String(),
		time.Monday.String(),
		time.Tuesday.String(),
		time.Wednesday.String(),
		time.Thursday.String(),
		time.Friday.String(),
		time.Saturday.String(),
	}
}

func GetWeekDayEnglishByTraditionalChinese(traditionalChinese string) string {
	switch traditionalChinese {
	case "星期日", "禮拜日":
		return GetWeekDays()[0]
	case "星期一", "禮拜一":
		return GetWeekDays()[1]
	case "星期二", "禮拜二":
		return GetWeekDays()[2]

	case "星期三", "禮拜三":
		return GetWeekDays()[3]

	case "星期四", "禮拜四":
		return GetWeekDays()[4]

	case "星期五", "禮拜五":
		return GetWeekDays()[5]

	case "星期六", "禮拜六":
		return GetWeekDays()[6]
	}
	return ""
}

func GetWeekDayTraditionalChineseByEnglish(english string) string {
	switch english {
	case "Sunday":
		return "星期日"
	case "Monday":
		return "星期一"
	case "Tuesday":
		return "星期二"
	case "Wednesday":
		return "星期三"
	case "Thursday":
		return "星期四"
	case "Friday":
		return "星期五"
	case "Saturday":
		return "星期六"
	}
	return ""
}

func GetWeekDayShortTraditionalChineseByEnglish(english string) string {
	switch english {
	case "Sunday":
		return "日"
	case "Monday":
		return "一"
	case "Tuesday":
		return "二"
	case "Wednesday":
		return "三"
	case "Thursday":
		return "四"
	case "Friday":
		return "五"
	case "Saturday":
		return "六"
	}
	return ""
}

func ShortDateIsEveryDay(tc string) bool {
	if tc == "每天" ||
		tc == "每日" ||
		tc == "every" ||
		tc == "every day" ||
		tc == "every-day" {
		return true
	}
	return false
}

func IsShortDateOrTraditionalChineseShortDate(tc string) bool {
	switch tc {
	case "每天", "每日", "every-day", "every day", "every_day", "今天", "今日", "today", "明天", "明日", "tomorrow":
		return true
	}
	return false
}

func GetDateTimeByTraditionalChinese[T string | time.Time](t T) (time.Time, error) {
	switch value := any(t).(type) {
	case string:
		parseDate := strings.Split(value, " ")
		shortTc := parseDate[0]
		dateTime, err := GetTimeByTimeString(parseDate[1])
		if err != nil {
			return time.Time{}, err
		}
		switch shortTc {
		case "每天", "每日":
			//check if the specified date time is greater than the current date time
			if CurrentDateTimeIsGreaterThanSpecifiedDateTime(dateTime) {
				return (*dateTime).AddDate(0, 0, 1), nil
			} else {
				return *dateTime, nil
			}
		case "今天", "今日":
			return *dateTime, nil
		case "明天", "明日":
			return (*dateTime).AddDate(0, 0, 1), nil
		case "昨天", "昨日":
			return (*dateTime).AddDate(0, 0, -1), nil
		}
	case time.Time:
		return value, nil
	}
	return time.Time{}, fmt.Errorf(
		"GetDateTimeByTraditionalChinese: %v  does not conform", t,
	)
}

func CurrentDateTimeIsGreaterThanSpecifiedDateTime(t *time.Time) bool {
	currentTime := time.Now()
	fmt.Printf("%v", t)
	return currentTime.After(*t)
}

func IsDateTime(s string) bool {
	_, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Now().Local().Location())
	return err == nil
}

func GetNowDateTimeByFormat(format string) string {
	if format == "" {
		format = Yyyymmddhhmmss
	}
	t := time.Now()
	return t.Format(format)
}

func GetSpecifyNextWeekDayDateTime(wdt time.Time) time.Time {
	snwddt := wdt.Weekday()
	nextWeekDay := time.Now()
	for nextWeekDay.Weekday() == snwddt {
		nextWeekDay = nextWeekDay.AddDate(0, 0, 1)
	}

	return nextWeekDay
}

func GetSpecifyNextWeekDayDateTimeByString(wdts string) time.Time {
	nextWeekDay := time.Now()
	for nextWeekDay.Weekday().String() == wdts {
		nextWeekDay = nextWeekDay.AddDate(0, 0, 1)
	}

	return nextWeekDay
}

func GetNextWeekDayDateTime(wdt time.Time) time.Time {
	diff := int(time.Now().Weekday()) - int(wdt.Weekday())
	if diff <= 0 {
		diff += 7
	}
	return wdt.AddDate(0, 0, diff)
}

func CalcTimeRange(fromDate string, toDate string) int64 {
	fromDateUnix := toUnix(fromDate)
	toDateUnix := toUnix(toDate)
	return toDateUnix - fromDateUnix
}

func toUnix(date string) int64 {
	t, err := time.Parse(Yyyymmddhhmmss, date)
	if err != nil {
		log.Println(err)
	}
	return t.Unix()
}
