package xyzmesh

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/alipsidikp/cast"
)

var (
	modelData  [][]*detailData
	isReachMax bool
)

func NewXyzMesh(xScale, yScale int) *XyzMesh {
	nxm := new(XyzMesh)

	nxm.isRounded = true
	nxm.blockRadius = 1
	nxm.returnOrigin = true

	nxm.xScaleDiv = xScale
	nxm.yScaleDiv = yScale

	if xScale == 0 {
		nxm.xScaleDiv = 10
	}

	if yScale == 0 {
		nxm.yScaleDiv = 10
	}
	// } else {
	// 	nxm.XScale = float64(xScale)
	// 	nxm.YScale = float64(yScale)
	// }

	nxm.itteratte = 10

	return nxm
}

// warning -1 will make it compare with max value in blok, and higher value will impact how long it calculate
func (x *XyzMesh) ChangeItterate(itt int) {
	x.itteratte = itt
}

func (x *XyzMesh) SetXScale(xScale float64) {
	x.XScale = xScale
}

func (x *XyzMesh) SetYScale(yScale float64) {
	x.YScale = yScale
}

//Set Source and calculate prepared data
func (x *XyzMesh) SetSource(source [][]float64) {
	x.Sources = source
	mxList, myList := make(map[float64]int, 0), make(map[float64]int, 0)

	for i, val := range source {
		if len(val) < 2 {
			continue
		}

		if i == 0 {
			x.xMin, x.xMax = val[0], val[0]
			x.yMin, x.yMax = val[1], val[1]
		}

		if val[0] > x.xMax {
			x.xMax = val[0]
		}

		if val[0] < x.xMin {
			x.xMin = val[0]
		}

		if val[1] > x.yMax {
			x.yMax = val[1]
		}

		if val[1] < x.yMin {
			x.yMin = val[1]
		}

		mxList[cast.ToF64(val[0], 2, cast.RoundingAuto)] = 1
		myList[cast.ToF64(val[1], 2, cast.RoundingAuto)] = 1
	}

	if x.XScale == 0 {
		x.XScale = (x.xMax - x.xMin) / float64(x.xScaleDiv)
	}

	x.XScale = cast.ToF64(x.XScale, 2, cast.RoundingAuto)

	if x.YScale == 0 {
		x.YScale = (x.yMax - x.yMin) / float64(x.yScaleDiv)
	}

	x.YScale = cast.ToF64(x.YScale, 2, cast.RoundingAuto)

	if x.isRounded {
		x.XScale = cast.ToF64(x.XScale, 0, cast.RoundingUp)
		x.YScale = cast.ToF64(x.YScale, 0, cast.RoundingUp)
	}

	for xVal, _ := range mxList {
		x.xList = append(x.xList, xVal)
	}

	for yVal, _ := range myList {
		x.yList = append(x.yList, yVal)
	}
}

func (x *XyzMesh) GetResult() (res [][]float64) {
	res = [][]float64{}

	t0 := time.Now()
	x.sortSource()
	fmt.Println("SORT : ", time.Since(t0).String())

	t0 = time.Now()
	x.setModelData()
	fmt.Println("SET MODEL : ", time.Since(t0).String())

	t0 = time.Now()
	x.processModelData()
	fmt.Println("Process MODEL : ", time.Since(t0).String())

	// for _, val := range modelData {
	// 	for _, xval := range val {
	// 		fmt.Println(" :: ", xval)
	// 	}
	// }

	t0 = time.Now()
	res = x.resultFromModel()
	fmt.Println("Generate Result : ", time.Since(t0).String())

	return
}

func (x *XyzMesh) sortSource() {
	sort.Slice(x.Sources, func(i, j int) bool {
		if len(x.Sources[i]) < 2 || len(x.Sources[j]) < 2 {
			return true
		}

		flag := x.Sources[i][0] < x.Sources[j][0]
		if !flag && x.Sources[i][0] == x.Sources[j][0] {
			flag = x.Sources[i][1] < x.Sources[j][1]
		}

		return flag
	})
}

