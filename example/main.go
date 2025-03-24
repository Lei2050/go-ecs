package main

import (
	"fmt"

	ecs "github.com/Lei2050/go-ecs"
)

func init() {
	ecs.RegisterComponentType[Position](4)
	ecs.RegisterComponentType[Speed](4)
	ecs.RegisterComponentType[Size](4)
	ecs.RegisterComponentType[Fly](4)
	ecs.RegisterComponentType[Walk](4)
	ecs.RegisterComponentType[Gender](4)
	ecs.RegisterComponentType[Name](4)
	ecs.RegisterComponentType[Kind](4)
	ecs.RegisterComponentType[Swim](4)
	ecs.RegisterComponentType[Age](4)
	ecs.RegisterComponentType[Immortal](4)
	ecs.RegisterComponentType[BreathInWater](4)
	ecs.RegisterComponentType[Identification](4)
}

type Position struct {
	X, Y int
}

type Speed struct {
	V float32
}

type Size struct {
	W, H int
}

func main() {
	world := ecs.NewWorld()
	testGroupFilter(world)
}

func testCommon(world *ecs.World) {
	positionFilter := ecs.RegisterFilter(world, ecs.NewFilter1[Position](world))
	speedFilter := ecs.RegisterFilter(world, ecs.NewFilter1[Speed](world))
	sizeFilter := ecs.RegisterFilter(world, ecs.NewFilter1[Size](world))

	allFilter := ecs.RegisterFilter(world, ecs.NewFilter4[Position, Size, Fly, Walk](world))

	player := world.NewEntity()
	ecs.Replace(player, Position{10, 25})
	pos, ok := ecs.TryGet[Position](player)
	fmt.Printf("player pos:%+v, ok:%+v\n", pos, ok)

	pos, ok = ecs.TryGetMayForWrite[Position](player)
	pos.X = 100
	pos.Y = 200
	fmt.Printf("player pos:%+v, ok:%+v\n", pos, ok)
	fmt.Printf("player has[Speed]=%v\n", ecs.Has[Speed](player))
	//pool := ecs.GetComponentPool[Position](world).(*ecs.ComponentPool[Position])
	//fmt.Printf("pool:%+v\n", pool)

	player2 := world.NewEntity()
	ecs.Replace(player2, Position{30, 15})
	pos, ok = ecs.TryGet[Position](player2)
	fmt.Printf("player2 pos:%+v, ok:%+v\n", pos, ok)

	pos, ok = ecs.TryGetMayForWrite[Position](player2)
	pos.X = 101
	pos.Y = 202
	fmt.Printf("player2 pos:%+v, ok:%+v\n", pos, ok)
	fmt.Printf("player2 has[Speed]=%v\n", ecs.Has[Speed](player2))
	ecs.Ensure[Speed](player2).V = 100
	fmt.Printf("player2 has[Speed]=%v\n", ecs.Get[Speed](player2))
	fmt.Println("========================================")
	positionFilter.Foreach(func(entity ecs.Entity, position Position) {
		fmt.Printf("entity:%+v, pos:%+v\n", entity, position)
	})
	speedFilter.Foreach(func(entity ecs.Entity, speed Speed) {
		fmt.Printf("entity:%+v, speed:%+v\n", entity, speed)
	})
	sizeFilter.Foreach(func(entity ecs.Entity, size Size) {
		fmt.Printf("entity:%+v, size:%+v\n", entity, size)
	})
	allFilter.Foreach(func(entity ecs.Entity, position Position, size Size, fly Fly, walk Walk) {
		fmt.Printf("entity:%+v, position:%+v, size:%+v, fly:%+v, walk:%+v\n", entity, position, size, fly, walk)
	})
}

// human
type Identification struct {
	Id      int
	SafeNum uint64
	Code    string
}

var _ ecs.IGroupKeyMap[Identification, ecs.Multi3Key[int, uint64, string]] = Identification{}

func (Identification) MapKey(id Identification) ecs.Multi3Key[int, uint64, string] {
	return ecs.Multi3Key[int, uint64, string]{Key1: id.Id, Key2: uint64(id.SafeNum), Key3: id.Code}
}

type Gender struct {
	Val int
}

type Name struct {
	First, Middle, Last string
}

func (Name) MapKey(name Name) string {
	return name.First
}

type Walk struct {
	Step int
}

// bird
type Fly struct {
}

type Kind struct {
	Name string
}

// fish
type Swim struct{}

type BreathInWater struct{}

// common
type Age struct {
	Val int
}

func (Age) MapKey(age Age) int {
	return age.Val
}

// super man
type Immortal struct{}

