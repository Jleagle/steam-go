package steamapi

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jleagle/unmarshal-go"
)

var (
	ErrAppNotFound      = errors.New("steam: store: app not found")
	ErrPackageNotFound  = errors.New("steam: store: package not found")
	ErrWishlistNotFound = errors.New("steam: store: wishlist not found")
	ErrNullResponse     = errors.New("steam: store: null response") // Probably being rate limited
	ErrHTMLResponse     = errors.New("steam: store: html response") // Probably down
)

func (c Client) GetAppDetails(id uint, cc ProductCC, language LanguageCode, filters []string) (app AppDetails, err error) {

	if id == 0 {
		return app, ErrAppNotFound // App 0 does exist but the API does not return it
	}

	resp, err := c.GetAppDetailsMulti([]uint{id}, cc, language, filters)
	if err != nil {
		return app, err
	}

	idx := strconv.FormatUint(uint64(id), 10)

	if resp[idx].Success == false {
		return app, ErrAppNotFound
	}

	return resp[idx], nil
}

func (c Client) GetAppDetailsMulti(ids []uint, cc ProductCC, language LanguageCode, filters []string) (resp map[string]AppDetails, err error) {

	var stringIDs []string
	for _, id := range ids {
		stringIDs = append(stringIDs, strconv.FormatUint(uint64(id), 10))
	}

	query := url.Values{}
	query.Set("appids", strings.Join(stringIDs, ","))
	query.Set("cc", string(cc))      // Country code (not from enum)
	query.Set("l", string(language)) // Text language
	if filters != nil && len(filters) > 0 {
		query.Set("filters", strings.Join(filters, ","))
	}

	b, err := c.getFromStore("api/appdetails", query)
	if err != nil {
		return resp, err
	}

	var bytesString = string(b)

	// Check invalid responses
	if bytesString == "null" || bytesString == "[]" {
		return resp, ErrNullResponse
	}
	if strings.HasPrefix(strings.TrimSpace(bytesString), "<") {
		return resp, ErrHTMLResponse
	}

	// Fix arrays that should be objects
	bytesString = strings.Replace(bytesString, `{"success":true,"data":[]}`, `{"success":true}`, 1)
	bytesString = strings.Replace(bytesString, `"pc_requirements":[]`, `"pc_requirements":{}`, 1)
	bytesString = strings.Replace(bytesString, `"mac_requirements":[]`, `"mac_requirements":{}`, 1)
	bytesString = strings.Replace(bytesString, `"linux_requirements":[]`, `"linux_requirements":{}`, 1)
	b = []byte(bytesString)

	// Unmarshal JSON
	resp = map[string]AppDetails{}
	err = json.Unmarshal(b, &resp)

	return resp, err
}

