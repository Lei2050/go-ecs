package ecs

const (
	groupKeyAdd groupKeyEventKind = iota + 1
	groupKeyRemove
)

type groupKeyEventKind = int

type groupKeyEventHandler interface {
	onGroupKeyEvent(key groupKeyEventKind, entity Entity)
}

type groupKeyEvent struct {
	eventKind groupKeyEventKind
	handler   groupKeyEventHandler
}

var _ groupKeyEventHandler = (*groupKeyEventProxy)(nil)

type groupKeyEventProxy struct {
	set    iEntitySet
	filter IFilter
}

func (g *groupKeyEventProxy) onGroupKeyEvent(eventKind groupKeyEventKind, entity Entity) {
	switch eventKind {
	case groupKeyAdd:
		if g.filter.isCompatibleAfterAddIncluded(entity.getEntityData()) {
			g.set.Add(entity)
		}
	case groupKeyRemove:
		g.set.Remove(entity)
	}
}

// 键值生成器：从Entity中提取Key
type groupKeyMaker[Key comparable] interface {
	makeKey(Entity) Key
}

// 实现groupKeyMaker[KeyComp]接口
type DirectlyKeyMaker[KeyComp comparable] struct{}

// 直接将entity的KeyComp组件数据作为key
func (d DirectlyKeyMaker[KeyComp]) makeKey(entity Entity) KeyComp {
	return *Get[KeyComp](entity)
}

type Directly2KeyMaker[KeyComp1, KeyComp2 comparable] struct{}
type Multi2Key[Key1, Key2 comparable] struct {
	Key1 Key1
	Key2 Key2
}

func (d Directly2KeyMaker[KeyComp1, KeyComp2]) makeKey(entity Entity) Multi2Key[KeyComp1, KeyComp2] {
	return Multi2Key[KeyComp1, KeyComp2]{*Get[KeyComp1](entity), *Get[KeyComp2](entity)}
}

type Directly3KeyMaker[KeyComp1, KeyComp2, KeyComp3 comparable] struct{}
type Multi3Key[Key1, Key2, Key3 comparable] struct {
	Key1 Key1
	Key2 Key2
	Key3 Key3
}

func (d Directly3KeyMaker[KeyComp1, KeyComp2, KeyComp3]) makeKey(entity Entity) Multi3Key[KeyComp1, KeyComp2, KeyComp3] {
	return Multi3Key[KeyComp1, KeyComp2, KeyComp3]{*Get[KeyComp1](entity), *Get[KeyComp2](entity), *Get[KeyComp3](entity)}
}

// 接口功能：将Source转换为Key
type IGroupKeyMap[Source any, Key comparable] interface {
	MapKey(Source) Key
}

// 将Source转换为Key的键值生成器
// 实现groupKeyMaker[Key]接口
type groupKeyMapper[Source any, Key comparable, Mapper IGroupKeyMap[Source, Key]] struct {
	mapper Mapper
}

// 直接将entity的Source组件数据，通过指定的mapper转化为key
func (g groupKeyMapper[Source, Key, Mapper]) makeKey(entity Entity) Key {
	return g.mapper.MapKey(*Get[Source](entity))
}

type group2KeyMapper[Source1, Source2 any, Key1, Key2 comparable, Mapper1 IGroupKeyMap[Source1, Key1], Mapper2 IGroupKeyMap[Source2, Key2]] struct {
	mapper1 Mapper1
	mapper2 Mapper2
}

func (g group2KeyMapper[Source1, Source2, Key1, Key2, Mapper1, Mapper2]) makeKey(entity Entity) Multi2Key[Key1, Key2] {
	return Multi2Key[Key1, Key2]{g.mapper1.MapKey(*Get[Source1](entity)), g.mapper2.MapKey(*Get[Source2](entity))}
}

type group3KeyMapper[
	Source1, Source2, Source3 any,
	Key1, Key2, Key3 comparable,
	Mapper1 IGroupKeyMap[Source1, Key1], Mapper2 IGroupKeyMap[Source2, Key2], Mapper3 IGroupKeyMap[Source3, Key3]] struct {
	mapper1 Mapper1
	mapper2 Mapper2
	mapper3 Mapper3
}

func (g group3KeyMapper[Source1, Source2, Source3, Key1, Key2, Key3, Mapper1, Mapper2, Mapper3]) makeKey(entity Entity) Multi3Key[Key1, Key2, Key3] {
	return Multi3Key[Key1, Key2, Key3]{g.mapper1.MapKey(*Get[Source1](entity)), g.mapper2.MapKey(*Get[Source2](entity)), g.mapper3.MapKey(*Get[Source3](entity))}
}
