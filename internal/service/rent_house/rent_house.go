package rent_house

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/programzheng/black-key/config"
	pb "github.com/programzheng/black-key/pkg/proto/rent-house"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RentHouse struct {
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Type holds the value of the "type" field.
	Type int `json:"type,omitempty"`
	// PostID holds the value of the "post_id" field.
	PostID int `json:"post_id,omitempty"`
	// KindName holds the value of the "kind_name" field.
	KindName string `json:"kind_name,omitempty"`
	// RoomStr holds the value of the "room_str" field.
	RoomStr string `json:"room_str,omitempty"`
	// FloorStr holds the value of the "floor_str" field.
	FloorStr string `json:"floor_str,omitempty"`
	// Community holds the value of the "community" field.
	Community string `json:"community,omitempty"`
	// Price holds the value of the "price" field.
	Price int `json:"price,omitempty"`
	// PriceUnit holds the value of the "price_unit" field.
	PriceUnit string `json:"price_unit,omitempty"`
	// PhotoList holds the value of the "photo_list" field.
	PhotoList []string `json:"photo_list,omitempty"`
	// RegionName holds the value of the "region_name" field.
	RegionName string `json:"region_name,omitempty"`
	// SectionName holds the value of the "section_name" field.
	SectionName string `json:"section_name,omitempty"`
	// StreetName holds the value of the "street_name" field.
	StreetName string `json:"street_name,omitempty"`
	// Location holds the value of the "location" field.
	Location string `json:"location,omitempty"`
	// Area holds the value of the "area" field.
	Area string `json:"area,omitempty"`
	// RoleName holds the value of the "role_name" field.
	RoleName string `json:"role_name,omitempty"`
	// Contact holds the value of the "contact" field.
	Contact string `json:"contact,omitempty"`
	// RefreshTime holds the value of the "refresh_time" field.
	RefreshTime string `json:"refresh_time,omitempty"`
	// YesterdayHit holds the value of the "yesterday_hit" field.
	YesterdayHit int `json:"yesterday_hit,omitempty"`
	// IsVip holds the value of the "is_vip" field.
	IsVip int `json:"is_vip,omitempty"`
	// IsCombine holds the value of the "is_combine" field.
	IsCombine int `json:"is_combine,omitempty"`
	// Hurry holds the value of the "hurry" field.
	Hurry int `json:"hurry,omitempty"`
	// IsSocial holds the value of the "is_social" field.
	IsSocial int `json:"is_social,omitempty"`
	// DiscountPriceStr holds the value of the "discount_price_str" field.
	DiscountPriceStr string `json:"discount_price_str,omitempty"`
	// CasesID holds the value of the "cases_id" field.
	CasesID int `json:"cases_id,omitempty"`
	// IsVideo holds the value of the "is_video" field.
	IsVideo int `json:"is_video,omitempty"`
	// Preferred holds the value of the "preferred" field.
	Preferred int `json:"preferred,omitempty"`
	// Cid holds the value of the "cid" field.
	Cid int `json:"cid,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DetailUrl string
}

type GetRentHousesConditions struct {
	Keyword string `json:"keyword,omitempty"`
	City    string `json:"city,omitempty"`
}

func GetGetRentHousesConditionsByJSONString(jsonString string) (*GetRentHousesConditions, error) {
	grhcs := &GetRentHousesConditions{}
	err := json.Unmarshal([]byte(jsonString), &GetRentHousesConditions{})
	if err != nil {
		return nil, err
	}
	return grhcs, nil
}

func ConvertGetRentHousesResponseToRentHouses(grhr *pb.GetRentHousesResponse) ([]*RentHouse, error) {
	rhs := make([]*RentHouse, 0, len(grhr.RentHouses))
	for _, rrh := range grhr.RentHouses {
		rhs = append(rhs, &RentHouse{
			Title:     rrh.Title,
			PhotoList: rrh.PhotoList,
			Price:     int(rrh.Price),
			PriceUnit: rrh.PriceUnit,
			DetailUrl: rrh.DetailUrl,
		})
	}

	return rhs, nil
}

func GetRentHousesByConditionsResponse(conditions *GetRentHousesConditions) (*pb.GetRentHousesResponse, error) {
	grpcRentHouseUrl := config.Cfg.GetString("GRPC_RENT_HOUSE_URL")
	if grpcRentHouseUrl == "" {
		return nil, nil
	}
	conn, err := grpc.Dial(grpcRentHouseUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProxyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	req := &pb.GetRentHousesByConditionsRequest{
		Keyword: conditions.Keyword,
		City:    conditions.City,
	}
	r, err := c.GetRentHousesByConditions(ctx, req)
	if err != nil {
		log.Fatalf("could not get proxy response: %v", err)
	}
	return r, nil
}