func testGroupFilter(world *ecs.World) {
	ecs.RegisterFilter(world, ecs.NewFilter1[Name](world))
	gf := ecs.NewGroupFilter[Name](world)
	ecs.RegisterFilter(world, ecs.NewFilter1[Identification](world))
	idGroupFilter := ecs.NewGroupFilterWithKeyMapper[Identification, ecs.Multi3Key[int, uint64, string], Identification](world)
	ecs.RegisterFilter(world, ecs.NewFilter2[Age, Name](world))
	ageAndNameGroupFilter := ecs.NewGroupFilter2WithKeyMapper[Age, Name, int, string, Age, Name](world)
	testFilters(world)
	entity, hasValue := gf.Find(Name{First: "香香"})
	fmt.Printf("entity:%+v, hasValue:%+v, name:%+v, gender:%+v, age:%+v\n", entity, hasValue, ecs.Get[Name](entity), ecs.Get[Gender](entity), ecs.Get[Age](entity))
	entity, hasValue = idGroupFilter.Find(ecs.Multi3Key[int, uint64, string]{Key1: 1002, Key2: 1992, Key3: "1314-2211"})
	fmt.Printf("entity:%+v, hasValue:%+v, name:%+v, gender:%+v, age:%+v\n", entity, hasValue, ecs.Get[Name](entity), ecs.Get[Gender](entity), ecs.Get[Age](entity))
	entity, hasValue = ageAndNameGroupFilter.Find(ecs.Multi2Key[int, string]{Key1: 34, Key2: "香香"})
	fmt.Printf("entity:%+v, hasValue:%+v, name:%+v, gender:%+v, age:%+v\n", entity, hasValue, ecs.Get[Name](entity), ecs.Get[Gender](entity), ecs.Get[Age](entity))
}

func testFilters(world *ecs.World) {
	humanFilter := ecs.RegisterFilter(world, ecs.NewFilter4Exclude4[Gender, Name, Walk, Age, Fly, Kind, Swim, BreathInWater](world))
	birdFilter := ecs.RegisterFilter(world, ecs.NewFilter2Exclude1[Fly, Kind, Swim](world))
	fishFilter := ecs.RegisterFilter(world, ecs.NewFilter2[Kind, BreathInWater](world))
	superManFilter := ecs.RegisterFilter(world, ecs.NewFilter3[Immortal, Age, Name](world))

	spawnHuman(world, 1001, 2005, "1314-2211", 0, "Jack", 18)
	spawnHuman(world, 1002, 1992, "1314-2211", 1, "香香", 34)
	spawnBird(world, "啄木鸟")
	spawnBird(world, "麻雀")
	spawnFish(world, "鲸鱼")
	spawnFish(world, "沙丁鱼")

	superMan := world.NewEntity()
	ecs.Replace(superMan, Immortal{})
	ecs.Replace(superMan, Age{Val: 1000})
	ecs.Replace(superMan, Name{First: "超人"})
	ecs.Replace(superMan, Fly{})
	ecs.Replace(superMan, BreathInWater{})

	birdMan := spawnHuman(world, 1003, 3333, "3333-3333", 2, "鸟人", 25)
	ecs.Replace(birdMan, Fly{})
	fishMan := spawnHuman(world, 1004, 4444, "3333-3333", 3, "鱼人", 25)
	ecs.Replace(fishMan, BreathInWater{})

	god := spawnHuman(world, 110, 1000, "0000-0000", 4, "玉皇大帝", 10000000000)
	ecs.Replace(god, Immortal{})
	ecs.Replace(god, Fly{})
	ecs.Replace(god, Kind{Name: "god"})
	ecs.Replace(god, Swim{})
	ecs.Replace(god, BreathInWater{})

	humanFilter.Foreach(func(entity ecs.Entity, gender Gender, name Name, walk Walk, age Age) {
		fmt.Printf("entity:%+v, walking, gender:%+v, name:%+v, age:%+v\n", entity, gender, name, age)
	})
	birdFilter.Foreach(func(entity ecs.Entity, fly Fly, kind Kind) {
		fmt.Printf("entity:%+v, flying, kind:%+v\n", entity, kind)
	})
	fishFilter.Foreach(func(entity ecs.Entity, kind Kind, breath BreathInWater) {
		fmt.Printf("entity:%+v, swimming, kind:%+v\n", entity, kind)
	})
	superManFilter.Foreach(func(entity ecs.Entity, immortal Immortal, age Age, name Name) {
		fmt.Printf("entity:%+v, super man, immortal:%+v, age:%+v, name:%+v\n", entity, immortal, age, name)
	})
}

func spawnHuman(world *ecs.World, id int, safeNum uint64, code string, gender int, name string, age int) ecs.Entity {
	people := world.NewEntity()
	ecs.Replace(people, Identification{Id: id, SafeNum: safeNum, Code: code})
	ecs.Replace(people, Gender{Val: gender})
	ecs.Replace(people, Name{First: name})
	ecs.Replace(people, Age{Val: age})
	ecs.Replace(people, Walk{})
	return people
}

func spawnBird(world *ecs.World, kind string) ecs.Entity {
	bird := world.NewEntity()
	ecs.Replace(bird, Fly{})
	ecs.Replace(bird, Kind{Name: kind})
	return bird
}

func spawnFish(world *ecs.World, kind string) ecs.Entity {
	fish := world.NewEntity()
	ecs.Replace(fish, Kind{Name: kind})
	ecs.Replace(fish, BreathInWater{})
	return fish
}

func testNewEntity(world *ecs.World) {
	entity1 := world.NewEntity()
	fmt.Printf("%+v\n", entity1)
	fmt.Printf("entity1 is Alive:%+v\n", entity1.IsAlive())
	entity2 := world.NewEntity()
	fmt.Printf("%+v\n", entity2)
	fmt.Printf("1==2 ? %+v\n", entity1.Equal(entity2))
	entity1.Destroy()
	fmt.Printf("entity1 is Alive:%+v\n", entity1.IsAlive())

	entity3 := world.NewEntity()
	fmt.Printf("entity3:%+v, 1==3 ? %v\n", entity3, entity1.Equal(entity3))

}
