package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	dataPool "github.com/Lei2050/array-pool"
)

// World 结构体代表一个实体组件系统（ECS）中的世界实例，它管理着所有的实体、组件和过滤器。
// 一个 World 实例可以包含多个实体，每个实体可以拥有多个组件，同时 World 还负责维护各种过滤器，
// 用于筛选符合特定条件的实体。
type World struct {
	// entityPool 是一个entityData对象池，用于存储和管理所有实体的基本数据。
	// 它使用 dataPool.Pool 来实现对象池化，提高实体创建和销毁的效率。
	entityPool *dataPool.Pool[EntityData]
	// componentPools 键为组件的反射类型，值为对应的组件对象池。
	componentPools map[reflect.Type]ComponentPooler //<Type, componentPool>
	// compTypeIndexPools 键为组件的类型索引，值为对应的组件对象池。
	compTypeIndexPools map[int]ComponentPooler //<typeIndex, componentPool>
	// filters 管理所有IFilter类型过滤器，这些过滤相当于是一个列表，不需要通过key来获取entity
	// <过滤器类型名称（由反射获取），过滤器实例>
	filters map[string]IFilter //<FilterTypeName, filter>
	// groupFilters 管理所有IGroupFilter类型过滤器，这些支持通过key来快速获取指定entity
	// <过滤器类型名称（由反射获取），过滤器实例>
	groupFilters map[string]IGroupFilter //<FilterTypeName, filter>
	// <组件类型索引，包含该组件的过滤器列表>
	//<typeIndex, []filter>
	filterByIncludedComps map[int][]IFilter
	// <组件类型索引，排斥该组件的过滤器列表>
	//<typeIndex, []filter>
	filterByExcludedComps map[int][]IFilter //<typeIndex, []filter> //key是排斥的comp的typeIndex
	// groupKeyEventReceivers 管理所有groupKey事件的接收者
	// 用于通知groupFilter过滤器，当entity的groupKey发生变化时，需要更新集合
	groupKeyEventReceivers map[int][]groupKeyEvent //<typeIndex, []handler> //key是comp的typeIndex
}

// 实列化一个World
func NewWorld() *World {
	return &World{
		entityPool:         dataPool.NewPool[EntityData](segmentSize),
		componentPools:     make(map[reflect.Type]ComponentPooler),
		compTypeIndexPools: make(map[int]ComponentPooler),

		filters:      make(map[string]IFilter),
		groupFilters: make(map[string]IGroupFilter),

		filterByIncludedComps: make(map[int][]IFilter),
		filterByExcludedComps: make(map[int][]IFilter),

		groupKeyEventReceivers: make(map[int][]groupKeyEvent),
	}
}

// 在当前世界中创建一个新的实体。
func (w *World) NewEntity() Entity {
	// 从实体池中分配一个新的实体数据，返回该实体在池中的索引 idx 和指向实体数据的指针 pe
	idx, pe := w.entityPool.Alloc()
	if pe.Gen == 0 {
		pe.Gen = 1
		pe.CompIndices = make(map[int]int)
	}
	return Entity{
		Id:       idx,
		Gen:      pe.Gen,
		WorldPtr: uintptr(unsafe.Pointer(w)),
	}
}

// 根据实体的 ID 获取对应的实体数据。
// 返回值是一个指向 EntityData 结构体的指针，该结构体包含了指定 ID 实体的相关数据。
func (w *World) getEntityData(id int) *EntityData {
	return w.entityPool.GetRef(id)
}

// freeEntityData 用于释放指定索引的entity数据，
// 参数 idx 表示要释放的entity在池中的索引。
func (w *World) freeEntityData(idx int) {
	entityData := w.entityPool.GetRef(idx)
	// 代数+1
	gen := entityData.Gen + 1
	// 调用实体池的 Free 方法释放指定索引的实体数据，该方法会重置数据
	w.entityPool.Free(idx)
	entityData.Gen = gen
}

