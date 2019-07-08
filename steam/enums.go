package steam

type CountryCode string

const (
	CountryAE CountryCode = "AE"
	CountryAR CountryCode = "AR"
	CountryAU CountryCode = "AU"
	CountryBR CountryCode = "BR"
	CountryCA CountryCode = "CA"
	CountryCH CountryCode = "CH"
	CountryCL CountryCode = "CL"
	CountryCN CountryCode = "CN"
	CountryCO CountryCode = "CO"
	CountryCR CountryCode = "CR"
	CountryDE CountryCode = "DE"
	CountryGB CountryCode = "GB"
	CountryHK CountryCode = "HK"
	CountryIL CountryCode = "IL"
	CountryID CountryCode = "ID"
	CountryIN CountryCode = "IN"
	CountryJP CountryCode = "JP"
	CountryKR CountryCode = "KR"
	CountryKW CountryCode = "KW"
	CountryKZ CountryCode = "KZ"
	CountryMX CountryCode = "MX"
	CountryMY CountryCode = "MY"
	CountryNO CountryCode = "NO"
	CountryNZ CountryCode = "NZ"
	CountryPE CountryCode = "PE"
	CountryPH CountryCode = "PH"
	CountryPL CountryCode = "PL"
	CountryQA CountryCode = "QA"
	CountryRU CountryCode = "RU"
	CountrySA CountryCode = "SA"
	CountrySG CountryCode = "SG"
	CountryTH CountryCode = "TH"
	CountryTR CountryCode = "TR"
	CountryTW CountryCode = "TW"
	CountryUA CountryCode = "UA"
	CountryUS CountryCode = "US"
	CountryUY CountryCode = "UY"
	CountryVN CountryCode = "VN"
	CountryZA CountryCode = "ZA"
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

func ValidCountryCode(cc CountryCode) bool {
	_, ok := Countries[cc]
	return ok
}

// https://partner.steamgames.com/doc/store/localization
type Language string

const (
	LanguageArabic              Language = "arabic"
	LanguageBulgarian           Language = "bulgarian"
	LanguageChineseSimplified   Language = "schinese"
	LanguageChineseTraditional  Language = "tchinese"
	LanguageCzech               Language = "czech"
	LanguageDanish              Language = "danish"
	LanguageDutch               Language = "dutch"
	LanguageEnglish             Language = "english"
	LanguageFinnish             Language = "finnish"
	LanguageFrench              Language = "french"
	LanguageGerman              Language = "german"
	LanguageGreek               Language = "greek"
	LanguageHungarian           Language = "hungarian"
	LanguageItalian             Language = "italian"
	LanguageJapanese            Language = "japanese"
	LanguageKorean              Language = "koreana"
	LanguageNorwegian           Language = "norwegian"
	LanguagePolish              Language = "polish"
	LanguagePortuguese          Language = "portuguese"
	LanguagePortugueseBrazil    Language = "brazilian"
	LanguageRomanian            Language = "romanian"
	LanguageRussian             Language = "russian"
	LanguageSpanishSpain        Language = "spanish"
	LanguageSpanishLatinAmerica Language = "latam"
	LanguageSwedish             Language = "swedish"
	LanguageThai                Language = "thai"
	LanguageTurkish             Language = "turkish"
	LanguageUkrainian           Language = "ukrainian"
	LanguageVietnamese          Language = "vietnamese"
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

// https://partner.steamgames.com/doc/store/pricing/currencies
type CurrencyCode string

const CurrencyAED CurrencyCode = "AED"
const CurrencyARS CurrencyCode = "ARS"
const CurrencyAUD CurrencyCode = "AUD"
const CurrencyBRL CurrencyCode = "BRL"
const CurrencyCAD CurrencyCode = "CAD"
const CurrencyCHF CurrencyCode = "CHF"
const CurrencyCLP CurrencyCode = "CLP"
const CurrencyCNY CurrencyCode = "CNY"
const CurrencyCOP CurrencyCode = "COP"
const CurrencyCRC CurrencyCode = "CRC"
const CurrencyEUR CurrencyCode = "EUR"
const CurrencyGBP CurrencyCode = "GBP"
const CurrencyHKD CurrencyCode = "HKD"
const CurrencyILS CurrencyCode = "ILS"
const CurrencyIDR CurrencyCode = "IDR"
const CurrencyINR CurrencyCode = "INR"
const CurrencyJPY CurrencyCode = "JPY"
const CurrencyKRW CurrencyCode = "KRW"
const CurrencyKWD CurrencyCode = "KWD"
const CurrencyKZT CurrencyCode = "KZT"
const CurrencyMXN CurrencyCode = "MXN"
const CurrencyMYR CurrencyCode = "MYR"
const CurrencyNOK CurrencyCode = "NOK"
const CurrencyNZD CurrencyCode = "NZD"
const CurrencyPEN CurrencyCode = "PEN"
const CurrencyPHP CurrencyCode = "PHP"
const CurrencyPLN CurrencyCode = "PLN"
const CurrencyQAR CurrencyCode = "QAR"
const CurrencyRUB CurrencyCode = "RUB"
const CurrencySAR CurrencyCode = "SAR"
const CurrencySGD CurrencyCode = "SGD"
const CurrencyTHB CurrencyCode = "THB"
const CurrencyTRY CurrencyCode = "TRY"
const CurrencyTWD CurrencyCode = "TWD"
const CurrencyUAH CurrencyCode = "UAH"
const CurrencyUSD CurrencyCode = "USD"
const CurrencyUYU CurrencyCode = "UYU"
const CurrencyVND CurrencyCode = "VND"
const CurrencyZAR CurrencyCode = "ZAR"

var Currencies = map[CurrencyCode]string{
	CurrencyAED: "United Arab Emirates Dirham",
	CurrencyARS: "Argentine Peso",
	CurrencyAUD: "Australian Dollars",
	CurrencyBRL: "Brazilian Reals",
	CurrencyCAD: "Canadian Dollars",
	CurrencyCHF: "Swiss Francs",
	CurrencyCLP: "Chilean Peso",
	CurrencyCNY: "Chinese Renminbi (yuan)",
	CurrencyCOP: "Colombian Peso",
	CurrencyCRC: "Costa Rican Colón",
	CurrencyEUR: "European Union Euro",
	CurrencyGBP: "United Kingdom Pound",
	CurrencyHKD: "Hong Kong Dollar",
	CurrencyILS: "Israeli New Shekel",
	CurrencyIDR: "Indonesian Rupiah",
	CurrencyINR: "Indian Rupee",
	CurrencyJPY: "Japanese Yen",
	CurrencyKRW: "South Korean Won",
	CurrencyKWD: "Kuwaiti Dinar",
	CurrencyKZT: "Kazakhstani Tenge",
	CurrencyMXN: "Mexican Peso",
	CurrencyMYR: "Malaysian Ringgit",
	CurrencyNOK: "Norwegian Krone",
	CurrencyNZD: "New Zealand Dollar",
	CurrencyPEN: "Peruvian Nuevo Sol",
	CurrencyPHP: "Philippine Peso",
	CurrencyPLN: "Polish Złoty",
	CurrencyQAR: "Qatari Riyal",
	CurrencyRUB: "Russian Rouble",
	CurrencySAR: "Saudi Riyal",
	CurrencySGD: "Singapore Dollar",
	CurrencyTHB: "Thai Baht",
	CurrencyTRY: "Turkish Lira",
	CurrencyTWD: "New Taiwan Dollar",
	CurrencyUAH: "Ukrainian Hryvnia",
	CurrencyUSD: "United States Dollar",
	CurrencyUYU: "Uruguayan Peso",
	CurrencyVND: "Vietnamese Dong",
	CurrencyZAR: "South African Rand",
}
