package main

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sqrt(x float64) float64 {
	y := 1.0
	for abs(x/y-y) > 0.00001 {
		y = (y + x/y) / 2
	}
	return y
}

func evalA(i, j int) int { return ((i+j)*(i+j+1)/2 + i + 1) }

func times(v, u []float64) {
	for i := 0; i < len(v); i++ {
		v[i] = 0
		for j := 0; j < len(u); j++ {
			v[i] = v[i] + (u[j] / float64(evalA(i, j)))
		}
	}
}

func timesTranspose(v, u []float64) {
	for i := 0; i < len(v); i++ {
		v[i] = 0
		for j := 0; j < len(u); j++ {
			v[i] = v[i] + (u[j] / float64(evalA(j, i)))
		}
	}
}

func transpose(v, u []float64) {
	x := make([]float64, len(u))
	times(x, u)
	timesTranspose(v, x)
}

func spectralnorm(n int) float64 {
	u := make([]float64, n)
	for i := 0; i < len(u); i++ {
		u[i] = 1
	}
	v := make([]float64, n)
	for i := 0; i < 10; i++ {
		transpose(v, u)
		transpose(u, v)
	}
	vBv := float64(0)
	vv := float64(0)
	for i := 0; i < n; i++ {
		vBv += u[i] * v[i]
		vv += v[i] * v[i]
	}
	return sqrt(vBv / vv)
}

func main() {
	println(spectralnorm(1))
	println(spectralnorm(10))
	println(spectralnorm(25))
}
