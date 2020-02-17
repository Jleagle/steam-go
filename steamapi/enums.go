package steamapi

import (
	"strings"
)

// https://partner.steamgames.com/doc/store/localization
type LanguageCode string

func (lc LanguageCode) Title() string {

	switch lc {
	case LanguageChineseSimplified:
		return "Chinese (Simplified)"
	case LanguageChineseTraditional:
		return "Chinese (Traditional)"
	case LanguageSpanishLatinAmerica:
		return "Spanish (Latin America)"
	default:
		return strings.Title(string(lc))
	}
}

const (
	LanguageArabic              LanguageCode = "arabic"
	LanguageBulgarian           LanguageCode = "bulgarian"
	LanguageChineseSimplified   LanguageCode = "schinese"
	LanguageChineseTraditional  LanguageCode = "tchinese"
	LanguageCzech               LanguageCode = "czech"
	LanguageDanish              LanguageCode = "danish"
	LanguageDutch               LanguageCode = "dutch"
	LanguageEnglish             LanguageCode = "english"
	LanguageFinnish             LanguageCode = "finnish"
	LanguageFrench              LanguageCode = "french"
	LanguageGerman              LanguageCode = "german"
	LanguageGreek               LanguageCode = "greek"
	LanguageHungarian           LanguageCode = "hungarian"
	LanguageItalian             LanguageCode = "italian"
	LanguageJapanese            LanguageCode = "japanese"
	LanguageKorean              LanguageCode = "koreana"
	LanguageNorwegian           LanguageCode = "norwegian"
	LanguagePolish              LanguageCode = "polish"
	LanguagePortuguese          LanguageCode = "portuguese"
	LanguagePortugueseBrazil    LanguageCode = "brazilian"
	LanguageRomanian            LanguageCode = "romanian"
	LanguageRussian             LanguageCode = "russian"
	LanguageSpanishSpain        LanguageCode = "spanish"
	LanguageSpanishLatinAmerica LanguageCode = "latam"
	LanguageSwedish             LanguageCode = "swedish"
	LanguageThai                LanguageCode = "thai"
	LanguageTurkish             LanguageCode = "turkish"
	LanguageUkrainian           LanguageCode = "ukrainian"
	LanguageVietnamese          LanguageCode = "vietnamese"
)

var LanguageCodes = []LanguageCode{LanguageArabic, LanguageBulgarian, LanguageChineseSimplified, LanguageChineseTraditional, LanguageCzech, LanguageDanish,
	LanguageDutch, LanguageEnglish, LanguageFinnish, LanguageFrench, LanguageGerman, LanguageGreek, LanguageHungarian, LanguageItalian, LanguageJapanese,
	LanguageKorean, LanguageNorwegian, LanguagePolish, LanguagePortuguese, LanguagePortugueseBrazil, LanguageRomanian, LanguageRussian, LanguageSpanishSpain,
	LanguageSpanishLatinAmerica, LanguageSwedish, LanguageThai, LanguageTurkish, LanguageUkrainian, LanguageVietnamese}

type LanguageCodeWeb string

