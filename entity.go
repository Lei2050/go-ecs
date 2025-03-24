package ecs

import (
	"unsafe"
)

type Entity struct {
	Id       int
	Gen      uint
	WorldPtr uintptr
}

func (e *Entity) World() *World {
	return (*World)(unsafe.Pointer(e.WorldPtr))
}

func (e *Entity) Equal(entity Entity) bool {
	return e.Id == entity.Id && e.Gen == entity.Gen
}

func (e *Entity) Destroy() {
	destroyEntity(*e)
}

func (e *Entity) IsAlive() bool {
	world := e.World()
	entityData := world.getEntityData(e.Id)
	return !entityData.IsDestroy && entityData.isCurrentEntityData(*e)
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
