package ecs

// GroupFilter[KeyComp comparable] 提供给用户使用的过滤器，
// 用于筛选持有KeyComp组件的Entity，同时支持通过KeyComp快速获取Entity。
// 相当于Filter1[KeyComp].GroupBy[KeyComp]，
// 注意：其依赖于Filter1[KeyComp]，要先注册Filter1[KeyComp]，否则无法使用GroupFilter[KeyComp]。
// TODO：未来可以实现自动注册，使用依赖注入项目可以在注入时自动注册Filter1[KeyComp]。
type GroupFilter[KeyComp comparable] struct {
	*groupFilterBase[KeyComp, DirectlyKeyMaker[KeyComp]]
}

func NewGroupFilter[KeyComp comparable](world *World) *GroupFilter[KeyComp] {
	filter := GetFilter[*Filter1[KeyComp]](world)
	gf := &GroupFilter[KeyComp]{
		newGroupFilterBase[KeyComp, DirectlyKeyMaker[KeyComp]](filter),
	}
	//KeyComp是groupKey，监听其增删
	registerGroupKeyEventByType[KeyComp](world, gf, filter)
	return gf
}

// 与GroupFilter类似，但是支持自定义Key类型以及自定义键值转换器。
// 用于可以指定key的类型，并需要指定转换器将KeyComp转换为Key。
// 相当于Filter1[KeyComp].GroupBy[KeyMapper(KeyComp)]，
// 注意：IGroupKeyMap[KeyComp, Key]必须是struct类型！！
// 相当于Filter1[KeyComp].GroupBy[Key]，
type GroupFilterWithKeyMapper[KeyComp any, Key comparable, KeyMapper IGroupKeyMap[KeyComp, Key]] struct {
	*groupFilterBase[Key, groupKeyMapper[KeyComp, Key, KeyMapper]]
}

func NewGroupFilterWithKeyMapper[KeyComp any, Key comparable, KeyMapper IGroupKeyMap[KeyComp, Key]](world *World) *GroupFilterWithKeyMapper[KeyComp, Key, KeyMapper] {
	filter := GetFilter[*Filter1[KeyComp]](world)
	gf := &GroupFilterWithKeyMapper[KeyComp, Key, KeyMapper]{
		newGroupFilterBase[Key, groupKeyMapper[KeyComp, Key, KeyMapper]](filter),
	}
	//KeyComp是groupKey，监听其增删
	registerGroupKeyEventByType[KeyComp](world, gf, filter)
	return gf
}

// 与GroupFilter类似，只是key是KeyComp1和KeyComp2构成的联合主键。
// 相当于Filter2[KeyComp1, KeyComp2].GroupBy[(KeyComp1, KeyComp2)]。
type GroupFilter2[KeyComp1, KeyComp2 comparable] struct {
	*groupFilterBase[Multi2Key[KeyComp1, KeyComp2], Directly2KeyMaker[KeyComp1, KeyComp2]]
}

func NewGroupFilter2[KeyComp1, KeyComp2 comparable](world *World) *GroupFilter2[KeyComp1, KeyComp2] {
	filter := GetFilter[*Filter2[KeyComp1, KeyComp2]](world)
	gf := &GroupFilter2[KeyComp1, KeyComp2]{
		newGroupFilterBase[Multi2Key[KeyComp1, KeyComp2], Directly2KeyMaker[KeyComp1, KeyComp2]](filter),
	}
	//KeyComp1、KeyComp2是groupKey，监听其增删
	registerGroupKeyEventByTypeAndHandler[KeyComp2](world,
		registerGroupKeyEventByType[KeyComp1](world, gf, filter))
	return gf
}

// 与GroupFilter2类似，但是支持自定义Key类型以及自定义键值转换器。
// 用于可以指定key的类型，并需要指定转换器将KeyComp转换为Key。
// 相当于Filter2[KeyComp1, KeyComp2].GroupBy[(KeyMapper1(KeyComp1), KeyMapper2(KeyComp2))]。
type GroupFilter2WithKeyMapper[KeyComp1, KeyComp2 any, Key1, Key2 comparable, KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2]] struct {
	*groupFilterBase[Multi2Key[Key1, Key2], group2KeyMapper[KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2]]
}