func (x *XyzMesh) setModelData() {
	modelData = make([][]*detailData, 0)
	//Buffer data source to map[float64][float64]float64{} use rounding for the keys
	for i := x.xMin; i <= x.xMax; i += x.XScale {
		xKey := cast.ToF64(i, 2, cast.RoundingAuto)
		if x.isRounded {
			xKey = cast.ToF64(i, 0, cast.RoundingAuto)
		}
		_mdata := make([]*detailData, 0)
		for ii := x.yMin; ii <= x.yMax; ii += x.YScale {
			// fmt.Println("Y Detail : ", ii, x.yMin, x.yMax, x.YScale)
			yKey := cast.ToF64(ii, 2, cast.RoundingAuto)
			if x.isRounded {
				yKey = cast.ToF64(ii, 0, cast.RoundingAuto)
			}
			_dData := new(detailData)
			_dData.XValue = xKey
			_dData.YValue = yKey

			//Set Value and Origin,
			_dData.Value, _dData.IsOrigin = x.getZOriginValue(xKey, yKey)

			_mdata = append(_mdata, _dData)
		}
		modelData = append(modelData, _mdata)
	}
}

func (x *XyzMesh) getZOriginValue(ix, iy float64) (float64, bool) {
	xVal, IsOrigin, Diff := float64(0), false, float64(0)
	for _, val := range x.Sources {
		if len(val) < 3 {
			continue
		}

		if val[0] >= (ix + x.XScale) {
			break
		}

		if val[0] <= (ix-x.XScale) || val[1] >= (iy+x.YScale) || val[1] <= (iy-x.YScale) {
			continue
		}

		iDiff := math.Abs(val[0]-ix) + math.Abs(val[1]-iy)
		if IsOrigin && iDiff < Diff {
			Diff = iDiff
			xVal = val[2]
		} else if !IsOrigin {
			Diff = iDiff
			IsOrigin = true
			xVal = val[2]
		}
	}
	return xVal, IsOrigin
}

func (x *XyzMesh) processModelData() {

	if len(modelData) == 0 {
		return
	}

	count, nFlag := int(0), int(0)
	for {
		count++
		nX, nY := len(modelData)-1, len(modelData[0])-1
		for xi, xval := range modelData {
			arrxi := []int{}
			for n := xi - x.blockRadius; n <= xi+x.blockRadius; n++ {
				if n >= 0 && n <= nX {
					arrxi = append(arrxi, n)
				}
			}
			for yi, yval := range xval {
				if yval.IsOrigin {
					continue
				}

				arryi := []int{}
				for n := yi - x.blockRadius; n <= yi+x.blockRadius; n++ {
					if n >= 0 && n <= nY {
						arryi = append(arryi, n)
					}
				}

				Sum, nDiv, maxOri := float64(0), float64(0), float64(0)
				for _, axi := range arrxi {
					for _, ayi := range arryi {
						nDiv += 1
						Sum += modelData[axi][ayi].Value
						if modelData[axi][ayi].IsOrigin && modelData[axi][ayi].Value > maxOri {
							maxOri = modelData[axi][ayi].Value
						}
					}
				}

				if nDiv > 0 {
					modelData[xi][yi].Value = Sum / nDiv
					modelData[xi][yi].MaxInBlok = maxOri
				}
			}
		}

		if count > x.itteratte || nFlag > len(modelData) {
			break
		}
	}
}

func (x *XyzMesh) resultFromModel() (res [][]float64) {
	res = make([][]float64, 0)

	sort.Float64s(x.xList)
	sort.Float64s(x.yList)

	for _, ix := range x.xList {
		for _, iy := range x.yList {
			res = append(res, []float64{ix, iy, getZValueModel(ix, iy, x.XScale, x.YScale)})
			// res[ix][iy] = getZValueModel(ix, iy, x.XScale, x.YScale)
		}
	}

	return
}

func getZValueModel(ix, iy, xscale, yscale float64) float64 {
	resVal, IsOrigin, Diff, IsBreak, isFill := float64(0), false, float64(0), false, false
	for _, xval := range modelData {
		for _, yval := range xval {
			if yval.XValue >= (ix + xscale) {
				IsBreak = true
				break
			}

			if yval.XValue <= (ix-xscale) || yval.YValue >= (iy+yscale) {
				break
			}

			if yval.YValue <= (iy - yscale) {
				continue
			}

			iDiff := math.Abs(yval.XValue-ix) + math.Abs(yval.YValue-iy)
			if (yval.IsOrigin && !IsOrigin) || (yval.IsOrigin && IsOrigin && iDiff < Diff) {
				Diff = iDiff
				resVal = yval.Value
				IsOrigin = true
			} else if (isFill && iDiff < Diff) || !isFill {
				Diff = iDiff
				resVal = yval.Value
			}

			isFill = true
		}

		if IsBreak {
			break
		}
	}
	return resVal
}
