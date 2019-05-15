package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Jleagle/unmarshal-go/ctypes"
)

var (
	ErrAppNotFound      = errors.New("steam: store: app not found")
	ErrPackageNotFound  = errors.New("steam: store: package not found")
	ErrWishlistNotFound = errors.New("steam: store: wishlist not found")
	ErrNullResponse     = errors.New("steam: store: null response") // Probably being rate limited
	ErrHTMLResponse     = errors.New("steam: store: html response") // Probably down
)

func (s Steam) GetAppDetails(id int, code CountryCode, language Language) (app AppDetailsBody, bytes []byte, err error) {

	idx := strconv.Itoa(id)

	query := url.Values{}
	query.Set("appids", idx)
	query.Set("cc", string(code))    // Price currency
	query.Set("l", string(language)) // Text

	bytes, err = s.getFromStore("api/appdetails", query)
	if err != nil {
		return app, bytes, err
	}

	// Check invalid responses
	if string(bytes) == "null" {
		return app, bytes, ErrNullResponse
	}
	if strings.HasPrefix(string(bytes), "<") {
		return app, bytes, ErrHTMLResponse
	}

	// Fix arrays that should be objects
	var str = string(bytes)
	str = strings.Replace(str, "\"pc_requirements\":[]", "\"pc_requirements\":{}", 1)
	str = strings.Replace(str, "\"mac_requirements\":[]", "\"mac_requirements\":{}", 1)
	str = strings.Replace(str, "\"linux_requirements\":[]", "\"linux_requirements\":{}", 1)
	bytes = []byte(str)

	// Unmarshal JSON
	resp := map[string]AppDetailsBody{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return app, bytes, err
	}

	if resp[idx].Success == false {
		return app, bytes, ErrAppNotFound
	}

	return resp[idx], bytes, nil
}

type AppDetailsBody struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string      `json:"type"`
		Name                string      `json:"name"`
		AppID               int         `json:"steam_appid"`
		RequiredAge         ctypes.CInt `json:"required_age"`
		IsFree              bool        `json:"is_free"`
		DLC                 []int       `json:"dlc"`
		ControllerSupport   string      `json:"controller_support"`
		DetailedDescription string      `json:"detailed_description"`
		AboutTheGame        string      `json:"about_the_game"`
		ShortDescription    string      `json:"short_description"`
		Fullgame            struct {
			AppID ctypes.CInt `json:"appid"`
			Name  string      `json:"name"`
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
		LegalNotice string   `json:"legal_notice"`
		Developers  []string `json:"developers"`
		Publishers  []string `json:"publishers"`
		Demos       []struct {
			AppID       ctypes.CInt `json:"appid"`
			Description string      `json:"description"`
		} `json:"demos"`
		PriceOverview struct {
			Currency        string `json:"currency"`
			Initial         int    `json:"initial"`
			Final           int    `json:"final"`
			DiscountPercent int    `json:"discount_percent"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string         `json:"name"`
			Title                   string         `json:"title"`
			Description             string         `json:"description"`
			SelectionText           string         `json:"selection_text"`
			SaveText                string         `json:"save_text"`
			DisplayType             ctypes.CString `json:"display_type"`
			IsRecurringSubscription string         `json:"is_recurring_subscription"`
			Subs                    []struct {
				PackageID                int          `json:"packageid"`
				PercentSavingsText       string       `json:"percent_savings_text"`
				PercentSavings           int          `json:"percent_savings"`
				OptionText               string       `json:"option_text"`
				OptionDescription        string       `json:"option_description"`
				CanGetFreeLicense        ctypes.CBool `json:"can_get_free_license"`
				IsFreeLicense            bool         `json:"is_free_license"`
				PriceInCentsWithDiscount int          `json:"price_in_cents_with_discount"`
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
		Categories []struct {
			ID          int8   `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          ctypes.CInt `json:"id"`
			Description string      `json:"description"`
		} `json:"genres"`
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
		Background string `json:"background"`
	} `json:"data"`
}

