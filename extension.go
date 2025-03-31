package ecs

import (
	"fmt"
	"reflect"
)

// checkEntity 用于检查Entity的有效性、是否存活，
// 返回值分别为Entity所在的世界实例指针、Entity数据指针和组件类型指针。
func checkEntity[T any](entity Entity) (*World, *EntityData, *ComponentType) {
	world := entity.World()
	entityData := world.getEntityData(entity.Id)
	if !entityData.isCurrentEntityData(entity) {
		//不允许使用非存活的entity
		panic("entity is not alive")
	}
	componentType := GetComponentType[T]()
	if componentType == nil {
		t := reflect.TypeOf((*T)(nil)).Elem()
		//不允许使用未注册的组件
		panic(fmt.Sprintf("component type %s is not registered", t.Name()))
	}
	return world, entityData, componentType
}

// Has 用于检查Entity是否拥有指定类型的组件。
func Has[T any](entity Entity) bool {
	_, entityData, componentType := checkEntity[T](entity)
	// 根据布隆过滤器的特性，如果标志位未设置，则说明必定不存在该组件
	if entityData.CompFlags&componentType.Flag == 0 {
		return false
	}
	// 标志位设置了，但目标组件可能不存在，需要进一步检查
	_, ok := entityData.CompIndices[componentType.TypeIndex]
	return ok
}

// Replace 用于附加/替换Entity的指定组件。
func Replace[T any](entity Entity, component T) {
	world, entityData, componentType := checkEntity[T](entity)

	// 检查Entity是否已经拥有该组件
	if entityData.CompFlags&componentType.Flag > 0 {
		dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
		if ok { //entity拥有该组件
			// 触发组件更新前的事件
			componentType.Events.BeforeUpdate.Invoke(entity)
			pool := getComponentPool[T](world)
			data := pool.GetRef(dataIdx)
			comp := data.(*T)
			//这里的用groupKey事件通知groupFilter操作entity，而不是用FilterEventListener等，
			//因为这里只是替换了Comp数据，entity必定还在filter中，只是其GroupKey发生变化，
			//只需要变更groupFilter中的集合数据。
			//触发一下groupKey移除的事件，通知相关groupFilter移除entitty
			world.fireGroupKeyEvent(componentType.TypeIndex, groupKeyRemove, entity)
			// 替换组件数据
			*comp = component
			//触发一下groupKey增加的事件，通知相关groupFilter移除entitty
			world.fireGroupKeyEvent(componentType.TypeIndex, groupKeyAdd, entity)
			return
		}
	}

	//组件不存在

	pool := getComponentPool[T](world)
	idx, data := pool.Alloc()
	comp := data.(*T)
	*comp = component
	// 应用新的组件到Entity
	applyComponent[T](world, entity, entityData, idx, componentType)
}

// applyComponent 用于将指定组件应用到Entity上。
// compPoolIdx 是组件在组件对象池中的索引。
func applyComponent[T any](world *World, entity Entity, entityData *EntityData, compPoolIdx int, componentType *ComponentType) {
	// 设置Entity的组件布隆过滤器
	entityData.CompFlags |= componentType.Flag
	// 记录Entity的组件索引信息
	entityData.CompIndices[componentType.TypeIndex] = compPoolIdx
	// 触发组件添加前的事件
	componentType.Events.BeforeAdd.Invoke(entity)
	// 触发组件添加前的事件，compPoolIdx可用于有关联的component的快速获取，
	// 比如典型的ParentComponent，如果是高频调用、可以缓存compPoolIdx以加速获取Parent。
	componentType.Events.BeforeAddWithPoolIdx.Invoke(entity, compPoolIdx)
	// 新组件添加，world通知相关过滤器执行更新
	world.updateFiltersAfterAdd(componentType.TypeIndex, entity, entityData)
	// 触发组件添加后的事件
	componentType.Events.AfterAdd.Invoke(entity)
	componentType.Events.AfterAddWithPoolIdx.Invoke(entity, compPoolIdx)
}

// TryGet 尝试获取Entity的指定组件。
// 返回值为组件指针和布尔类型，若获取成功则返回组件指针和 true，否则返回 nil 和 false。
// 注意，不要长期持有返回的指针，指向的对象可能频繁地被回收/变更/复用；
// 注意，不要再堆上持有返回的指针，除非你明确了解其生命周期。
func TryGet[T any](entity Entity) (*T, bool) {
	return TryGetMayForWrite[T](entity)
}

// TryGetMayForWrite 尝试获取Entity的指定组件，可能用于写入操作。
// 返回值为组件指针和布尔类型，若获取成功则返回组件指针和 true，否则返回 nil 和 false。
// 注意，不要长期持有返回的指针，指向的对象可能频繁地被回收/变更/复用；
// 注意，不要再堆上持有返回的指针，除非你明确了解其生命周期。
func TryGetMayForWrite[T any](entity Entity) (*T, bool) {
	world, entityData, componentType := checkEntity[T](entity)

	// 检查Entity是否拥有该组件
	if entityData.CompFlags&componentType.Flag == 0 { //没有该comp
		return nil, false
	}

	typeIndex := componentType.TypeIndex
	dataIdx, ok := entityData.CompIndices[typeIndex]
	if ok {
		pool := getComponentPool[T](world)
		data := pool.GetRef(dataIdx)
		return data.(*T), true
	}

	return nil, false
}

