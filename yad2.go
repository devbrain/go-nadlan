package gonadlan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func FetchYad2Page(page int, city int, forSale bool) ([]byte, int, error) {
	tp := "forsale"
	if !forSale {
		tp = "rent"
	}
	yad2Url, err := url.Parse("https://gw.yad2.co.il/feed-search-legacy/realestate/" + tp)
	if err != nil {
		return nil, 0, err
	}
	params := yad2Url.Query()
	params.Add("city", strconv.Itoa(city))
	params.Add("page", strconv.Itoa(page))
	params.Add("forceLoad", "true")
	yad2Url.RawQuery = params.Encode()
	targetUrl := yad2Url.String()
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		return nil, 0, err
	}
	SetStandardHeaders(req, "gw.yad2.co.il")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	body, err := ReadHTTPResponse(res, err)
	if err != nil {
		return nil, 0, err
	}
	//f, err := os.Create(fmt.Sprintf("%d-%d-%d-%v.json", city, page, res.StatusCode, forSale))
	//defer f.Close()
	//f.Write(body)

	return body, res.StatusCode, nil
}

func GetYad2Data(page int, city int, forSale bool) ([]Yad2Data, int, error) {
	body, statusCode, err := FetchYad2Page(page, city, forSale)
	if err != nil {
		return nil, 0, err
	}
	if statusCode != 200 {
		return nil, 0, fmt.Errorf("bad status code %d", statusCode)
	}
	var yad2Data Yad2RawData
	err = json.Unmarshal(body, &yad2Data)
	if err != nil {
		return nil, 0, err
	}
	data, lastPage := ParseYad2RawData(&yad2Data, forSale)
	return data, lastPage, nil
}

func anyToInt(x any) (int, error) {
	switch v := x.(type) {
	case int:
		return x.(int), nil
	case float64:
		return int(x.(float64)), nil
	case float32:
		return int(x.(float32)), nil
	default:
		return 0, fmt.Errorf("anyToInt dont know how to convert [%v]", v)
	}
}

func anyToFloat(x any) (float32, error) {
	switch v := x.(type) {
	case int:
		return float32(x.(int)), nil
	case float64:
		return float32(x.(float64)), nil
	case float32:
		return x.(float32), nil
	case string:
		strVal := x.(string)
		if len(strVal) == 0 {
			return -1000.0, nil
		}
		if strVal == "קרקע" {
			return 0, nil
		} else if strVal == "לא צוין" {
			return -1000.0, nil
		} else {
			intVal, err := strconv.Atoi(strVal)
			if err != nil {
				return -1000.0, fmt.Errorf("anyToFloat dont know how to convert [%v]", v)
			}
			return float32(intVal), nil
		}
	default:
		return -1000.0, fmt.Errorf("anyToFloat dont know how to convert [%v]", v)
	}
}

func anyToString(x any) string {
	switch x.(type) {
	case int:
	case float64:
	case float32:
		return fmt.Sprintf("%v", x)
	case string:
		return x.(string)
	default:
		return ""
	}
	return ""
}

func parsePriceStr(price string, currency string) (float64, error) {
	f := strings.Trim(strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(price, currency, ""),
			",", ""),
		".", ""), " ")
	atoi, err := strconv.Atoi(f)
	if err == nil {
		return float64(atoi), nil
	}
	return 0, nil
}

func parsePrice(price any, currency string) (float64, error) {
	switch v := price.(type) {
	case int:
		return float64(price.(int)), nil
	case float64:
		return price.(float64), nil
	case float32:
		return float64(price.(float32)), nil
	case int64:
		return float64(price.(int64)), nil
	case string:
		return parsePriceStr(price.(string), currency)
	default:
		return 0, fmt.Errorf("parsePrice dont know how to parse [%v]", v)
	}
}

func minDates(d1 *time.Time, d2 *time.Time) *time.Time {
	if d1 != nil && d2 != nil {
		if d1.Before(*d2) {
			return d1
		}
		return d2
	}
	if d1 != nil {
		return d1
	}
	if d2 != nil {
		return d2
	}
	return nil
}

