// BSD 3-Clause No Military License
// Copyright © 2024, Piotr Bajdek. All Rights Reserved.

// gccgo -Ofast -march=native theil_sen.go -o theil_sen

package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type DataPoint struct {
	Date   string
	Time   string
	O, H, L, C, V float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input file selected!")
		return
	}

	arg1 := os.Args[1]	

	fileName := filepath.Base(arg1)
	fmt.Printf("8-Day Theil–Sen Estimator for %s\n", fileName)

	file, err := os.Open(arg1)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	if len(records) < 8 {
		fmt.Println("Not enough records in the CSV file!")
		return
	}

	var data []DataPoint
	for _, record := range records {
		if len(record) != 7 {
			fmt.Println("Invalid record format")
			return
		}
		o, _ := strconv.ParseFloat(record[2], 64)
		h, _ := strconv.ParseFloat(record[3], 64)
		l, _ := strconv.ParseFloat(record[4], 64)
		c, _ := strconv.ParseFloat(record[5], 64)
		v, _ := strconv.ParseFloat(record[6], 64)
		data = append(data, DataPoint{record[0], record[1], o, h, l, c, v})
	}

	lastRecord := data[len(data)-1]
	fmt.Printf("Date: %s, Time: %s\n", lastRecord.Date, lastRecord.Time)

	predictAndPrint("O", getColumn(data, func(d DataPoint) float64 { return d.O }))
	predictAndPrint("H", getColumn(data, func(d DataPoint) float64 { return d.H }))
	predictAndPrint("L", getColumn(data, func(d DataPoint) float64 { return d.L }))
	predictAndPrint("C", getColumn(data, func(d DataPoint) float64 { return d.C }))
	predictAndPrint("V", getColumn(data, func(d DataPoint) float64 { return d.V }))
}

func getColumn(data []DataPoint, accessor func(DataPoint) float64) []float64 {
	result := make([]float64, len(data))
	for i, d := range data {
		result[i] = accessor(d)
	}
	return result
}

func predictAndPrint(columnName string, sequence []float64) {
	seqLen := len(sequence)
	startIdx := int(math.Max(0, float64(seqLen-8)))
	var slopes []float64
	for i := startIdx; i < seqLen-1; i++ {
		for j := i + 1; j < seqLen; j++ {
			slopes = append(slopes, (sequence[j]-sequence[i])/float64(j-i))
		}
	}

	sort.Float64s(slopes)
	b := median(slopes)
	yPred := b*float64(seqLen-1) + sequence[0]

	fmt.Printf("Predicted value for %s = %.6f\n", columnName, yPred)
}

func median(array []float64) float64 {
	n := len(array)
	if n%2 == 1 {
		return array[n/2]
	}
	return (array[n/2-1] + array[n/2]) / 2.0
}
