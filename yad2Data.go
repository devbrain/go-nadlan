package gonadlan

import "time"

type AssetType int

const (
	UnknownAT AssetType = iota
	TotalNewAT
	NewAT
	RenovatedAT
	RenovationNeededAT
	PreservedAT
)

type HomeType int

const (
	UnknownHT HomeType = iota
	GenericHT
	PratiHT
	RoofHT
	TriplexHT
	DualHT
	StudioHT
	GardenHT
	DuplexHT
	FlatHT
	MigrashHT
	YehidaHT
	ParkingHT
	WholeBuildingHT
	KvuzatRehishaHT
	CellarHT
	StorageHT
	TourismHT
	PensionHT
	SubletHT
	ExchangeHT
)

type Yad2Data struct {
	ForSale bool `json:"for_sale"`
	ExtInfo struct {
		Text          string   `json:"text"`
		Title         string   `json:"title"`
		Title1        string   `json:"title1"`
		Title2        string   `json:"title2"`
		Images        []string `json:"images"`
		PrimaryArea   string   `json:"PrimaryArea"`
		PrimaryAreaID int      `json:"PrimaryAreaID"`
		AreaIDText    string   `json:"AreaID_text"`
		SecondaryArea string   `json:"SecondaryArea"`
		AreaID        int      `json:"area_id"`
		Street        string   `json:"street"`
	} `json:"ext_info"`
	CityCode    int `json:"city_code"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	AdNumber  int    `json:"ad_number"`
	ID        string `json:"id"`
	LinkToken string `json:"link_token"`
	Merchant  bool   `json:"merchant"`
	RecordID  int    `json:"record_id"`

	Price       float64    `json:"price"`
	DateAdded   *time.Time `json:"date_added"`
	DateUpdated *time.Time `json:"date"`

	HoodID       int       `json:"hood_id"`
	Neighborhood string    `json:"neighborhood"`
	Asset        AssetType `json:"Asset"`
	Home         HomeType  `json:"Home"`
	Properties   struct {
		Rooms       float32 `json:"rooms"`
		Floor       float32 `json:"floor"`
		SquareMeter float32 `json:"squareMeter"`
	} `json:"properties"`
}
