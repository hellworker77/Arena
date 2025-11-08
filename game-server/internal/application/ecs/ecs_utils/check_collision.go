package ecs_utils

import (
	"game-server/internal/application/ecs/ecs_signatures/runtime"
	"game-server/internal/application/ecs/ecs_signatures/static"
)

func abs(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}

func intersectsRectVsRect(aPos, bPos runtime.Position, a, b *static.Collider) bool {
	ax, ay := aPos.X, aPos.Y
	aw, ah := a.Params["width"], a.Params["height"]

	bx, by := bPos.X, bPos.Y
	bw, bh := b.Params["width"], b.Params["height"]

	return ax < bx+bw && ax+aw > bx &&
		ay < by+bh && ay+ah > by
}

func intersectsEllipseVsEllipse(aPos, bPos runtime.Position, a, b *static.Collider) bool {
	rx := a.Params["rx"]
	ry := a.Params["ry"]
	sx := b.Params["rx"]
	sy := b.Params["ry"]

	dx := bPos.X - aPos.X
	dy := bPos.Y - aPos.Y

	nx := dx / (rx + sx)
	ny := dy / (ry + sy)

	return nx*nx+ny*ny <= 1
}

func intersectsEllipseVsRect(ePos runtime.Position, rPos runtime.Position, e *static.Collider, r *static.Collider) bool {
	rx, ry := r.Params["width"]/2, r.Params["height"]/2
	ex, ey := e.Params["rx"], e.Params["ry"]

	cx, cy := rPos.X+rx, rPos.Y+ry
	dx := abs(ePos.X - cx)
	dy := abs(ePos.Y - cy)

	if dx > rx+ex || dy > ry+ey {
		return false
	}

	if dx <= rx || dy <= ry {
		return true
	}

	ux := dx - rx
	uy := dy - ry
	return (ux*ux)/(ex*ex)+(uy*uy)/(ey*ey) <= 1
}

func checkGeometry(aPos, bPos runtime.Position, a, b static.Collider) bool {
	switch a.Type {
	case static.ColliderRectangle:
		switch b.Type {
		case static.ColliderRectangle:
			return intersectsRectVsRect(aPos, bPos, &a, &b)
		case static.ColliderEllipse:
			return intersectsEllipseVsRect(bPos, aPos, &a, &b)
		}
	case static.ColliderEllipse:
		switch b.Type {
		case static.ColliderRectangle:
			return intersectsEllipseVsRect(aPos, bPos, &a, &b)
		case static.ColliderEllipse:
			return intersectsEllipseVsEllipse(aPos, bPos, &a, &b)
		}
	}
	return false
}

func CheckCollision(aPos, bPos runtime.Position, a, b static.Collider) bool {
	if (a.Layer&b.Mask) == 0 || (b.Layer&a.Mask) == 0 {
		return false
	}

	return checkGeometry(aPos, bPos, a, b)
}
