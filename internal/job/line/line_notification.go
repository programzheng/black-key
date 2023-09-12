package line

import model "github.com/programzheng/black-key/internal/model/bot"

type lineNotificationModel interface {
	*model.LineNotification | *model.LineFeatureNotification
}

func getPushID[LNM lineNotificationModel](ln interface{}) string {
	switch v := ln.(type) {
	case *model.LineNotification:
		if v.RoomID != "" {
			return v.RoomID
		}
		if v.GroupID != "" {
			return v.GroupID
		}
		if v.UserID != "" {
			return v.UserID
		}
	case *model.LineFeatureNotification:
		if v.RoomID != "" {
			return v.RoomID
		}
		if v.GroupID != "" {
			return v.GroupID
		}
		if v.UserID != "" {
			return v.UserID
		}
	}

	return ""
}

func afterPushLineNotification[LNM lineNotificationModel](ln interface{}) error {
	var err error

	switch v := ln.(type) {
	case *model.LineNotification:
		//Limit < 0 is unlimited
		if v.Limit > 0 {
			v.Limit -= 1
			err := v.Save()
			if err != nil {
				return err
			}
		}

		// if v.PushCycle != "specify" && v.Limit == -1 {
		// 	//push weekday is next weekday
		// 	nextDateTime := v.PushDateTime.AddDate(0, 0, 1)
		// 	wds := nextDateTime.Weekday().String()
		// 	weekDayCycle := strings.Split(v.PushCycle, ",")
		// 	if slices.Contains(weekDayCycle, wds) {
		// 		v.PushDateTime = nextDateTime
		// 		err := v.Save()
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }

		if v.Limit == 0 {
			err = v.Delete()
			if err != nil {
				return err
			}
		}
	case *model.LineFeatureNotification:
		//Limit < 0 is unlimited
		if v.Limit > 0 {
			v.Limit -= 1
			err := v.Save()
			if err != nil {
				return err
			}
		}

		// if v.PushCycle != "specify" && v.Limit == -1 {
		// 	//push weekday is next weekday
		// 	nextDateTime := v.PushDateTime.AddDate(0, 0, 1)
		// 	wds := nextDateTime.Weekday().String()
		// 	weekDayCycle := strings.Split(v.PushCycle, ",")
		// 	if slices.Contains(weekDayCycle, wds) {
		// 		v.PushDateTime = nextDateTime
		// 		err := v.Save()
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }

		if v.Limit == 0 {
			err = v.Delete()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
