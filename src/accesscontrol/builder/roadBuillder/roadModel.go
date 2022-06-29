package roadBuillder

type RoadModel[T any] struct {
	roads map[string]*RoadModel[T]
	level int
	data  *T
}
