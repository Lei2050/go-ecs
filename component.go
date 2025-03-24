package ecs

import (
	"fmt"
	"reflect"

	dataPool "github.com/Lei2050/array-pool"
)

var TypeIndex int

var componentTypeMap = make(map[reflect.Type]*ComponentType)

type ComponentType struct {
	TypeIndex int
	Type      reflect.Type
	Flag      uint64

	PoolSegmentSize int

	Events EntityEvents
}

func RegisterComponentType[T any](poolSegmentSize int) *ComponentType {
	t := reflect.TypeOf((*T)(nil)).Elem()
	TypeIndex++
	ct := &ComponentType{
		TypeIndex:       TypeIndex,
		Type:            t,
		Flag:            1 << (TypeIndex % 64),
		PoolSegmentSize: poolSegmentSize,
	}
	componentTypeMap[t] = ct
	return ct
}

func GetComponentType[T any]() *ComponentType {
	t := reflect.TypeOf((*T)(nil)).Elem()
	componentType, ok := componentTypeMap[t]
	if !ok {
		panic(fmt.Sprintf("component:%+v not register", t.Name()))
	}
	return componentType
}

type ComponentPooler interface {
	Alloc() (int, any)
	GetRef(id int) any
	Free(id int)
}

type ComponentPool[T any] struct {
	pool *dataPool.Pool[T]
}

func (cp *ComponentPool[T]) Alloc() (int, any) {
	return cp.pool.Alloc()
}

func (cp *ComponentPool[T]) GetRef(id int) any {
	return cp.pool.GetRef(id)
}

func (cp *ComponentPool[T]) Free(id int) {
	cp.pool.Free(id)
}
