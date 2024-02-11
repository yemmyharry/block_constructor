package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MempoolTransaction struct {
	Txid    string
	Fee     int
	Weight  int
	Parents map[string]bool
}

type Block []MempoolTransaction

func parseMempoolCSV(filename string) (Block, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions Block

	for _, line := range lines {
		txid := line[0]
		fee, _ := strconv.Atoi(line[1])
		weight, _ := strconv.Atoi(line[2])
		parents := make(map[string]bool)

		if len(line) > 3 && len(line[3]) != 0 {
			parentTxids := strings.Split(line[3], ";")
			for _, parentTxid := range parentTxids {
				parents[parentTxid] = true
			}
		}

		transactions = append(transactions, MempoolTransaction{
			Txid:    txid,
			Fee:     fee,
			Weight:  weight,
			Parents: parents,
		})
	}

	return transactions, nil
}

func calculateMaxFeeBlock(transactions []MempoolTransaction) []string {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Fee > transactions[j].Fee
	})

	var includedTransactions []string
	includedSet := make(map[string]bool)
	totalWeight := 0

	for _, transaction := range transactions {

		parentsIncluded := true
		for parent := range transaction.Parents {
			if !includedSet[parent] {
				parentsIncluded = false
				break
			}
		}

		if parentsIncluded {
			includedSet[transaction.Txid] = true
			includedTransactions = append(includedTransactions, transaction.Txid)
			totalWeight += transaction.Weight

			if totalWeight >= 4000000 {
				break
			}
		}
	}

	return includedTransactions
}

func main() {
	mempool, err := parseMempoolCSV("mempool.csv")
	if err != nil {
		fmt.Println("Error when parsing file :", err.Error())
		return
	}

	maxFeeBlock := calculateMaxFeeBlock(mempool)

	for _, txid := range maxFeeBlock {
		fmt.Println(txid)
	}
}