func maxDates(d1 *time.Time, d2 *time.Time) *time.Time {
	if d1 != nil && d2 != nil {
		if d1.Before(*d2) {
			return d2
		}
		return d1
	}
	if d1 != nil {
		return d1
	}
	if d2 != nil {
		return d2
	}
	return nil
}

func parseDates(date string, dateAdded string, dateUpdated string) (*time.Time, *time.Time) {
	d, errD := time.Parse(time.DateTime, date)
	dA, errA := time.Parse(time.DateTime, dateAdded)
	dU, errU := time.Parse(time.DateTime, dateUpdated)

	var dPtr *time.Time = nil
	var dAPtr *time.Time = nil
	var dUPtr *time.Time = nil

	if errD == nil {
		dPtr = &d
	}
	if errA == nil {
		dAPtr = &dA
	}
	if errU == nil {
		dUPtr = &dU
	}

	return minDates(dPtr, minDates(dAPtr, dUPtr)), maxDates(dPtr, maxDates(dAPtr, dUPtr))
}

func parseAssetType(at string) AssetType {
	if len(at) == 0 {
		return UnknownAT
	}
	if strings.Index(at, "חדש מקבלן") != -1 {
		return TotalNewAT
	}
	if strings.Index(at, "גרו בנכס") != -1 {
		return NewAT
	}
	if strings.Index(at, "דרוש שיפוץ") != -1 {
		return RenovationNeededAT
	}
	if strings.Index(at, "במצב שמור") != -1 {
		return PreservedAT
	}
	if strings.Index(at, "משופץ") != -1 {
		return RenovatedAT
	}
	return UnknownAT
}

func parseHomeType(ht string) HomeType {
	if len(ht) == 0 {
		return UnknownHT
	}
	if strings.Index(ht, "כללי") != -1 {
		return GenericHT
	}
	if strings.Index(ht, "פרטי") != -1 {
		return PratiHT
	}
	if strings.Index(ht, "גג") != -1 {
		return RoofHT
	}
	if strings.Index(ht, "טריפלקס") != -1 {
		return TriplexHT
	}
	if strings.Index(ht, "משפחתי") != -1 {
		return DualHT
	}
	if strings.Index(ht, "סטודיו") != -1 {
		return StudioHT
	}
	if strings.Index(ht, "דירת גן") != -1 {
		return GardenHT
	}
	if strings.Index(ht, "דופלקס") != -1 {
		return DuplexHT
	}
	if strings.Index(ht, "מגרשים") != -1 {
		return MigrashHT
	}
	if strings.Index(ht, "יחידת דיור") != -1 {
		return YehidaHT
	}
	if strings.Index(ht, "חניה") != -1 {
		return ParkingHT
	}
	if strings.Index(ht, "דירה") != -1 {
		return FlatHT
	}
	if strings.Index(ht, "בניין") != -1 {
		return WholeBuildingHT
	}
	if strings.Index(ht, "קב' רכישה") != -1 {
		return KvuzatRehishaHT
	}
	if strings.Index(ht, "מרתף") != -1 {
		return CellarHT
	}
	if strings.Index(ht, "מחסן") != -1 {
		return StorageHT
	}
	if strings.Index(ht, "תיירות") != -1 {
		return TourismHT
	}
	if strings.Index(ht, "דיור מוגן") != -1 {
		return PensionHT
	}
	if strings.Index(ht, "סאבלט") != -1 {
		return SubletHT
	}
	if strings.Index(ht, "החלפת") != -1 {
		return ExchangeHT
	}
	return UnknownHT
}

