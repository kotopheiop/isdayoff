package isdayoff

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestIsLeap(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		expected bool
	}{
		{
			name:     "2020 is leap year",
			year:     2020,
			expected: true,
		},
		{
			name:     "2021 is not leap year",
			year:     2021,
			expected: false,
		},
		{
			name:     "2024 is leap year",
			year:     2024,
			expected: true,
		},
		{
			name:     "1900 is not leap year (century)",
			year:     1900,
			expected: false,
		},
		{
			name:     "2000 is leap year (century leap)",
			year:     2000,
			expected: true,
		},
	}

	client := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leap, err := client.IsLeap(tt.year)
			if err != nil {
				t.Fatalf("IsLeap(%d) failed: %v", tt.year, err)
			}
			if leap != tt.expected {
				t.Errorf("IsLeap(%d) = %v, expected %v", tt.year, leap, tt.expected)
			}
		})
	}
}

func TestGetByYear(t *testing.T) {
	tests := []struct {
		name         string
		year         int
		expectedDays int
		description  string
	}{
		{
			name:         "2020 is leap year - 366 days",
			year:         2020,
			expectedDays: 366,
			description:  "Leap year should have 366 days",
		},
		{
			name:         "2021 is not leap year - 365 days",
			year:         2021,
			expectedDays: 365,
			description:  "Non-leap year should have 365 days",
		},
		{
			name:         "2024 is leap year - 366 days",
			year:         2024,
			expectedDays: 366,
			description:  "Leap year should have 366 days",
		},
	}

	client := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, err := client.GetBy(Params{Year: tt.year})
			if err != nil {
				t.Fatalf("GetBy(Year: %d) failed: %v", tt.year, err)
			}
			if len(days) != tt.expectedDays {
				t.Errorf("GetBy(Year: %d) returned %d days, expected %d. %s", tt.year, len(days), tt.expectedDays, tt.description)
			}
			// Проверяем, что все дни имеют валидный тип
			validTypes := map[DayType]bool{
				DayTypeWorking:      true,
				DayTypeNonWorking:   true,
				DayTypeHalfHoliday:  true,
				DayTypeWorkingCovid: true,
			}
			for i, day := range days {
				if !validTypes[day] {
					t.Errorf("GetBy(Year: %d) returned invalid day type at index %d: %v", tt.year, i, day)
				}
			}
		})
	}
}

func TestGetByDay(t *testing.T) {
	tests := []struct {
		name         string
		year         int
		month        time.Month
		day          int
		countryCode  CountryCode
		pre          bool
		covid        bool
		expectedDays int
		description  string
	}{
		{
			name:         "New Year 2020 in Kazakhstan",
			year:         2020,
			month:        time.January,
			day:          1,
			countryCode:  CountryCodeKazakhstan,
			pre:          false,
			covid:        false,
			expectedDays: 1,
			description:  "Should return single day for specific date",
		},
		{
			name:         "New Year 2021 in Russia",
			year:         2021,
			month:        time.January,
			day:          1,
			countryCode:  CountryCodeRussia,
			pre:          false,
			covid:        false,
			expectedDays: 1,
			description:  "Should return single day for specific date",
		},
		{
			name:         "May 1st 2024 in Belarus",
			year:         2024,
			month:        time.May,
			day:          1,
			countryCode:  CountryCodeBelarus,
			pre:          false,
			covid:        false,
			expectedDays: 1,
			description:  "Should return single day for specific date",
		},
	}

	client := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, err := client.GetBy(Params{
				Year:        tt.year,
				Month:       &tt.month,
				Day:         &tt.day,
				CountryCode: &tt.countryCode,
				Pre:         &tt.pre,
				Covid:       &tt.covid,
			})
			if err != nil {
				t.Fatalf("GetBy() failed: %v", err)
			}
			if len(days) != tt.expectedDays {
				t.Errorf("GetBy() returned %d days, expected %d. %s", len(days), tt.expectedDays, tt.description)
			}
			// Проверяем, что день имеет валидный тип
			if len(days) > 0 {
				validTypes := map[DayType]bool{
					DayTypeWorking:      true,
					DayTypeNonWorking:   true,
					DayTypeHalfHoliday:  true,
					DayTypeWorkingCovid: true,
				}
				if !validTypes[days[0]] {
					t.Errorf("GetBy() returned invalid day type: %v", days[0])
				}
				t.Logf("Day type for %d-%02d-%02d in %s: %v", tt.year, tt.month, tt.day, tt.countryCode, days[0])
			}
		})
	}
}

