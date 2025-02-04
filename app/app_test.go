package app

import (
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rtequida/Receipt_Processor/api"
)

func TestGenerateId(t *testing.T) {
	receipt := api.Receipt{}
	id := GenerateId(receipt)
	_, err := uuid.Parse(id)
	if err != nil {
		t.Errorf("Generated Id: %v is invalid", id)
	}
}

func TestAlphanumbericCount(t *testing.T) {
	retailer_score_map := make(map[string]int)
	retailer_score_map["Just A Store"] = 10
	retailer_score_map["    Test &*Random--*&^Symbols"] = 17
	for retailer, score := range retailer_score_map {
		f_score := alphanumeric_count(retailer)
		if f_score != score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, score)
		}
	}
}

func TestIsRoundDollarAmount(t *testing.T) {
	total_score_map := make(map[string]int)
	total_score_map["test"] = 0
	total_score_map["10.00"] = 50
	total_score_map["10.100"] = 0
	total_score_map["10.10"] = 0
	total_score_map["10"] = 50
	for total, score := range total_score_map {
		f_score := is_round_dollar_amount(total)
		if f_score != score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, score)
		}
	}
}

func TestIsMultOf25(t *testing.T) {
	total_score_map := make(map[string]int)
	total_score_map["test"] = 0
	total_score_map["10.75"] = 25
	total_score_map["10.100"] = 0
	total_score_map["10.10"] = 0
	total_score_map["10"] = 25
	for total, score := range total_score_map {
		f_score := is_mult_of_25(total)
		if f_score != score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, score)
		}
	}
}

func TestForEveryPair(t *testing.T) {
	test_items := []struct {
		items []api.Item
		score int
	}{
		{[]api.Item{}, 0},
		{[]api.Item{{}, {}, {}}, 5},
		{[]api.Item{{}}, 0},
		{[]api.Item{{}, {}, {}, {}}, 10},
	}
	for _, ti := range test_items {
		f_score := for_every_pair(ti.items)
		if f_score != ti.score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, ti.score)
		}
	}
}

func TestIsTrimmedLengthMultOf3(t *testing.T) {
	test_items := []struct {
		items []api.Item
		score int
	}{
		{[]api.Item{{ShortDescription: "   123123     ", Price: "10.00"}, {ShortDescription: "testing", Price: "10000.25"}}, 2},
		{[]api.Item{{ShortDescription: "asdfg asdfgh", Price: "500.15"}}, 101},
		{[]api.Item{{}}, 0},
		{[]api.Item{{ShortDescription: "asd", Price: "2.00"}, {ShortDescription: "asd", Price: "100.00"}, {ShortDescription: "asdf", Price: "75.50"}}, 21},
	}
	for _, ti := range test_items {
		f_score := is_trimmed_length_mult_of_3(ti.items)
		if f_score != ti.score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, ti.score)
		}
	}
}

func TestDayIsOdd(t *testing.T) {
	test_dates := []struct {
		purchaseDate openapi_types.Date
		score        int
	}{
		{openapi_types.Date{Time: time.Date(2001, 05, 15, 0, 0, 0, 0, time.UTC)}, 6},
		{openapi_types.Date{Time: time.Date(2001, 05, 14, 0, 0, 0, 0, time.UTC)}, 0},
	}
	for _, td := range test_dates {
		f_score := day_is_odd(td.purchaseDate)
		if f_score != td.score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, td.score)
		}
	}
}

func TestCheckIfBetweenTime(t *testing.T) {
	time_score_map := make(map[string]int)
	time_score_map["test"] = 0
	time_score_map["14:01"] = 10
	time_score_map["14:00"] = 0
	time_score_map["15:59"] = 10
	time_score_map["16:00"] = 0
	for time, score := range time_score_map {
		f_score := check_if_between_time(time)
		if f_score != score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, score)
		}
	}
}

func TestGetPoints(t *testing.T) {
	test_receipts := []struct {
		receipt api.Receipt
		score   int
	}{
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, 106},
	}
	for _, tr := range test_receipts {
		f_score := GetPoints(tr.receipt)
		if f_score != tr.score {
			t.Errorf("Score generated = %v, expected value is = %v", f_score, tr.score)
		}
	}
}

func TestValidateReceipt(t *testing.T) {
	test_receipts := []struct {
		receipt     api.Receipt
		flag        bool
		err_message string
	}{
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, true, "No Errors."},
		{api.Receipt{PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Retailer is missing."},
		{api.Receipt{Retailer: "Test# World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Retailer name is not allowed."},
		{api.Receipt{Retailer: "Test World", PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Purchase Date is missing."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Purchase Time is missing."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:68", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Purchase Time name is not allowed."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{}}, false, "Items is missing."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Total is missing."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.0", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Total is not allowed."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{Price: "5.00"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Short Description for Item[0] is missing."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bol*ogna", Price: "48.00"}}}, false, "Short Description for Item[1] name is not allowed."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.00"}, {ShortDescription: "bologna"}}}, false, "Price for Item[1] is missing."},
		{api.Receipt{Retailer: "Test World", PurchaseDate: openapi_types.Date{Time: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)}, PurchaseTime: "15:48", Total: "53.00", Items: []api.Item{{ShortDescription: "lemons", Price: "5.500"}, {ShortDescription: "bologna", Price: "48.00"}}}, false, "Price for Item[0] name is not allowed."},
	}
	for _, tr := range test_receipts {
		f_flag, f_err_message := ValidateReceipt(tr.receipt)
		if f_flag != tr.flag || f_err_message != tr.err_message {
			t.Errorf("Flag generated = %v, expected value is = %v. Message generated = %s, expected value is = %s", f_flag, tr.flag, f_err_message, tr.err_message)
		}
	}
}

func TestValidateID(t *testing.T) {
	test_id := []struct {
		id          string
		flag        bool
		err_message string
	}{
		{id: "fc465ee8-29dc-47c0-adba-2845126b4d05", flag: true, err_message: "No Errors."},
		{id: "asdf asdf", flag: false, err_message: "ID is not allowed."},
	}
	for _, ti := range test_id {
		f_flag, f_err_message := ValidateID(ti.id)
		if f_flag != ti.flag || f_err_message != ti.err_message {
			t.Errorf("Flag generated = %v, expected value is = %v. Message generated = %s, expected value is = %s", f_flag, ti.flag, f_err_message, ti.err_message)
		}
	}
}
