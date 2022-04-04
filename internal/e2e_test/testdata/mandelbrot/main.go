package main

const limit = 4.0  // abs(z) < 2
const maxIter = 50 // number of iterations

func renderRow(initialR, initialI []float64, row []byte, y0 int) {
	i := 0
	j := 0
	x := 0
	res := byte(0)
	b := byte(0)
	Zr1 := float64(0)
	Zr2 := float64(0)
	Zi1 := float64(0)
	Zi2 := float64(0)
	Tr1 := float64(0)
	Tr2 := float64(0)
	Ti1 := float64(0)
	Ti2 := float64(0)

	// TODO: use range stmt.
	for xByte := 0; xByte < len(row); xByte++ {
		res = 0
		Ci := initialI[y0]
		for i = 0; i < 8; i += 2 {
			x = xByte << 3
			Cr1 := initialR[x+i]
			Cr2 := initialR[x+i+1]

			Zr1 = Cr1
			Zi1 = Ci

			Zr2 = Cr2
			Zi2 = Ci

			b = 0

			for j = 0; j < maxIter; j++ {
				Tr1 = Zr1 * Zr1
				Ti1 = Zi1 * Zi1
				Zi1 = 2*Zr1*Zi1 + Ci
				Zr1 = Tr1 - Ti1 + Cr1

				if Tr1+Ti1 > limit {
					b |= 2
					if b == 3 {
						break
					}
				}

				Tr2 = Zr2 * Zr2
				Ti2 = Zi2 * Zi2
				Zi2 = 2*Zr2*Zi2 + Ci
				Zr2 = Tr2 - Ti2 + Cr2

				if Tr2+Ti2 > limit {
					b |= 1
					if b == 3 {
						break
					}
				}
			}
			res = (res << 2) | b
		}
		row[xByte] = ^res
	}
}

func mandelbrot(size int) []byte {
	numRows := size

	initialR := make([]float64, size)
	initialI := make([]float64, size)

	inv := 2.0 / float64(size)
	for xy := 0; xy < size; xy++ {
		i := inv * float64(xy)
		initialR[xy] = i - 1.5
		initialI[xy] = i - 1.0
	}

	bytesPerRow := size >> 3
	rowsData := make([]byte, size*bytesPerRow)
	rowOffset := 0
	for i := 0; i < numRows; i++ {
		rowBytes := rowsData[rowOffset : rowOffset+bytesPerRow]
		renderRow(initialR, initialI, rowBytes, i)
		rowOffset += bytesPerRow
	}
	return rowsData
}

func main() {
	data := mandelbrot(80)
	for i := 0; i < len(data); i++ {
		println(data[i])
	}
}
