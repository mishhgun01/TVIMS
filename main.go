package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"sort"
)

type Num interface {
	int | float64
}

func main() {

	var sampleFirst = []int{3, 2, 4, 2, 3, 1, 2, 1, 7, 3, 6, 3, 6, 5, 3, 7, 3, 3, 7, 0, 4, 3, 5, 4, 7,
		1, 9, 6, 2, 4, 4, 3, 7, 7, 6, 1, 6, 6, 3, 2, 3, 4, 5, 3, 3, 1, 4,
		5, 8, 4}
	var sampleSecond = []float64{1.36, 1.08, 0.09, 0.79, 3.17, 1.73, 2.19, 5.16, 16.35, 4.10, 4.09,
		1.50, 2.03, 2.93, 12.81, 0.01, 4.14, 8.34, 15.41, 3.16, 1.35, 7.45, 1.69, 3.34, 0.06,
		3.63, 5.19, 0.00, 2.56, 9.03, 2.47, 8.46, 8.37, 9.97, 10.46, 3.71, 2.99, 0.76, 5.79, 5.50, 3.28, 0.11, 3.06, 24.33, 4.70, 0.60,
		10.03, 2.27, 14.75, 0.96}
	histPlot(sampleFirst)
	freqPolygon(sampleSecond)
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

func freqPolygon(sample []float64) {
	p := plot.New()

	p.Title.Text = "Полигон частот"
	p.X.Label.Text = "Xi"
	p.Y.Label.Text = "Frequency"
	p.X.Max = 30
	p.Y.Max = 30

	f := make(map[float64]float64)
	sort.Slice(sample, func(i, j int) bool {
		return sample[i] < sample[j]
	})
	fmt.Println(sample)
	k := 4.0
	for _, v := range sample {

		if v <= k {
			f[k] += 1
		} else {
			k += 4
			f[k] = 1
		}
	}
	i := 0
	var keys []float64
	for k := range f {
		fmt.Println(k)
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	fmt.Println(keys)
	pts := make(plotter.XYs, len(keys))
	for _, key := range keys {
		pts[i].X = key
		pts[i].Y = f[key]
		i++
	}
	fmt.Println(pts)
	err := plotutil.AddLinePoints(p, pts)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		log.Fatal(err.Error())
	}
}
func distFunction[T Num](sample []T) map[T]float64 {
	f := make(map[T]float64)
	sort.Slice(sample, func(i, j int) bool {
		return sample[i] < sample[j]
	})
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
	s := removeDuplicate(sample)
	for _, v := range s {
		if v == 0 {
			continue
		}
		f[v] += f[v-1]
	}
	return f

}

func expectation[T Num](sample []T) float64 {
	n := T(len(sample))
	var sum T
	for _, v := range sample {
		sum += v
	}
	return float64(sum / n)
}
func dispersion[T Num](sample []T) float64 {
	n := T(len(sample))
	sum := float64(0)
	exp := expectation(sample)
	for _, v := range sample {
		sum += math.Pow(float64(v)-exp, 2)
	}
	return sum / float64(n-1)
}

func median[T Num](sample []T) T {
	if len(sample)%2 == 0 {
		return (sample[(len(sample)/2)-1] + sample[(len(sample)/2)]) / 2
	} else {
		return sample[(len(sample) / 2)]
	}
}

func asymmetry[T Num](sample []T) float64 {
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

func excess[T Num](sample []T) float64 {
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

func Prob[T Num](sample []T, a, b float64) float64 {
	f := distFunction(sample)
	s := removeDuplicate(sample)
	var vA, vB float64
	for i, v := range s {
		if i == 0 {
			continue
		}

		if a > float64(s[i-1]) && a < float64(v) {
			vA = f[s[i-1]]
		}
		if b > float64(s[i-1]) && b < float64(v) {
			vB = f[s[i-1]]
		}
		i++
	}
	return vB - vA
}

func removeDuplicate[T Num](slice []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
