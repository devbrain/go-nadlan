package gonadlan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	SetStandardHeaders(req, "https://yad2.co.il")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	body, err := ReadHTTPResponse(res, err)
	if err != nil {
		return nil, 0, err
	}

	return body, res.StatusCode, nil
}

func GetYad2Data(page int, city int, forSale bool) (*Yad2RawData, error) {
	body, statusCode, err := FetchYad2Page(page, city, forSale)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("bad status code %d", statusCode)
	}
	var yad2Data Yad2RawData
	err = json.Unmarshal(body, &yad2Data)
	if err != nil {
		return nil, err
	}
	return &yad2Data, nil
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
				Row5       any    `json:"row_5,omitempty"`
				SearchText string `json:"search_text,omitempty"`
				Title1     string `json:"title_1,omitempty"`
				Title2     string `json:"title_2,omitempty"`
				Images     struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
				} `json:"images,omitempty"`
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
				Street        string   `json:"street,omitempty"`
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
				ContactName                 string `json:"contact_name,omitempty"`
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
				AddressMore                 string `json:"address_more,omitempty"`
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
				Images0                     struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
				} `json:"images,omitempty"`
				Images1 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
					Image11 struct {
						Src string `json:"src"`
					} `json:"Image11"`
					Image12 struct {
						Src string `json:"src"`
					} `json:"Image12"`
					Image13 struct {
						Src string `json:"src"`
					} `json:"Image13"`
					Image14 struct {
						Src string `json:"src"`
					} `json:"Image14"`
				} `json:"images,omitempty"`
				Images2 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
				} `json:"images,omitempty"`
				Images3 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
				} `json:"images,omitempty"`
				Images4 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
				} `json:"images,omitempty"`
				Images5 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
				} `json:"images,omitempty"`
				Images6 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
				} `json:"images,omitempty"`
				Images7 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
				} `json:"images,omitempty"`
				Images8 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
				} `json:"images,omitempty"`
				Images9 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
				} `json:"images,omitempty"`
				Images10 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
				} `json:"images,omitempty"`
				Images11 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
				} `json:"images,omitempty"`
				Images12 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
				} `json:"images,omitempty"`
				Images13 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
				} `json:"images,omitempty"`
				Images14 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
				} `json:"images,omitempty"`
				Title    string `json:"title,omitempty"`
				Images15 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
				} `json:"images,omitempty"`
				Images16 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
				} `json:"images,omitempty"`
				Images17 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
				} `json:"images,omitempty"`
				Images18 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
				} `json:"images,omitempty"`
				Images19 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
					Image11 struct {
						Src string `json:"src"`
					} `json:"Image11"`
				} `json:"images,omitempty"`
				Images20 struct {
				} `json:"images,omitempty"`
				Images21 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
				} `json:"images,omitempty"`
				Images22 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
				} `json:"images,omitempty"`
				Images23 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
				} `json:"images,omitempty"`
				Images24 struct {
				} `json:"images,omitempty"`
				Images25 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
				} `json:"images,omitempty"`
				Images26 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
					Image11 struct {
						Src string `json:"src"`
					} `json:"Image11"`
					Image12 struct {
						Src string `json:"src"`
					} `json:"Image12"`
				} `json:"images,omitempty"`
				Images27 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
				} `json:"images,omitempty"`
				Images28 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
				} `json:"images,omitempty"`
				Images29 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
				} `json:"images,omitempty"`
				Images30 struct {
					Image1 struct {
						Src string `json:"src"`
					} `json:"Image1"`
					Image2 struct {
						Src string `json:"src"`
					} `json:"Image2"`
					Image3 struct {
						Src string `json:"src"`
					} `json:"Image3"`
					Image4 struct {
						Src string `json:"src"`
					} `json:"Image4"`
					Image5 struct {
						Src string `json:"src"`
					} `json:"Image5"`
					Image6 struct {
						Src string `json:"src"`
					} `json:"Image6"`
					Image7 struct {
						Src string `json:"src"`
					} `json:"Image7"`
					Image8 struct {
						Src string `json:"src"`
					} `json:"Image8"`
					Image9 struct {
						Src string `json:"src"`
					} `json:"Image9"`
					Image10 struct {
						Src string `json:"src"`
					} `json:"Image10"`
					Image11 struct {
						Src string `json:"src"`
					} `json:"Image11"`
					Image12 struct {
						Src string `json:"src"`
					} `json:"Image12"`
					Image13 struct {
						Src string `json:"src"`
					} `json:"Image13"`
					Image14 struct {
						Src string `json:"src"`
					} `json:"Image14"`
					Image15 struct {
						Src string `json:"src"`
					} `json:"Image15"`
				} `json:"images,omitempty"`
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
