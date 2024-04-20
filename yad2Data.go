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

type Yad2AdditionalData struct {
	AdNumber          int     `json:"ad_number"`
	TotalFloor        int     `json:"total_floor"`
	AssetExclusive    bool    `json:"asset_exclusive"`
	AirConditioner    bool    `json:"air_conditioner"`
	Bars              bool    `json:"bars"`
	Boiler            bool    `json:"boiler"`
	Elevator          bool    `json:"elevator"`
	Accessibility     bool    `json:"accessibility"`
	Renovated         bool    `json:"renovated"`
	Shelter           bool    `json:"shelter"`
	Warehouse         bool    `json:"warehouse"`
	Pets              bool    `json:"pets"`
	RavBariach        bool    `json:"ravBariach"`
	Tornado           bool    `json:"tornado"`
	Furniture         bool    `json:"furniture"`
	FlexibleEnterDate bool    `json:"flexibleEnterDate"`
	LongTerm          bool    `json:"longTerm"`
	Balconies         float32 `json:"balconies"`
	GardenArea        float32 `json:"gardenArea"`
	Parking           float32 `json:"parking"`
}
