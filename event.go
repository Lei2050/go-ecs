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

type EntityEvents struct {
	BeforeAdd Delegate
	AfterAdd  Delegate

	BeforeUpdate Delegate

	BeforeDelete Delegate
	AfterDelete  Delegate
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
