package ecs

type iEntitySet interface {
	Add(e Entity)
	Remove(e Entity)
}

// 带有key的过滤器接口，方便World中管理
type IGroupFilter interface {
	iamGroupFilter()
}

// groupFilterBase 是一个Entity集合的集合，支持按Key值查询Entity。
// 需要指定一个键值生成器groupKeyMaker[Key]，用以从Entity中提取Key值。
// groupFilterBase依赖于一个IFilter，过滤器的功能由所依赖的IFilter提供，
// 当IFilter中的Entity发生变动时，groupFilterBase会自动更新自己的Entity集合（
//
//	通过注册FilterEventListener接口）。
type groupFilterBase[Key comparable, KeyMaker groupKeyMaker[Key]] struct {
	keyMaker KeyMaker
	entities map[Key]Set[Entity] //set大部分情况可能只有一个元素，可以用数组链表优化
}

// 实例化一个groupFilterBase，需要传入一个所以依赖的IFilter。
func newGroupFilterBase[Key comparable, KeyMaker groupKeyMaker[Key]](filter IFilter) *groupFilterBase[Key, KeyMaker] {
	gf := &groupFilterBase[Key, KeyMaker]{
		entities: make(map[Key]Set[Entity]),
	}
	//监听filter的Entity增删事件
	filter.AddListener(gf)
	return gf
}

func (gf *groupFilterBase[Key, KeyMaker]) iamGroupFilter() {}

// 实现FilterEventListener接口
// 实现当所依赖的filter中的Entity发生变动时，groupFilterBase会自动更新自己的Entity集合
func (gf *groupFilterBase[Key, KeyMaker]) OnEntityAdded(entity Entity) {
	key := gf.keyMaker.makeKey(entity)
	set, ok := gf.entities[key]
	if !ok {
		set = make(Set[Entity])
		gf.entities[gf.keyMaker.makeKey(entity)] = set
	}
	set.Add(entity)
}

// 实现FilterEventListener接口
// 实现当所依赖的filter中的Entity发生变动时，groupFilterBase会自动更新自己的Entity集合
func (gf *groupFilterBase[Key, KeyMaker]) OnEntityRemoved(entity Entity) {
	set, ok := gf.entities[gf.keyMaker.makeKey(entity)]
	if !ok {
		return
	}
	delete(set, entity)
}

// 实现iEntitySet接口
func (gf *groupFilterBase[Key, KeyMaker]) Add(entity Entity) {
	gf.OnEntityAdded(entity)
}

// 实现iEntitySet接口
func (gf *groupFilterBase[Key, KeyMaker]) Remove(entity Entity) {
	gf.OnEntityRemoved(entity)
}

// 根据key找到任意一个Entity
func (gf *groupFilterBase[Key, KeyMaker]) FindOne(key Key) (Entity, bool) {
	set, ok := gf.entities[key]
	if ok {
		for e := range set {
			return e, true
		}
	}
	return Entity{}, false
}

// 遍历keyMaker(entity)==key的所有Entity
func (gf *groupFilterBase[Key, KeyMaker]) Foreach(key Key, f func(Entity)) {
	set, ok := gf.entities[key]
	if ok {
		for e := range set {
			f(e)
		}
	}
}
