package steam

type CountryCode string

const (
	CountryAE CountryCode = "AE"
	CountryAR             = "AR"
	CountryAU             = "AU"
	CountryBR             = "BR"
	CountryCA             = "CA"
	CountryCH             = "CH"
	CountryCL             = "CL"
	CountryCN             = "CN"
	CountryCO             = "CO"
	CountryCR             = "CR"
	CountryDE             = "DE"
	CountryGB             = "GB"
	CountryHK             = "HK"
	CountryIL             = "IL"
	CountryID             = "ID"
	CountryIN             = "IN"
	CountryJP             = "JP"
	CountryKR             = "KR"
	CountryKW             = "KW"
	CountryKZ             = "KZ"
	CountryMX             = "MX"
	CountryMY             = "MY"
	CountryNO             = "NO"
	CountryNZ             = "NZ"
	CountryPE             = "PE"
	CountryPH             = "PH"
	CountryPL             = "PL"
	CountryQA             = "QA"
	CountryRU             = "RU"
	CountrySA             = "SA"
	CountrySG             = "SG"
	CountryTH             = "TH"
	CountryTR             = "TR"
	CountryTW             = "TW"
	CountryUA             = "UA"
	CountryUS             = "US"
	CountryUY             = "UY"
	CountryVN             = "VN"
	CountryZA             = "ZA"
)

var Countries = map[CountryCode]string{
	CountryAE: "United Arab Emirates Dirham",
	CountryAR: "Argentine Peso",
	CountryAU: "Australian Dollars",
	CountryBR: "Brazilian Reals",
	CountryCA: "Canadian Dollars",
	CountryCH: "Swiss Francs",
	CountryCL: "Chilean Peso",
	CountryCN: "Chinese Renminbi (yuan)",
	CountryCO: "Colombian Peso",
	CountryCR: "Costa Rican Colón",
	CountryDE: "European Union Euro",
	CountryGB: "United Kingdom Pound",
	CountryHK: "Hong Kong Dollar",
	CountryIL: "Israeli New Shekel",
	CountryID: "Indonesian Rupiah",
	CountryIN: "Indian Rupee",
	CountryJP: "Japanese Yen",
	CountryKR: "South Korean Won",
	CountryKW: "Kuwaiti Dinar",
	CountryKZ: "Kazakhstani Tenge",
	CountryMX: "Mexican Peso",
	CountryMY: "Malaysian Ringgit",
	CountryNO: "Norwegian Krone",
	CountryNZ: "New Zealand Dollar",
	CountryPE: "Peruvian Nuevo Sol",
	CountryPH: "Philippine Peso",
	CountryPL: "Polish Złoty",
	CountryQA: "Qatari Riyal",
	CountryRU: "Russian Rouble",
	CountrySA: "Saudi Riyal",
	CountrySG: "Singapore Dollar",
	CountryTH: "Thai Baht",
	CountryTR: "Turkish Lira",
	CountryTW: "New Taiwan Dollar",
	CountryUA: "Ukrainian Hryvnia",
	CountryUS: "United States Dollar",
	CountryUY: "Uruguayan Peso",
	CountryVN: "Vietnamese Dong",
	CountryZA: "South African Rand",
}

type Language string

const (
	LanguageArabic              Language = "arabic"
	LanguageBulgarian                    = "bulgarian"
	LanguageChineseSimplified            = "schinese"
	LanguageChineseTraditional           = "tchinese"
	LanguageCzech                        = "czech"
	LanguageDanish                       = "danish"
	LanguageDutch                        = "dutch"
	LanguageEnglish                      = "english"
	LanguageFinnish                      = "finnish"
	LanguageFrench                       = "french"
	LanguageGerman                       = "german"
	LanguageGreek                        = "greek"
	LanguageHungarian                    = "hungarian"
	LanguageItalian                      = "italian"
	LanguageJapanese                     = "japanese"
	LanguageKorean                       = "koreana"
	LanguageNorwegian                    = "norwegian"
	LanguagePolish                       = "polish"
	LanguagePortuguese                   = "portuguese"
	LanguagePortugueseBrazil             = "brazilian"
	LanguageRomanian                     = "romanian"
	LanguageRussian                      = "russian"
	LanguageSpanishSpain                 = "spanish"
	LanguageSpanishLatinAmerica          = "latam"
	LanguageSwedish                      = "swedish"
	LanguageThai                         = "thai"
	LanguageTurkish                      = "turkish"
	LanguageUkrainian                    = "ukrainian"
	LanguageVietnamese                   = "vietnamese"
)

var Languages = map[Language]string{
	LanguageArabic:              "Arabic",
	LanguageBulgarian:           "Bulgarian",
	LanguageChineseSimplified:   "Chinese (Simplified)",
	LanguageChineseTraditional:  "Chinese (Traditional)",
	LanguageCzech:               "Czech",
	LanguageDanish:              "Danish",
	LanguageDutch:               "Dutch",
	LanguageEnglish:             "English",
	LanguageFinnish:             "Finnish",
	LanguageFrench:              "French",
	LanguageGerman:              "German",
	LanguageGreek:               "Greek",
	LanguageHungarian:           "Hungarian",
	LanguageItalian:             "Italian",
	LanguageJapanese:            "Japanese",
	LanguageKorean:              "Korean",
	LanguageNorwegian:           "Norwegian",
	LanguagePolish:              "Polish",
	LanguagePortuguese:          "Portuguese",
	LanguagePortugueseBrazil:    "Portuguese-Brazil",
	LanguageRomanian:            "Romanian",
	LanguageRussian:             "Russian",
	LanguageSpanishSpain:        "Spanish-Spain",
	LanguageSpanishLatinAmerica: "Spanish-Latin America",
	LanguageSwedish:             "Swedish",
	LanguageThai:                "Thai",
	LanguageTurkish:             "Turkish",
	LanguageUkrainian:           "Ukrainian",
	LanguageVietnamese:          "Vietnamese",
}