const (
	LanguageWebAR    LanguageCodeWeb = "ar"
	LanguageWebBG    LanguageCodeWeb = "bg"
	LanguageWebZHCN  LanguageCodeWeb = "zh-CN"
	LanguageWebZHTW  LanguageCodeWeb = "zh-TW"
	LanguageWebCS    LanguageCodeWeb = "cs"
	LanguageWebDA    LanguageCodeWeb = "da"
	LanguageWebNL    LanguageCodeWeb = "nl"
	LanguageWebEN    LanguageCodeWeb = "en"
	LanguageWebFI    LanguageCodeWeb = "fi"
	LanguageWebFR    LanguageCodeWeb = "fr"
	LanguageWebDE    LanguageCodeWeb = "de"
	LanguageWebEL    LanguageCodeWeb = "el"
	LanguageWebHU    LanguageCodeWeb = "hu"
	LanguageWebIT    LanguageCodeWeb = "it"
	LanguageWebJA    LanguageCodeWeb = "ja"
	LanguageWebKO    LanguageCodeWeb = "ko"
	LanguageWebNO    LanguageCodeWeb = "no"
	LanguageWebPL    LanguageCodeWeb = "pl"
	LanguageWebPT    LanguageCodeWeb = "pt"
	LanguageWebPTBR  LanguageCodeWeb = "pt-BR"
	LanguageWebRO    LanguageCodeWeb = "ro"
	LanguageWebRU    LanguageCodeWeb = "ru"
	LanguageWebES    LanguageCodeWeb = "es"
	LanguageWebES419 LanguageCodeWeb = "es-419"
	LanguageWebSV    LanguageCodeWeb = "sv"
	LanguageWebTH    LanguageCodeWeb = "th"
	LanguageWebTR    LanguageCodeWeb = "tr"
	LanguageWebUK    LanguageCodeWeb = "uk"
	LanguageWebVN    LanguageCodeWeb = "vn"
)

type Language struct {
	EnglishName        string
	NativeName         string
	LanguageCode       LanguageCode
	WebAPILanguageCode LanguageCodeWeb
}

