package ecs

type Filter1[Comp1 any] struct {
	*filterBase1[Comp1]
}

func NewFilter1[Comp1 any](world *World) *Filter1[Comp1] {
	return &Filter1[Comp1]{
		filterBase1: newFilterBase1[Comp1](world, GetComponentType[Comp1]().TypeIndex),
	}
}

type Filter1Exclude1[Comp1, ExcComp1 any] struct {
	*filterBase1[Comp1]
}

func NewFilter1Exclude[Comp1, ExcComp1 any](world *World) *Filter1Exclude1[Comp1, ExcComp1] {
	f := &Filter1Exclude1[Comp1, ExcComp1]{
		filterBase1: newFilterBase1[Comp1](world, GetComponentType[Comp1]().TypeIndex),
	}
	initMask1[ExcComp1](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter1Exclude2[Comp1, ExcComp1, ExcComp2 any] struct {
	*filterBase1[Comp1]
}

func NewFilter1Exclude2[Comp1, ExcComp1, ExcComp2 any](world *World) *Filter1Exclude2[Comp1, ExcComp1, ExcComp2] {
	f := &Filter1Exclude2[Comp1, ExcComp1, ExcComp2]{
		filterBase1: newFilterBase1[Comp1](world, GetComponentType[Comp1]().TypeIndex),
	}
	initMask2[ExcComp1, ExcComp2](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter1Exclude3[Comp1, ExcComp1, ExcComp2, ExcComp3 any] struct {
	*filterBase1[Comp1]
}

func NewFilter1Exclude3[Comp1, ExcComp1, ExcComp2, ExcComp3 any](world *World) *Filter1Exclude3[Comp1, ExcComp1, ExcComp2, ExcComp3] {
	f := &Filter1Exclude3[Comp1, ExcComp1, ExcComp2, ExcComp3]{
		filterBase1: newFilterBase1[Comp1](world, GetComponentType[Comp1]().TypeIndex),
	}
	initMask3[ExcComp1, ExcComp2, ExcComp3](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter1Exclude4[Comp1, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any] struct {
	*filterBase1[Comp1]
}

func NewFilter1Exclude4[Comp1, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any](world *World) *Filter1Exclude4[Comp1, ExcComp1, ExcComp2, ExcComp3, ExcComp4] {
	f := &Filter1Exclude4[Comp1, ExcComp1, ExcComp2, ExcComp3, ExcComp4]{
		filterBase1: newFilterBase1[Comp1](world, GetComponentType[Comp1]().TypeIndex),
	}
	initMask4[ExcComp1, ExcComp2, ExcComp3, ExcComp4](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter2[Comp1, Comp2 any] struct {
	*filterBase2[Comp1, Comp2]
}

func NewFilter2[Comp1, Comp2 any](world *World) *Filter2[Comp1, Comp2] {
	return &Filter2[Comp1, Comp2]{
		filterBase2: newFilterBase2[Comp1, Comp2](world, GetComponentType[Comp1]().TypeIndex, GetComponentType[Comp2]().TypeIndex),
	}
}

type Filter2Exclude1[Comp1, Comp2, ExcComp1 any] struct {
	*filterBase2[Comp1, Comp2]
}

func NewFilter2Exclude1[Comp1, Comp2, ExcComp1 any](world *World) *Filter2Exclude1[Comp1, Comp2, ExcComp1] {
	f := &Filter2Exclude1[Comp1, Comp2, ExcComp1]{
		filterBase2: newFilterBase2[Comp1, Comp2](world, GetComponentType[Comp1]().TypeIndex, GetComponentType[Comp2]().TypeIndex),
	}
	initMask1[ExcComp1](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter2Exclude2[Comp1, Comp2, ExcComp1, ExcComp2 any] struct {
	*filterBase2[Comp1, Comp2]
}

func NewFilter2Exclude2[Comp1, Comp2, ExcComp1, ExcComp2 any](world *World) *Filter2Exclude2[Comp1, Comp2, ExcComp1, ExcComp2] {
	f := &Filter2Exclude2[Comp1, Comp2, ExcComp1, ExcComp2]{
		filterBase2: newFilterBase2[Comp1, Comp2](world, GetComponentType[Comp1]().TypeIndex, GetComponentType[Comp2]().TypeIndex),
	}
	initMask2[ExcComp1, ExcComp2](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter2Exclude3[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3 any] struct {
	*filterBase2[Comp1, Comp2]
}

func NewFilter2Exclude3[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3 any](world *World) *Filter2Exclude3[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3] {
	f := &Filter2Exclude3[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3]{
		filterBase2: newFilterBase2[Comp1, Comp2](world, GetComponentType[Comp1]().TypeIndex, GetComponentType[Comp2]().TypeIndex),
	}
	initMask3[ExcComp1, ExcComp2, ExcComp3](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter2Exclude4[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any] struct {
	*filterBase2[Comp1, Comp2]
}

func NewFilter2Exclude4[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any](world *World) *Filter2Exclude4[
	Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3, ExcComp4] {
	f := &Filter2Exclude4[Comp1, Comp2, ExcComp1, ExcComp2, ExcComp3, ExcComp4]{
		filterBase2: newFilterBase2[Comp1, Comp2](world, GetComponentType[Comp1]().TypeIndex, GetComponentType[Comp2]().TypeIndex),
	}
	initMask4[ExcComp1, ExcComp2, ExcComp3, ExcComp4](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter3[Comp1, Comp2, Comp3 any] struct {
	*filterBase3[Comp1, Comp2, Comp3]
}

func NewFilter3[Comp1, Comp2, Comp3 any](world *World) *Filter3[Comp1, Comp2, Comp3] {
	return &Filter3[Comp1, Comp2, Comp3]{
		filterBase3: newFilterBase3[Comp1, Comp2, Comp3](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
		),
	}
}

type Filter3Exclude1[Comp1, Comp2, Comp3, ExcComp1 any] struct {
	*filterBase3[Comp1, Comp2, Comp3]
}

func NewFilter3Exclude1[Comp1, Comp2, Comp3, ExcComp1 any](world *World) *Filter3Exclude1[Comp1, Comp2, Comp3, ExcComp1] {
	f := &Filter3Exclude1[Comp1, Comp2, Comp3, ExcComp1]{
		filterBase3: newFilterBase3[Comp1, Comp2, Comp3](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
		),
	}
	initMask1[ExcComp1](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter3Exclude2[Comp1, Comp2, Comp3, ExcComp1, ExcComp2 any] struct {
	*filterBase3[Comp1, Comp2, Comp3]
}

func NewFilter3Exclude2[Comp1, Comp2, Comp3, ExcComp1, ExcComp2 any](world *World) *Filter3Exclude2[Comp1, Comp2, Comp3, ExcComp1, ExcComp2] {
	f := &Filter3Exclude2[Comp1, Comp2, Comp3, ExcComp1, ExcComp2]{
		filterBase3: newFilterBase3[Comp1, Comp2, Comp3](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
		),
	}
	initMask2[ExcComp1, ExcComp2](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter3Exclude3[Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3 any] struct {
	*filterBase3[Comp1, Comp2, Comp3]
}

func NewFilter3Exclude3[Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3 any](world *World) *Filter3Exclude3[
	Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3] {
	f := &Filter3Exclude3[Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3]{
		filterBase3: newFilterBase3[Comp1, Comp2, Comp3](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
		),
	}
	initMask3[ExcComp1, ExcComp2, ExcComp3](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter3Exclude4[Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any] struct {
	*filterBase3[Comp1, Comp2, Comp3]
}

func NewFilter3Exclude4[Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any](world *World) *Filter3Exclude4[
	Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3, ExcComp4] {
	f := &Filter3Exclude4[Comp1, Comp2, Comp3, ExcComp1, ExcComp2, ExcComp3, ExcComp4]{
		filterBase3: newFilterBase3[Comp1, Comp2, Comp3](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
		),
	}
	initMask4[ExcComp1, ExcComp2, ExcComp3, ExcComp4](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter4[Comp1, Comp2, Comp3, Comp4 any] struct {
	*filterBase4[Comp1, Comp2, Comp3, Comp4]
}

func NewFilter4[Comp1, Comp2, Comp3, Comp4 any](world *World) *Filter4[Comp1, Comp2, Comp3, Comp4] {
	return &Filter4[Comp1, Comp2, Comp3, Comp4]{
		filterBase4: newFilterBase4[Comp1, Comp2, Comp3, Comp4](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
			GetComponentType[Comp4]().TypeIndex,
		),
	}
}

type Filter4Exclude1[Comp1, Comp2, Comp3, Comp4, ExcComp1 any] struct {
	*filterBase4[Comp1, Comp2, Comp3, Comp4]
}

func NewFilter4Exclude1[Comp1, Comp2, Comp3, Comp4, ExcComp1 any](world *World) *Filter4Exclude1[Comp1, Comp2, Comp3, Comp4, ExcComp1] {
	f := &Filter4Exclude1[Comp1, Comp2, Comp3, Comp4, ExcComp1]{
		filterBase4: newFilterBase4[Comp1, Comp2, Comp3, Comp4](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
			GetComponentType[Comp4]().TypeIndex,
		),
	}
	initMask1[ExcComp1](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter4Exclude2[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2 any] struct {
	*filterBase4[Comp1, Comp2, Comp3, Comp4]
}

func NewFilter4Exclude2[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2 any](world *World) *Filter4Exclude2[
	Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2] {
	f := &Filter4Exclude2[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2]{
		filterBase4: newFilterBase4[Comp1, Comp2, Comp3, Comp4](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
			GetComponentType[Comp4]().TypeIndex,
		),
	}
	initMask2[ExcComp1, ExcComp2](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter4Exclude3[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3 any] struct {
	*filterBase4[Comp1, Comp2, Comp3, Comp4]
}

func NewFilter4Exclude3[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3 any](world *World) *Filter4Exclude3[
	Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3] {
	f := &Filter4Exclude3[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3]{
		filterBase4: newFilterBase4[Comp1, Comp2, Comp3, Comp4](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
			GetComponentType[Comp4]().TypeIndex,
		),
	}
	initMask3[ExcComp1, ExcComp2, ExcComp3](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}

type Filter4Exclude4[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any] struct {
	*filterBase4[Comp1, Comp2, Comp3, Comp4]
}

func NewFilter4Exclude4[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3, ExcComp4 any](world *World) *Filter4Exclude4[
	Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3, ExcComp4] {
	f := &Filter4Exclude4[Comp1, Comp2, Comp3, Comp4, ExcComp1, ExcComp2, ExcComp3, ExcComp4]{
		filterBase4: newFilterBase4[Comp1, Comp2, Comp3, Comp4](
			world,
			GetComponentType[Comp1]().TypeIndex,
			GetComponentType[Comp2]().TypeIndex,
			GetComponentType[Comp3]().TypeIndex,
			GetComponentType[Comp4]().TypeIndex,
		),
	}
	initMask4[ExcComp1, ExcComp2, ExcComp3, ExcComp4](&f.ExcludeTypeIndices, &f.ExcludeMask)
	return f
}
