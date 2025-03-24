package ecs

type iEntitySet interface {
	Add(e Entity)
	Remove(e Entity)
}

type groupFilterBase[Key comparable, KeyMaker groupKeyMaker[Key]] struct {
	keyMaker KeyMaker
	entities map[Key]Entity
}

func newGroupFilterBase[Key comparable, KeyMaker groupKeyMaker[Key]](filter IFilter) *groupFilterBase[Key, KeyMaker] {
	gf := &groupFilterBase[Key, KeyMaker]{
		entities: make(map[Key]Entity),
	}
	filter.AddListener(gf)
	return gf
}

func (gf *groupFilterBase[Key, KeyMaker]) OnEntityAdded(entity Entity) {
	gf.entities[gf.keyMaker.makeKey(entity)] = entity
}
func (gf *groupFilterBase[Key, KeyMaker]) OnEntityRemoved(entity Entity) {
	delete(gf.entities, gf.keyMaker.makeKey(entity))
}

func (gf *groupFilterBase[Key, KeyMaker]) Find(comp Key) (Entity, bool) {
	v, ok := gf.entities[comp]
	return v, ok
}
