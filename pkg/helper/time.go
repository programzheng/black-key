package helper

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const Iso8601 = "2006-01-02"
const Yyyymmddhhmmss = "2006/01/02 15:04:05"
const Rfc2822 = "Mon Jan 02 15:04:05 -0700 2006"

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

func GetWeekDayByTraditionalChinese(tC string) string {
	switch tC {
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

func IsDateTime(s string) bool {
	_, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Now().Local().Location())
	return err != nil
}

func GetNowDateTimeByFormat(format string) string {
	if format == "" {
		format = Yyyymmddhhmmss
	}
	t := time.Now()
	return t.Format(format)
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