var Languages = []Language{
	{
		EnglishName:        "Arabic",
		NativeName:         "العربية",
		LanguageCode:       LanguageArabic,
		WebAPILanguageCode: LanguageWebAR,
	},
	{
		EnglishName:        "Bulgarian",
		NativeName:         "български език",
		LanguageCode:       LanguageBulgarian,
		WebAPILanguageCode: LanguageWebBG,
	},
	{
		EnglishName:        "Chinese (Simplified)",
		NativeName:         "简体中文",
		LanguageCode:       LanguageChineseSimplified,
		WebAPILanguageCode: LanguageWebZHCN,
	},
	{
		EnglishName:        "Chinese (Traditional)",
		NativeName:         "繁體中文",
		LanguageCode:       LanguageChineseTraditional,
		WebAPILanguageCode: LanguageWebZHTW,
	},
	{
		EnglishName:        "Czech",
		NativeName:         "čeština",
		LanguageCode:       LanguageCzech,
		WebAPILanguageCode: LanguageWebCS,
	},
	{
		EnglishName:        "Danish",
		NativeName:         "Dansk",
		LanguageCode:       LanguageDanish,
		WebAPILanguageCode: LanguageWebDA,
	},
	{
		EnglishName:        "Dutch",
		NativeName:         "Nederlands",
		LanguageCode:       LanguageDutch,
		WebAPILanguageCode: LanguageWebNL,
	},
	{
		EnglishName:        "English",
		NativeName:         "English",
		LanguageCode:       LanguageEnglish,
		WebAPILanguageCode: LanguageWebEN,
	},
	{
		EnglishName:        "Finnish",
		NativeName:         "Suomi",
		LanguageCode:       LanguageFinnish,
		WebAPILanguageCode: LanguageWebFI,
	},
	{
		EnglishName:        "French",
		NativeName:         "Français",
		LanguageCode:       LanguageFrench,
		WebAPILanguageCode: LanguageWebFR,
	},
	{
		EnglishName:        "German",
		NativeName:         "Deutsch",
		LanguageCode:       LanguageGerman,
		WebAPILanguageCode: LanguageWebDE,
	},
	{
		EnglishName:        "Greek",
		NativeName:         "Ελληνικά",
		LanguageCode:       LanguageGreek,
		WebAPILanguageCode: LanguageWebEL,
	},
	{
		EnglishName:        "Hungarian",
		NativeName:         "Magyar",
		LanguageCode:       LanguageHungarian,
		WebAPILanguageCode: LanguageWebHU,
	},
	{
		EnglishName:        "Italian",
		NativeName:         "Italiano",
		LanguageCode:       LanguageItalian,
		WebAPILanguageCode: LanguageWebIT,
	},
	{
		EnglishName:        "Japanese",
		NativeName:         "日本語",
		LanguageCode:       LanguageJapanese,
		WebAPILanguageCode: LanguageWebJA,
	},
	{
		EnglishName:        "Korean",
		NativeName:         "한국어",
		LanguageCode:       LanguageKorean,
		WebAPILanguageCode: LanguageWebKO,
	},
	{
		EnglishName:        "Norwegian",
		NativeName:         "Norsk",
		LanguageCode:       LanguageNorwegian,
		WebAPILanguageCode: LanguageWebNO,
	},
	{
		EnglishName:        "Polish",
		NativeName:         "Polski",
		LanguageCode:       LanguagePolish,
		WebAPILanguageCode: LanguageWebPL,
	},
	{
		EnglishName:        "Portuguese",
		NativeName:         "Português",
		LanguageCode:       LanguagePortuguese,
		WebAPILanguageCode: LanguageWebPT,
	},
	{
		EnglishName:        "Portuguese-Brazil",
		NativeName:         "Português-Brasil",
		LanguageCode:       LanguagePortugueseBrazil,
		WebAPILanguageCode: LanguageWebPTBR,
	},
	{
		EnglishName:        "Romanian",
		NativeName:         "Română",
		LanguageCode:       LanguageRomanian,
		WebAPILanguageCode: LanguageWebRO,
	},
	{
		EnglishName:        "Russian",
		NativeName:         "Русский",
		LanguageCode:       LanguageRussian,
		WebAPILanguageCode: LanguageWebRU,
	},
	{
		EnglishName:        "Spanish-Spain",
		NativeName:         "Español-España",
		LanguageCode:       LanguageSpanishSpain,
		WebAPILanguageCode: LanguageWebES,
	},
	{
		EnglishName:        "Spanish-Latin America",
		NativeName:         "Español-Latinoamérica",
		LanguageCode:       LanguageSpanishLatinAmerica,
		WebAPILanguageCode: LanguageWebES419,
	},
	{
		EnglishName:        "Swedish",
		NativeName:         "Svenska",
		LanguageCode:       LanguageSwedish,
		WebAPILanguageCode: LanguageWebSV,
	},
	{
		EnglishName:        "Thai",
		NativeName:         "ไทย",
		LanguageCode:       LanguageThai,
		WebAPILanguageCode: LanguageWebTH,
	},
	{
		EnglishName:        "Turkish",
		NativeName:         "Türkçe",
		LanguageCode:       LanguageTurkish,
		WebAPILanguageCode: LanguageWebTR,
	},
	{
		EnglishName:        "Ukrainian",
		NativeName:         "Українська",
		LanguageCode:       LanguageUkrainian,
		WebAPILanguageCode: LanguageWebUK,
	},
	{
		EnglishName:        "Vietnamese",
		NativeName:         "Tiếng Việt",
		LanguageCode:       LanguageVietnamese,
		WebAPILanguageCode: LanguageWebVN,
	},
}

// https://partner.steamgames.com/doc/store/pricing/currencies
type CurrencyCode string