func ParseYad2RawData(rawData *Yad2RawData, forSale bool) ([]Yad2Data, int) {
	out := make([]Yad2Data, 0)
	for _, f := range rawData.Data.Feed.FeedItems {
		if f.Type != "ad" {
			continue
		}
		var x Yad2Data
		x.ForSale = forSale
		var errPrice error
		x.Price, errPrice = parsePrice(f.Price, f.CurrencyText)
		if errPrice != nil {
			continue
		}
		x.Properties.Floor = -1000.0
		x.Properties.Rooms = -1000.0
		x.Properties.SquareMeter = -1000.0

		for _, obj := range f.Row4 {
			val, e := anyToFloat(obj.Value)
			if e != nil {
				continue
			}
			if obj.Key == "rooms" {
				x.Properties.Rooms = val
			} else if obj.Key == "floor" {
				x.Properties.Floor = val
			} else if obj.Key == "SquareMeter" {
				x.Properties.SquareMeter = val
			}
		}

		if x.Properties.Floor == -1000.0 || x.Properties.Rooms == -1000.0 || x.Properties.SquareMeter == -1000.0 {
			continue
		}
		var parseError error
		x.ExtInfo.Text = f.SearchText
		x.ExtInfo.Images = f.ImagesUrls
		x.ExtInfo.PrimaryArea = f.PrimaryArea
		x.ExtInfo.PrimaryAreaID = f.PrimaryAreaID
		x.ExtInfo.AreaIDText = f.AreaIDText
		x.ExtInfo.SecondaryArea = f.SecondaryArea
		x.ExtInfo.AreaID = f.AreaID
		x.ExtInfo.Title = f.Title
		x.ExtInfo.Title1 = f.Title1
		x.ExtInfo.Title2 = f.Title2
		x.CityCode, parseError = anyToInt(f.CityCode)
		if parseError != nil {
			continue
		}

		x.ExtInfo.Street = anyToString(f.Street)
		x.Coordinates.Latitude = f.Coordinates.Latitude
		x.Coordinates.Longitude = f.Coordinates.Longitude

		x.AdNumber = f.AdNumber
		x.ID = f.ID
		x.LinkToken = f.LinkToken
		x.RecordID = f.RecordID
		x.Merchant = f.Merchant

		x.HoodID = f.HoodID
		x.Neighborhood = f.Neighborhood
		x.Asset = parseAssetType(f.AssetClassificationIDText)
		x.Home = parseHomeType(f.HomeTypeIDText)
		x.DateAdded, x.DateUpdated = parseDates(f.Date, f.DateAdded, f.UpdatedAt)

		out = append(out, x)
	}
	return out, rawData.Data.Pagination.LastPage
}

