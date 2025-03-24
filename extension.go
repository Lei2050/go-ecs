package ecs

import (
	"fmt"
	"reflect"
)

func checkEntity[T any](entity Entity) (*World, *EntityData, *ComponentType) {
	world := entity.World()
	entityData := world.getEntityData(entity.Id)
	if entityData.IsDestroy || !entityData.isCurrentEntityData(entity) {
		panic("entity is not alive")
	}
	componentType := GetComponentType[T]()
	if componentType == nil {
		t := reflect.TypeOf((*T)(nil)).Elem()
		panic(fmt.Sprintf("component type %s is not registered", t.Name()))
	}
	return world, entityData, componentType
}

func Has[T any](entity Entity) bool {
	_, entityData, componentType := checkEntity[T](entity)
	if entityData.CompFlags&componentType.Flag == 0 {
		return false
	}
	_, ok := entityData.CompIndices[componentType.TypeIndex]
	return ok
}

func Replace[T any](entity Entity, component T) {
	world, entityData, componentType := checkEntity[T](entity)

	if entityData.CompFlags&componentType.Flag > 0 {
		dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
		if ok {
			componentType.Events.BeforeUpdate.Invoke(entity)
			pool := getComponentPool[T](world)
			data := pool.GetRef(dataIdx)
			comp := data.(*T)
			// TODO ... process remove filter key
			*comp = component
			// TODO ... process add filter key
		}
	}

	pool := getComponentPool[T](world)
	idx, data := pool.Alloc()
	comp := data.(*T)
	*comp = component
	entityData.CompFlags |= componentType.Flag
	applyComponent[T](world, entity, entityData, idx, componentType)
}

func applyComponent[T any](world *World, entity Entity, entityData *EntityData, compPoolIdx int, componentType *ComponentType) {
	entityData.CompFlags |= componentType.Flag
	entityData.CompIndices[componentType.TypeIndex] = compPoolIdx
	componentType.Events.BeforeAdd.Invoke(entity)
	world.updateFiltersAfterAdd(componentType.TypeIndex, entity, entityData)
	componentType.Events.AfterAdd.Invoke(entity)
}

func TryGet[T any](entity Entity) (*T, bool) {
	return TryGetMayForWrite[T](entity)
}

func TryGetMayForWrite[T any](entity Entity) (*T, bool) {
	world, entityData, componentType := checkEntity[T](entity)

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

func Get[T any](entity Entity) *T {
	return GetMayForWrite[T](entity)
}

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

func GetForWrite[T any](entity Entity) *T {
	world, entityData, componentType := checkEntity[T](entity)
	dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	if !ok {
		t := reflect.TypeOf((*T)(nil)).Elem()
		panic(fmt.Sprintf("entity:%+v not has component:%s", entity, t.Name()))
	}
	componentType.Events.BeforeUpdate.Invoke(entity) //更新通知
	pool := getComponentPool[T](world)
	data := pool.GetRef(dataIdx)
	return data.(*T)
}

func Ensure[T any](entity Entity) *T {
	return EnsureMayForWrite[T](entity)
}

func EnsureMayForWrite[T any](entity Entity) *T {
	world, entityData, componentType := checkEntity[T](entity)
	dataIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	pool := getComponentPool[T](world)
	if ok {
		componentType.Events.BeforeUpdate.Invoke(entity) //更新通知
		data := pool.GetRef(dataIdx)
		return data.(*T)
	}

	idx, data := pool.Alloc()
	entityData.CompFlags |= componentType.Flag
	entityData.CompIndices[componentType.TypeIndex] = idx
	componentType.Events.BeforeAdd.Invoke(entity)
	world.updateFiltersAfterAdd(componentType.TypeIndex, entity, entityData)
	componentType.Events.AfterAdd.Invoke(entity)
	return data.(*T)
}

func MarkDirty[T any](entity Entity) {
	if !Has[T](entity) {
		return
	}
	_, _, componentType := checkEntity[T](entity)
	componentType.Events.BeforeUpdate.Invoke(entity)
}

func Del[T any](entity Entity) bool {
	world, entityData, componentType := checkEntity[T](entity)
	if entityData.CompFlags&componentType.Flag == 0 {
		return false
	}

	gen := entityData.Gen
	_, ok := entityData.CompIndices[componentType.TypeIndex]
	if !ok {
		return false
	}

	world.updateFiltersBeforeRemove(componentType.TypeIndex, entity, entityData)
	componentType.Events.BeforeDelete.Invoke(entity)
	if gen != entity.Gen {
		return false
	}
	//重新寻找一下idx，防止前面的事件执行时改变了idx
	compPoolIdx, ok := entityData.CompIndices[componentType.TypeIndex]
	if ok {
		pool := getComponentPool[T](world)
		pool.Free(compPoolIdx)
		delete(entityData.CompIndices, componentType.TypeIndex)
	}
	componentType.Events.AfterDelete.Invoke(entity)
	return true
}

func destroyEntity(entity Entity) {
	world := entity.World()
	entityData := world.getEntityData(entity.Id)
	if entityData.isCurrentEntityData(entity) {
		return
	}
	if entityData.IsDestroy {
		return
	}
	entityData.IsDestroy = true
	var saveEntity Entity
	saveEntity.Id = entity.Id
	saveEntity.Gen = entity.Gen
	saveEntity.WorldPtr = entity.WorldPtr

	for typeIndex, compPoolIdx := range entityData.CompIndices {
		world.updateFiltersBeforeRemove(typeIndex, saveEntity, entityData)
		pool := world.getComponentPoolByTypeIndex(typeIndex)
		pool.Free(compPoolIdx)
	}

	world.freeEntityData(entity.Id)
}