//func (w *World) isCurrentEntityData(entity Entity) bool {
//	entityData := w.entityPool.GetRef(entity.Id)
//	return entityData.Gen == entity.Gen
//}

// getComponentPool 函数用于根据指定的组件类型获取对应的组件池。
// 如果组件池不存在，则会创建一个新的组件池并进行注册。
// 参数 w 是 World 结构体的指针，代表当前的世界实例。
func getComponentPool[T any](w *World) ComponentPooler {
	t := reflect.TypeOf((*T)(nil)).Elem()
	//componentPoolMapLocker.Lock()
	//defer componentPoolMapLocker.Unlock()
	pool, ok := w.componentPools[t]
	if !ok {
		componentType, ok := componentTypeMap[t]
		if !ok {
			panic(fmt.Sprintf("component type %s is not registered", t.Name()))
		}
		pool = &ComponentPool[T]{
			pool: dataPool.NewPool[T](segmentSize),
		}
		w.componentPools[t] = pool
		w.compTypeIndexPools[componentType.TypeIndex] = pool
	}
	return pool
}

// 根据组件类型索引获取组件池
func (w *World) getComponentPoolByTypeIndex(typeIndex int) ComponentPooler {
	return w.compTypeIndexPools[typeIndex]
}

// updateFiltersAfterAdd 用于在entity添加指定组件后更新过滤器。
// 该方法会检查所有与指定组件类型索引相关的包含过滤器和排除过滤器，
// 根据实体添加组件后的兼容性，决定是否将实体添加到过滤器中或从过滤器中移除。
// 参数 typeIndex 表示新添加组件的类型索引。
// 参数 entity 表示添加了组件的实体。
// 参数 entityData 表示该实体的数据。
func (w *World) updateFiltersAfterAdd(typeIndex int, entity Entity, entityData *EntityData) {
	// 从 filterByIncludedComps 中获取与该组件类型索引相关的包含过滤器列表
	// 这些过滤器要求实体包含该组件
	filters, ok := w.filterByIncludedComps[typeIndex]
	if ok {
		for _, filter := range filters {
			// 检查实体添加该组件后是否满足过滤器的包含条件
			// 如果满足条件，则将该实体添加到过滤器中
			if filter.isCompatibleAfterAddIncluded(entityData) {
				filter.addEntity(entity)
			}
		}
	}
	// 从 filterByExcludedComps 中获取与该组件类型索引相关的排斥过滤器列表
	// 这些过滤器要求实体不包含该组件
	filters, ok = w.filterByExcludedComps[typeIndex]
	if ok {
		for _, filter := range filters {
			// 检查实体添加该组件后是否满足过滤器的排除条件
			// 如果满足条件，则将该实体从过滤器中移除
			if filter.isCompatibleAfterAddExcluded(entityData, typeIndex) {
				filter.removeEntity(entity)
			}
		}
	}
}

// updateFiltersBeforeRemove 用于在entity移除指定组件之前更新过滤器。
// 该方法会检查所有与指定组件类型索引相关的包含过滤器和排除过滤器，
// 根据实体移除组件前的兼容性，决定是否将实体从过滤器中移除或添加到过滤器中。
// 参数 typeIndex 表示即将移除组件的类型索引。
// 参数 entity 表示即将移除组件的实体。
// 参数 entityData 表示该实体的数据。
func (w *World) updateFiltersBeforeRemove(typeIndex int, entity Entity, entityData *EntityData) {
	// 从 filterByIncludedComps 中获取与该组件类型索引相关的包含过滤器列表
	// 这些过滤器要求实体包含该组件
	filters, ok := w.filterByIncludedComps[typeIndex]
	if ok {
		for _, filter := range filters {
			// 检查实体在移除该组件前是否满足过滤器的包含条件
			// 如果满足条件，说明所需的组件被移除，则将该实体从过滤器中移除
			if filter.isCompatibleBeforeRemoveIncluded(entityData) {
				filter.removeEntity(entity)
			}
		}
	}
	// 从 filterByExcludedComps 中获取与该组件类型索引相关的排除过滤器列表
	// 这些过滤器要求实体不包含该组件
	filters, ok = w.filterByExcludedComps[typeIndex]
	if ok {
		for _, filter := range filters {
			// 如果移除该组件，entity满足过滤器条件，则加入过滤器
			if filter.isCompatibleBeforeRemoveExcluded(entityData, typeIndex) {
				filter.addEntity(entity)
			}
		}
	}
}

