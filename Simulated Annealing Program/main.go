package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"net/http"
	"html/template"
	"path"
)

func main() {
	// Melakukan deklarasi dan memanggil fungsi simulated anealing

    rand.Seed(time.Now().UnixNano())							// Digunakan untuk mendapatkan nilai random yang uniq
    var temperature, final_temperature, alpha, x1, x2 float64 	
    var initial_state, cost_solutions []float64
    var x_solutions [][]float64        							

    temperature = 100.0
    final_temperature = 0.000001
    alpha = 0.99999
    x1 = randFloats(-10, 10)
    x2 = randFloats(-10, 10)
    initial_state = []float64{x1, x2}

    // Pemanggilan fungsi dan menampung hasil returnnya kedalan dua varible, untuk mendapatkan nilai x1, x2, dan nilai minimummnya
    x_solutions, cost_solutions = simulated_annealing(initial_state, temperature, final_temperature, alpha)		
    nilai_minimum := cost_solutions[len(cost_solutions)-1]
    nilai_x := x_solutions[len(x_solutions)-1]
    nilai_x1 := nilai_x[0]
    nilai_x2 := nilai_x[1]
    akurasi := count_accuracy(nilai_minimum)

    // Melakukan proses untuk rendering file html, untuk dijadikan tempat output program

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var filepath = path.Join("views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
		    http.Error(w, err.Error(), http.StatusInternalServerError)
		    return
		}

		var data = map[string]interface{}{
		    "nilai_minimum":  nilai_minimum,
		    "nilai_x1": nilai_x1,
		    "nilai_x2": nilai_x2,
		    "akurasi": akurasi,
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
	return (-((math.Sin(x1) * math.Cos(x2)) + ((4.0/5.0) * math.Exp(1 - math.Sqrt(math.Pow(x1, 2) + math.Pow(x2, 2))))))
}

func probability_acceptance(new_cost, current_cost, temperature float64) float64 {
	/*
		I.S : Telah terdefinisi parameter untuk dihitung
		F.S : Mengembalikan probability acceptance
	*/
	return math.Exp(-(new_cost - current_cost) / temperature)
}

func simulated_annealing(initial_state []float64, temperature, final_temperature, alpha float64) ([][]float64, []float64){
	/*
		I.S : Telah terdefinisi paremeter untuk dilakukan perhitungan dengan menggunakan metode simulated anealing
		F.S : Telah didapatkan nilai minimun, x1, dan x2
	*/
	var x_solutions [][]float64
	var cost_solutions []float64
	var x1, x2, new_cost, prob_acc float64
	
	var current_state = initial_state
	var current_cost = fungsi(initial_state[0], initial_state[1])

	var best_solution = initial_state
	var best_cost = current_cost

	for temperature > final_temperature {
		var i = 0
		for i < 100 {
			x1 = randFloats(-10, 10)
			x2 = randFloats(-10, 10)
			new_cost = fungsi(x1, x2)
			i = i + 1
		}
		if new_cost < current_cost {
			current_state = []float64{x1, x2}
			current_cost = new_cost
			if current_cost < best_cost {
				best_solution = current_state
				best_cost = current_cost
			}
		} else {
			prob_acc = probability_acceptance(new_cost, current_cost, temperature)
			if prob_acc > randFloats(0, 1) {
				current_state = []float64{x1, x2}
				current_cost = new_cost
			}
		}
		x_solutions = append(x_solutions, best_solution)
		cost_solutions = append(cost_solutions, best_cost)
		temperature = temperature * alpha
	}
	return x_solutions, cost_solutions
}

func randFloats(min, max float64) float64 {
	/*
		I.S : Telah terdefinisi range nilainya
		F.S : Telah dihasilkan nilai random dari rentang nilai yang sudah didefinisikan
	*/
    var random float64
    random = min + rand.Float64() * (max - min)
    return random
}

func count_accuracy(sa_solution float64) float64 {
	/*
		I.S : Telah terdefinisi nilai dari hasil perhitungan dengan metode simulated anealing
		F.S : Telah terhitung akurasi dari hasil perhitungan dengan metode simulated anealing dengan solusi eksaknya
	*/
	exact_solution := fungsi(0, 0)
	return (1.0 - math.Abs((sa_solution - exact_solution) / exact_solution)) * 100.0
}