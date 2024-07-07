// BSD 3-Clause No Military License
// Copyright © 2024, Piotr Bajdek. All Rights Reserved.

// gccgo -Ofast -march=native strata.go -o strata

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Użycie: strata <prawdopodobieństwo_straty_w_kwartale>")
		os.Exit(1)
	}

	probKwartal, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil || probKwartal < 0 || probKwartal > 1 {
		fmt.Println("Błąd: Prawdopodobieństwo musi być liczbą z zakresu 0-1")
		os.Exit(1)
	}

	probRok := 1 - math.Pow(1-probKwartal, 4)
	fmt.Printf("Prawdopodobieństwo wystąpienia kwartalnej straty w skali rocznej: %.2f%%\n", probRok*100)

	maxProb := 0.0
	maxKwartaly := 0
	for k := 0; k <= 4; k++ {
		prob := prawdopodobienstwoDwumianowe(4, k, probKwartal)
		if prob > maxProb {
			maxProb = prob
			maxKwartaly = k
		}
	}

	fmt.Printf("Najbardziej prawdopodobna liczba kwartałów ze stratą: %d\n", maxKwartaly)
}

func prawdopodobienstwoDwumianowe(n, k int, p float64) float64 {
	return float64(wspolczynnikDwumianowy(n, k)) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

func wspolczynnikDwumianowy(n, k int) int {
	return silnia(n) / (silnia(k) * silnia(n-k))
}

func silnia(n int) int {
	if n <= 1 {
		return 1
	}
	return n * silnia(n-1)
}
