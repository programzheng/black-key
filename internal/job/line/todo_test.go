package line

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertTimeToPushDateTime(t *testing.T) {
	testDateTime := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s|13:00:00",
		time.Sunday.String(),
		time.Monday.String(),
		time.Tuesday.String(),
		time.Wednesday.String(),
		time.Thursday.String(),
		time.Friday.String(),
		time.Saturday.String(),
	)
	pushDateTime := convertTimeToPushDateTime(testDateTime)
	nowDateTime := time.Now().Format("2006-01-02") + " 13:00:00"
	if pushDateTime != nowDateTime {
		t.Errorf("pushDateTime %s, nowDateTime %s", pushDateTime, nowDateTime)
		t.Fail()
		return
	}
	t.Log("success")
}

func TestRunPushLineNotificationSchedule(t *testing.T) {
	RunPushLineNotificationSchedule()
}