func TestGetByMonth(t *testing.T) {
	tests := []struct {
		name         string
		year         int
		month        time.Month
		countryCode  CountryCode
		expectedDays int
		description  string
	}{
		{
			name:         "January 2020 in Kazakhstan",
			year:         2020,
			month:        time.January,
			countryCode:  CountryCodeKazakhstan,
			expectedDays: 31,
			description:  "January has 31 days",
		},
		{
			name:         "February 2020 in Russia (leap year)",
			year:         2020,
			month:        time.February,
			countryCode:  CountryCodeRussia,
			expectedDays: 29,
			description:  "February in leap year has 29 days",
		},
		{
			name:         "February 2021 in Russia (non-leap year)",
			year:         2021,
			month:        time.February,
			countryCode:  CountryCodeRussia,
			expectedDays: 28,
			description:  "February in non-leap year has 28 days",
		},
	}

	client := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, err := client.GetBy(Params{
				Year:        tt.year,
				Month:       &tt.month,
				CountryCode: &tt.countryCode,
			})
			if err != nil {
				t.Fatalf("GetBy() failed: %v", err)
			}
			if len(days) != tt.expectedDays {
				t.Errorf("GetBy() returned %d days, expected %d. %s", len(days), tt.expectedDays, tt.description)
			}
			// Проверяем, что все дни имеют валидный тип
			validTypes := map[DayType]bool{
				DayTypeWorking:      true,
				DayTypeNonWorking:   true,
				DayTypeHalfHoliday:  true,
				DayTypeWorkingCovid: true,
			}
			for i, day := range days {
				if !validTypes[day] {
					t.Errorf("GetBy() returned invalid day type at index %d: %v", i, day)
				}
			}
			// Подсчитываем рабочие и нерабочие дни
			working := 0
			nonWorking := 0
			for _, day := range days {
				if day == DayTypeWorking || day == DayTypeWorkingCovid {
					working++
				} else if day == DayTypeNonWorking {
					nonWorking++
				}
			}
			t.Logf("Month %d/%d in %s: %d working days, %d non-working days", tt.month, tt.year, tt.countryCode, working, nonWorking)
		})
	}
}

func TestToday(t *testing.T) {
	client := New()

	tests := []struct {
		name        string
		countryCode CountryCode
		pre         bool
		covid       bool
		description string
	}{
		{
			name:        "Today in Kazakhstan",
			countryCode: CountryCodeKazakhstan,
			pre:         false,
			covid:       false,
			description: "Should return valid day type for today",
		},
		{
			name:        "Today in Russia",
			countryCode: CountryCodeRussia,
			pre:         false,
			covid:       false,
			description: "Should return valid day type for today",
		},
		{
			name:        "Today in Belarus",
			countryCode: CountryCodeBelarus,
			pre:         false,
			covid:       false,
			description: "Should return valid day type for today",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := client.Today(Params{
				CountryCode: &tt.countryCode,
				Pre:         &tt.pre,
				Covid:       &tt.covid,
			})
			if err != nil {
				t.Fatalf("Today() failed: %v", err)
			}
			if day == nil {
				t.Fatal("Today() returned nil day")
			}

			// Проверяем, что получили валидный тип дня
			validTypes := map[DayType]bool{
				DayTypeWorking:      true,
				DayTypeNonWorking:   true,
				DayTypeHalfHoliday:  true,
				DayTypeWorkingCovid: true,
			}
			if !validTypes[*day] {
				t.Errorf("Today() returned invalid day type: %v. %s", *day, tt.description)
			}

			dayName := map[DayType]string{
				DayTypeWorking:      "рабочий день",
				DayTypeNonWorking:   "нерабочий день",
				DayTypeHalfHoliday:  "сокращенный день",
				DayTypeWorkingCovid: "рабочий день (COVID)",
			}
			t.Logf("Сегодня в %s: %s (%v)", tt.countryCode, dayName[*day], *day)
		})
	}
}

