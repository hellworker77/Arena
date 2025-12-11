package static

type ColliderType int

type CollisionGroup uint64

const (
	LayerNone   CollisionGroup = 0
	LayerPlayer CollisionGroup = 1 << iota
	LayerEnemy
	LayerProjectile
	LayerWorld
	LayerItem
)

const (
	ColliderEllipse ColliderType = iota
	ColliderRectangle
)

type Collider struct {
	Type   ColliderType
	Params map[string]float32

	Layer CollisionGroup
	Mask  CollisionGroup
}