const (
	CurrencyAED CurrencyCode = "AED"
	CurrencyARS CurrencyCode = "ARS"
	CurrencyAUD CurrencyCode = "AUD"
	CurrencyBRL CurrencyCode = "BRL"
	CurrencyCAD CurrencyCode = "CAD"
	CurrencyCHF CurrencyCode = "CHF"
	CurrencyCLP CurrencyCode = "CLP"
	CurrencyCNY CurrencyCode = "CNY"
	CurrencyCOP CurrencyCode = "COP"
	CurrencyCRC CurrencyCode = "CRC"
	CurrencyEUR CurrencyCode = "EUR"
	CurrencyGBP CurrencyCode = "GBP"
	CurrencyHKD CurrencyCode = "HKD"
	CurrencyILS CurrencyCode = "ILS"
	CurrencyIDR CurrencyCode = "IDR"
	CurrencyINR CurrencyCode = "INR"
	CurrencyJPY CurrencyCode = "JPY"
	CurrencyKRW CurrencyCode = "KRW"
	CurrencyKWD CurrencyCode = "KWD"
	CurrencyKZT CurrencyCode = "KZT"
	CurrencyMXN CurrencyCode = "MXN"
	CurrencyMYR CurrencyCode = "MYR"
	CurrencyNOK CurrencyCode = "NOK"
	CurrencyNZD CurrencyCode = "NZD"
	CurrencyPEN CurrencyCode = "PEN"
	CurrencyPHP CurrencyCode = "PHP"
	CurrencyPLN CurrencyCode = "PLN"
	CurrencyQAR CurrencyCode = "QAR"
	CurrencyRUB CurrencyCode = "RUB"
	CurrencySAR CurrencyCode = "SAR"
	CurrencySGD CurrencyCode = "SGD"
	CurrencyTHB CurrencyCode = "THB"
	CurrencyTRY CurrencyCode = "TRY"
	CurrencyTWD CurrencyCode = "TWD"
	CurrencyUAH CurrencyCode = "UAH"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyUYU CurrencyCode = "UYU"
	CurrencyVND CurrencyCode = "VND"
	CurrencyZAR CurrencyCode = "ZAR"
)

type Currency struct {
	CurrencyCode CurrencyCode
	Description  string
}

var Currencies = []Currency{
	{CurrencyCode: CurrencyARS, Description: "Argentine Peso"},
	{CurrencyCode: CurrencyAUD, Description: "Australian Dollars"},
	{CurrencyCode: CurrencyBRL, Description: "Brazilian Reals"},
	{CurrencyCode: CurrencyCAD, Description: "Canadian Dollars"},
	{CurrencyCode: CurrencyCLP, Description: "Chilean Peso"},
	{CurrencyCode: CurrencyCNY, Description: "Chinese Renminbi (yuan)"},
	{CurrencyCode: CurrencyCOP, Description: "Colombian Peso"},
	{CurrencyCode: CurrencyCRC, Description: "Costa Rican Colón"},
	{CurrencyCode: CurrencyEUR, Description: "European Union Euro"},
	{CurrencyCode: CurrencyHKD, Description: "Hong Kong Dollar"},
	{CurrencyCode: CurrencyINR, Description: "Indian Rupee"},
	{CurrencyCode: CurrencyIDR, Description: "Indonesian Rupiah"},
	{CurrencyCode: CurrencyILS, Description: "Israeli New Shekel"},
	{CurrencyCode: CurrencyJPY, Description: "Japanese Yen"},
	{CurrencyCode: CurrencyKZT, Description: "Kazakhstani Tenge"},
	{CurrencyCode: CurrencyKWD, Description: "Kuwaiti Dinar"},
	{CurrencyCode: CurrencyMYR, Description: "Malaysian Ringgit"},
	{CurrencyCode: CurrencyMXN, Description: "Mexican Peso"},
	{CurrencyCode: CurrencyTWD, Description: "New Taiwan Dollar"},
	{CurrencyCode: CurrencyNZD, Description: "New Zealand Dollar"},
	{CurrencyCode: CurrencyNOK, Description: "Norwegian Krone"},
	{CurrencyCode: CurrencyPEN, Description: "Peruvian Nuevo Sol"},
	{CurrencyCode: CurrencyPHP, Description: "Philippine Peso"},
	{CurrencyCode: CurrencyPLN, Description: "Polish Złoty"},
	{CurrencyCode: CurrencyQAR, Description: "Qatari Riyal"},
	{CurrencyCode: CurrencyRUB, Description: "Russian Rouble"},
	{CurrencyCode: CurrencySAR, Description: "Saudi Riyal"},
	{CurrencyCode: CurrencySGD, Description: "Singapore Dollar"},
	{CurrencyCode: CurrencyZAR, Description: "South African Rand"},
	{CurrencyCode: CurrencyKRW, Description: "South Korean Won"},
	{CurrencyCode: CurrencyCHF, Description: "Swiss Francs"},
	{CurrencyCode: CurrencyTHB, Description: "Thai Baht"},
	{CurrencyCode: CurrencyTRY, Description: "Turkish Lira"},
	{CurrencyCode: CurrencyUAH, Description: "Ukrainian Hryvnia"},
	{CurrencyCode: CurrencyAED, Description: "United Arab Emirates Dirham"},
	{CurrencyCode: CurrencyGBP, Description: "United Kingdom Pound"},
	{CurrencyCode: CurrencyUSD, Description: "United States Dollar"},
	{CurrencyCode: CurrencyUYU, Description: "Uruguayan Peso"},
	{CurrencyCode: CurrencyVND, Description: "Vietnamese Dong"},
}

