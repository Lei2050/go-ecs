package ecs

import "reflect"

type Delegate struct {
	callbacks []func(Entity)
}

func (d *Delegate) AddCallback(callback func(Entity)) {
	d.callbacks = append(d.callbacks, callback)
}

func (d *Delegate) RemoveCallback(callback func(Entity)) {
	pf := reflect.ValueOf(callback).Pointer()
	for i, cb := range d.callbacks {
		if reflect.ValueOf(cb).Pointer() == pf {
			d.callbacks = append(d.callbacks[:i], d.callbacks[i+1:]...)
			break
		}
	}
}

func (d *Delegate) Invoke(entity Entity) {
	for _, cb := range d.callbacks {
		cb(entity)
	}
}

type DelegateWithParam struct {
	callbacks []func(Entity, ...any)
}

func (d *DelegateWithParam) AddCallback(callback func(Entity, ...any)) {
	d.callbacks = append(d.callbacks, callback)
}

func (d *DelegateWithParam) RemoveCallback(callback func(Entity, ...any)) {
	pf := reflect.ValueOf(callback).Pointer()
	for i, cb := range d.callbacks {
		if reflect.ValueOf(cb).Pointer() == pf {
			d.callbacks = append(d.callbacks[:i], d.callbacks[i+1:]...)
			break
		}
	}
}

func (d *DelegateWithParam) Invoke(entity Entity, params ...any) {
	for _, cb := range d.callbacks {
		cb(entity, params...)
	}
}

//与Entity有关的各种事件，可以用来监听entity的数据变化。
//目前主要是用来Component数据的变化。
type EntityEvents struct {
	BeforeAdd Delegate
	AfterAdd  Delegate

	BeforeUpdate Delegate

	BeforeDelete Delegate
	AfterDelete  Delegate

	BeforeAddWithPoolIdx DelegateWithParam
	AfterAddWithPoolIdx  DelegateWithParam
}

type FilterEventListener interface {
	OnEntityAdded(entity Entity)
	OnEntityRemoved(entity Entity)
}

var _ FilterEventListener = &FilterEventListen{}

type FilterEventListen struct {
	EntityAdded   Delegate
	EntityRemoved Delegate
}

func newFilterEventListener() *FilterEventListen {
	return &FilterEventListen{}
}

func (f *FilterEventListen) OnEntityAdded(entity Entity) {
	f.EntityAdded.Invoke(entity)
}

func (f *FilterEventListen) OnEntityRemoved(entity Entity) {
	f.EntityRemoved.Invoke(entity)
}
