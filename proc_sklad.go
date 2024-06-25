// BSD 3-Clause No Military License
// Copyright © 2024, Piotr Bajdek. All Rights Reserved.

// gccgo -Ofast -march=native proc_sklad.go -o proc_sklad

package main

import (
    "fmt"
    "math"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Użycie: proc_sklad <procent_zwrotu_za_okres> <liczba_okresów>")
        os.Exit(1)
    }

    periodReturn, err := strconv.ParseFloat(os.Args[1], 64)
    if err != nil {
        fmt.Println("Błąd: Nieprawidłowy format procentu zwrotu za okres")
        os.Exit(1)
    }

    periods, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("Błąd: Nieprawidłowy format liczby okresów")
        os.Exit(1)
    }

    periodReturnDecimal := periodReturn / 100
    totalGrowth := math.Pow(1+periodReturnDecimal, float64(periods)) - 1

    fmt.Printf("Zwrot za okres: %.2f%%\n", periodReturn)
    fmt.Printf("Liczba okresów: %d\n", periods)
    fmt.Printf("Całkowity wzrost: %.2f%%\n", totalGrowth*100)
}
