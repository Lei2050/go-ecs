package ecs

import (
	"unsafe"

	dataPool "github.com/Lei2050/array-pool"
)

type Include[T any] struct {
	compTypeIndex int
	// 对应T的组件池
	pool ComponentPooler
	// get中的id所指向的是entity的T组件在pool中的idx
	// 而该id和FilterBase.entities中的id是一一对应的
	get *dataPool.ArrayList[int]
}

func newInclude[T any](world *World, compTypeIndex int) *Include[T] {
	return &Include[T]{
		compTypeIndex: compTypeIndex,
		pool:          getComponentPool[T](world),
		get:           dataPool.NewArrayList[int](128),
	}
}

func (inc *Include[T]) FastRemoveAt(idx int) {
	inc.get.FastRemoveAt(idx)
}

// 就是获取池中第idx个数据
func (inc *Include[T]) GetItem(idx int) *T {
	return inc.pool.GetRef(inc.get.Get(idx)).(*T)
}

func (inc *Include[T]) AddIdx(typeIndex int, idxInPool int) {
	inc.get.Add(idxInPool)
}

type afterAddEntityProcesser interface {
	afterAddEntity(compIndices map[int]int)
}
type afterRemoveEntityProcesser interface {
	afterRemoveEntity(idx int) //idx是在filter.entities中的id
}

type filterBase struct {
	world       *World
	eventListen *FilterEventListen
	listeners   []FilterEventListener

	IncludeTypeIndices []int
	ExcludeTypeIndices []int

	IncludeMask uint64
	ExcludeMask uint64

	entities    *dataPool.ArrayList[EntityId] //这里的entity位置和include中的get会是一一对应的
	entitiesMap map[EntityId]int              //<entityId, 在entities中的id>

	afterAddEntityProcesser    afterAddEntityProcesser
	afterRemoveEntityProcesser afterRemoveEntityProcesser
}

func newFilterBase(world *World, ap afterAddEntityProcesser, rp afterRemoveEntityProcesser) *filterBase {
	return &filterBase{
		world:                      world,
		entities:                   dataPool.NewArrayList[EntityId](128),
		entitiesMap:                make(map[EntityId]int),
		afterAddEntityProcesser:    ap,
		afterRemoveEntityProcesser: rp,
	}
}

// 外部想监听filter中增加/移除entity事件，可用此方法注册
func (f *filterBase) AddListener(listener FilterEventListener) {
	f.listeners = append(f.listeners, listener)
}

func (f *filterBase) RemoveListener(listener FilterEventListener) {
	for i, l := range f.listeners {
		if l == listener {
			f.listeners = append(f.listeners[:i], f.listeners[i+1:]...)
			break
		}
	}
}

func (f *filterBase) notifyAdd(entity Entity) {
	for _, listener := range f.listeners {
		listener.OnEntityAdded(entity)
	}
}

func (f *filterBase) notifyRemove(entity Entity) {
	for _, listener := range f.listeners {
		listener.OnEntityRemoved(entity)
	}
}

// 外部想单独监听filter中增加entity事件，可用此方法注册
func (f *filterBase) OnAdd(cb func(entity Entity)) {
	if f.eventListen == nil {
		f.eventListen = newFilterEventListener()
		f.AddListener(f.eventListen)
	}
	f.eventListen.EntityAdded.AddCallback(cb)
}

// 外部想单独监听filter中移除entity事件，可用此方法注册
func (f *filterBase) OnRemove(cb func(entity Entity)) {
	if f.eventListen == nil {
		f.eventListen = newFilterEventListener()
		f.AddListener(f.eventListen)
	}
	f.eventListen.EntityRemoved.AddCallback(cb)
}

func initMask1[T any](typeIndices *[]int, mask *uint64) {
	componentType := GetComponentType[T]()
	*typeIndices = append(*typeIndices, componentType.TypeIndex)
	*mask |= componentType.Flag
}

func initMask2[T1, T2 any](typeIndices *[]int, mask *uint64) {
	componentType1 := GetComponentType[T1]()
	componentType2 := GetComponentType[T2]()
	*typeIndices = append(*typeIndices, componentType1.TypeIndex, componentType2.TypeIndex)
	*mask |= componentType1.Flag | componentType2.Flag
}

