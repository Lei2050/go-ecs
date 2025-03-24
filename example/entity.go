package main

import ecs "github.com/Lei2050/go-ecs"

func spawnHuman(world *ecs.World, id int, safeNum uint64, code string, gender int, name string, age int) ecs.Entity {
	people := world.NewEntity()
	ecs.Replace(people, IdCardComponent{Id: id, Security: safeNum, Province: code})
	ecs.Replace(people, GenderComponent{Val: gender})
	ecs.Replace(people, NameComponent{First: name})
	ecs.Replace(people, AgeComponent{Val: age})
	ecs.Replace(people, WalkComponent{})
	return people
}

func spawnBird(world *ecs.World, kind string) ecs.Entity {
	bird := world.NewEntity()
	ecs.Replace(bird, FlyComponent{})
	ecs.Replace(bird, KindComponent{Name: kind})
	return bird
}

func spawnFish(world *ecs.World, kind string) ecs.Entity {
	fish := world.NewEntity()
	ecs.Replace(fish, KindComponent{Name: kind})
	ecs.Replace(fish, BreathInWaterComponent{})
	return fish
}