func TestTomorrow(t *testing.T) {
	client := New()

	tests := []struct {
		name        string
		countryCode CountryCode
		pre         bool
		covid       bool
		description string
	}{
		{
			name:        "Tomorrow in Kazakhstan",
			countryCode: CountryCodeKazakhstan,
			pre:         false,
			covid:       false,
			description: "Should return valid day type for tomorrow",
		},
		{
			name:        "Tomorrow in Russia",
			countryCode: CountryCodeRussia,
			pre:         false,
			covid:       false,
			description: "Should return valid day type for tomorrow",
		},
		{
			name:        "Tomorrow in Ukraine",
			countryCode: CountryCodeUkraine,
			pre:         false,
			covid:       false,
			description: "Should return valid day type for tomorrow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := client.Tomorrow(Params{
				CountryCode: &tt.countryCode,
				Pre:         &tt.pre,
				Covid:       &tt.covid,
			})
			if err != nil {
				t.Fatalf("Tomorrow() failed: %v", err)
			}
			if day == nil {
				t.Fatal("Tomorrow() returned nil day")
			}

			// Проверяем, что получили валидный тип дня
			validTypes := map[DayType]bool{
				DayTypeWorking:      true,
				DayTypeNonWorking:   true,
				DayTypeHalfHoliday:  true,
				DayTypeWorkingCovid: true,
			}
			if !validTypes[*day] {
				t.Errorf("Tomorrow() returned invalid day type: %v. %s", *day, tt.description)
			}

			dayName := map[DayType]string{
				DayTypeWorking:      "рабочий день",
				DayTypeNonWorking:   "нерабочий день",
				DayTypeHalfHoliday:  "сокращенный день",
				DayTypeWorkingCovid: "рабочий день (COVID)",
			}
			t.Logf("Завтра в %s: %s (%v)", tt.countryCode, dayName[*day], *day)
		})
	}
}

func TestNewWithClient(t *testing.T) {
	t.Run("Create client with custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 10 * time.Second,
		}
		client := NewWithClient(customClient)
		if client == nil {
			t.Fatal("NewWithClient() returned nil")
		}
		if client.httpClient != customClient {
			t.Error("NewWithClient() did not use provided HTTP client")
		}
	})
}

func TestCountryCodes(t *testing.T) {
	tests := []struct {
		name        string
		countryCode CountryCode
		description string
	}{
		{"Belarus", CountryCodeBelarus, "BY"},
		{"Kazakhstan", CountryCodeKazakhstan, "KZ"},
		{"Russia", CountryCodeRussia, "RU"},
		{"Ukraine", CountryCodeUkraine, "UA"},
		{"USA", CountryCodeUSA, "US"},
		{"Uzbekistan", CountryCodeUzbekistan, "UZ"},
		{"Turkey", CountryCodeTurkey, "TR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.countryCode) == "" {
				t.Errorf("CountryCode for %s is empty", tt.name)
			}
			t.Logf("Country code %s: %s", tt.name, tt.countryCode)
		})
	}
}

func TestDayTypes(t *testing.T) {
	tests := []struct {
		name        string
		dayType     DayType
		description string
	}{
		{"Working", DayTypeWorking, "0"},
		{"NonWorking", DayTypeNonWorking, "1"},
		{"HalfHoliday", DayTypeHalfHoliday, "2"},
		{"WorkingCovid", DayTypeWorkingCovid, "4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.dayType) != tt.description {
				t.Errorf("DayType %s should be %s, got %s", tt.name, tt.description, tt.dayType)
			}
			t.Logf("Day type %s: %s", tt.name, tt.dayType)
		})
	}
}

func ExampleClient_IsLeap() {
	client := New()
	isLeap, err := client.IsLeap(2020)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("2020 is leap year: %v\n", isLeap)
	// Output: 2020 is leap year: true
}