func initMask3[T1, T2, T3 any](typeIndices *[]int, mask *uint64) {
	componentType1 := GetComponentType[T1]()
	componentType2 := GetComponentType[T2]()
	componentType3 := GetComponentType[T3]()
	*typeIndices = append(*typeIndices, componentType1.TypeIndex, componentType2.TypeIndex, componentType3.TypeIndex)
	*mask |= componentType1.Flag | componentType2.Flag | componentType3.Flag
}

func initMask4[T1, T2, T3, T4 any](typeIndices *[]int, mask *uint64) {
	componentType1 := GetComponentType[T1]()
	componentType2 := GetComponentType[T2]()
	componentType3 := GetComponentType[T3]()
	componentType4 := GetComponentType[T4]()
	*typeIndices = append(*typeIndices, componentType1.TypeIndex, componentType2.TypeIndex, componentType3.TypeIndex, componentType4.TypeIndex)
	*mask |= componentType1.Flag | componentType2.Flag | componentType3.Flag | componentType4.Flag
}

// entity新增一个所要求的component后，是否满足filter的过滤条件
func (f *filterBase) isCompatibleAfterAddIncluded(entityData *EntityData) bool {
	return f.isCompatibleBeforeRemoveIncluded(entityData)
}

// entity删除一个所要求的component时，是否应该从filter中移除
// 就是判断entity当前是否满足filter的过滤条件
func (f *filterBase) isCompatibleBeforeRemoveIncluded(entityData *EntityData) bool {
	if entityData.CompFlags&f.IncludeMask != f.IncludeMask {
		//entity不完全包含filter中所需的component
		return false
	}
	for _, typeIndex := range f.IncludeTypeIndices {
		//entity不包含任意filter中所需的component
		if _, ok := entityData.CompIndices[typeIndex]; !ok {
			return false
		}
	}
	if entityData.CompFlags&f.ExcludeMask == 0 {
		//entity完全不包含filter中所排斥的component
		return true
	}
	for _, typeIndex := range f.ExcludeTypeIndices {
		//entity包含任意filter中所排斥的component
		if _, ok := entityData.CompIndices[typeIndex]; ok {
			return false
		}
	}
	return true
}

// entity新增一个所排斥的component后，是否满足filter的过滤条件
func (f *filterBase) isCompatibleAfterAddExcluded(entityData *EntityData, removeTypeIndex int) bool {
	return f.isCompatibleBeforeRemoveExcluded(entityData, removeTypeIndex)
}

// entity删除一个所排斥的component时，是否满足filter中移除
func (f *filterBase) isCompatibleBeforeRemoveExcluded(entityData *EntityData, removeTypeIndex int) bool {
	if entityData.CompFlags&f.IncludeMask != f.IncludeMask {
		//entity不完全包含filter中所需的component
		return false
	}
	for _, typeIndex := range f.IncludeTypeIndices {
		//entity不包含任意filter中所需的component
		if _, ok := entityData.CompIndices[typeIndex]; !ok {
			return false
		}
	}
	if entityData.CompFlags&f.ExcludeMask == 0 {
		//entity完全不包含filter中所排斥的component
		return true
	}
	for _, typeIndex := range f.ExcludeTypeIndices {
		if removeTypeIndex == typeIndex { //removeTypeIndex是即将移除的comp
			continue
		}
		//entity包含任意filter中所排斥的component
		if _, ok := entityData.CompIndices[typeIndex]; ok {
			return false
		}
	}
	return true
}

func (f *filterBase) addEntity(entity Entity) {
	world := entity.World()
	entityData := world.getEntityData(entity.Id)
	if !entityData.isCurrentEntityData(entity) {
		return
	}

	entityId := entity.GetId()
	if _, ok := f.entitiesMap[entityId]; ok {
		return
	}

	f.afterAddEntityProcesser.afterAddEntity(entityData.CompIndices)
	idx := f.entities.Count()
	f.entities.Add(entityId)
	f.entitiesMap[entityId] = idx
	f.notifyAdd(entity)
}

func (f *filterBase) removeEntity(entity Entity) {
	world := entity.World()
	entityData := world.getEntityData(entity.Id)
	if !entityData.isCurrentEntityData(entity) {
		return
	}

	entityId := entity.GetId()
	idx, ok := f.entitiesMap[entityId]
	if !ok {
		return
	}

	f.afterRemoveEntityProcesser.afterRemoveEntity(idx)
	delete(f.entitiesMap, entityId)
	f.entities.FastRemoveAt(idx)
	if idx < f.entities.Count() && f.entities.Count() > 0 {
		f.entitiesMap[f.entities.Get(idx)] = idx
	}
	f.notifyRemove(entity)
}

