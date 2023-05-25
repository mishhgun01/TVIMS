package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"sort"
)

func main() {
	
	var sample = []int{3, 2, 4, 2, 3, 1, 2, 1, 7, 3, 6, 3, 6, 5, 3, 7, 3, 3, 7, 0, 4, 3, 5, 4, 7,
		1, 9, 6, 2, 4, 4, 3, 7, 7, 6, 1, 6, 6, 3, 2, 3, 4, 5, 3, 3, 1, 4,
		5, 8, 4}
	histPlot(sample)
	fmt.Println(distFunction(sample))
	fmt.Println(median(sample))
}

func histPlot(sample []int) {
	var vars plotter.Values
	
	sort.Ints(sample)
	for _, v := range sample {
		vars = append(vars, float64(v))
	}
	
	p := plot.New()
	p.Title.Text = "Гистограмма частот"
	hist, err := plotter.NewHist(vars, len(sample))
	if err != nil {
		log.Fatal(err.Error())
	}
	
	p.Add(hist)
	p.X.Max = 10
	p.Y.Max = 15
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "hist.png"); err != nil {
		log.Fatal(err.Error())
	}
}

func distFunction(sample []int) map[int]float64 {
	f := make(map[int]float64)
	sort.Ints(sample)
	fmt.Println(sample)
	for i, v := range sample {
		if i == 0 {
			f[v] = 1
			continue
		}
		_, ok := f[v]
		if !ok {
			f[v] = 1
		}
		if v == sample[i-1] {
			f[v] += 1
		}
	}
	for k, _ := range f {
		f[k] = f[k] / float64(len(sample))
	}
	s := removeDuplicateInt(sample)
	for _, v := range s {
		if v == 0 {
			continue
		}
		f[v] += f[v-1]
	}
	return f
	
}

func expectation(sample []int) float64 {
	n := len(sample)
	sum := 0
	for _, v := range sample {
		sum += v
	}
	return float64(sum / n)
}
func dispersion(sample []int) float64 {
	n := len(sample)
	sum := float64(0)
	exp := expectation(sample)
	for _, v := range sample {
		sum += float64(v) - exp
		sum *= sum
	}
	return sum / float64(n)
}

func median(sample []int) int {
	if len(sample)%2 == 0 {
		return (sample[(len(sample)/2)-1] + sample[(len(sample)/2)]) / 2
	} else {
		return sample[(len(sample) / 2)]
	}
}

func asymmetry(sample []int) float64 {
	exp := expectation(sample)
	disp := dispersion(sample)
	n := len(sample)
	helper := float64(n / ((n - 1) * (n - 2)))
	var sum float64
	for _, v := range sample {
		sum += math.Pow(float64(v)-exp, 3)
	}
	return (sum / math.Pow(math.Sqrt(disp), 3)) * helper
}

func excess(sample []int) float64 {
	exp := expectation(sample)
	disp := dispersion(sample)
	n := len(sample)
	helper := float64((n) * (n + 1) / ((n - 1) * (n - 2) * (n - 3)))
	helper2 := float64(3) * math.Pow(float64(n-1), 2) / float64((n-2)*(n-3))
	var sum float64
	for _, v := range sample {
		sum += math.Pow(float64(v)-exp, 4)
	}
	return (sum/math.Pow(disp, 2))*helper - helper2
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	var list []int
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