type AppDetails struct {
	Success bool `json:"success"`
	Data    *struct {
		Type                string        `json:"type"`
		Name                string        `json:"name"`
		AppID               int           `json:"steam_appid"`
		RequiredAge         unmarshal.Int `json:"required_age"`
		IsFree              bool          `json:"is_free"`
		DLC                 []int         `json:"dlc"`
		ControllerSupport   string        `json:"controller_support"`
		DetailedDescription string        `json:"detailed_description"`
		AboutTheGame        string        `json:"about_the_game"`
		ShortDescription    string        `json:"short_description"`
		Fullgame            struct {
			AppID unmarshal.Int `json:"appid"`
			Name  string        `json:"name"`
		} `json:"fullgame"`
		SupportedLanguages string `json:"supported_languages"`
		Reviews            string `json:"reviews"`
		HeaderImage        string `json:"header_image"`
		Website            string `json:"website"`
		PcRequirements     struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"pc_requirements"`
		MacRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"mac_requirements"`
		LinuxRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"linux_requirements"`
		LegalNotice          string   `json:"legal_notice"`
		ExtUserAccountNotice string   `json:"ext_user_account_notice"`
		DRMNotice            string   `json:"drm_notice"`
		Developers           []string `json:"developers"`
		Publishers           []string `json:"publishers"`
		Demos                []struct {
			AppID       unmarshal.Int `json:"appid"`
			Description string        `json:"description"`
		} `json:"demos"`
		PriceOverview *struct {
			Currency         CurrencyCode `json:"currency"`
			Initial          int          `json:"initial"`
			Final            int          `json:"final"`
			DiscountPercent  int          `json:"discount_percent"`
			InitialFormatted string       `json:"initial_formatted"`
			FinalFormatted   string       `json:"final_formatted"`
			RecurringSub     interface{}  `json:"recurring_sub"` // Either "false" or a sub id int
			RecurringSubDesc string       `json:"recurring_sub_desc"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string           `json:"name"`
			Title                   string           `json:"title"`
			Description             string           `json:"description"`
			SelectionText           string           `json:"selection_text"`
			SaveText                string           `json:"save_text"`
			DisplayType             unmarshal.String `json:"display_type"`
			IsRecurringSubscription unmarshal.Bool   `json:"is_recurring_subscription"`
			Subs                    []struct {
				PackageID                int            `json:"packageid"`
				PercentSavingsText       string         `json:"percent_savings_text"`
				PercentSavings           int            `json:"percent_savings"`
				OptionText               string         `json:"option_text"`
				OptionDescription        string         `json:"option_description"`
				CanGetFreeLicense        unmarshal.Bool `json:"can_get_free_license"`
				IsFreeLicense            bool           `json:"is_free_license"`
				PriceInCentsWithDiscount int            `json:"price_in_cents_with_discount"`
			} `json:"subs"`
		} `json:"package_groups"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int8   `json:"score"`
			URL   string `json:"url"`
		} `json:"metacritic"`
		Categories  AppDetailsCategory `json:"categories"`
		Genres      AppDetailsGenre    `json:"genres"`
		Screenshots []struct {
			ID            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		ContentDescriptors struct {
			IDs   interface{}
			Notes interface{}
		} `json:"content_descriptors"`
	} `json:"data"`
}

type AppDetailsGenre []struct {
	ID          unmarshal.Int `json:"id"`
	Description string        `json:"description"`
}

func (g AppDetailsGenre) IDs() (IDs []int) {
	for _, v := range g {
		IDs = append(IDs, int(v.ID))
	}
	return IDs
}

func (g AppDetailsGenre) Names() (names []string) {
	for _, v := range g {
		names = append(names, v.Description)
	}
	return names
}

type AppDetailsCategory []struct {
	ID          int8   `json:"id"`
	Description string `json:"description"`
}

func (c AppDetailsCategory) IDs() (IDs []int) {
	for _, v := range c {
		IDs = append(IDs, int(v.ID))
	}
	return IDs
}

func (c AppDetailsCategory) Names() (names []string) {
	for _, v := range c {
		names = append(names, v.Description)
	}
	return names
}

func (c Client) GetPackageDetails(id uint, code ProductCC, language LanguageCode) (pack PackageDetailsBody, err error) {

	if id == 0 {
		return pack, ErrPackageNotFound // Package 0 does exist but the API does not return it
	}

	idx := strconv.FormatUint(uint64(id), 10)

	query := url.Values{}
	query.Set("packageids", idx)
	query.Set("cc", string(code))    // Price currency
	query.Set("l", string(language)) // Text

	b, err := c.getFromStore("api/packagedetails", query)
	if err != nil {
		return pack, err
	}

	var bytesString = string(b)

	// Check invalid responses
	if bytesString == "null" {
		return pack, ErrNullResponse
	}
	if strings.HasPrefix(strings.TrimSpace(bytesString), "<") {
		return pack, ErrHTMLResponse
	}

	// Unmarshal JSON
	resp := map[string]PackageDetailsBody{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return pack, err
	}

	if resp[idx].Success == false {
		return pack, ErrPackageNotFound
	}

	return resp[idx], nil
}

type PackageDetailsBody struct {
	Success bool `json:"success"`
	Data    struct {
		Name        string `json:"name"`
		PageImage   string `json:"page_image"`
		HeaderImage string `json:"header_image"`
		SmallLogo   string `json:"small_logo"`
		Apps        []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"apps"`
		Price struct {
			Currency        CurrencyCode `json:"currency"`
			Initial         int          `json:"initial"`
			Final           int          `json:"final"`
			DiscountPercent int          `json:"discount_percent"`
			Individual      int          `json:"individual"`
		} `json:"price"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Controller  map[string]bool `json:"controller"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
	} `json:"data"`
}

func (c Client) GetTags() (tags Tags, err error) {

	b, err := c.getFromStore("tagdata/populartags/english", url.Values{})
	if err != nil {
		return tags, err
	}

	var resp []Tag
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return tags, err
	}

	return Tags{Tags: resp}, nil
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

func (t Tags) GetSlice() (ids []int) {

	for _, v := range t.Tags {
		ids = append(ids, v.TagID)
	}
	return ids
}

func (t Tags) GetMap() (tags map[int]string) {

	tags = map[int]string{}
	for _, v := range t.Tags {
		tags[v.TagID] = v.Name
	}
	return tags
}

type Tag struct {
	TagID int    `json:"tagid"`
	Name  string `json:"name"`
}