func (f *filterBase) getIncludeTypeIndices() []int {
	return f.IncludeTypeIndices
}
func (f *filterBase) getExcludeTypeIndices() []int {
	return f.ExcludeTypeIndices
}

type IFilter interface {
	getIncludeTypeIndices() []int
	getExcludeTypeIndices() []int
	isCompatibleAfterAddIncluded(entityData *EntityData) bool
	isCompatibleBeforeRemoveIncluded(entityData *EntityData) bool
	isCompatibleAfterAddExcluded(entityData *EntityData, removeTypeIndex int) bool
	isCompatibleBeforeRemoveExcluded(entityData *EntityData, removeTypeIndex int) bool
	addEntity(entity Entity)
	removeEntity(entity Entity)

	AddListener(listener FilterEventListener)
	RemoveListener(listener FilterEventListener)
}

type filterBase1[Include1 any] struct {
	*filterBase
	include1 *Include[Include1]
}

func newFilterBase1[Include1 any](world *World, compTypeIndex int) *filterBase1[Include1] {
	f := &filterBase1[Include1]{
		include1: newInclude[Include1](world, compTypeIndex),
	}
	f.filterBase = newFilterBase(world, f, f)
	initMask1[Include1](&f.filterBase.IncludeTypeIndices, &f.filterBase.IncludeMask)
	return f
}

func (f *filterBase1[Include1]) afterAddEntity(compIndices map[int]int) {
	for k, v := range compIndices {
		if f.include1.compTypeIndex == k {
			f.include1.AddIdx(k, v)
			break
		}
	}
}

// idx是在filter.entities中的id
func (f *filterBase1[Include1]) afterRemoveEntity(idx int) {
	f.include1.FastRemoveAt(idx)
}

func (f *filterBase1[Include1]) Foreach(callback func(entity Entity, comp Include1)) {
	count := f.include1.get.Count()
	for i := range count {
		entityId := f.entities.Get(i)
		entity := Entity{
			Id:       entityId.Id,
			Gen:      entityId.Gen,
			WorldPtr: uintptr(unsafe.Pointer(f.world)),
		}
		callback(entity, *f.include1.GetItem(i))
	}
}

type filterBase2[Include1, Include2 any] struct {
	*filterBase
	include1 *Include[Include1]
	include2 *Include[Include2]
}

func newFilterBase2[Include1, Include2 any](world *World, compTypeIndex1 int, compTypeIndex2 int) *filterBase2[Include1, Include2] {
	f := &filterBase2[Include1, Include2]{
		include1: newInclude[Include1](world, compTypeIndex1),
		include2: newInclude[Include2](world, compTypeIndex2),
	}
	f.filterBase = newFilterBase(world, f, f)
	initMask2[Include1, Include2](&f.filterBase.IncludeTypeIndices, &f.filterBase.IncludeMask)
	return f
}

func (f *filterBase2[Include1, Include2]) afterAddEntity(compIndices map[int]int) {
	for k, v := range compIndices {
		if f.include1.compTypeIndex == k {
			f.include1.AddIdx(k, v)
		}
		if f.include2.compTypeIndex == k {
			f.include2.AddIdx(k, v)
		}
	}
}

func (f *filterBase2[Include1, Include2]) afterRemoveEntity(idx int) {
	f.include1.FastRemoveAt(idx)
	f.include2.FastRemoveAt(idx)
}

func (f *filterBase2[Include1, Include2]) Foreach(callback func(entity Entity, comp1 Include1, comp2 Include2)) {
	count := f.include1.get.Count()
	for i := range count {
		entityId := f.entities.Get(i)
		entity := Entity{
			Id:       entityId.Id,
			Gen:      entityId.Gen,
			WorldPtr: uintptr(unsafe.Pointer(f.world)),
		}
		callback(entity, *f.include1.GetItem(i), *f.include2.GetItem(i))
	}
}

type filterBase3[Include1, Include2, Include3 any] struct {
	*filterBase
	include1 *Include[Include1]
	include2 *Include[Include2]
	include3 *Include[Include3]
}

