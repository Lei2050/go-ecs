package ecs

import (
	"fmt"
	"reflect"

	dataPool "github.com/Lei2050/array-pool"
)

var TypeIndex int

// 所有已注册的组件，//<组件反射类型, 组件类型数据>
// 目前，使用组件之前必须先使用RegisterComponentType注册组件类型
var componentTypeMap = make(map[reflect.Type]*ComponentType)

// 组件类型数据
type ComponentType struct {
	// 组件类型索引，也可以理解为组件类型的ID，组件类型的索引是唯一的
	TypeIndex int
	// 组件反射类型
	Type reflect.Type
	// 组件类型的标志位，用于entity标识布隆过滤器
	// 等于 1 << (TypeIndex % 64)
	// 所以，不同组件可能会有相同的Flag
	Flag uint64
	// 对象池的段大小，这个参数开放给上层好像很迷惑
	PoolSegmentSize int
	// 与该组件变更的相关事件，外部可以通过它来监听具体组件的变更
	Events EntityEvents
}

// 注册组件类型
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

// 获取组件类型数据
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

// 基于数组的组件对象池
// 类型T需要是一个结构体才能发挥作用。如果是一个指针，没啥意义。
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