// Get 用于获取Entity的指定组件，其假定Entity拥有该组件；
// 若Entity不拥有该组件则触发 panic。
func Get[T any](entity Entity) *T {
	return GetMayForWrite[T](entity)
}

// GetMayForWrite 用于获取Entity的指定组件，可能用于写入操作，其假定Entity拥有该组件；
// 返回值为组件指针，若Entity不拥有该组件则触发 panic。
func GetMayForWrite[T any](entity Entity) *T {
	world, entityData, componentType := checkEntity[T](entity)
	dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem().Elem()
		panic(fmt.Sprintf("entity:%+v not has component:%s", entity, t.Name()))
	}
	pool := getComponentPool[T](world)
	data := pool.GetRef(dataIdx)
	return data.(*T)
}

// GetForWrite 用于获取Entity的指定组件，用于写入操作。
// 其假定Entity拥有该组件，并且假定用户调用后必定会修改该组件的数据。
// 返回值为组件指针，若Entity不拥有该组件则触发 panic。
func GetForWrite[T any](entity Entity) *T {
	world, entityData, componentType := checkEntity[T](entity)
	dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem()
		panic(fmt.Sprintf("entity:%+v not has component:%s", entity, t.Name()))
	}
	// 触发组件更新前的事件
	componentType.Events.BeforeUpdate.Invoke(entity) //更新通知
	pool := getComponentPool[T](world)
	data := pool.GetRef(dataIdx)
	return data.(*T)
}

// Ensure 确保Entity拥有指定组件。
// 返回值为组件指针，若Entity没有该组件则添加并返回。
func Ensure[T any](entity Entity) *T {
	return EnsureMayForWrite[T](entity)
}

// EnsureMayForWrite 确保Entity拥有指定组件，可能用于写入操作。
// 返回值为组件指针，若Entity没有该组件则添加并返回。
func EnsureMayForWrite[T any](entity Entity) *T {
	world, entityData, componentType := checkEntity[T](entity)
	dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	pool := getComponentPool[T](world)
	if ok {
		data := pool.GetRef(dataIdx)
		return data.(*T)
	}

	idx, data := pool.Alloc()
	// 应用新的组件到Entity
	applyComponent[T](world, entity, entityData, idx, componentType)

	return data.(*T)
}

// MarkDirty 标记Entity的指定组件为脏数据。
// 若Entity拥有该组件，则触发组件更新前的事件。
func MarkDirty[T any](entity Entity) {
	if !Has[T](entity) {
		return
	}
	_, _, componentType := checkEntity[T](entity)
	// 触发组件更新前的事件
	componentType.Events.BeforeUpdate.Invoke(entity)
}

// Del 用于删除Entity的指定组件。
// 删除成功则返回 true，否则返回 false。
func Del[T any](entity Entity) bool {
	world, entityData, componentType := checkEntity[T](entity)
	// 检查Entity是否拥有该组件
	if entityData.CompFlags&componentType.Flag == 0 {
		return false
	}

	gen := entityData.Gen
	_, ok := entityData.CompIndices[componentType.TypeIndex]
	if !ok {
		return false
	}

	// 组件删除，world通知相关过滤器执行更新
	world.updateFiltersBeforeRemove(componentType.TypeIndex, entity, entityData)
	// 触发组件删除前的事件
	componentType.Events.BeforeDelete.Invoke(entity)
	if gen != entity.Gen { //期间执行事件导致entity销毁过了？重复删除？
		return false
	}
	//重新寻找一下idx，防止前面的事件执行时改变了idx
	compPoolIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	if ok {
		pool := getComponentPool[T](world)
		pool.Free(compPoolIdx)
		delete(entityData.CompIndices, componentType.TypeIndex)
	}
	// 触发组件删除后的事件
	componentType.Events.AfterDelete.Invoke(entity)
	return true
}

// destroyEntity 函数用于销毁指定的Entity。
func destroyEntity(entity Entity) {
	world := entity.World()
	entityData := world.getEntityData(entity.Id)
	if !entityData.isCurrentEntityData(entity) {
		return
	}
	if entityData.IsDestroying {
		return
	}
	// 标记Entity为已销毁状态
	entityData.IsDestroying = true
	var saveEntity Entity
	saveEntity.Id = entity.Id
	saveEntity.Gen = entity.Gen
	saveEntity.WorldPtr = entity.WorldPtr

	// 遍历Entity的所有组件索引
	for typeIndex, compPoolIdx := range entityData.CompIndices {
		// 更新所有相关过滤器
		world.updateFiltersBeforeRemove(typeIndex, saveEntity, entityData)
		pool := world.getComponentPoolByTypeIndex(typeIndex)
		pool.Free(compPoolIdx)
	}

	// 回收EntityData
	world.freeEntityData(entity.Id)
}