func NewGroupFilter2WithKeyMapper[KeyComp1, KeyComp2 any, Key1, Key2 comparable,
	KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2]](world *World) *GroupFilter2WithKeyMapper[
	KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2] {
	//
	filter := GetFilter[*Filter2[KeyComp1, KeyComp2]](world)
	gf := &GroupFilter2WithKeyMapper[KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2]{
		newGroupFilterBase[Multi2Key[Key1, Key2], group2KeyMapper[KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2]](filter),
	}
	//KeyComp1、KeyComp2是groupKey，监听其增删
	registerGroupKeyEventByTypeAndHandler[KeyComp2](world,
		registerGroupKeyEventByType[KeyComp1](world, gf, filter))
	return gf
}

// 参考上述注释，不再赘述。
type GroupFilter3[KeyComp1, KeyComp2, KeyComp3 comparable] struct {
	*groupFilterBase[Multi3Key[KeyComp1, KeyComp2, KeyComp3], Directly3KeyMaker[KeyComp1, KeyComp2, KeyComp3]]
}

func NewGroupFilter3[KeyComp1, KeyComp2, KeyComp3 comparable](world *World) *GroupFilter3[KeyComp1, KeyComp2, KeyComp3] {
	filter := GetFilter[*Filter3[KeyComp1, KeyComp2, KeyComp3]](world)
	gf := &GroupFilter3[KeyComp1, KeyComp2, KeyComp3]{
		newGroupFilterBase[Multi3Key[KeyComp1, KeyComp2, KeyComp3], Directly3KeyMaker[KeyComp1, KeyComp2, KeyComp3]](filter),
	}
	handler := registerGroupKeyEventByType[KeyComp1](world, gf, filter)
	registerGroupKeyEventByTypeAndHandler[KeyComp2](world, handler)
	registerGroupKeyEventByTypeAndHandler[KeyComp3](world, handler)
	return gf
}

// 参考上述注释，不再赘述。
type GroupFilter3WithKeyMapper[KeyComp1, KeyComp2, KeyComp3 any, Key1, Key2, Key3 comparable,
	KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2], KeyMapper3 IGroupKeyMap[KeyComp3, Key3]] struct {
	//
	*groupFilterBase[Multi3Key[Key1, Key2, Key3], group3KeyMapper[KeyComp1, KeyComp2, KeyComp3, Key1, Key2, Key3, KeyMapper1, KeyMapper2, KeyMapper3]]
}

func NewGroupFilter3WithKeyMapper[KeyComp1, KeyComp2, KeyComp3 any, Key1, Key2, Key3 comparable,
	KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2], KeyMapper3 IGroupKeyMap[KeyComp3, Key3]](world *World) *GroupFilter3WithKeyMapper[
	KeyComp1, KeyComp2, KeyComp3, Key1, Key2, Key3, KeyMapper1, KeyMapper2, KeyMapper3] {
	//
	filter := GetFilter[*Filter3[KeyComp1, KeyComp2, KeyComp3]](world)
	gf := &GroupFilter3WithKeyMapper[KeyComp1, KeyComp2, KeyComp3, Key1, Key2, Key3, KeyMapper1, KeyMapper2, KeyMapper3]{
		newGroupFilterBase[Multi3Key[Key1, Key2, Key3],
			group3KeyMapper[KeyComp1, KeyComp2, KeyComp3,
				Key1, Key2, Key3,
				KeyMapper1, KeyMapper2, KeyMapper3]](filter),
	}
	handler := registerGroupKeyEventByType[KeyComp1](world, gf, filter)
	registerGroupKeyEventByTypeAndHandler[KeyComp2](world, handler)
	registerGroupKeyEventByTypeAndHandler[KeyComp3](world, handler)
	return gf
}
