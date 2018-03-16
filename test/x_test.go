package xyzmesh

import (
	"fmt"
	"testing"

	. "github.com/alipsidikp/xyzmesh"
)

var (
	dataTest1 = [][]float64{[]float64{7, 6, 31.978}, []float64{3, 4, 59.935}, []float64{7, 8, 87.426}, []float64{3, 6, 44.879},
		[]float64{9, 10, 94.789}, []float64{1, 2, 33.08}, []float64{5, 10, 15.978},
	}
)

func TestInitialData(t *testing.T) {
	xyz := NewXyzMesh(10, 10)
	xyz.SetSource(dataTest1)
	result := xyz.GetResult()

	for _, val := range result {
		for _, xval := range val {
			fmt.Printf("%v,", xval)
		}
		fmt.Println()
	}
}
