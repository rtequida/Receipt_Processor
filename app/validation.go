package app

import (
	"regexp"
	"strconv"

	"github.com/rtequida/Receipt_Processor/api"
)

// ValidateReceipt ensures that all required fields in the passed in receipt are included
// as well verifying that each field meets its defined pattern.
func ValidateReceipt(receipt api.Receipt) (bool, string) {
	if receipt.Retailer == "" {
		return false, "Retailer is missing."
	}
	// Retailer should only contain alphanumeric characters, whitespaces, hyphens, and ampersands
	// and should have at least 1 character.
	if match, err := regexp.MatchString(`^[\w\s\-&]+$`, receipt.Retailer); err != nil || !match {
		return false, "Retailer name is not allowed."
	}
	if receipt.PurchaseDate.Time.IsZero() {
		return false, "Purchase Date is missing."
	}
	if receipt.PurchaseTime == "" {
		return false, "Purchase Time is missing."
	}
	// PurchaseTime should be in the format of '13:00' ranging from '00:00' to '23:59'
	if match, err := regexp.MatchString(`^(?:[01]\d|2[0-3]):[0-5]\d$`, receipt.PurchaseTime); err != nil || !match {
		return false, "Purchase Time name is not allowed."
	}
	if len(receipt.Items) == 0 {
		return false, "Items is missing."
	}
	if receipt.Total == "" {
		return false, "Total is missing."
	}
	// Total should be a float with two decimal places
	if match, err := regexp.MatchString(`^\d+\.\d{2}$`, receipt.Total); err != nil || !match {
		return false, "Total is not allowed."
	}

	// Compiled regex for item's short description. Short description should only contain
	// alphanumeric characters, whitespaces, and hyphens
	re_short_description := regexp.MustCompile(`^[\w\s\-]+$`)
	// Complied regex for item's price. Price should be a float with two decimal places
	re_price := regexp.MustCompile(`^\d+\.\d{2}$`)
	for i, item := range receipt.Items {
		if item.ShortDescription == "" {
			return false, "Short Description for Item[" + strconv.Itoa(i) + "] is missing."
		}
		if !re_short_description.MatchString(item.ShortDescription) {
			return false, "Short Description for Item[" + strconv.Itoa(i) + "] name is not allowed."
		}
		if item.Price == "" {
			return false, "Price for Item[" + strconv.Itoa(i) + "] is missing."
		}
		if !re_price.MatchString(item.Price) {
			return false, "Price for Item[" + strconv.Itoa(i) + "] name is not allowed."
		}
	}
	return true, "No Errors."
}

// ValidateID ensures the passed ID contains 1 or more non-whitespace characters
func ValidateID(id string) (bool, string) {
	if match, err := regexp.MatchString(`^\S+$`, id); err != nil || !match {
		return false, "ID is not allowed."
	}
	return true, "No Errors."
}
