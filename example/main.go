package main

import (
	"fmt"

	ecs "github.com/Lei2050/go-ecs"
)

func main() {
	//component.init()
	world := ecs.NewWorld()
	initFilters(world)

	//筛选所有正常人，有身份证、名字、年龄，但不会飞、不能在水里呼吸
	humanFilter := ecs.GetFilter[*ecs.Filter3Exclude2[IdCardComponent, NameComponent, AgeComponent, FlyComponent, BreathInWaterComponent]](world)
	//监听humanFilter里的entity变化
	humanFilter.AddListener(&humanFilterListener{})
	humanFilter.OnAdd(func(entity ecs.Entity) {
		idCard := ecs.Get[IdCardComponent](entity)
		name := ecs.Get[NameComponent](entity)
		fmt.Printf("    -human entity:%+v add, id:%+v, name:%+v\n", entity, idCard, name)
	})

	//监听IdCardComponent变化
	registerCompChangeEvent[IdCardComponent](func(e ecs.Entity) {
		fmt.Printf("    EVENT: beforeAdd IdCardComponent entity:%+v change\n", e)
	},
		nil,
		//func(e ecs.Entity) {
		//	fmt.Printf("    EVENT: afterAdd IdCardComponent entity:%+v change\n", e)
		//},
		func(e ecs.Entity) {
			fmt.Printf("    EVENT: beforeUpdate IdCardComponent entity:%+v change\n", e)
		},
		func(e ecs.Entity) {
			fmt.Printf("    EVENT: beforeDelete IdCardComponent entity:%+v change\n", e)
		},
		func(e ecs.Entity) {
			fmt.Printf("    EVENT: afterDelete IdCardComponent entity:%+v change\n", e)
		},
		//func(e ecs.Entity, params ...any) {
		//	compPoolIdx := params[0].(int)
		//	fmt.Printf("    EVENT: beforeAddWithPoolIdx IdCardComponent entity:%+v change, compPoolIdx:%+v\n", e, compPoolIdx)
		//},
		nil,
		func(e ecs.Entity, params ...any) {
			compPoolIdx := params[0].(int)
			fmt.Printf("    EVENT: afterAddWithPoolIdx IdCardComponent entity:%+v change, compPoolIdx:%+v\n", e, compPoolIdx)
		},
	)

	//生成entity，并添加组件
	leilei := spawnHuman(world, 9527, 123456, "浙江", 0, "雷雷", 25)
	xx := spawnHuman(world, 9381, 567890, "湖北", 1, "香香", 18)
	bx := spawnHuman(world, 9382, 567891, "湖北", 1, "不香", 20)
	xx2 := spawnHuman(world, 9383, 567892, "湖北", 0, "香香", 25)
	spawnBird(world, "啄木鸟")
	spawnBird(world, "麻雀")
	spawnFish(world, "鲸鱼")
	spawnFish(world, "沙丁鱼")
	//超人啥组件都有
	superMan := world.NewEntity()
	ecs.Replace(superMan, ImmortalComponent{})
	ecs.Replace(superMan, AgeComponent{Val: 1000})
	ecs.Replace(superMan, NameComponent{First: "超人"})
	ecs.Replace(superMan, FlyComponent{})
	ecs.Replace(superMan, BreathInWaterComponent{})
	//变异鸟人
	birdMan := spawnHuman(world, 1003, 3333, "上海", 2, "鸟人", 25)
	ecs.Replace(birdMan, FlyComponent{})
	//变异鱼人
	fishMan := spawnHuman(world, 1004, 4444, "海南", 3, "鱼人", 25)
	ecs.Replace(fishMan, BreathInWaterComponent{})
	//玉皇大帝
	god := spawnHuman(world, 110, 1000, "天庭", 0, "玉皇大帝", 10000000000)
	ecs.Replace(god, ImmortalComponent{})
	ecs.Replace(god, FlyComponent{})
	ecs.Replace(god, KindComponent{Name: "god"})
	ecs.Replace(god, SwimComponent{})
	ecs.Replace(god, BreathInWaterComponent{})

	//修改超人的年龄
	ecs.GetForWrite[AgeComponent](superMan).Val = 100
	kindComp, ok := ecs.TryGetMayForWrite[KindComponent](superMan) //尝试修改
	//false，超人没有kind组件
	if ok {
		kindComp.Name = "god"
	} else {
		fmt.Println("superMan has no kind component")
	}
	//玉皇大帝没有身份证
	fmt.Printf("god has IdCardComponent? %+v\n", ecs.Has[IdCardComponent](god))
	fmt.Println("---------------------")

	//目前所有Filter在使用之前都要先Register，并且要在world.NewEntity()之前注册
	//如果在world.NewEntity()之后注册，会导致Filter中的entity数量不准确

	//筛选所有能飞的entity
	flyFilter := ecs.GetFilter[*ecs.Filter1[FlyComponent]](world)
	flyFilter.Foreach(func(entity ecs.Entity, fly FlyComponent) {
		fmt.Printf("entity:%+v, can flying\n", entity)
	})
	fmt.Println("---------------------")
	humanFilter.Foreach(func(entity ecs.Entity, idCard IdCardComponent, name NameComponent, age AgeComponent) {
		fmt.Printf("entity:%+v, human, idCard:%+v, name:%+v, age:%+v\n", entity, idCard, name, age)
	})
	fmt.Println("---------------------")

	//检索身份证id为9381的人
	idGroupFilter := ecs.GetGroupFilter[*ecs.GroupFilterWithKeyMapper[IdCardComponent, int, IdCardGroupIdMapper]](world)
	entity, ok := idGroupFilter.FindOne(9381)
	if ok {
		fmt.Printf("entity:%+v, human, id:%+v, name:%+v, age:%+v\n", entity, ecs.Get[IdCardComponent](entity), ecs.Get[NameComponent](entity), ecs.Get[AgeComponent](entity))
	}
	fmt.Println("after changing xx's id ---------------------")
	ecs.Replace(xx, IdCardComponent{9681, 567890, "湖北"})
	entity, ok = idGroupFilter.FindOne(9381)
	if ok {
		fmt.Printf("entity:%+v, human, id:%+v, name:%+v, age:%+v\n", entity, ecs.Get[IdCardComponent](entity), ecs.Get[NameComponent](entity), ecs.Get[AgeComponent](entity))
	} else {
		fmt.Printf("not find id=9381\n")
	}
	entity, ok = idGroupFilter.FindOne(9681)
	if ok {
		fmt.Printf("entity:%+v, human, id:%+v, name:%+v, age:%+v\n", entity, ecs.Get[IdCardComponent](entity), ecs.Get[NameComponent](entity), ecs.Get[AgeComponent](entity))
	} else {
		fmt.Printf("not find id=9681\n")
	}
	fmt.Println("---------------------")
	//检索所有省份为湖北、性别为0的人
	xxGroupFilter := ecs.GetGroupFilter[*ecs.GroupFilter2WithKeyMapper[IdCardComponent, GenderComponent, string, int, IdCardGroupProvinceMapper, GenderGroupMapper]](world)
	xxGroupFilter.Foreach(ecs.Multi2Key[string, int]{Key1: "湖北", Key2: 0}, func(entity ecs.Entity) {
		fmt.Printf("entity:%+v, human, id:%+v, gender:%+v, name:%+v, age:%+v\n",
			entity, ecs.Get[IdCardComponent](entity), ecs.Get[GenderComponent](entity), ecs.Get[NameComponent](entity), ecs.Get[AgeComponent](entity))
	})

	fmt.Println("going to destroy ---------------------")
	leilei.Destroy()
	ecs.Del[IdCardComponent](xx2)
	xx2.Destroy()
	ecs.Del[AgeComponent](bx)
	xx.Destroy()
}

type humanFilterListener struct{}

func (h *humanFilterListener) OnEntityAdded(entity ecs.Entity) {
	idCard := ecs.Get[IdCardComponent](entity)
	name := ecs.Get[NameComponent](entity)
	fmt.Printf("    =human entity:%+v add, id:%+v, name:%+v\n", entity, idCard, name)
}

func (h *humanFilterListener) OnEntityRemoved(entity ecs.Entity) {
	fmt.Printf("    =human entity:%+v remove\n", entity)
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