func ExampleClient_GetBy() {
	client := New()
	month := time.January
	days, err := client.GetBy(Params{
		Year:  2024,
		Month: &month,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("January 2024 has %d days\n", len(days))
	// Output: January 2024 has 31 days
}

func TestGetByPeriod(t *testing.T) {
	tests := []struct {
		name        string
		date1       string
		date2       string
		countryCode CountryCode
		pre         bool
		covid       bool
		sixDayWeek  bool
		expectedMin int
		expectedMax int
		description string
	}{
		{
			name:        "New Year week 2024 in Russia",
			date1:       "20240101",
			date2:       "20240107",
			countryCode: CountryCodeRussia,
			pre:         false,
			covid:       false,
			sixDayWeek:  false,
			expectedMin: 7,
			expectedMax: 7,
			description: "Should return 7 days for a week",
		},
		{
			name:        "January 2024 in Kazakhstan",
			date1:       "20240101",
			date2:       "20240131",
			countryCode: CountryCodeKazakhstan,
			pre:         false,
			covid:       false,
			sixDayWeek:  false,
			expectedMin: 31,
			expectedMax: 31,
			description: "Should return 31 days for January",
		},
		{
			name:        "Short period in Belarus",
			date1:       "20240201",
			date2:       "20240205",
			countryCode: CountryCodeBelarus,
			pre:         false,
			covid:       false,
			sixDayWeek:  false,
			expectedMin: 5,
			expectedMax: 5,
			description: "Should return 5 days",
		},
	}

	client := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, err := client.GetByPeriod(tt.date1, tt.date2, Params{
				CountryCode: &tt.countryCode,
				Pre:         &tt.pre,
				Covid:       &tt.covid,
				SixDayWeek:  &tt.sixDayWeek,
			})
			if err != nil {
				t.Fatalf("GetByPeriod() failed: %v", err)
			}
			if len(days) < tt.expectedMin || len(days) > tt.expectedMax {
				t.Errorf("GetByPeriod() returned %d days, expected between %d and %d. %s", len(days), tt.expectedMin, tt.expectedMax, tt.description)
			}
			// Проверяем, что все дни имеют валидный тип
			validTypes := map[DayType]bool{
				DayTypeWorking:      true,
				DayTypeNonWorking:   true,
				DayTypeHalfHoliday:  true,
				DayTypeWorkingCovid: true,
			}
			for i, day := range days {
				if !validTypes[day] {
					t.Errorf("GetByPeriod() returned invalid day type at index %d: %v", i, day)
				}
			}
			t.Logf("Period %s to %s in %s: %d days", tt.date1, tt.date2, tt.countryCode, len(days))
		})
	}
}

func TestGetByWithSixDayWeek(t *testing.T) {
	client := New()
	sixDayWeek := true
	month := time.January
	countryCode := CountryCodeRussia

	days, err := client.GetBy(Params{
		Year:        2024,
		Month:       &month,
		CountryCode: &countryCode,
		SixDayWeek:  &sixDayWeek,
	})
	if err != nil {
		t.Fatalf("GetBy() with SixDayWeek failed: %v", err)
	}
	if len(days) != 31 {
		t.Errorf("GetBy() with SixDayWeek returned %d days, expected 31", len(days))
	}

	// Проверяем, что все дни имеют валидный тип
	validTypes := map[DayType]bool{
		DayTypeWorking:      true,
		DayTypeNonWorking:   true,
		DayTypeHalfHoliday:  true,
		DayTypeWorkingCovid: true,
	}
	for i, day := range days {
		if !validTypes[day] {
			t.Errorf("GetBy() with SixDayWeek returned invalid day type at index %d: %v", i, day)
		}
	}
	t.Logf("January 2024 in Russia with six-day week: %d days", len(days))
}

func TestAPIError(t *testing.T) {
	client := New()

	// Тест на обработку ошибки неправильной даты
	// Используем несуществующую дату
	invalidDate := "20240230" // 30 февраля не существует
	validDate := "20240228"

	_, err := client.GetByPeriod(invalidDate, validDate, Params{})
	if err == nil {
		t.Error("GetByPeriod() should return error for invalid date")
	} else {
		apiErr, ok := err.(*APIError)
		if ok {
			t.Logf("Got API error: %s (Code: %s, Status: %d)", apiErr.Message, apiErr.Code, apiErr.Status)
			if apiErr.Code != ErrorCodeWrongDate && apiErr.Code != ErrorCodeNotFound {
				t.Errorf("Expected ErrorCodeWrongDate or ErrorCodeNotFound, got %s", apiErr.Code)
			}
		} else {
			t.Logf("Got error (not APIError): %v", err)
		}
	}
}
