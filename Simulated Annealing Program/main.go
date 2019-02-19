package main

import (
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"path"
	"time"
)

func main() {
	// Melakukan deklarasi dan memanggil fungsi simulated anealing

	rand.Seed(time.Now().UnixNano()) // Digunakan untuk mendapatkan nilai random yang uniq
	var temperature, finalTemperature, alpha, x1, x2 float64
	var initialState, costSolutions []float64
	var xSolutions [][]float64

	temperature = 100.0
	finalTemperature = 0.000001
	alpha = 0.99999
	x1 = randFloats(-10, 10)
	x2 = randFloats(-10, 10)
	initialState = []float64{x1, x2}

	// Pemanggilan fungsi dan menampung hasil returnnya kedalan dua varible, untuk mendapatkan nilai x1, x2, dan nilai minimummnya
	xSolutions, costSolutions = simulatedAnnealing(initialState, temperature, finalTemperature, alpha)
	nilaiMinimum := costSolutions[len(costSolutions)-1]
	nilaiX := xSolutions[len(xSolutions)-1]
	nilaiX1 := nilaiX[0]
	nilaiX2 := nilaiX[1]
	akurasi := countAccuracy(nilaiMinimum)

	// Melakukan proses untuk rendering file html, untuk dijadikan tempat output program

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var filepath = path.Join("views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = map[string]interface{}{
			"nilai_minimum": nilaiMinimum,
			"nilai_x1":      nilaiX1,
			"nilai_x2":      nilaiX2,
			"akurasi":       akurasi,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	fmt.Println("Buka pada browser => localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func fungsi(x1, x2 float64) float64 {
	/*
		I.S : Telah terdefinisi fungsi berdasarkan soal yang akan didekati oleh metode simulated anealing
		F.S : mengembalikan nilai dari fungsi
	*/
	return (-((math.Sin(x1) * math.Cos(x2)) + ((4.0 / 5.0) * math.Exp(1-math.Sqrt(math.Pow(x1, 2)+math.Pow(x2, 2))))))
}

func probabilityAcceptance(newCost, currentCost, temperature float64) float64 {
	/*
		I.S : Telah terdefinisi parameter untuk dihitung
		F.S : Mengembalikan probability acceptance
	*/
	return math.Exp(-(newCost - currentCost) / temperature)
}

func simulatedAnnealing(initialState []float64, temperature, finalTemperature, alpha float64) ([][]float64, []float64) {
	/*
		I.S : Telah terdefinisi paremeter untuk dilakukan perhitungan dengan menggunakan metode simulated anealing
		F.S : Telah didapatkan nilai minimun, x1, dan x2
	*/
	var xSolutions [][]float64
	var costSolutions []float64
	var x1, x2, newCost, probAcc float64

	var currentState = initialState
	var currentCost = fungsi(initialState[0], initialState[1])

	var bestSolution = initialState
	var bestCost = currentCost

	for temperature > finalTemperature {
		var i = 0
		for i < 100 {
			x1 = randFloats(-10, 10)
			x2 = randFloats(-10, 10)
			newCost = fungsi(x1, x2)
			i = i + 1
		}
		if newCost < currentCost {
			currentState = []float64{x1, x2}
			currentCost = newCost
			if currentCost < bestCost {
				bestSolution = currentState
				bestCost = currentCost
			}
		} else {
			probAcc = probabilityAcceptance(newCost, currentCost, temperature)
			if probAcc > randFloats(0, 1) {
				currentState = []float64{x1, x2}
				currentCost = newCost
			}
		}
		xSolutions = append(xSolutions, bestSolution)
		costSolutions = append(costSolutions, bestCost)
		temperature = temperature * alpha
	}
	return xSolutions, costSolutions
}

func randFloats(min, max float64) float64 {
	/*
		I.S : Telah terdefinisi range nilainya
		F.S : Telah dihasilkan nilai random dari rentang nilai yang sudah didefinisikan
	*/
	var random float64
	random = min + rand.Float64()*(max-min)
	return random
}

func countAccuracy(saSolution float64) float64 {
	/*
		I.S : Telah terdefinisi nilai dari hasil perhitungan dengan metode simulated anealing
		F.S : Telah terhitung akurasi dari hasil perhitungan dengan metode simulated anealing dengan solusi eksaknya
	*/
	exactSolution := fungsi(0, 0)
	return (1.0 - math.Abs((saSolution-exactSolution)/exactSolution)) * 100.0
}
