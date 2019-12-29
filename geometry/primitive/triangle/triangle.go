package triangle

import (
	"fluorescence/geometry"
	"fluorescence/geometry/primitive"
	"fluorescence/geometry/primitive/aabb"
	"fluorescence/shading/material"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Triangle is an internal representation of a Triangle geometry contruct
type Triangle struct {
	A        mgl64.Vec3 `json:"a"`
	B        mgl64.Vec3 `json:"b"`
	C        mgl64.Vec3 `json:"c"`
	normal   mgl64.Vec3 // normal of the Triangle's surface
	IsCulled bool       `json:"is_culled"` // whether or not the Triangle is culled, or single-sided
	mat      material.Material
}

// Data holds information needed to contruct a Triangle
// type Data struct {
// 	A       mgl64.Vec3 `json:"a"`
// 	B       mgl64.Vec3 `json:"b"`
// 	C       mgl64.Vec3 `json:"c"`
// 	IsCulled bool           `json:"is_culled"`
// }

// Setup fills calculated fields in an Triangle
func (t *Triangle) Setup() (*Triangle, error) {
	if t.A == t.B || t.A == t.C || t.B == t.C {
		return nil, fmt.Errorf("Triangle resolves to line or point")
	}
	t.normal = t.B.Sub(t.A).Cross(t.C.Sub(t.A)).Normalize()
	return t, nil
}

// Intersection computer the intersection of this object and a given ray if it exists
func (t *Triangle) Intersection(ray geometry.Ray, tMin, tMax float64) (*material.RayHit, bool) {
	ab := t.B.Sub(t.A)
	ac := t.C.Sub(t.A)
	pVector := ray.Direction.Cross(ac)
	determinant := ab.Dot(pVector)
	if t.IsCulled && determinant < 1e-7 {
		// This ray is parallel to this Triangle or back-facing.
		return nil, false
	} else if determinant > -1e-7 && determinant < 1e-7 {
		return nil, false
	}

	inverseDeterminant := 1.0 / determinant

	tVector := ray.Origin.Sub(t.A)
	u := inverseDeterminant * (tVector.Dot(pVector))
	if u < 0.0 || u > 1.0 {
		return nil, false
	}

	qVector := tVector.Cross(ab)
	v := inverseDeterminant * (ray.Direction.Dot(qVector))
	if v < 0.0 || u+v > 1.0 {
		return nil, false
	}

	// At this stage we can compute time to find out where the intersection point is on the line.
	time := inverseDeterminant * (ac.Dot(qVector))
	if time >= tMin && time <= tMax {
		// ray intersection
		return &material.RayHit{
			Ray:         ray,
			NormalAtHit: t.normal,
			Time:        time,
			U:           0,
			V:           0,
			Material:    t.mat,
		}, true
	}
	return nil, false
}

// BoundingBox returns an AABB for this object
func (t *Triangle) BoundingBox(t0, t1 float64) (*aabb.AABB, bool) {
	return &aabb.AABB{
		A: mgl64.Vec3{
			math.Min(math.Min(t.A.X(), t.A.X()), t.C.X()) - 1e-7,
			math.Min(math.Min(t.A.Y(), t.A.Y()), t.C.Y()) - 1e-7,
			math.Min(math.Min(t.A.Z(), t.A.Z()), t.C.Z()) - 1e-7,
		},
		B: mgl64.Vec3{
			math.Max(math.Max(t.A.X(), t.B.X()), t.C.X()) + 1e-7,
			math.Max(math.Max(t.A.Y(), t.B.Y()), t.C.Y()) + 1e-7,
			math.Max(math.Max(t.A.Z(), t.B.Z()), t.C.Z()) + 1e-7,
		},
	}, true
}

// SetMaterial sets the material of this object
func (t *Triangle) SetMaterial(m material.Material) {
	t.mat = m
}

// IsInfinite returns whether this object is infinite
func (t *Triangle) IsInfinite() bool {
	return false
}

// IsClosed returns whether this object is closed
func (t *Triangle) IsClosed() bool {
	return false
}

// Copy returns a shallow copy of this object
func (t *Triangle) Copy() primitive.Primitive {
	newT := *t
	return &newT
}

// Unit creates a unit Triangle.
// The points of this Triangle are:
// A: (0, 0, 0),
// B: (1, 0, 0),
// C: (0, 1, 0).
func Unit(xOffset, yOffset, zOffset float64) *Triangle {
	t, _ := (&Triangle{
		A: mgl64.Vec3{
			0.0 + xOffset,
			0.0 + yOffset,
			0.0 + zOffset,
		},
		B: mgl64.Vec3{
			1.0 + xOffset,
			0.0 + yOffset,
			0.0 + zOffset,
		},
		C: mgl64.Vec3{
			0.0 + xOffset,
			1.0 + yOffset,
			0.0 + zOffset,
		},
		IsCulled: true,
	}).Setup()
	return t
}
