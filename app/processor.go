// Package app implements the logic used to generate ids and determine point values for receipts.
// It also implements validation for incoming receipts to make sure they match the defined schema.
package app

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rtequida/Receipt_Processor/api"
)

// GenerateId generates and returns a new ID for each receipt.
func GenerateId(receipt api.Receipt) string {
	return uuid.NewString()
}

// alphanumeric_count determines how many alphanumeric characters are within the retailer
// string and return 1 point for each of those characters. Returns 0 otherwise.
func alphanumeric_count(retailer string) int {
	count := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}

// is_round_dollar_amount determines if the total of the receipt is a round dollar amount
// with no cents and returns 50 points to the receipt if it does. Returns 0 otherwise.
func is_round_dollar_amount(total string) int {
	total_float, err := strconv.ParseFloat(total, 64)
	if err == nil && math.Mod(total_float, 1) == 0 {
		return 50
	}
	return 0
}

// is_mult_of_25 determines if the total of the receipt is a multiple of 0.25 and returns
// 25 points if so. Returns 0 otherwise.
func is_mult_of_25(total string) int {
	total_float, err := strconv.ParseFloat(total, 64)
	if err == nil && math.Mod(total_float, 0.25) == 0 {
		return 25
	}
	return 0
}

// for_every_pair returns 5 points for every two items found in the receipt's items array.
func for_every_pair(items []api.Item) int {
	//5 points for every two items on the receipt.
	count := 0
	count += 5 * (len(items) / 2)
	return count

}

// is_trimmed_length_mult_of_3 for each item in the receipt's item array, trims any beginning
// or trailing whitespace and then determines if the remaining characters have a length that
// is divisible 3. If it does, multiply the price of the item by 0.2 and round up to the nearst
// whole number. Otherwise award 0 points for that item.
func is_trimmed_length_mult_of_3(items []api.Item) int {
	count := 0
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				count += int(math.Ceil(price * 0.2))
			}
		}
	}
	return count
}

// day_is_odd determines if the purchase day of the receipt is odd. Returns 6 points if the day
// is odd and returns 0 otherwise.
func day_is_odd(purchaseDate openapi_types.Date) int {
	day := purchaseDate.Day()
	if day%2 == 1 {
		return 6
	}
	return 0
}

// check_if_between_time determines if the purchase time of the receipt is after 2:00pm and
// before 4:00pm. If it falls between those times returns 10 points or returns 0 points otherwise.
func check_if_between_time(purchaseTime string) int {
	timeFormat := "15:04"
	parsedTime, err := time.Parse(timeFormat, purchaseTime)
	// Generate new time objects to compare with.
	after, _ := time.Parse(timeFormat, "14:00")
	before, _ := time.Parse(timeFormat, "16:00")
	if err == nil && parsedTime.After(after) && parsedTime.Before(before) {
		return 10
	}
	return 0
}

// GetPoints calls all the helper functions that determine the points for certain sections of
// the receipt and tallies them together to get the overall receipt score.
func GetPoints(receipt api.Receipt) int {
	points := 0

	points += alphanumeric_count(receipt.Retailer)
	points += is_round_dollar_amount(receipt.Total)
	points += is_mult_of_25(receipt.Total)
	points += for_every_pair(receipt.Items)
	points += is_trimmed_length_mult_of_3(receipt.Items)
	points += day_is_odd(receipt.PurchaseDate)
	points += check_if_between_time(receipt.PurchaseTime)

	return points
}
