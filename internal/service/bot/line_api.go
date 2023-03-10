package bot

import (
	log "github.com/sirupsen/logrus"
)

type LineMember struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Language      string `json:"language"`
}

func GetGroupMemberCount(groupID string) int {
	groupMemberCount, err := BotClient.GetGroupMemberCount(groupID).Do()
	if err != nil {
		log.Fatal("line messaging api get group member count error:", err)
	}
	return groupMemberCount.Count
}

func GetGroupMemberProfile(groupID string, userID string) (*LineMember, error) {
	userProfileResponse, err := BotClient.GetGroupMemberProfile(groupID, userID).Do()
	if err != nil {
		log.Errorf("line messaging api get group member profile error:", err)
		return nil, err
	}
	lineMember := LineMember(*userProfileResponse)
	return &lineMember, nil
}
