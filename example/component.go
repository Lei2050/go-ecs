package main

import ecs "github.com/Lei2050/go-ecs"

func init() {
	ecs.RegisterComponentType[FlyComponent](4)
	ecs.RegisterComponentType[WalkComponent](4)
	ecs.RegisterComponentType[GenderComponent](4)
	ecs.RegisterComponentType[NameComponent](4)
	ecs.RegisterComponentType[KindComponent](4)
	ecs.RegisterComponentType[SwimComponent](4)
	ecs.RegisterComponentType[AgeComponent](4)
	ecs.RegisterComponentType[ImmortalComponent](4)
	ecs.RegisterComponentType[BreathInWaterComponent](4)
	ecs.RegisterComponentType[IdCardComponent](4)
}

// human
type IdCardComponent struct {
	Id       int
	Security uint64
	Province string
}

var _ ecs.IGroupKeyMap[IdCardComponent, ecs.Multi3Key[int, uint64, string]] = IdCardComponent{}

func (IdCardComponent) MapKey(id IdCardComponent) ecs.Multi3Key[int, uint64, string] {
	return ecs.Multi3Key[int, uint64, string]{Key1: id.Id, Key2: uint64(id.Security), Key3: id.Province}
}

type IdCardGroupIdMapper struct{}

func (IdCardGroupIdMapper) MapKey(id IdCardComponent) int {
	return id.Id
}

type IdCardGroupProvinceMapper struct{}

func (IdCardGroupProvinceMapper) MapKey(id IdCardComponent) string {
	return id.Province
}

type GenderComponent struct {
	Val int
}
type GenderGroupMapper struct{}

func (GenderGroupMapper) MapKey(gender GenderComponent) int {
	return gender.Val
}

type NameComponent struct {
	First, Middle, Last string
}

func (NameComponent) MapKey(name NameComponent) string {
	return name.First
}

type WalkComponent struct {
	Step int
}

// bird
type FlyComponent struct {
}

type KindComponent struct {
	Name string
}

// fish
type SwimComponent struct{}

type BreathInWaterComponent struct{}

// common
type AgeComponent struct {
	Val int
}

func (AgeComponent) MapKey(age AgeComponent) int {
	return age.Val
}

// super man
type ImmortalComponent struct{}

/*
	BeforeAdd Delegate
	AfterAdd  Delegate

	BeforeUpdate Delegate

	BeforeDelete Delegate
	AfterDelete  Delegate

	BeforeAddWithPoolIdx DelegateWithParam
	AfterAddWithPoolIdx  DelegateWithParam
*/
func registerCompChangeEvent[T any](beforeAdd func(ecs.Entity),
	afterAdd func(ecs.Entity),
	beforeUpdate func(ecs.Entity),
	beforeDelete func(ecs.Entity),
	afterDelete func(ecs.Entity),
	beforeAddWithPoolIdx func(ecs.Entity, ...any),
	afterAddWithPoolIdx func(ecs.Entity, ...any),
) {
	componentType := ecs.GetComponentType[T]()
	if beforeAdd != nil {
		componentType.Events.BeforeAdd.AddCallback(beforeAdd)
	}
	if afterAdd != nil {
		componentType.Events.AfterAdd.AddCallback(afterAdd)
	}
	if beforeUpdate != nil {
		componentType.Events.BeforeUpdate.AddCallback(beforeUpdate)
	}
	if beforeDelete != nil {
		componentType.Events.BeforeDelete.AddCallback(beforeDelete)
	}
	if afterDelete != nil {
		componentType.Events.AfterDelete.AddCallback(afterDelete)
	}
	if beforeAddWithPoolIdx != nil {
		componentType.Events.BeforeAddWithPoolIdx.AddCallback(beforeAddWithPoolIdx)
	}
	if afterAddWithPoolIdx != nil {
		componentType.Events.AfterAddWithPoolIdx.AddCallback(afterAddWithPoolIdx)
	}
}
