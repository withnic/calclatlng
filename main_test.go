package main

import (
	"math"
	"testing"
)

func BenchmarkDistanceOnTrigonometry(b *testing.B) {
	from := geo{
		Lat: 36.10056,
		Lng: 140.09111,
	}

	to := geo{
		Lat: 35.65500,
		Lng: 139.74472,
	}
	for i := 0; i < b.N; i++ {
		distanceByTrigonometry(from, to)
	}
}

func BenchmarkDistanceHubenyWGS(b *testing.B) {
	from := geo{
		Lat: 36.10056,
		Lng: 140.09111,
	}

	to := geo{
		Lat: 35.65500,
		Lng: 139.74472,
	}
	for i := 0; i < b.N; i++ {
		wgs(from, to)
	}
}

func BenchmarkDistanceHubenyGRS(b *testing.B) {
	from := geo{
		Lat: 36.10056,
		Lng: 140.09111,
	}

	to := geo{
		Lat: 35.65500,
		Lng: 139.74472,
	}
	for i := 0; i < b.N; i++ {
		grs(from, to)
	}
}

func BenchmarkDistanceHubenyBessel(b *testing.B) {
	from := geo{
		Lat: 36.10056,
		Lng: 140.09111,
	}

	to := geo{
		Lat: 35.65500,
		Lng: 139.74472,
	}
	for i := 0; i < b.N; i++ {
		bessel(from, to)
	}
}

func TestDistanceOnTrigonometry(t *testing.T) {
	from := geo{
		Lat: 36.10056,
		Lng: 140.09111,
	}

	to := geo{
		Lat: 35.65500,
		Lng: 139.74472,
	}

	d := distanceByTrigonometry(from, to)
	expected := 58.0 // km
	if math.Abs(d-expected) > 1.0 {
		t.Fatalf("expedted: %f, got %f", expected, d)
	}
}

func TestHubeny(t *testing.T) {
	from := geo{
		Lat: 36.10056,
		Lng: 140.09111,
	}

	to := geo{
		Lat: 35.65500,
		Lng: 139.74472,
	}

	distance := hubeny(from, to, w_a, w_e2, w_m)
	expected := 58502.45

	if math.Abs(distance-expected) > 0.1 {
		t.Fatalf("expedted: %f, got %f", expected, distance)
	}
}

func TestAzimuthByTrigonometry(t *testing.T) {
	from := geo{
		Lat: 36.06135892,
		Lng: 140.05162781,
	}

	to := geo{
		Lat: 35.39181025,
		Lng: 139.44411016,
	}

	azimuth := azimuthByTrigonometry(from, to)
	expected := 210.0

	if math.Abs(azimuth-expected) > 50 {
		t.Fatalf("expedted: %f, got %f", expected, azimuth)
	}
}