func (c Client) GetReviews(appID int, language LanguageCode) (reviews ReviewsResponse, err error) {

	query := url.Values{}
	query.Set("json", "1")
	query.Set("language", string(language))
	query.Set("l", string(language))
	query.Set("filter", "all")
	query.Set("purchase_type", "all")
	query.Set("date_range_type", "all")
	query.Set("review_type", "all")
	query.Set("start_date", "-1")
	query.Set("end_date", "-1")
	query.Set("cursor", "*")

	b, err := c.getFromStore("appreviews/"+strconv.Itoa(appID), query)
	if err != nil {
		return reviews, err
	}

	// Unmarshal JSON
	err = json.Unmarshal(b, &reviews)
	if err != nil {
		return reviews, err
	}

	return reviews, nil
}

type ReviewsResponse struct {
	Success      int `json:"success"`
	QuerySummary struct {
		NumReviews      int     `json:"num_reviews"`
		ReviewScore     float64 `json:"review_score"`
		ReviewScoreDesc string  `json:"review_score_desc"`
		TotalPositive   int     `json:"total_positive"`
		TotalNegative   int     `json:"total_negative"`
		TotalReviews    int     `json:"total_reviews"`
	} `json:"query_summary"`
	Reviews []struct {
		Recommendationid string `json:"recommendationid"`
		Author           struct {
			SteamID              unmarshal.Int64 `json:"steamid"`
			NumGamesOwned        int             `json:"num_games_owned"`
			NumReviews           int             `json:"num_reviews"`
			PlaytimeForever      int             `json:"playtime_forever"`
			PlaytimeLastTwoWeeks int             `json:"playtime_last_two_weeks"`
			LastPlayed           int             `json:"last_played"`
		} `json:"author"`
		Language                 string            `json:"language"`
		Review                   string            `json:"review"`
		TimestampCreated         int64             `json:"timestamp_created"`
		TimestampUpdated         int64             `json:"timestamp_updated"`
		VotedUp                  bool              `json:"voted_up"`
		VotesUp                  int               `json:"votes_up"`
		VotesFunny               int               `json:"votes_funny"`
		WeightedVoteScore        unmarshal.Float64 `json:"weighted_vote_score"`
		CommentCount             int               `json:"comment_count"`
		SteamPurchase            bool              `json:"steam_purchase"`
		ReceivedForFree          bool              `json:"received_for_free"`
		WrittenDuringEarlyAccess bool              `json:"written_during_early_access"`
	} `json:"reviews"`
}

func (r ReviewsResponse) GetPositivePercent() float64 {
	return float64(r.QuerySummary.TotalPositive) / float64(r.QuerySummary.TotalReviews) * 100
}

func (r ReviewsResponse) GetNegativePercent() float64 {
	return float64(r.QuerySummary.TotalNegative) / float64(r.QuerySummary.TotalReviews) * 100
}

func (c Client) GetWishlist(playerID int64) (wishlist Wishlist, err error) {

	query := url.Values{}
	query.Set("p", "0")

	b, err := c.getFromStore("wishlist/profiles/"+strconv.FormatInt(playerID, 10)+"/wishlistdata", query)
	if err != nil {
		return wishlist, err
	}

	// No items
	if strings.TrimSpace(string(b)) == "[]" {
		return wishlist, err
	}

	// Check for fail response
	failResp := WishlistFail{}
	err = json.Unmarshal(b, &failResp)
	if err == nil && failResp.Success > 0 {
		return wishlist, ErrWishlistNotFound
	}

	// Unmarshal JSON
	err = json.Unmarshal(b, &wishlist.Items)
	return wishlist, err
}

type WishlistFail struct {
	Success int `json:"success"`
}

type Wishlist struct {
	Items map[unmarshal.Int]WishlistItem
}

type WishlistItem struct {
	Name           string          `json:"name"`
	Capsule        string          `json:"capsule"`
	ReviewScore    int             `json:"review_score"`
	ReviewDesc     string          `json:"review_desc"`
	ReviewsTotal   string          `json:"reviews_total"`
	ReviewsPercent int             `json:"reviews_percent"`
	ReleaseDate    unmarshal.Int64 `json:"release_date"`
	ReleaseString  string          `json:"release_string"`
	PlatformIcons  string          `json:"platform_icons"`
	Subs           []struct {
		ID            int           `json:"id"`
		DiscountBlock string        `json:"discount_block"`
		DiscountPct   int           `json:"discount_pct"`
		Price         unmarshal.Int `json:"price"`
	} `json:"subs"`
	Type        string         `json:"type"`
	Screenshots []string       `json:"screenshots"`
	ReviewCSS   string         `json:"review_css"`
	Priority    int            `json:"priority"`
	Added       int            `json:"added"`
	Background  string         `json:"background"`
	Rank        unmarshal.Int  `json:"rank"`
	Tags        []string       `json:"tags"`
	EarlyAccess int            `json:"early_access"`
	IsFreeGame  bool           `json:"is_free_game"`
	Win         unmarshal.Bool `json:"win"`
	Mac         unmarshal.Bool `json:"mac"`
	Linux       unmarshal.Bool `json:"linux"`
}