func (s Steam) GetPackageDetails(id int, code CountryCode, language Language) (pack PackageDetailsBody, bytes []byte, err error) {

	idx := strconv.Itoa(id)

	query := url.Values{}
	query.Set("packageids", idx)
	query.Set("cc", string(code))    // Price currency
	query.Set("l", string(language)) // Text

	bytes, err = s.getFromStore("api/packagedetails", query)
	if err != nil {
		return pack, bytes, err
	}

	// Check invalid responses
	if string(bytes) == "null" {
		return pack, bytes, ErrNullResponse
	}
	if strings.HasPrefix(string(bytes), "<") {
		return pack, bytes, ErrHTMLResponse
	}

	// Unmarshal JSON
	resp := map[string]PackageDetailsBody{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return pack, bytes, err
	}

	if resp[idx].Success == false {
		return pack, bytes, ErrPackageNotFound
	}

	return resp[idx], bytes, nil
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
			Currency        string `json:"currency"`
			Initial         int    `json:"initial"`
			Final           int    `json:"final"`
			DiscountPercent int    `json:"discount_percent"`
			Individual      int    `json:"individual"`
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

func (s Steam) GetTags() (tags Tags, bytes []byte, err error) {

	bytes, err = s.getFromStore("tagdata/populartags/english", url.Values{})
	if err != nil {
		return tags, bytes, err
	}

	var resp []Tag
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return tags, bytes, err
	}

	return Tags{Tags: resp}, bytes, nil
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

func (s Steam) GetReviews(appID int) (reviews ReviewsResponse, bytes []byte, err error) {

	query := url.Values{}
	query.Set("json", "1")
	query.Set("filter", "all") // all / summary
	query.Set("language", string(LanguageEnglish))
	query.Set("day_range", "all")
	query.Set("start_offset", "0")
	query.Set("review_type", "all")
	query.Set("purchase_type", "all")
	// query.Set("start_date", "")
	// query.Set("end_date", "")
	// query.Set("date_range_type", "include")
	// query.Set("review_beta_enabled", "1")
	// query.Set("summary_num_positive_reviews", "1")
	// query.Set("summary_num_reviews", "1")

	bytes, err = s.getFromStore("appreviews/"+strconv.Itoa(appID), query)
	if err != nil {
		return reviews, bytes, err
	}

	// Unmarshal JSON
	err = json.Unmarshal(bytes, &reviews)
	if err != nil {
		return reviews, bytes, err
	}

	return reviews, bytes, nil
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
			SteamID              ctypes.CInt64 `json:"steamid"`
			NumGamesOwned        int           `json:"num_games_owned"`
			NumReviews           int           `json:"num_reviews"`
			PlaytimeForever      int           `json:"playtime_forever"`
			PlaytimeLastTwoWeeks int           `json:"playtime_last_two_weeks"`
			LastPlayed           int           `json:"last_played"`
		} `json:"author"`
		Language                 string          `json:"language"`
		Review                   string          `json:"review"`
		TimestampCreated         int64           `json:"timestamp_created"`
		TimestampUpdated         int64           `json:"timestamp_updated"`
		VotedUp                  bool            `json:"voted_up"`
		VotesUp                  int             `json:"votes_up"`
		VotesFunny               int             `json:"votes_funny"`
		WeightedVoteScore        ctypes.CFloat64 `json:"weighted_vote_score"`
		CommentCount             int             `json:"comment_count"`
		SteamPurchase            bool            `json:"steam_purchase"`
		ReceivedForFree          bool            `json:"received_for_free"`
		WrittenDuringEarlyAccess bool            `json:"written_during_early_access"`
	} `json:"reviews"`
}

func (r ReviewsResponse) GetPositivePercent() float64 {
	return float64(r.QuerySummary.TotalPositive) / float64(r.QuerySummary.TotalReviews) * 100
}

func (r ReviewsResponse) GetNegativePercent() float64 {
	return float64(r.QuerySummary.TotalNegative) / float64(r.QuerySummary.TotalReviews) * 100
}

func (s Steam) GetWishlist(playerID int64) (wishlist Wishlist, bytes []byte, err error) {

	query := url.Values{}
	query.Set("p", "0")

	bytes, err = s.getFromStore("wishlist/profiles/"+strconv.FormatInt(playerID, 10)+"/wishlistdata", query)
	if err != nil {
		return wishlist, bytes, err
	}

	// No items
	if strings.TrimSpace(string(bytes)) == "[]" {
		return wishlist, bytes, err
	}

	// Check for fail response
	failResp := WishlistFail{}
	err = json.Unmarshal(bytes, &failResp)
	if err == nil && failResp.Success > 0 {
		return wishlist, bytes, ErrWishlistNotFound
	}

	// Unmarshal JSON
	err = json.Unmarshal(bytes, &wishlist.Items)
	return wishlist, bytes, err
}

type WishlistFail struct {
	Success int `json:"success"`
}

type Wishlist struct {
	Items map[string]WishlistItem
}

type WishlistItem struct {
	Name           string        `json:"name"`
	Capsule        string        `json:"capsule"`
	ReviewScore    int           `json:"review_score"`
	ReviewDesc     string        `json:"review_desc"`
	ReviewsTotal   string        `json:"reviews_total"`
	ReviewsPercent int           `json:"reviews_percent"`
	ReleaseDate    ctypes.CInt64 `json:"release_date"`
	ReleaseString  string        `json:"release_string"`
	PlatformIcons  string        `json:"platform_icons"`
	Subs           []struct {
		ID            int    `json:"id"`
		DiscountBlock string `json:"discount_block"`
		DiscountPct   int    `json:"discount_pct"`
		Price         int    `json:"price"`
	} `json:"subs"`
	Type        string   `json:"type"`
	Screenshots []string `json:"screenshots"`
	ReviewCSS   string   `json:"review_css"`
	Priority    int      `json:"priority"`
	Added       int      `json:"added"`
	Background  string   `json:"background"`
	Rank        string   `json:"rank"`
	Tags        []string `json:"tags"`
	EarlyAccess int      `json:"early_access"`
	Win         int      `json:"win"`
}
