package ecs

import (
	"unsafe"
)

// 一个实体，包含一个ID和一个版本号
type Entity struct {
	Id  int
	Gen uint
	// 指向实体所属世界的指针
	WorldPtr uintptr
}

func (e *Entity) World() *World {
	return (*World)(unsafe.Pointer(e.WorldPtr))
}

// 比较两个实体是否相等
// 只有ID和版本号都相等时，两个实体才相等
func (e *Entity) Equal(entity Entity) bool {
	return e.Id == entity.Id && e.Gen == entity.Gen
}

// 销毁实体
func (e *Entity) Destroy() {
	destroyEntity(*e)
}

// 检查实体是否存活
// 如果实体的版本号与当前实体数据的版本号不相等，说明实体已经被销毁过，视为不存活
func (e *Entity) IsAlive() bool {
	world := e.World()
	entityData := world.getEntityData(e.Id)
	return !entityData.IsDestroying && entityData.isCurrentEntityData(*e)
}

func (e *Entity) getEntityData() *EntityData {
	world := e.World()
	return world.getEntityData(e.Id)
}

func (e *Entity) GetId() EntityId {
	return EntityId{
		Id:  e.Id,
		Gen: e.Gen,
	}
}

type EntityId struct {
	Id  int
	Gen uint
}
