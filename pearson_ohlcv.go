// BSD 3-Clause No Military License
// Copyright © 2024, Piotr Bajdek. All Rights Reserved.

// gccgo -Ofast -march=native pearson_ohlcv.go -o pearson_ohlcv

package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Użycie: program <plik1.csv> <plik2.csv>")
		os.Exit(1)
	}

	columns := []struct {
		index int
		name  string
	}{
		{2, "O"},
		{3, "H"},
		{4, "L"},
		{5, "C"},
		{6, "V"},
	}

	fmt.Println("Obliczone korelacje Pearsona:")
	for _, col := range columns {
		seq1 := readCSVColumn(os.Args[1], col.index)
		seq2 := readCSVColumn(os.Args[2], col.index)

		if len(seq1) != len(seq2) {
			fmt.Printf("Błąd: Sekwencje dla kolumny %s mają różne długości\n", col.name)
			continue
		}

		correlation := calculateCorrelation(seq1, seq2)

		var formattedCorrelation string
		if correlation >= 0 {
			formattedCorrelation = fmt.Sprintf(" %.6f", correlation)
		} else {
			formattedCorrelation = fmt.Sprintf("%.6f", correlation)
		}

		fmt.Printf("Współczynnik korelacji dla %s: %s\n", col.name, formattedCorrelation)
	}
}

func readCSVColumn(filename string, columnIndex int) []float64 {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Błąd przy otwieraniu pliku %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Błąd przy czytaniu pliku CSV %s: %v\n", filename, err)
		os.Exit(1)
	}

	var sequence []float64
	for _, record := range records {
		if len(record) > columnIndex {
			value, err := strconv.ParseFloat(record[columnIndex], 64)
			if err != nil {
				fmt.Printf("Błąd przy konwersji wartości %s: %v\n", record[columnIndex], err)
				continue
			}
			sequence = append(sequence, value)
		}
	}
	return sequence
}

func calculateCorrelation(x, y []float64) float64 {
	n := float64(len(x))
	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))

	return numerator / denominator
}
