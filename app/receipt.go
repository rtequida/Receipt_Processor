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

func GenerateId(receipt api.Receipt) string {
	return uuid.NewString()
}

func alphanumeric_count(retailer string) int {
	//One point for every alphanumeric character in the retailer name.
	count := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}

func is_round_dollar_amount(total string) int {
	//50 points if the total is a round dollar amount with no cents.
	total_float, err := strconv.ParseFloat(total, 64)
	if err == nil && math.Mod(total_float, 1) == 0 {
		return 50
	}
	return 0
}

func is_mult_of_25(total string) int {
	//25 points if the total is a multiple of 0.25.
	total_float, err := strconv.ParseFloat(total, 64)
	if err == nil && math.Mod(total_float, 0.25) == 0 {
		return 25
	}
	return 0
}

func for_every_pair(items []api.Item) int {
	//5 points for every two items on the receipt.
	count := 0
	count += 5 * (len(items) / 2)
	return count

}

func is_trimmed_length_mult_of_3(items []api.Item) int {
	//If the trimmed length of the item description is a multiple of 3,
	//multiply the price by 0.2 and round up to the nearest integer.
	//The result is the number of points earned.
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

func day_is_odd(purchaseDate openapi_types.Date) int {
	//6 points if the day in the purchase date is odd.
	day := purchaseDate.Day()
	if day%2 == 1 {
		return 6
	}
	return 0
}

func check_if_between_time(purchaseTime string) int {
	//10 points if the time of purchase is after 2:00pm and before 4:00pm.
	timeFormat := "15:04"
	parsedTime, err := time.Parse(timeFormat, purchaseTime)
	after, _ := time.Parse(timeFormat, "14:00")
	before, _ := time.Parse(timeFormat, "16:00")
	if err == nil && parsedTime.After(after) && parsedTime.Before(before) {
		return 10
	}
	return 0
}

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
