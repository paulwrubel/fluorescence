package infinitecylinder

import (
	"fluorescence/geometry"
	"fluorescence/geometry/primitive"
	"fluorescence/geometry/primitive/aabb"
	"fluorescence/shading/material"
	"fmt"
	"math"
)

type infiniteCylinder struct {
	ray                geometry.Ray
	radius             float64
	hasInvertedNormals bool
	mat                material.Material
}

type Data struct {
	Ray                geometry.Ray `json:"ray"`
	Radius             float64      `json:"radius"`
	HasInvertedNormals bool         `json:"has_inverted_normals"`
}

func New(icd *Data) (*infiniteCylinder, error) {
	// if icd.Ray == nil {
	// 	return nil, fmt.Errorf("infiniteCylinder ray is nil")
	// }
	// if icd.Ray.Origin == nil || icd.Ray.Direction == nil {
	// 	return nil, fmt.Errorf("infiniteCylinder ray origin or ray direction is nil")
	// }
	if icd.Ray.Direction.Magnitude() == 0 {
		return nil, fmt.Errorf("infiniteCylinder ray direction is zero vector")
	}
	if icd.Radius <= 0.0 {
		return nil, fmt.Errorf("infiniteCylinder radius is 0 or negative")
	}
	icd.Ray.Direction.Unit()
	return &infiniteCylinder{
		ray:                icd.Ray,
		radius:             icd.Radius,
		hasInvertedNormals: icd.HasInvertedNormals,
	}, nil
}

func (ic *infiniteCylinder) Intersection(ray geometry.Ray, tMin, tMax float64) (*material.RayHit, bool) {
	deltaP := ic.ray.Origin.To(ray.Origin)
	preA := ray.Direction.Sub(ic.ray.Direction.MultScalar(ray.Direction.Dot(ic.ray.Direction)))
	preB := deltaP.Sub(ic.ray.Direction.MultScalar(deltaP.Dot(ic.ray.Direction)))

	// terms of the quadratic equation we are solving
	a := preA.Dot(preA)
	b := preA.Dot(preB)
	c := preB.Dot(preB) - (ic.radius * ic.radius)

	preDiscriminant := b*b - a*c

	if preDiscriminant > 0 {
		root := math.Sqrt(preDiscriminant)
		// evaluate first solution, which will be smaller
		t1 := (-b - root) / a
		// return if within range
		if t1 >= tMin && t1 <= tMax {
			return &material.RayHit{
				Ray:         ray,
				NormalAtHit: ic.normalAt(ray.PointAt(t1)),
				Time:        t1,
				Material:    ic.mat,
			}, true
		}
		// evaluate and return second solution if in range
		t2 := (-b + root) / a
		if t2 >= tMin && t2 <= tMax {
			return &material.RayHit{
				Ray:         ray,
				NormalAtHit: ic.normalAt(ray.PointAt(t2)),
				Time:        t2,
				Material:    ic.mat,
			}, true
		}
	}

	return nil, false
}

func (ic *infiniteCylinder) BoundingBox(t0, t1 float64) (*aabb.AABB, bool) {
	return nil, false
}

func (ic *infiniteCylinder) SetMaterial(m material.Material) {
	ic.mat = m
}

func (ic *infiniteCylinder) IsInfinite() bool {
	return true
}

func (ic *infiniteCylinder) IsClosed() bool {
	return true
}

func (ic *infiniteCylinder) Copy() primitive.Primitive {
	newIC := *ic
	return &newIC
}

func (ic *infiniteCylinder) normalAt(p geometry.Point) geometry.Vector {
	if ic.hasInvertedNormals {
		return ic.ray.ClosestPoint(p).To(p).Unit().Negate()
	}
	return ic.ray.ClosestPoint(p).To(p).Unit()
}

func UnitInfiniteCylinder(xOffset, yOffset, zOffset float64) *infiniteCylinder {
	icd := Data{
		Ray: geometry.Ray{
			Origin: geometry.Point{
				X: 0.0 + xOffset,
				Y: 0.0 + yOffset,
				Z: 0.0 + zOffset,
			},
			Direction: geometry.Vector{
				X: 0.0 + xOffset,
				Y: 1.0 + yOffset,
				Z: 0.0 + zOffset,
			},
		},
		Radius: 1.0,
	}
	ic, _ := New(&icd)
	return ic
}