package steam

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrAppNotFound     = errors.New("steam: store: app not found")
	ErrPackageNotFound = errors.New("steam: store: package not found")
	ErrNullResponse    = errors.New("steam: store: null response")
)

func (s Steam) GetAppDetails(id int) (app AppDetailsBody, bytes []byte, err error) {

	idx := strconv.Itoa(id)

	query := url.Values{}
	query.Set("appids", idx)
	query.Set("cc", "us") // 2 letter country code
	query.Set("l", "en")  // Language

	bytes, err = s.getFromStore("api/appdetails", query)
	if err != nil {
		return app, bytes, err
	}

	// Check for no app
	if string(bytes) == "null" {
		return app, bytes, ErrNullResponse
	}

	// Fix values that can change type, causing unmarshal errors
	var regex *regexp.Regexp
	var str = string(bytes)

	// Convert strings to ints
	regex = regexp.MustCompile(`:\s?"(\d+)"`) // After colon
	str = regex.ReplaceAllString(str, `:$1`)

	regex = regexp.MustCompile(`,\s?"(\d+)"`) // After comma
	str = regex.ReplaceAllString(str, `,$1`)

	regex = regexp.MustCompile(`"(\d+)",`) // Before comma
	str = regex.ReplaceAllString(str, `$1,`)

	regex = regexp.MustCompile(`"packages":\s?\["(\d+)"]`) // Package array with single int
	str = regex.ReplaceAllString(str, `"packages":[$1]`)

	// Make some its strings again
	regex = regexp.MustCompile(`"date":\s?(\d+)`)
	str = regex.ReplaceAllString(str, `"date":"$1"`)

	regex = regexp.MustCompile(`"name":\s?(\d+)`)
	str = regex.ReplaceAllString(str, `"name":"$1"`)

	regex = regexp.MustCompile(`"description":\s?(\d+)`)
	str = regex.ReplaceAllString(str, `"description":"$1"`)

	regex = regexp.MustCompile(`"display_type":\s?(\d+)`)
	str = regex.ReplaceAllString(str, `"display_type":"$1"`)

	regex = regexp.MustCompile(`"legal_notice":\s?(\d+)`)
	str = regex.ReplaceAllString(str, `"legal_notice":"$1"`)

	// Fix arrays that should be objects
	str = strings.Replace(str, "\"pc_requirements\":[]", "\"pc_requirements\":null", 1)
	str = strings.Replace(str, "\"mac_requirements\":[]", "\"mac_requirements\":null", 1)
	str = strings.Replace(str, "\"linux_requirements\":[]", "\"linux_requirements\":null", 1)

	bytes = []byte(str)

	// Unmarshal JSON
	resp := map[string]AppDetailsBody{}
	if err := json.Unmarshal(bytes, &resp); err != nil {
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
		Type                string `json:"type"`
		Name                string `json:"name"`
		AppID               int    `json:"steam_appid"`
		RequiredAge         int    `json:"required_age"`
		IsFree              bool   `json:"is_free"`
		DLC                 []int  `json:"dlc"`
		ControllerSupport   string `json:"controller_support"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		Fullgame            struct {
			AppID int    `json:"appid"`
			Name  string `json:"name"`
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
			AppID       int    `json:"appid"`
			Description string `json:"description"`
		} `json:"demos"`
		PriceOverview struct {
			Currency        string `json:"currency"`
			Initial         int    `json:"initial"`
			Final           int    `json:"final"`
			DiscountPercent int    `json:"discount_percent"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string `json:"name"`
			Title                   string `json:"title"`
			Description             string `json:"description"`
			SelectionText           string `json:"selection_text"`
			SaveText                string `json:"save_text"`
			DisplayType             string `json:"display_type"` // Can be string or int
			IsRecurringSubscription string `json:"is_recurring_subscription"`
			Subs                    []struct {
				PackageID                int    `json:"packageid"`
				PercentSavingsText       string `json:"percent_savings_text"`
				PercentSavings           int    `json:"percent_savings"`
				OptionText               string `json:"option_text"`
				OptionDescription        string `json:"option_description"`
				CanGetFreeLicense        int    `json:"can_get_free_license"`
				IsFreeLicense            bool   `json:"is_free_license"`
				PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
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
		Categories  []AppDetailsCategory   `json:"categories"`
		Genres      []AppDetailsGenre      `json:"genres"`
		Screenshots []AppDetailsScreenshot `json:"screenshots"`
		Movies      []struct {
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
		Achievements AppDetailsAchievements `json:"achievements"`
		ReleaseDate  struct {
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

type AppDetailsScreenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type AppDetailsAchievements struct {
	Total       int `json:"total"`
	Highlighted []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"highlighted"`
}

type AppDetailsGenre struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type AppDetailsCategory struct {
	ID          int8   `json:"id"`
	Description string `json:"description"`
}

func (s Steam) GetPackageDetails(id int) (pack PackageDetailsBody, bytes []byte, err error) {

	idx := strconv.Itoa(id)

	query := url.Values{}
	query.Set("packageids", idx)
	query.Set("cc", "us") // 2 letter country code
	query.Set("l", "en")  // Language

	bytes, err = s.getFromStore("api/packagedetails", query)
	if err != nil {
		return pack, bytes, err
	}

	// Check for no pack
	if string(bytes) == "null" {
		return pack, bytes, ErrNullResponse
	}

	// Unmarshal JSON
	resp := map[string]PackageDetailsBody{}
	if err := json.Unmarshal(bytes, &resp); err != nil {
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
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return tags, bytes, err
	}

	return Tags{Tags: resp,}, bytes, nil
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

func (t Tags) GetIDs() (ids []int) {

	for _, v := range t.Tags {
		ids = append(ids, v.TagID)
	}
	return ids
}

type Tag struct {
	TagID int    `json:"tagid"`
	Name  string `json:"name"`
}

func (s Steam) GetReviews(appID int) (reviews ReviewsResponse, bytes []byte, err error) {

	query := url.Values{}
	query.Set("json", "1")
	query.Set("filter", "all")
	//query.Set("language", "all")
	query.Set("day_range", "all")
	query.Set("start_offset", "0")
	query.Set("review_type", "all")
	query.Set("purchase_type", "all")

	bytes, err = s.getFromStore("appreviews/"+strconv.Itoa(appID), query)
	if err != nil {
		return reviews, bytes, err
	}

	str := string(bytes)

	regex := regexp.MustCompile(`"comment_count":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"comment_count": $1`)

	regex = regexp.MustCompile(`"steamid":\s?"(\d+)"`)
	str = regex.ReplaceAllString(str, `"steamid": $1`)

	bytes = []byte(str)

	// Unmarshal JSON
	if err := json.Unmarshal(bytes, &reviews); err != nil {
		return reviews, bytes, err
	}

	return reviews, bytes, nil
}

type ReviewsResponse struct {
	Success      int                    `json:"success"`
	QuerySummary ReviewsSummaryResponse `json:"query_summary"`
	Reviews      []struct {
		Recommendationid string `json:"recommendationid"`
		Author           struct {
			SteamID              int64 `json:"steamid"`
			NumGamesOwned        int   `json:"num_games_owned"`
			NumReviews           int   `json:"num_reviews"`
			PlaytimeForever      int   `json:"playtime_forever"`
			PlaytimeLastTwoWeeks int   `json:"playtime_last_two_weeks"`
			LastPlayed           int   `json:"last_played"`
		} `json:"author"`
		Language                 string `json:"language"`
		Review                   string `json:"review"`
		TimestampCreated         int64  `json:"timestamp_created"`
		TimestampUpdated         int64  `json:"timestamp_updated"`
		VotedUp                  bool   `json:"voted_up"`
		VotesUp                  int    `json:"votes_up"`
		VotesFunny               int    `json:"votes_funny"`
		WeightedVoteScore        string `json:"weighted_vote_score"`
		CommentCount             int    `json:"comment_count"`
		SteamPurchase            bool   `json:"steam_purchase"`
		ReceivedForFree          bool   `json:"received_for_free"`
		WrittenDuringEarlyAccess bool   `json:"written_during_early_access"`
	} `json:"reviews"`
}

type ReviewsSummaryResponse struct {
	NumReviews      int     `json:"num_reviews"`
	ReviewScore     float64 `json:"review_score"`
	ReviewScoreDesc string  `json:"review_score_desc"`
	TotalPositive   int     `json:"total_positive"`
	TotalNegative   int     `json:"total_negative"`
	TotalReviews    int     `json:"total_reviews"`
}

func (r ReviewsSummaryResponse) GetPositivePercent() float64 {
	return float64(r.TotalPositive) / float64(r.TotalReviews) * 100
}

func (r ReviewsSummaryResponse) GetNegativePercent() float64 {
	return float64(r.TotalNegative) / float64(r.TotalReviews) * 100
}
