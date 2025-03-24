package ecs

type GroupFilter[KeyComp comparable] struct {
	*groupFilterBase[KeyComp, DirectlyKeyMaker[KeyComp]]
}

func NewGroupFilter[KeyComp comparable](world *World) *GroupFilter[KeyComp] {
	return &GroupFilter[KeyComp]{
		newGroupFilterBase[KeyComp, DirectlyKeyMaker[KeyComp]](GetFilter[*Filter1[KeyComp]](world)),
	}
}

type GroupFilterWithKeyMapper[KeyComp any, Key comparable, KeyMapper IGroupKeyMap[KeyComp, Key]] struct {
	*groupFilterBase[Key, groupKeyMapper[KeyComp, Key, KeyMapper]]
}

func NewGroupFilterWithKeyMapper[KeyComp any, Key comparable, KeyMapper IGroupKeyMap[KeyComp, Key]](world *World) *GroupFilterWithKeyMapper[KeyComp, Key, KeyMapper] {
	return &GroupFilterWithKeyMapper[KeyComp, Key, KeyMapper]{
		newGroupFilterBase[Key, groupKeyMapper[KeyComp, Key, KeyMapper]](GetFilter[*Filter1[KeyComp]](world)),
	}
}

type GroupFilter2[KeyComp1, KeyComp2 comparable] struct {
	*groupFilterBase[Multi2Key[KeyComp1, KeyComp2], Directly2KeyMaker[KeyComp1, KeyComp2]]
}

func NewGroupFilter2[KeyComp1, KeyComp2 comparable](world *World) *GroupFilter2[KeyComp1, KeyComp2] {
	return &GroupFilter2[KeyComp1, KeyComp2]{
		newGroupFilterBase[Multi2Key[KeyComp1, KeyComp2], Directly2KeyMaker[KeyComp1, KeyComp2]](GetFilter[*Filter2[KeyComp1, KeyComp2]](world)),
	}
}

type GroupFilter2WithKeyMapper[KeyComp1, KeyComp2 any, Key1, Key2 comparable, KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2]] struct {
	*groupFilterBase[Multi2Key[Key1, Key2], group2KeyMapper[KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2]]
}

func NewGroupFilter2WithKeyMapper[KeyComp1, KeyComp2 any, Key1, Key2 comparable,
	KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2]](world *World) *GroupFilter2WithKeyMapper[
	KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2] {
	//
	return &GroupFilter2WithKeyMapper[KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2]{
		newGroupFilterBase[Multi2Key[Key1, Key2], group2KeyMapper[KeyComp1, KeyComp2, Key1, Key2, KeyMapper1, KeyMapper2]](GetFilter[*Filter2[KeyComp1, KeyComp2]](world)),
	}
}

type GroupFilter3[KeyComp1, KeyComp2, KeyComp3 comparable] struct {
	*groupFilterBase[Multi3Key[KeyComp1, KeyComp2, KeyComp3], Directly3KeyMaker[KeyComp1, KeyComp2, KeyComp3]]
}

func NewGroupFilter3[KeyComp1, KeyComp2, KeyComp3 comparable](world *World) *GroupFilter3[KeyComp1, KeyComp2, KeyComp3] {
	return &GroupFilter3[KeyComp1, KeyComp2, KeyComp3]{
		newGroupFilterBase[Multi3Key[KeyComp1, KeyComp2, KeyComp3], Directly3KeyMaker[KeyComp1, KeyComp2, KeyComp3]](GetFilter[*Filter3[KeyComp1, KeyComp2, KeyComp3]](world)),
	}
}

type GroupFilter3WithKeyMapper[KeyComp1, KeyComp2, KeyComp3 any, Key1, Key2, Key3 comparable,
	KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2], KeyMapper3 IGroupKeyMap[KeyComp3, Key3]] struct {
	//
	*groupFilterBase[Multi3Key[Key1, Key2, Key3], group3KeyMapper[KeyComp1, KeyComp2, KeyComp3, Key1, Key2, Key3, KeyMapper1, KeyMapper2, KeyMapper3]]
}

func NewGroupFilter3WithKeyMapper[KeyComp1, KeyComp2, KeyComp3 any, Key1, Key2, Key3 comparable,
	KeyMapper1 IGroupKeyMap[KeyComp1, Key1], KeyMapper2 IGroupKeyMap[KeyComp2, Key2], KeyMapper3 IGroupKeyMap[KeyComp3, Key3]](world *World) *GroupFilter3WithKeyMapper[
	KeyComp1, KeyComp2, KeyComp3, Key1, Key2, Key3, KeyMapper1, KeyMapper2, KeyMapper3] {
	//
	return &GroupFilter3WithKeyMapper[KeyComp1, KeyComp2, KeyComp3, Key1, Key2, Key3, KeyMapper1, KeyMapper2, KeyMapper3]{
		newGroupFilterBase[Multi3Key[Key1, Key2, Key3],
			group3KeyMapper[KeyComp1, KeyComp2, KeyComp3,
				Key1, Key2, Key3,
				KeyMapper1, KeyMapper2, KeyMapper3]](GetFilter[*Filter3[KeyComp1, KeyComp2, KeyComp3]](world)),
	}
}
