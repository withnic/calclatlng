package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

// radian
type rad float64

type geo struct {
	Lng float64 //経度 x
	Lat float64 //緯度 y
}

func (point *geo) radian() (rad, rad) {
	return degToRad(point.Lat), degToRad(point.Lng)
}

// 地球の半径(km)
const r = 6378.137

// 地球の半径(m)
const rm = 6378137

//ベッセル楕円体（旧日本測地系）
const (
	bA  = 6377397.155
	bB  = 6356079.000000
	bE2 = 0.00667436061028297
	bM  = 6334832.10663254
)

//GRS80（世界測地系）
const (
	gA  = 6378137.000
	gB  = 6356752.314140
	gE2 = 0.00669438002301188
	gM  = 6335439.32708317
)

//WGS84 (GPS)
const (
	wA  = 6378137.000
	wB  = 6356752.314245
	wE2 = 0.00669437999019758
	wM  = 6335439.32729246
)

var (
	// 緯度
	lat = flag.Float64("lat", 0.0, "latitude")
	// 経度
	lng = flag.Float64("lng", 0.0, "longitude")

	// 距離を知りたい点の緯度 （デフォルト:日本)
	tlat = flag.Float64("tlat", 11.1, "latitude")
	// 距離を知りたい点の経度 (デフォルト:日本)
	tlng = flag.Float64("tlng", 135.5, "longitude")

	//距離演算のモード
	mode = flag.String("mode", "g", "mode")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of  this:
		-mode calcmode: default grs80 (g:grs80 | w:wgs84 | b:bessel | t:torigonometory)
		-lat latitude: origin latitude. default 0.0
		-lng longitude: origin longitude. default 0.0
		-tlat target latitude: target latitude. default 11.1
		-tlng target longitude: target longitude. default 135.5 `)
	}
	flag.Parse()

	from := geo{
		Lat: *lat,
		Lng: *lng,
	}

	to := geo{
		Lat: *tlat,
		Lng: *tlng,
	}
	os.Exit(run(from, to, *mode))
}

// degToRad returns radian
func degToRad(pos float64) rad {
	return rad(pos * math.Pi / 180.0)
}

// sin resturns result of math.Sin
func sin(r rad) float64 {
	return math.Sin(float64(r))
}

// cos returns result of math.Cos
func cos(r rad) float64 {
	return math.Cos(float64(r))
}

func grs(from geo, to geo) float64 {
	return hubeny(from, to, gA, gE2, gM)
}

func bessel(from geo, to geo) float64 {
	return hubeny(from, to, bA, bE2, bM)
}

func wgs(from geo, to geo) float64 {
	return hubeny(from, to, wA, wE2, wM)
}

// hubeny returns distance. (m)
func hubeny(from geo, to geo, a float64, e2 float64, m float64) float64 {
	my := degToRad((from.Lat + to.Lat) / 2.0)
	dy := degToRad(from.Lat - to.Lat)
	dx := degToRad(from.Lng - to.Lng)
	s := sin(my)
	w := math.Sqrt(1.0 - e2*s*s)

	// zero divide
	if w == 0 {
		return 0
	}

	mm := m / (w * w * w)
	n := a / w

	dym := float64(dy) * mm
	dxncos := float64(dx) * n * cos(my)

	return math.Sqrt(dym*dym + dxncos*dxncos)
}

// distance returns distance(km). 球面三角法
func distanceByTrigonometry(from geo, to geo) float64 {
	lat1, lng1 := from.radian()
	lat2, lng2 := to.radian()
	return r * math.Acos(sin(lat1)*sin(lat2)+cos(lat1)*cos(lat2)*cos(lng2-lng1))
}

// direction returns azimuth. 球面三角法
func azimuthByTrigonometry(from geo, to geo) float64 {
	lat1, lng1 := from.radian()
	lat2, lng2 := to.radian()
	dy := cos(lat2) * sin(lng2-lng1)
	dx := cos(lat1)*sin(lat2) - sin(lat1)*cos(lat2)*cos(lng2-lng1)
	rad := math.Atan2(dy, dx)
	if dx == 0 && dy == 0 {
		return 0
	}
	deg := rad * 180 / math.Pi

	if deg < 0.0 {
		deg += 360.0
	}
	return deg
}

func calcDistanceAndName(from geo, to geo, m string) (float64, string) {
	switch m {
	case "t":
		return distanceByTrigonometry(from, to), "球面三角法"
	case "w":
		return wgs(from, to), "Hudeny (WGS84)"
	case "b":
		return bessel(from, to), "Hudeny (ベッセル楕円体)"
	case "g":
		return grs(from, to), "Hudeny (GRS80（世界測地系）)"
	default:
		return 0.0, "No selected"
	}
}

func run(from geo, to geo, m string) int {
	var d float64
	d, name := calcDistanceAndName(from, to, m)
	deg := azimuthByTrigonometry(from, to)

	fmt.Fprintf(os.Stdout, "距離計算方法: %s\n", name)
	fmt.Fprintf(os.Stdout, "from Point(%f, %f) to Point (%f, %f)\n", from.Lat, from.Lng, to.Lat, to.Lng)
	fmt.Fprintf(os.Stdout, "距離: %f, 方位角: %f\n", d, deg)
	return 0
}
