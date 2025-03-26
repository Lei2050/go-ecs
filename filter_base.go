package ecs

import (
	"unsafe"

	dataPool "github.com/Lei2050/array-pool"
)

// Include 存储一系列实体的组件的索引，该索引是在组件对象池中的索引
type Include[T any] struct {
	compTypeIndex int
	// 对应T的组件池
	pool ComponentPooler
	// 可以理解get是一个数组，存储的是T组件在pool中的idx，
	// 数组的下标就是在filter.entities中的下标是一一对应的。
	// 比如，get[2] = 10，那么pool.Get(10)就是filter.entities[2]的T组件对象
	get *dataPool.ArrayList[int]
}

func newInclude[T any](world *World, compTypeIndex int) *Include[T] {
	return &Include[T]{
		compTypeIndex: compTypeIndex,
		pool:          getComponentPool[T](world),
		get:           dataPool.NewArrayList[int](segmentSize),
	}
}

// 移除数组中的指定下标的数据，
// 直接将元素交换到数组的末尾，然后移除数组的末尾元素，所以非常快。
func (inc *Include[T]) FastRemoveAt(idx int) {
	inc.get.FastRemoveAt(idx)
}

// 获取数组中第idx个组件对象
func (inc *Include[T]) GetItem(idx int) *T {
	return inc.pool.GetRef(inc.get.Get(idx)).(*T)
}

// 增加一个组件对象索引（对象池中的索引）到数组末尾
func (inc *Include[T]) AddIdx(typeIndex int, idxInPool int) {
	if typeIndex != inc.compTypeIndex {
		return
	}
	inc.get.Add(idxInPool)
}

type afterAddEntityProcesser interface {
	afterAddEntity(compIndices map[int]int)
}
type afterRemoveEntityProcesser interface {
	afterRemoveEntity(idx int) //idx是在filter.entities中的id
}

var _ IFilter = &filterBase{}

// 过滤器的基类，实现了IFilter接口
type filterBase struct {
	world *World
	// 过滤器新增/删除entity时的事件监听
	// 外部可以通过OnAdd/OnRemove方法注册/移除监听
	eventListen *FilterEventListen
	// 外部可以通过AddListener/RemoveListener方法注册/移除监听
	listeners []FilterEventListener
	// 过滤器包含的组件类型索引id列表，即过滤器中的entity必定拥有这些组件
	// 比如，包含了A、B、C三个组件，他们的TypeIndex分别为1、3、7，那么IncludeTypeIndices就是[1, 3, 7]
	IncludeTypeIndices []int
	// 过滤器排斥的组件类型索引id列表，即过滤器中的entity必定不会拥有这些组件
	ExcludeTypeIndices []int
	// 过滤器包含的组件列表的掩码，用于快速判断一个entity是否包含这些组件
	IncludeMask uint64
	// 过滤器排斥的组件列表的掩码，用于快速判断一个entity是否包含这些组件
	ExcludeMask uint64
	// 满足过滤器要求的所有entity，
	// 可理解为一个数组，数组的位置和include中的get会是一一对应的
	entities *dataPool.ArrayList[EntityId] //
	// 记录EntityId与entity在数组中位置的映射
	entitiesMap map[EntityId]int //<entityId, 在entities中的id>
	// 用于组合类（子类）实现的接口，用于在添加/移除entity时做一些额外的处理
	afterAddEntityProcesser    afterAddEntityProcesser
	afterRemoveEntityProcesser afterRemoveEntityProcesser
}

func newFilterBase(world *World, ap afterAddEntityProcesser, rp afterRemoveEntityProcesser) *filterBase {
	return &filterBase{
		world:                      world,
		entities:                   dataPool.NewArrayList[EntityId](segmentSize),
		entitiesMap:                make(map[EntityId]int),
		afterAddEntityProcesser:    ap,
		afterRemoveEntityProcesser: rp,
	}
}

// 外部想监听filter中增加entity事件，可用此方法注册
func (f *filterBase) AddListener(listener FilterEventListener) {
	f.listeners = append(f.listeners, listener)
}

// 外部想监听filter中移除entity事件，可用此方法注册
func (f *filterBase) RemoveListener(listener FilterEventListener) {
	for i, l := range f.listeners {
		if l == listener {
			f.listeners = append(f.listeners[:i], f.listeners[i+1:]...)
			break
		}
	}
}

// 通知事件监听：有entity被添加到filter中
func (f *filterBase) notifyAdd(entity Entity) {
	for _, listener := range f.listeners {
		listener.OnEntityAdded(entity)
	}
}

// 通知事件监听：有entity从filter中被移除
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

// 过滤器中加入新的entity
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
	//通知子类执行额外的处理
	f.afterAddEntityProcesser.afterAddEntity(entityData.CompIndices)
	//相当于是直接添加到数组末尾
	idx := f.entities.Count()
	f.entities.Add(entityId)
	f.entitiesMap[entityId] = idx
	f.notifyAdd(entity)
}

// 过滤器中移除entity
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
	//通知子类执行额外的处理
	f.afterRemoveEntityProcesser.afterRemoveEntity(idx)
	delete(f.entitiesMap, entityId)
	//移除数组中的指定下标的数据，
	//直接将元素交换到数组的末尾，然后移除数组的末尾元素，所以非常快。
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

// filterBase1 是一个过滤器，它包含一个Include，用于快速过滤器中Entity的组件数据。
// 它依赖于*filterBase，自动实现了IFilter接口。
// filterBase1[Include1] 表示该过滤器过滤的Entity必须拥有Include1组件。
type filterBase1[Include1 any] struct {
	*filterBase
	// include1 用于快速获取过滤器中Entity的组件数据。具体可参考Include的注释。
	// Include中维持一个数组，存储的是Include1组件在pool中的idx，
	// 数组的下标就是在filterBase.entities中的下标是一一对应的。
	// 比如，get[2] = 10，那么pool.Get(10)就是filter.entities[2]的T组件对象
	include1 *Include[Include1]
}

func newFilterBase1[Include1 any](world *World, compTypeIndex int) *filterBase1[Include1] {
	f := &filterBase1[Include1]{
		include1: newInclude[Include1](world, compTypeIndex),
	}
	f.filterBase = newFilterBase(world, f, f)
	//初始化filterBase.IncludeTypeIndices和filterBase.IncludeMask
	initMask1[Include1](&f.filterBase.IncludeTypeIndices, &f.filterBase.IncludeMask)
	return f
}

// 实现afterAddEntityProcesser接口
func (f *filterBase1[Include1]) afterAddEntity(compIndices map[int]int) {
	for k, v := range compIndices {
		if f.include1.compTypeIndex == k {
			f.include1.AddIdx(k, v)
			break
		}
	}
}

// 实现afterRemoveEntityProcesser接口
// idx是在filter.entities中的id
func (f *filterBase1[Include1]) afterRemoveEntity(idx int) {
	f.include1.FastRemoveAt(idx)
}

// 遍历过滤器中所有entity
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

// 与filterBase1类似，filterBase2要求Entity必须同时包含Include1和Include2组件，
// 即实现过滤同时包含Include1和Include2组件的Entity。
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

// filterBase3实现过滤同时包含Include1、Include2、Include3的Entity。
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

// filterBase3实现过滤同时包含Include1、Include2、Include3、Include4的Entity。
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
