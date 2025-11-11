package isdayoff

// CountryCode type
type CountryCode string

const (
	// CountryCodeBelarus BY - Белоруссия
	CountryCodeBelarus CountryCode = "by"
	// CountryCodeKazakhstan KZ - Казахстан
	CountryCodeKazakhstan CountryCode = "kz"
	// CountryCodeRussia RU - Россия (по умолчанию)
	CountryCodeRussia CountryCode = "ru"
	// CountryCodeUkraine UA - Украина
	CountryCodeUkraine CountryCode = "ua"
	// CountryCodeUSA USA - США (расширенный код)
	CountryCodeUSA CountryCode = "us"
	// CountryCodeUzbekistan UZ - Узбекистан (расширенный код)
	CountryCodeUzbekistan CountryCode = "uz"
	// CountryCodeTurkey TR - Турция (расширенный код)
	CountryCodeTurkey CountryCode = "tr"
)

// DayType type
type DayType string

// YearType type
type YearType string

// ErrorCode type
type ErrorCode string

const (
	// DayTypeWorking working day
	DayTypeWorking DayType = "0"
	// DayTypeNonWorking non working day
	DayTypeNonWorking DayType = "1"
	// DayTypeHalfHoliday half holiday
	DayTypeHalfHoliday DayType = "2"
	// DayTypeWorkingCovid working day for Covid
	DayTypeWorkingCovid DayType = "4"

	// YearTypeNotLeap leap year
	YearTypeNotLeap YearType = "0"
	// YearTypeLeap non leap year
	YearTypeLeap YearType = "1"

	// ErrorCodeWrongDate wrong date err
	ErrorCodeWrongDate ErrorCode = "100"
	// ErrorCodeNotFound not found err
	ErrorCodeNotFound ErrorCode = "101"
	// ErrorCodeInternalError internal error
	ErrorCodeInternalError ErrorCode = "199"
)