func newFilterBase3[Include1, Include2, Include3 any](world *World, compTypeIndex1 int, compTypeIndex2 int, compTypeIndex3 int) *filterBase3[Include1, Include2, Include3] {
	f := &filterBase3[Include1, Include2, Include3]{
		include1: newInclude[Include1](world, compTypeIndex1),
		include2: newInclude[Include2](world, compTypeIndex2),
		include3: newInclude[Include3](world, compTypeIndex3),
	}
	f.filterBase = newFilterBase(world, f, f)
	initMask3[Include1, Include2, Include3](&f.filterBase.IncludeTypeIndices, &f.filterBase.IncludeMask)
	return f
}

func (f *filterBase3[Include1, Include2, Include3]) afterAddEntity(compIndices map[int]int) {
	for k, v := range compIndices {
		if f.include1.compTypeIndex == k {
			f.include1.AddIdx(k, v)
		}
		if f.include2.compTypeIndex == k {
			f.include2.AddIdx(k, v)
		}
		if f.include3.compTypeIndex == k {
			f.include3.AddIdx(k, v)
		}
	}
}

func (f *filterBase3[Include1, Include2, Include3]) afterRemoveEntity(idx int) {
	f.include1.FastRemoveAt(idx)
	f.include2.FastRemoveAt(idx)
	f.include3.FastRemoveAt(idx)
}

func (f *filterBase3[Include1, Include2, Include3]) Foreach(callback func(entity Entity, comp1 Include1, comp2 Include2, comp3 Include3)) {
	count := f.include1.get.Count()
	for i := range count {
		entityId := f.entities.Get(i)
		entity := Entity{
			Id:       entityId.Id,
			Gen:      entityId.Gen,
			WorldPtr: uintptr(unsafe.Pointer(f.world)),
		}
		callback(entity, *f.include1.GetItem(i), *f.include2.GetItem(i), *f.include3.GetItem(i))
	}
}

type filterBase4[Include1, Include2, Include3, Include4 any] struct {
	*filterBase
	include1 *Include[Include1]
	include2 *Include[Include2]
	include3 *Include[Include3]
	include4 *Include[Include4]
}

func newFilterBase4[Include1, Include2, Include3, Include4 any](world *World, compTypeIndex1 int, compTypeIndex2 int, compTypeIndex3 int, compTypeIndex4 int) *filterBase4[Include1, Include2, Include3, Include4] {
	f := &filterBase4[Include1, Include2, Include3, Include4]{
		include1: newInclude[Include1](world, compTypeIndex1),
		include2: newInclude[Include2](world, compTypeIndex2),
		include3: newInclude[Include3](world, compTypeIndex3),
		include4: newInclude[Include4](world, compTypeIndex4),
	}
	f.filterBase = newFilterBase(world, f, f)
	initMask4[Include1, Include2, Include3, Include4](&f.filterBase.IncludeTypeIndices, &f.filterBase.IncludeMask)
	return f
}

func (f *filterBase4[Include1, Include2, Include3, Include4]) afterAddEntity(compIndices map[int]int) {
	for k, v := range compIndices {
		if f.include1.compTypeIndex == k {
			f.include1.AddIdx(k, v)
		}
		if f.include2.compTypeIndex == k {
			f.include2.AddIdx(k, v)
		}
		if f.include3.compTypeIndex == k {
			f.include3.AddIdx(k, v)
		}
		if f.include4.compTypeIndex == k {
			f.include4.AddIdx(k, v)
		}
	}
}

func (f *filterBase4[Include1, Include2, Include3, Include4]) afterRemoveEntity(idx int) {
	f.include1.FastRemoveAt(idx)
	f.include2.FastRemoveAt(idx)
	f.include3.FastRemoveAt(idx)
	f.include4.FastRemoveAt(idx)
}

func (f *filterBase4[Include1, Include2, Include3, Include4]) Foreach(callback func(entity Entity, comp1 Include1, comp2 Include2, comp3 Include3, comp4 Include4)) {
	count := f.include1.get.Count()
	for i := range count {
		entityId := f.entities.Get(i)
		entity := Entity{
			Id:       entityId.Id,
			Gen:      entityId.Gen,
			WorldPtr: uintptr(unsafe.Pointer(f.world)),
		}
		callback(entity, *f.include1.GetItem(i), *f.include2.GetItem(i), *f.include3.GetItem(i), *f.include4.GetItem(i))
	}
}