// RegisterFilter 向指定的world注册一个filter
// 目前过滤器在使用之前都要先Register，并且要在world.NewEntity()之前注册
// 如果在world.NewEntity()之后注册，会导致Filter中的entity数量不准确
func RegisterFilter[T IFilter](w *World, filter T) T {
	t := reflect.TypeOf(filter)
	if t.Kind() != reflect.Pointer {
		panic("filter must be a pointer")
	}

	t = t.Elem()
	filterName := t.Name()
	_, ok := w.filters[filterName]
	if ok {
		panic("repeat register filter")
	}

	w.filters[filterName] = filter
	//filter所有包含的组件类型索引
	for _, typeIndex := range filter.getIncludeTypeIndices() {
		w.filterByIncludedComps[typeIndex] = append(w.filterByIncludedComps[typeIndex], filter)
	}
	//filter所有排除的组件类型索引
	for _, typeIndex := range filter.getExcludeTypeIndices() {
		w.filterByExcludedComps[typeIndex] = append(w.filterByExcludedComps[typeIndex], filter)
	}

	return filter
}

// RegisterGroupFilter 向指定的world注册一个groupFilter。
// 目前过滤器在使用之前都要先Register，并且要在world.NewEntity()之前注册。
// 如果在world.NewEntity()之后注册，会导致Filter中的entity数量不准确。
// GroupFilter类型依赖一个对应的Filter类型，该Filter类型必须在GroupFilter类型注册之前注册；
// TODO: 未来可以考虑实现自动注册对应的Filter类型
func RegisterGroupFilter[T IGroupFilter](w *World, filter T) T {
	t := reflect.TypeOf(filter)
	if t.Kind() != reflect.Pointer {
		panic("filter must be a pointer")
	}

	t = t.Elem()
	filterName := t.Name()
	_, ok := w.groupFilters[filterName]
	if ok {
		panic("repeat register filter")
	}

	w.groupFilters[filterName] = filter
	return filter
}

// GetFilter 是一个泛型函数，用于从指定的 World 实例中获取已注册的过滤器。
// 它接收一个 World 结构体指针和一个实现了 IFilter 接口的泛型类型 T。
// 函数会通过反射获取过滤器的名称，并从 World 实例的 filters 映射中查找对应的过滤器。
// 如果过滤器不是指针类型或者未注册，函数将触发 panic。
// 参数 w 是 World 结构体的指针，代表当前的世界实例。
// 返回值是实现了 IFilter 接口的泛型类型 T 对应的过滤器实例。
func GetFilter[T IFilter](w *World) T {
	var dt T
	t := reflect.TypeOf(dt)
	if t.Kind() != reflect.Pointer { //要求过滤器都是指针
		panic("filter must be a pointer")
	}

	t = t.Elem()
	filterName := t.Name()
	filter, ok := w.filters[filterName]
	if !ok {
		// 不允许使用未注册的过滤器
		panic(fmt.Sprintf("filter:%s not registered", filterName))
	}
	return filter.(T)
}

