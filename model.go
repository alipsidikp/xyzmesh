package xyzmesh

type detailData struct {
	XValue    float64
	YValue    float64
	Value     float64
	IsOrigin  bool
	MaxInBlok float64 //only for origin data
}

type XyzMesh struct {
	//input
	Sources      [][]float64 //[][x,y,z][x,y,z][x,y,z]...
	isRounded    bool        //using rounded xy, true will use without comma decimal false will use 2 comma, default is true
	xScaleDiv    int         //Scale Div >0, for more detail result (*higher value will make process more), default is 10
	yScaleDiv    int         //Scale Div >0, for more detail result (*higher value will make process more), default is 10
	returnOrigin bool        //return value for mesh, default is true
	blockRadius  int         //Block 1-5, will itterate calculation until this detail. 1 Block radius 9 cell, 2 Block 25 cell, default is 1

	//optional
	itteratte int     //itterate until defined value or using threshold with max value, default is 5, -1 will using max value per blok radius
	XScale    float64 //scale for x value, if not set will calculate with XScaleDiv
	YScale    float64 //scale for y value, if not set will calculate with YScaleDiv

	//calculate
	xMax float64
	yMax float64
	xMin float64
	yMin float64

	xList []float64
	yList []float64
}
