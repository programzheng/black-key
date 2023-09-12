package line

import (
	"fmt"
	"strings"
	"time"

	model "github.com/programzheng/black-key/internal/model/bot"
	"golang.org/x/exp/slices"
)

func checkCanPushLineFeatureNotification(lfn *model.LineFeatureNotification) bool {
	pushDateTime := lfn.PushDateTime

	if lfn.PushCycle != "specify" {
		nowWeekDay := time.Now().Weekday()
		pcs := strings.Split(lfn.PushCycle, ",")
		if slices.Contains(pcs, nowWeekDay.String()) {
			nowDate := time.Now().Format("2006-01-02")
			st := pushDateTime.Format("15:04:05")
			pdts := fmt.Sprintf("%s %s", nowDate, st)
			pushDateTime, err := time.ParseInLocation("2006-01-02 15:04:05", pdts, time.Now().Local().Location())
			if err != nil {
				return false
			}
			minTolerantDateTime := time.Now().Add(-30 * time.Second)
			maxTolerantDateTime := time.Now().Add(30 * time.Second)
			return minTolerantDateTime.Before(pushDateTime) && maxTolerantDateTime.After(pushDateTime)
		}
	}
	minTolerantDateTime := time.Now().Add(-30 * time.Second)
	maxTolerantDateTime := time.Now().Add(30 * time.Second)
	return minTolerantDateTime.Before(pushDateTime) && maxTolerantDateTime.After(pushDateTime)
}