// GetGroupFilter 是一个泛型函数，用于从指定的 World 实例中获取已注册的分组过滤器。
// 它接收一个 World 结构体指针和一个实现了 IGroupFilter 接口的泛型类型 T。
// 函数会通过反射获取分组过滤器的名称，并从 World 实例的 groupFilters 映射中查找对应的分组过滤器。
// 如果分组过滤器不是指针类型或者未注册，函数将触发 panic。
// 参数 w 是 World 结构体的指针，代表当前的世界实例。
// 返回值是实现了 IGroupFilter 接口的泛型类型 T 对应的分组过滤器实例。
func GetGroupFilter[T IGroupFilter](w *World) T {
	var dt T
	t := reflect.TypeOf(dt)
	if t.Kind() != reflect.Pointer { //要求过滤器都是指针
		panic("filter must be a pointer")
	}

	t = t.Elem()
	filterName := t.Name()
	filter, ok := w.groupFilters[filterName]
	if !ok {
		// 不允许使用未注册的过滤器
		panic(fmt.Sprintf("filter:%s not registered", filterName))
	}
	return filter.(T)
}

// 注册相关的groupKey事件
func registerGroupKeyEventByType[T any](world *World, entitySet iEntitySet, filter IFilter) *groupKeyEventProxy {
	componentType := GetComponentType[T]()
	proxy := &groupKeyEventProxy{set: entitySet, filter: filter}
	world.registerGroupKeyEvent(componentType.TypeIndex, groupKeyAdd, proxy)
	world.registerGroupKeyEvent(componentType.TypeIndex, groupKeyRemove, proxy)
	return proxy
}

// 注册相关的groupKey事件
func registerGroupKeyEventByTypeAndHandler[T any](world *World, handler groupKeyEventHandler) groupKeyEventHandler {
	componentType := GetComponentType[T]()
	world.registerGroupKeyEvent(componentType.TypeIndex, groupKeyAdd, handler)
	world.registerGroupKeyEvent(componentType.TypeIndex, groupKeyRemove, handler)
	return handler
}

func (w *World) registerGroupKeyEvent(typeIndex int, eventEnum groupKeyEventKind, handler groupKeyEventHandler) {
	w.groupKeyEventReceivers[typeIndex] = append(w.groupKeyEventReceivers[typeIndex], groupKeyEvent{
		eventKind: eventEnum,
		handler:   handler,
	})
}

// 触发groupKey变更事件
func (w *World) fireGroupKeyEvent(typeIndex int, eventEnum groupKeyEventKind, entity Entity) {
	receivers, ok := w.groupKeyEventReceivers[typeIndex]
	if ok {
		for _, receiver := range receivers {
			receiver.handler.onGroupKeyEvent(eventEnum, entity)
		}
	}
}

// EntityData 结构体表示一个entity的基本数据，它在World中池化存储，
// World通过entity的id来访问EntityData，entity的id是World中EntityData池中的索引。
type EntityData struct {
	// Gen 表示第几代，也可理解未是entity的版本。
	// 每当实体被重新分配时，生成序号会递增，以此来区分不同生命周期的同名实体。
	// 比如某个Entity可能被某个系统长期记录，但在该系统取消记录前，该Entity可能被销毁，
	// 需要通过isCurrentEntityData函数来确定是否是经过回收的数据。
	Gen uint
	// IsDestroying 表示实体是否正在被销毁。
	IsDestroying bool
	// CompFlags 是一个布隆过滤器，用于快速判断实体是否包含某些组件。
	// 其中的每个位位置可能代表某个组件类型的存在与否。
	// 假设AComponent组件的类型索引是2，当第2位为0时，表示实体必定不包含AComponent组件；
	// 当第2位为1时，表示实体可能包含AComponent组件，仍需要进一步验证，
	// 因为可能存在其他组件也占用了第2位的情况，比如BComponent组件的类型索引是66。
	CompFlags uint64
	// CompIndices 键为组件类型索引，值为该组件在组件池中的索引。
	// 它记录了entity所包含的每个组件在对应的组件池中的位置。
	//<componentTypeIndex, componentPoolIndex>
	CompIndices map[int]int
}

// isCurrentEntityData 用于判断当前的实体数据是否与传入的实体匹配。
// 由于实体在被销毁并重新分配后，其生成序号（Gen）会递增，因此可以通过比较版本序号来确定
func (ed *EntityData) isCurrentEntityData(entity Entity) bool {
	return ed.Gen == entity.Gen
}
