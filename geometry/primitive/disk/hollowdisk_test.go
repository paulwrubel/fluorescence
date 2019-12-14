package disk

import (
	"fluorescence/geometry"
	"testing"
)

var hollowDiskHit bool

func TestHollowDiskIntersectionHit(t *testing.T) {
	hollowDisk := BasicHollowDisk(0.0, 0.0, 0.0)
	r := &geometry.Ray{
		Origin: &geometry.Point{
			X: 0.75,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: &geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := hollowDisk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308)
	if !h {
		t.Errorf("Expected true (hit) but got %t\n", h)
	}
}

func BenchmarkHollowDiskIntersectionHit(b *testing.B) {
	hollowDisk := BasicHollowDisk(0.0, 0.0, 0.0)
	r := &geometry.Ray{
		Origin: &geometry.Point{
			X: 0.75,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: &geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = hollowDisk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308)
	}
	hollowDiskHit = h
}

func TestHollowDiskIntersectionMiss(t *testing.T) {
	hollowDisk := BasicHollowDisk(0.0, 0.0, 0.0)
	r := &geometry.Ray{
		Origin: &geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: &geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	_, h := hollowDisk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308)
	if h {
		t.Errorf("Expected false (miss) but got %t\n", h)
	}
}

func BenchmarkHollowDiskIntersectionMiss(b *testing.B) {
	hollowDisk := BasicHollowDisk(0.0, 0.0, 0.0)
	r := &geometry.Ray{
		Origin: &geometry.Point{
			X: 0.0,
			Y: 0.0,
			Z: 1.0,
		},
		Direction: &geometry.Vector{
			X: 0.0,
			Y: 0.0,
			Z: -1.0,
		},
	}
	var h bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, h = hollowDisk.Intersection(r, 1e-7, 1.797693134862315708145274237317043567981e+308)
	}
	hollowDiskHit = h
}