type Yad2RawData struct {
	Data struct {
		Feed struct {
			CatID      int    `json:"cat_id"`
			SubcatID   int    `json:"subcat_id"`
			TitleText  string `json:"title_text"`
			SortValues []struct {
				Title    string `json:"title"`
				Value    int    `json:"value"`
				Selected int    `json:"selected"`
			} `json:"sort_values"`
			FeedItems []struct {
				Line1 string   `json:"line_1,omitempty"`
				Line2 string   `json:"line_2,omitempty"`
				Line3 string   `json:"line_3,omitempty"`
				Row1  string   `json:"row_1,omitempty"`
				Row2  string   `json:"row_2,omitempty"`
				Row3  []string `json:"row_3,omitempty"`
				Row4  []struct {
					Key   string `json:"key"`
					Label string `json:"label"`
					Value any    `json:"value"`
				} `json:"row_4,omitempty"`
				Row5          any      `json:"row_5,omitempty"`
				SearchText    string   `json:"search_text,omitempty"`
				Title1        string   `json:"title_1,omitempty"`
				Title2        string   `json:"title_2,omitempty"`
				ImagesCount   int      `json:"images_count,omitempty"`
				ImgURL        string   `json:"img_url,omitempty"`
				ImagesUrls    []string `json:"images_urls,omitempty"`
				Mp4VideoURL   any      `json:"mp4_video_url,omitempty"`
				VideoURL      any      `json:"video_url,omitempty"`
				PrimaryArea   string   `json:"PrimaryArea,omitempty"`
				PrimaryAreaID int      `json:"PrimaryAreaID,omitempty"`
				AreaIDText    string   `json:"AreaID_text,omitempty"`
				SecondaryArea string   `json:"SecondaryArea,omitempty"`
				AreaID        int      `json:"area_id,omitempty"`
				City          string   `json:"city,omitempty"`
				CityCode      any      `json:"city_code,omitempty"`
				Street        any      `json:"street,omitempty"`
				Coordinates   struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"coordinates,omitempty"`
				Geohash                     string `json:"geohash,omitempty"`
				AdHighlightType             string `json:"ad_highlight_type,omitempty"`
				BackgroundColor             string `json:"background_color,omitempty"`
				HighlightText               string `json:"highlight_text,omitempty"`
				OrderTypeID                 int    `json:"order_type_id,omitempty"`
				AdNumber                    int    `json:"ad_number,omitempty"`
				CatID                       int    `json:"cat_id,omitempty"`
				CustomerID                  int    `json:"customer_id,omitempty"`
				FeedSource                  string `json:"feed_source,omitempty"`
				ID                          string `json:"id,omitempty"`
				LinkToken                   string `json:"link_token,omitempty"`
				Merchant                    bool   `json:"merchant,omitempty"`
				ContactName                 any    `json:"contact_name,omitempty"`
				MerchantName                string `json:"merchant_name,omitempty"`
				RecordID                    int    `json:"record_id,omitempty"`
				SubcatID                    string `json:"subcat_id,omitempty"`
				Currency                    string `json:"currency,omitempty"`
				CurrencyText                string `json:"currency_text,omitempty"`
				Price                       any    `json:"price,omitempty"`
				DealInfo                    any    `json:"deal_info,omitempty"`
				Date                        string `json:"date,omitempty"`
				DateAdded                   string `json:"date_added,omitempty"`
				UpdatedAt                   string `json:"updated_at,omitempty"`
				IsVisibleForReco            bool   `json:"IsVisibleForReco,omitempty"`
				AdType                      string `json:"ad_type,omitempty"`
				CanChangeLayout             int    `json:"can_change_layout,omitempty"`
				CanHide                     int    `json:"can_hide,omitempty"`
				DefaultLayout               string `json:"default_layout,omitempty"`
				External                    []any  `json:"external,omitempty"`
				IsHidden                    int    `json:"is_hidden,omitempty"`
				IsLiked                     int    `json:"is_liked,omitempty"`
				IsTradeInButton             bool   `json:"is_trade_in_button,omitempty"`
				LikeCount                   int    `json:"like_count,omitempty"`
				Line1TextColor              string `json:"line_1_text_color,omitempty"`
				Line2TextColor              string `json:"line_2_text_color,omitempty"`
				PromotionalAd               int    `json:"promotional_ad,omitempty"`
				RemoveOnUnlike              bool   `json:"remove_on_unlike,omitempty"`
				Type                        string `json:"type"`
				UID                         any    `json:"uid,omitempty"`
				AddressMore                 any    `json:"address_more,omitempty"`
				BrokerAvatar                string `json:"broker_avatar,omitempty"`
				HoodID                      int    `json:"hood_id,omitempty"`
				OfficeAbout                 string `json:"office_about,omitempty"`
				OfficeLogoURL               string `json:"office_logo_url,omitempty"`
				SquareMeters                int    `json:"square_meters,omitempty"`
				HomeTypeIDText              string `json:"HomeTypeID_text,omitempty"`
				Neighborhood                string `json:"neighborhood,omitempty"`
				AssetClassificationIDText   string `json:"AssetClassificationID_text,omitempty"`
				RoomsText                   any    `json:"Rooms_text,omitempty"`
				IsPrivateCommercialMinisite bool   `json:"is_private_commercial_minisite,omitempty"`
				AbovePrice                  string `json:"abovePrice,omitempty"`
				Priority                    int    `json:"priority,omitempty"`
				BackgroundType              int    `json:"background_type,omitempty"`
				IsPlatinum                  bool   `json:"is_platinum,omitempty"`
				IsMobilePlatinum            bool   `json:"is_mobile_platinum,omitempty"`
				Title                       string `json:"title,omitempty"`
			} `json:"feed_items"`
			CurrentPage  int    `json:"current_page"`
			PageSize     int    `json:"page_size"`
			TotalItems   int    `json:"total_items"`
			TotalPages   int    `json:"total_pages"`
			BreadCrumbs  []any  `json:"breadCrumbs"`
			Canonical    string `json:"canonical"`
			LeftColumn   []any  `json:"left_column"`
			SearchParams struct {
				City string `json:"city"`
			} `json:"search_params"`
			SeoParams struct {
				TotalAdCount int `json:"totalAdCount"`
			} `json:"seo_params"`
			AssociatedLinks []struct {
				TitleText string `json:"title_text"`
				URL       string `json:"url"`
			} `json:"associated_links"`
			HeaderText  string `json:"header_text"`
			FeedLiteral struct {
				City []struct {
					Title string `json:"title"`
					ID    string `json:"id"`
				} `json:"city"`
			} `json:"feedLiteral"`
			NhoodKingPackage []struct {
				CustID        int    `json:"CustID"`
				Total         int    `json:"total"`
				OfficeName    string `json:"office_name"`
				OfficeLogoURL string `json:"office_logo_url"`
				AgencyURL     struct {
					City             string `json:"city"`
					DealerID         string `json:"dealerID"`
					RedirectPathOnly string `json:"redirect_path_only"`
				} `json:"agency_url"`
			} `json:"nhood_king_package"`
			ThreeInFeedPackage []struct {
				CustID        int    `json:"CustID"`
				Total         int    `json:"total"`
				OfficeName    string `json:"office_name"`
				OfficeLogoURL string `json:"office_logo_url"`
				AgencyURL     struct {
					City             string `json:"city"`
					DealerID         string `json:"dealerID"`
					RedirectPathOnly string `json:"redirect_path_only"`
				} `json:"agency_url"`
			} `json:"three_in_feed_package"`
		} `json:"feed"`
		Title   string `json:"title"`
		Filters []struct {
			Title    string `json:"title"`
			Value    int    `json:"value"`
			Selected int    `json:"selected"`
		} `json:"filters"`
		Pagination struct {
			CurrentPage        int `json:"current_page"`
			ItemsInCurrentPage int `json:"items_in_current_page"`
			LastPage           int `json:"last_page"`
			MaxItemsPerPage    int `json:"max_items_per_page"`
			TotalItems         int `json:"total_items"`
		} `json:"pagination"`
		CatTitle   string `json:"catTitle"`
		LeftColumn bool   `json:"left_column"`
		Address    struct {
			TopArea struct {
				Level string `json:"level"`
				ID    int    `json:"id"`
				Name  string `json:"name"`
			} `json:"topArea"`
			Area struct {
				Level string `json:"level"`
				ID    int    `json:"id"`
				Name  string `json:"name"`
			} `json:"area"`
			City struct {
				Level string `json:"level"`
				ID    string `json:"id"`
				Name  string `json:"name"`
			} `json:"city"`
			Neighborhood struct {
				Level string `json:"level"`
				ID    any    `json:"id"`
				Name  any    `json:"name"`
			} `json:"neighborhood"`
			Street struct {
				Level string `json:"level"`
				ID    any    `json:"id"`
				Name  any    `json:"name"`
			} `json:"street"`
		} `json:"address"`
		Yad1Ads struct {
			TopGallery []struct {
				CityNeighborhood string `json:"CityNeighborhood"`
				Image            string `json:"Image"`
				ProjectName      string `json:"projectName"`
				ProjectID        int    `json:"projectID"`
				PromotionText    string `json:"promotion_text"`
				Neighborhood     string `json:"Neighborhood"`
				SalePic          string `json:"SalePic"`
			} `json:"top_gallery"`
			BottomGallery []struct {
				CityNeighborhood string `json:"CityNeighborhood"`
				Image            string `json:"Image"`
				ProjectName      string `json:"projectName"`
				ProjectID        int    `json:"projectID"`
				PromotionText    string `json:"promotion_text"`
				Neighborhood     string `json:"Neighborhood"`
				SalePic          string `json:"SalePic"`
			} `json:"bottom_gallery"`
		} `json:"yad1Ads"`
		Yad1Listing []any `json:"yad1Listing"`
	} `json:"data"`
	Message string `json:"message"`
}
