package main

import "github.com/Lei2050/go-ecs"

func initFilters(world *ecs.World) {
	ecs.RegisterFilter(world, ecs.NewFilter4Exclude4[GenderComponent, NameComponent, WalkComponent, AgeComponent, FlyComponent, KindComponent, SwimComponent, BreathInWaterComponent](world))
	ecs.RegisterFilter(world, ecs.NewFilter2Exclude1[FlyComponent, KindComponent, SwimComponent](world))
	ecs.RegisterFilter(world, ecs.NewFilter2[KindComponent, BreathInWaterComponent](world))
	ecs.RegisterFilter(world, ecs.NewFilter3[ImmortalComponent, AgeComponent, NameComponent](world))

	ecs.RegisterFilter(world, ecs.NewFilter1[FlyComponent](world))
	ecs.RegisterFilter(world, ecs.NewFilter1[NameComponent](world))
	ecs.RegisterGroupFilter(world, ecs.NewGroupFilter[NameComponent](world))
	ecs.RegisterFilter(world, ecs.NewFilter1[IdCardComponent](world))
	ecs.RegisterGroupFilter(world, ecs.NewGroupFilterWithKeyMapper[IdCardComponent, ecs.Multi3Key[int, uint64, string], IdCardComponent](world))
	ecs.RegisterFilter(world, ecs.NewFilter2[AgeComponent, NameComponent](world))
	ecs.RegisterGroupFilter(world, ecs.NewGroupFilter2WithKeyMapper[AgeComponent, NameComponent, int, string, AgeComponent, NameComponent](world))

	ecs.RegisterFilter(world, ecs.NewFilter3Exclude2[IdCardComponent, NameComponent, AgeComponent, FlyComponent, BreathInWaterComponent](world))
	ecs.RegisterGroupFilter(world, ecs.NewGroupFilterWithKeyMapper[IdCardComponent, int, IdCardGroupIdMapper](world))
	ecs.RegisterFilter(world, ecs.NewFilter2[IdCardComponent, GenderComponent](world))
	ecs.RegisterGroupFilter(world, ecs.NewGroupFilter2WithKeyMapper[IdCardComponent, GenderComponent, string, int, IdCardGroupProvinceMapper, GenderGroupMapper](world))
}
