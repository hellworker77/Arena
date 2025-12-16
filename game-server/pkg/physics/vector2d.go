package physics

import "math"

type Vector2D struct {
	X float32
	Y float32
}

func NewVector2D(x, y float32) Vector2D {
	return Vector2D{X: x, Y: y}
}

// -------------------------------- Metrics --------------------------------

func (v Vector2D) Magnitude() float32 {
	return float32(math.Hypot(float64(v.X), float64(v.Y))) // math.Hypot всё ещё float64
}

func (v Vector2D) Dot(o Vector2D) float32 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vector2D) DistanceTo(o Vector2D) float32 {
	dx := v.X - o.X
	dy := v.Y - o.Y
	return float32(math.Hypot(float64(dx), float64(dy)))
}

// -------------------------------- New Vector creators --------------------------------

func (v Vector2D) Add(o Vector2D) Vector2D {
	return Vector2D{v.X + o.X, v.Y + o.Y}
}

func (v Vector2D) Sub(o Vector2D) Vector2D {
	return Vector2D{v.X - o.X, v.Y - o.Y}
}

func (v Vector2D) Mul(s float32) Vector2D {
	return Vector2D{v.X * s, v.Y * s}
}

func (v Vector2D) Div(s float32) Vector2D {
	return Vector2D{v.X / s, v.Y / s}
}

func (v Vector2D) Normalized() Vector2D {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2D{}
	}
	return v.Div(mag)
}

// -------------------------------- Changes vector --------------------------------

func (v *Vector2D) Normalize() {
	mag := v.Magnitude()
	if mag != 0 {
		v.X /= mag
		v.Y /= mag
	}
}

func (v *Vector2D) MoveTowards(o Vector2D, d float32) {
	dx := o.X - v.X
	dy := o.Y - v.Y
	mag := float32(math.Hypot(float64(dx), float64(dy)))
	if mag == 0 || d == 0 {
		return
	}
	scale := d / mag
	v.X += dx * scale
	v.Y += dy * scale
}

// -------------------------------- For Projectiles and AI --------------------------------

func (v Vector2D) Lerp(to Vector2D, t float32) Vector2D {
	return Vector2D{
		X: v.X + (to.X-v.X)*t,
		Y: v.Y + (to.Y-v.Y)*t,
	}
}

func (v Vector2D) Reflect(normal Vector2D) Vector2D {
	n := normal.Normalized()
	return v.Sub(n.Mul(2 * v.Dot(n)))
}

func (v Vector2D) Project(onto Vector2D) Vector2D {
	ontoNorm := onto.Normalized()
	return ontoNorm.Mul(v.Dot(ontoNorm))
}

func (v Vector2D) Neg() Vector2D {
	return Vector2D{-v.X, -v.Y}
}

//TODO: SIMD optimizations for vector operations
