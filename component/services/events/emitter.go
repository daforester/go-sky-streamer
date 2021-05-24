package events

import "reflect"

type EventHandlerFunc = func(map[string]interface{})

type EventEmitter struct {
	handlers map[string][]EventHandlerFunc
}

func (E *EventEmitter) Emit(event string, data map[string]interface{}) {
	if E.handlers == nil {
		return
	}
	if len(E.handlers[event]) == 0 {
		return
	}
	for _, h := range E.handlers[event] {
		go h(data)
	}
}

func (E *EventEmitter) Off(event string, f EventHandlerFunc) {
	if E.handlers == nil || E.handlers[event] == nil || len(E.handlers[event]) == 0 {
		return
	}
	for i, ef := range E.handlers[event] {
		if reflect.ValueOf(f).Pointer() == reflect.ValueOf(ef).Pointer() {
			E.handlers[event][i] = E.handlers[event][len(E.handlers)-1]
			E.handlers[event] = E.handlers[event][:len(E.handlers)-1]
		}
	}
}

func (E *EventEmitter) On(event string, f EventHandlerFunc) {
	if E.handlers == nil {
		E.handlers = make(map[string][]EventHandlerFunc)
	}
	if E.handlers[event] == nil {
		E.handlers[event] = make([]EventHandlerFunc, 0)
	}
	E.handlers[event] = append(E.handlers[event], f)
}
