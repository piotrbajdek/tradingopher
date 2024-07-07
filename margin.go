// BSD 3-Clause No Military License
// Copyright © 2024, Piotr Bajdek. All Rights Reserved.

// gccgo -Ofast -march=native margin.go -o margin

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instrument struct {
	Nazwa  string
	Depozyt float64
}

type Konfiguracja struct {
	AktualneEquity float64
	AktualnyMargin float64
	Instrumenty    []Instrument
}

func wczytajKonfiguracje(nazwaPliku string) (Konfiguracja, error) {
	plik, err := os.Open(nazwaPliku)
	if err != nil {
		return Konfiguracja{}, err
	}
	defer plik.Close()

	var konfiguracja Konfiguracja
	skaner := bufio.NewScanner(plik)
	numerLinii := 0

	for skaner.Scan() {
		linia := strings.TrimSpace(skaner.Text())
		if linia == "" || strings.HasPrefix(linia, "#") {
			continue
		}

		czesci := strings.SplitN(linia, "=", 2)
		if len(czesci) != 2 {
			return Konfiguracja{}, fmt.Errorf("Nieprawidłowa linia konfiguracji: %s", linia)
		}

		klucz := strings.TrimSpace(czesci[0])
		wartosc := strings.TrimSpace(czesci[1])

		switch numerLinii {
		case 0:
			konfiguracja.AktualneEquity, err = strconv.ParseFloat(wartosc, 64)
		case 1:
			konfiguracja.AktualnyMargin, err = strconv.ParseFloat(wartosc, 64)
		default:
			instrument := Instrument{
				Nazwa:  klucz,
				Depozyt: parsujFloat(wartosc),
			}
			konfiguracja.Instrumenty = append(konfiguracja.Instrumenty, instrument)
		}

		if err != nil {
			return Konfiguracja{}, fmt.Errorf("Błąd parsowania wartości w linii %d: %v", numerLinii+1, err)
		}

		numerLinii++
	}

	if err := skaner.Err(); err != nil {
		return Konfiguracja{}, err
	}

	return konfiguracja, nil
}

func parsujFloat(s string) float64 {
	wartosc, _ := strconv.ParseFloat(s, 64)
	return wartosc
}

func obliczNowyPoziomMargin(konfiguracja Konfiguracja, nowyInstrument string, rozmiarNowejPozycji float64) (float64, error) {
	var depozytInstrumentu float64
	for _, inst := range konfiguracja.Instrumenty {
		if inst.Nazwa == nowyInstrument {
			depozytInstrumentu = inst.Depozyt
			break
		}
	}

	if depozytInstrumentu == 0 {
		return 0, fmt.Errorf("Nie znaleziono instrumentu: %s", nowyInstrument)
	}

	nowyMargin := konfiguracja.AktualnyMargin + (depozytInstrumentu * rozmiarNowejPozycji)
	nowyPoziomMargin := (konfiguracja.AktualneEquity / nowyMargin) * 100

	return nowyPoziomMargin, nil
}

func main() {
	konfiguracja, err := wczytajKonfiguracje("margin.conf")
	if err != nil {
		fmt.Printf("Błąd odczytu konfiguracji: %v\n", err)
		return
	}

	var instrument string
	var rozmiarPozycji float64
	var stopLoss float64

	fmt.Print("Podaj nazwę instrumentu: ")
	fmt.Scanln(&instrument)

	fmt.Print("Podaj rozmiar pozycji w walucie bazowej: ")
	fmt.Scanln(&rozmiarPozycji)

	fmt.Print("Podaj wartość Stop Loss (jako kwotę straty): ")
	fmt.Scanln(&stopLoss)

	nowyMargin, err := obliczNowyMargin(konfiguracja, instrument, rozmiarPozycji)
	if err != nil {
		fmt.Printf("Błąd: %v\n", err)
		return
	}

	nowyPoziomMargin := (konfiguracja.AktualneEquity / nowyMargin) * 100

	noweEquityNaStopLoss := konfiguracja.AktualneEquity - stopLoss

	poziomMarginNaStopLoss := (noweEquityNaStopLoss / nowyMargin) * 100

	fmt.Printf("Nowy Poziom Margin: %.2f%%\n", nowyPoziomMargin)

	if poziomMarginNaStopLoss < 200.0 {
		fmt.Printf("Poziom Margin na Stop Loss: \033[31m%.2f%%\033[0m\n", poziomMarginNaStopLoss)
	} else {
		fmt.Printf("Poziom Margin na Stop Loss: %.2f%%\n", poziomMarginNaStopLoss)
	}
}

func obliczNowyMargin(konfiguracja Konfiguracja, nowyInstrument string, rozmiarNowejPozycji float64) (float64, error) {
	var depozytInstrumentu float64
	for _, inst := range konfiguracja.Instrumenty {
		if inst.Nazwa == nowyInstrument {
			depozytInstrumentu = inst.Depozyt
			break
		}
	}
	if depozytInstrumentu == 0 {
		return 0, fmt.Errorf("nie znaleziono instrumentu: %s", nowyInstrument)
	}
	nowyMargin := konfiguracja.AktualnyMargin + (depozytInstrumentu * rozmiarNowejPozycji)
	return nowyMargin, nil
}