// Country codes to get info on apps/packages
type ProductCC string

const (
	ProductCCAE ProductCC = "ae"
	ProductCCAR ProductCC = "ar"
	ProductCCAU ProductCC = "au"
	ProductCCAZ ProductCC = "az"
	ProductCCBR ProductCC = "br"
	ProductCCCA ProductCC = "ca"
	ProductCCCH ProductCC = "ch"
	ProductCCCL ProductCC = "cl"
	ProductCCCN ProductCC = "cn"
	ProductCCCO ProductCC = "co"
	ProductCCCR ProductCC = "cr"
	ProductCCEU ProductCC = "eu"
	ProductCCHK ProductCC = "hk"
	ProductCCID ProductCC = "id"
	ProductCCIL ProductCC = "il"
	ProductCCIN ProductCC = "in"
	ProductCCJP ProductCC = "jp"
	ProductCCKR ProductCC = "kr"
	ProductCCKW ProductCC = "kw"
	ProductCCKZ ProductCC = "kz"
	ProductCCMX ProductCC = "mx"
	ProductCCMY ProductCC = "my"
	ProductCCNO ProductCC = "no"
	ProductCCNZ ProductCC = "nz"
	ProductCCPE ProductCC = "pe"
	ProductCCPH ProductCC = "ph"
	ProductCCPK ProductCC = "pk"
	ProductCCPL ProductCC = "pl"
	ProductCCQA ProductCC = "qa"
	ProductCCRU ProductCC = "ru"
	ProductCCSA ProductCC = "sa"
	ProductCCSG ProductCC = "sg"
	ProductCCTH ProductCC = "th"
	ProductCCTR ProductCC = "tr"
	ProductCCTW ProductCC = "tw"
	ProductCCUA ProductCC = "ua"
	ProductCCUK ProductCC = "uk"
	ProductCCUS ProductCC = "us"
	ProductCCUY ProductCC = "uy"
	ProductCCVN ProductCC = "vn"
	ProductCCZA ProductCC = "za"
)

var ProductCCs = []ProductCC{ProductCCAE, ProductCCAR, ProductCCAU, ProductCCAZ, ProductCCBR, ProductCCCA, ProductCCCH, ProductCCCL,
	ProductCCCN, ProductCCCO, ProductCCCR, ProductCCEU, ProductCCHK, ProductCCID, ProductCCIL, ProductCCIN, ProductCCJP, ProductCCKR,
	ProductCCKW, ProductCCKZ, ProductCCMX, ProductCCMY, ProductCCNO, ProductCCNZ, ProductCCPE, ProductCCPH, ProductCCPK, ProductCCPL,
	ProductCCQA, ProductCCRU, ProductCCSA, ProductCCSG, ProductCCTH, ProductCCTR, ProductCCTW, ProductCCUA, ProductCCUK, ProductCCUS,
	ProductCCUY, ProductCCVN, ProductCCZA}

func IsProductCC(i string) bool {
	for _, v := range ProductCCs {
		if string(v) == i {
			return true
		}
	}
	return false
}
