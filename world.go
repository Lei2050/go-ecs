package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	dataPool "github.com/Lei2050/array-pool"
)

type World struct {
	entityPool         *dataPool.Pool[EntityData]
	componentPools     map[reflect.Type]ComponentPooler //<Type, componentPool>
	compTypeIndexPools map[int]ComponentPooler          //<typeIndex, componentPool>

	filters      map[string]IFilter      //<FilterTypeName, filter>
	groupFilters map[string]IGroupFilter //<FilterTypeName, filter>

	filterByIncludedComps map[int][]IFilter //<typeIndex, []filter> //key是所需comp的typeIndex
	filterByExcludedComps map[int][]IFilter //<typeIndex, []filter> //key是排斥的comp的typeIndex

	groupKeyEventReceivers map[int][]groupKeyEvent //<typeIndex, []handler> //key是comp的typeIndex
}

func NewWorld() *World {
	return &World{
		entityPool:         dataPool.NewPool[EntityData](128),
		componentPools:     make(map[reflect.Type]ComponentPooler),
		compTypeIndexPools: make(map[int]ComponentPooler),

		filters:      make(map[string]IFilter),
		groupFilters: make(map[string]IGroupFilter),

		filterByIncludedComps: make(map[int][]IFilter),
		filterByExcludedComps: make(map[int][]IFilter),

		groupKeyEventReceivers: make(map[int][]groupKeyEvent),
	}
}

func (w *World) NewEntity() Entity {
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

func (w *World) getEntityData(id int) *EntityData {
	return w.entityPool.GetRef(id)
}

func (w *World) freeEntityData(idx int) {
	entityData := w.entityPool.GetRef(idx)
	gen := entityData.Gen + 1
	w.entityPool.Free(idx) // 这里会重置数据
	entityData.Gen = gen
}

//func (w *World) isCurrentEntityData(entity Entity) bool {
//	entityData := w.entityPool.GetRef(entity.Id)
//	return entityData.Gen == entity.Gen
//}

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
			pool: dataPool.NewPool[T](128),
		}
		w.componentPools[t] = pool
		w.compTypeIndexPools[componentType.TypeIndex] = pool
	}
	return pool
}

func (w *World) getComponentPoolByTypeIndex(typeIndex int) ComponentPooler {
	return w.compTypeIndexPools[typeIndex]
}

func (w *World) updateFiltersAfterAdd(typeIndex int, entity Entity, entityData *EntityData) {
	filters, ok := w.filterByIncludedComps[typeIndex] //comp是filter所需的
	if ok {
		for _, filter := range filters {
			//增加了comp后，filter是否满足条件，满足则加入filter
			if filter.isCompatibleAfterAddIncluded(entityData) {
				filter.addEntity(entity)
			}
		}
	}
	filters, ok = w.filterByExcludedComps[typeIndex] //comp是filter排斥的
	if ok {
		for _, filter := range filters {
			if filter.isCompatibleAfterAddExcluded(entityData, typeIndex) {
				filter.removeEntity(entity)
			}
		}
	}
}

func (w *World) updateFiltersBeforeRemove(typeIndex int, entity Entity, entityData *EntityData) {
	filters, ok := w.filterByIncludedComps[typeIndex]
	if ok {
		for _, filter := range filters {
			if filter.isCompatibleBeforeRemoveIncluded(entityData) {
				filter.removeEntity(entity)
			}
		}
	}
	filters, ok = w.filterByExcludedComps[typeIndex]
	if ok {
		for _, filter := range filters {
			if filter.isCompatibleBeforeRemoveExcluded(entityData, typeIndex) {
				filter.addEntity(entity)
			}
		}
	}
}

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
	for _, typeIndex := range filter.getIncludeTypeIndices() {
		w.filterByIncludedComps[typeIndex] = append(w.filterByIncludedComps[typeIndex], filter)
	}
	for _, typeIndex := range filter.getExcludeTypeIndices() {
		w.filterByExcludedComps[typeIndex] = append(w.filterByExcludedComps[typeIndex], filter)
	}

	return filter
}

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

func GetFilter[T IFilter](w *World) T {
	var dt T
	t := reflect.TypeOf(dt)
	if t.Kind() != reflect.Pointer {
		panic("filter must be a pointer")
	}

	t = t.Elem()
	filterName := t.Name()
	filter, ok := w.filters[filterName]
	if !ok {
		panic(fmt.Sprintf("filter:%s not registered", filterName))
	}
	return filter.(T)
}

func GetGroupFilter[T IGroupFilter](w *World) T {
	var dt T
	t := reflect.TypeOf(dt)
	if t.Kind() != reflect.Pointer {
		panic("filter must be a pointer")
	}

	t = t.Elem()
	filterName := t.Name()
	filter, ok := w.groupFilters[filterName]
	if !ok {
		panic(fmt.Sprintf("filter:%s not registered", filterName))
	}
	return filter.(T)
}

func registerGroupKeyEventByType[T any](world *World, entitySet iEntitySet, filter IFilter) {
	componentType := GetComponentType[T]()
	proxy := &groupKeyEventProxy{set: entitySet, filter: filter}
	world.registerGroupKeyEvent(componentType.TypeIndex, groupKeyAdd, proxy)
	world.registerGroupKeyEvent(componentType.TypeIndex, groupKeyRemove, proxy)
}

func (w *World) registerGroupKeyEvent(typeIndex int, eventEnum groupKeyEventKind, handler groupKeyEventHandler) {
	w.groupKeyEventReceivers[typeIndex] = append(w.groupKeyEventReceivers[typeIndex], groupKeyEvent{
		eventKind: eventEnum,
		handler:   handler,
	})
}

type EntityData struct {
	Gen       uint
	IsDestroy bool

	CompFlags uint64
	//<componentTypeIndex, componentPoolIndex>
	CompIndices map[int]int
}

func (ed *EntityData) isCurrentEntityData(entity Entity) bool {
	return ed.Gen == entity.Gen
}
