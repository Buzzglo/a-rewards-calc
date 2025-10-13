package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Date        time.Time
	Description string
	Amount      float64
	Balance     float64
	AccountType string
	Retailer    string
	Basket      string
}

const (
	SpendCapPerCategory  = 3000.0
	RewardRate           = 0.30
	MaxRewardPerCategory = SpendCapPerCategory * RewardRate
)

var basketRetailers = map[string][]string{
	"Grocery":      {"WOOLWORTHS", "FOOD LOVERS", "FLM", "PICK N PAY", "PNP", "CHECKERS"},
	"HealthBeauty": {"DISCHEM", "DIS-CHEM"},
	"Fuel":         {"SASOL"},
}

var totalSpendThresholds = []struct {
	Min    float64
	Max    float64
	Reward float64
}{
	{0, 9999.99, 0},
	{10000, 19999.99, 150},
	{20000, 29999.99, 300},
	{30000, 39999.99, 450},
	{40000, 49999.99, 800},
	{50000, 1e12, 1500},
}

func main() {
	chequeFile := "cheque.csv"
	creditFile := "credit.csv"

	chequeTxs, err := parseCSV(chequeFile, "Cheque")
	if err != nil {
		fmt.Println("Error reading cheque file:", err)
		return
	}

	creditTxs, err := parseCSV(creditFile, "Credit")
	if err != nil {
		fmt.Println("Error reading credit file:", err)
		return
	}

	// Rewards cycle: 16 June 2025 to 15 July 2025
	cycleStart := time.Date(2025, 9, 16, 0, 0, 0, 0, time.Local)
	cycleEnd := time.Date(2025, 10, 15, 23, 59, 59, 0, time.Local)

	chequeTxs = filterByDateRange(chequeTxs, cycleStart, cycleEnd)
	creditTxs = filterByDateRange(creditTxs, cycleStart, cycleEnd)

	// Tag credit transactions with retailer and basket
	for i := range creditTxs {
		creditTxs[i].Retailer = extractRetailer(creditTxs[i].Description)
		creditTxs[i].Basket = assignBasket(creditTxs[i].Retailer)
	}

	// Tag cheque transactions with retailer
	for i := range chequeTxs {
		chequeTxs[i].Retailer = extractRetailer(chequeTxs[i].Description)
	}

	// Calculate digital voucher reward (only "DIGITAL VOUCHERS")
	voucherSpend := sumDigitalVoucherSpend(chequeTxs)
	if voucherSpend > SpendCapPerCategory {
		voucherSpend = SpendCapPerCategory
	}
	chequeReward := voucherSpend * RewardRate

	// Credit Card Basket Spend
	basketSpend := map[string]float64{}
	for _, tx := range creditTxs {
		if tx.Amount < 0 {
			basketSpend[tx.Basket] += -tx.Amount
		}
	}

	qualifyingBaskets := []string{"Grocery", "HealthBeauty", "Fuel"}
	basketRewards := map[string]float64{}
	totalBasketReward := 0.0

	for _, basket := range qualifyingBaskets {
		spend := basketSpend[basket]
		if spend > SpendCapPerCategory {
			spend = SpendCapPerCategory
		}
		reward := spend * RewardRate
		basketRewards[basket] = reward
		totalBasketReward += reward
	}

	// Total Spend Reward
	totalCreditSpend := sumSpendAmounts(creditTxs)
	totalSpendReward := rewardForTotalSpend(totalCreditSpend)

	// Output
	fmt.Println("==== Rewards Estimate ====")
	fmt.Printf("Digital Voucher Reward (Cheque Account): R %.2f\n", chequeReward)

	fmt.Println("\nCredit Card Basket Rewards:")
	for _, basket := range qualifyingBaskets {
		if reward, ok := basketRewards[basket]; ok {
			fmt.Printf(" - %s: R %.2f\n", basket, reward)
		}
	}

	fmt.Printf("\nCredit Card Total Spend: R %.2f\n", totalCreditSpend)
	fmt.Printf("Total Spend Reward: R %.2f\n", totalSpendReward)

	grandTotal := chequeReward + totalBasketReward + totalSpendReward
	fmt.Printf("\n*** Total Estimated Rewards: R %.2f ***\n", grandTotal)
}

// === Helpers ===

func parseCSV(file string, accountType string) ([]Transaction, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	_, _ = reader.Read() // skip header

	var txs []Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 4 {
			continue
		}

		date, err := time.Parse("20060102", strings.TrimSpace(record[0]))
		if err != nil {
			continue
		}

		amount, err := strconv.ParseFloat(strings.ReplaceAll(record[2], ",", ""), 64)
		if err != nil {
			continue
		}

		balance, _ := strconv.ParseFloat(strings.ReplaceAll(record[3], ",", ""), 64)

		txs = append(txs, Transaction{
			Date:        date,
			Description: record[1],
			Amount:      amount,
			Balance:     balance,
			AccountType: accountType,
		})
	}
	return txs, nil
}

func filterByDateRange(txs []Transaction, start, end time.Time) []Transaction {
	var filtered []Transaction
	for _, tx := range txs {
		if !tx.Date.Before(start) && !tx.Date.After(end) {
			filtered = append(filtered, tx)
		}
	}
	return filtered
}

func extractRetailer(desc string) string {
	return strings.ToUpper(desc)
}

func assignBasket(retailer string) string {
	for basket, names := range basketRetailers {
		for _, name := range names {
			if strings.Contains(retailer, name) {
				return basket
			}
		}
	}
	return "Others"
}

func sumDigitalVoucherSpend(txs []Transaction) float64 {
	total := 0.0
	for _, tx := range txs {
		if tx.Amount < 0 && strings.Contains(strings.ToUpper(tx.Description), "DIGITAL VOUCHERS") {
			total += -tx.Amount
		}
	}
	return total
}

func sumSpendAmounts(txs []Transaction) float64 {
	total := 0.0
	for _, tx := range txs {
		if tx.Amount < 0 {
			total += -tx.Amount
		}
	}
	return total
}

func rewardForTotalSpend(spend float64) float64 {
	for _, tier := range totalSpendThresholds {
		if spend >= tier.Min && spend <= tier.Max {
			return tier.Reward
		}
	}
	return 0
}
