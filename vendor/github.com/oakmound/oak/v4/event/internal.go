package event

import (
	"sync"
)

type bindableList map[BindID]UnsafeBindable

func (eb *Bus) getBindableList(eventID UnsafeEventID, callerID CallerID) bindableList {
	if m := eb.bindingMap[eventID]; m == nil {
		eb.bindingMap[eventID] = make(map[CallerID]bindableList)
		bl := make(bindableList)
		eb.bindingMap[eventID][callerID] = bl
		return bl
	}
	bl := eb.bindingMap[eventID][callerID]
	if bl == nil {
		bl = make(bindableList)
		eb.bindingMap[eventID][callerID] = bl
	}
	return bl
}

func (bus *Bus) trigger(binds bindableList, eventID UnsafeEventID, callerID CallerID, data interface{}) {
	wg := &sync.WaitGroup{}
	wg.Add(len(binds))
	for bindID, bnd := range binds {
		bindID := bindID
		bnd := bnd
		go func() {
			if callerID == Global || bus.callerMap.HasEntity(callerID) {
				response := bnd(callerID, bus, data)
				switch response {
				case ResponseUnbindThisBinding:
					// Q: Why does this call bus.Unbind when it already has the event index to delete?
					// A: This goroutine does not own a write lock on the bus, and should therefore
					//    not modify its contents. We do not have a simple way of promoting our read lock
					//    to a write lock.
					bus.Unbind(Binding{EventID: eventID, CallerID: callerID, BindID: bindID, busResetCount: bus.resetCount})
				case ResponseUnbindThisCaller:
					bus.UnbindAllFrom(callerID)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
