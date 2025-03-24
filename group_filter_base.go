package ecs

type iEntitySet interface {
	Add(e Entity)
	Remove(e Entity)
}

type IGroupFilter interface {
	iamGroupFilter()
}

type groupFilterBase[Key comparable, KeyMaker groupKeyMaker[Key]] struct {
	keyMaker KeyMaker
	entities map[Key]Set[Entity] //set大部分情况可能只有一个元素，可以用数组链表优化
}

func newGroupFilterBase[Key comparable, KeyMaker groupKeyMaker[Key]](filter IFilter) *groupFilterBase[Key, KeyMaker] {
	gf := &groupFilterBase[Key, KeyMaker]{
		entities: make(map[Key]Set[Entity]),
	}
	filter.AddListener(gf)
	return gf
}

func (gf *groupFilterBase[Key, KeyMaker]) iamGroupFilter() {}

func (gf *groupFilterBase[Key, KeyMaker]) OnEntityAdded(entity Entity) {
	key := gf.keyMaker.makeKey(entity)
	set, ok := gf.entities[key]
	if !ok {
		set = make(Set[Entity])
		gf.entities[gf.keyMaker.makeKey(entity)] = set
	}
	set.Add(entity)
}

func (gf *groupFilterBase[Key, KeyMaker]) OnEntityRemoved(entity Entity) {
	set, ok := gf.entities[gf.keyMaker.makeKey(entity)]
	if !ok {
		return
	}
	delete(set, entity)
}

func (gf *groupFilterBase[Key, KeyMaker]) FindOne(key Key) (Entity, bool) {
	set, ok := gf.entities[key]
	if ok {
		for e := range set {
			return e, true
		}
	}
	return Entity{}, false
}

func (gf *groupFilterBase[Key, KeyMaker]) Foreach(key Key, f func(Entity)) {
	set, ok := gf.entities[key]
	if ok {
		for e := range set {
			f(e)
		}
	}
}
