package model

//#include <stdint.h>
//#include <stdlib.h>
//#include <string.h>
//#include "moc.h"
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/therecipe/qt"
	std_core "github.com/therecipe/qt/core"
)

func cGoUnpackString(s C.struct_Moc_PackedString) string {
	if len := int(s.len); len == -1 {
		return C.GoString(s.data)
	}
	return C.GoStringN(s.data, C.int(s.len))
}

type ClientQT_ITF interface {
	std_core.QObject_ITF
	ClientQT_PTR() *ClientQT
}

func (ptr *ClientQT) ClientQT_PTR() *ClientQT {
	return ptr
}

func (ptr *ClientQT) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *ClientQT) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromClientQT(ptr ClientQT_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.ClientQT_PTR().Pointer()
	}
	return nil
}

func NewClientQTFromPointer(ptr unsafe.Pointer) *ClientQT {
	var n *ClientQT
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(ClientQT)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *ClientQT:
			n = deduced

		case *std_core.QObject:
			n = &ClientQT{QObject: *deduced}

		default:
			n = new(ClientQT)
			n.SetPointer(ptr)
		}
	}
	return n
}

//export callbackClientQT_Constructor
func callbackClientQT_Constructor(ptr unsafe.Pointer) {
	gPtr := NewClientQTFromPointer(ptr)
	qt.Register(ptr, gPtr)
}

//export callbackClientQT_ReloadUI
func callbackClientQT_ReloadUI(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "reloadUI"); signal != nil {
		signal.(func())()
	}

}

func (ptr *ClientQT) ConnectReloadUI(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "reloadUI") {
			C.ClientQT_ConnectReloadUI(ptr.Pointer())
		}

		if signal := qt.LendSignal(ptr.Pointer(), "reloadUI"); signal != nil {
			qt.ConnectSignal(ptr.Pointer(), "reloadUI", func() {
				signal.(func())()
				f()
			})
		} else {
			qt.ConnectSignal(ptr.Pointer(), "reloadUI", f)
		}
	}
}

func (ptr *ClientQT) DisconnectReloadUI() {
	if ptr.Pointer() != nil {
		C.ClientQT_DisconnectReloadUI(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "reloadUI")
	}
}

func (ptr *ClientQT) ReloadUI() {
	if ptr.Pointer() != nil {
		C.ClientQT_ReloadUI(ptr.Pointer())
	}
}

func ClientQT_QRegisterMetaType() int {
	return int(int32(C.ClientQT_ClientQT_QRegisterMetaType()))
}

func (ptr *ClientQT) QRegisterMetaType() int {
	return int(int32(C.ClientQT_ClientQT_QRegisterMetaType()))
}

func ClientQT_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ClientQT_ClientQT_QRegisterMetaType2(typeNameC)))
}

func (ptr *ClientQT) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ClientQT_ClientQT_QRegisterMetaType2(typeNameC)))
}

func ClientQT_QmlRegisterType() int {
	return int(int32(C.ClientQT_ClientQT_QmlRegisterType()))
}

func (ptr *ClientQT) QmlRegisterType() int {
	return int(int32(C.ClientQT_ClientQT_QmlRegisterType()))
}

func ClientQT_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ClientQT_ClientQT_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *ClientQT) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ClientQT_ClientQT_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *ClientQT) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		var tmpValue = std_core.NewQByteArrayFromPointer(C.ClientQT___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		runtime.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *ClientQT) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *ClientQT) __dynamicPropertyNames_newList() unsafe.Pointer {
	return unsafe.Pointer(C.ClientQT___dynamicPropertyNames_newList(ptr.Pointer()))
}

func (ptr *ClientQT) __findChildren_atList2(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		var tmpValue = std_core.NewQObjectFromPointer(C.ClientQT___findChildren_atList2(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ClientQT) __findChildren_setList2(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT___findChildren_setList2(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ClientQT) __findChildren_newList2() unsafe.Pointer {
	return unsafe.Pointer(C.ClientQT___findChildren_newList2(ptr.Pointer()))
}

func (ptr *ClientQT) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		var tmpValue = std_core.NewQObjectFromPointer(C.ClientQT___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ClientQT) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ClientQT) __findChildren_newList3() unsafe.Pointer {
	return unsafe.Pointer(C.ClientQT___findChildren_newList3(ptr.Pointer()))
}

func (ptr *ClientQT) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		var tmpValue = std_core.NewQObjectFromPointer(C.ClientQT___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ClientQT) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ClientQT) __findChildren_newList() unsafe.Pointer {
	return unsafe.Pointer(C.ClientQT___findChildren_newList(ptr.Pointer()))
}

func (ptr *ClientQT) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		var tmpValue = std_core.NewQObjectFromPointer(C.ClientQT___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ClientQT) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ClientQT) __children_newList() unsafe.Pointer {
	return unsafe.Pointer(C.ClientQT___children_newList(ptr.Pointer()))
}

func NewClientQT(parent std_core.QObject_ITF) *ClientQT {
	var tmpValue = NewClientQTFromPointer(C.ClientQT_NewClientQT(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackClientQT_DestroyClientQT
func callbackClientQT_DestroyClientQT(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~ClientQT"); signal != nil {
		signal.(func())()
	} else {
		NewClientQTFromPointer(ptr).DestroyClientQTDefault()
	}
}

func (ptr *ClientQT) ConnectDestroyClientQT(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~ClientQT"); signal != nil {
			qt.ConnectSignal(ptr.Pointer(), "~ClientQT", func() {
				signal.(func())()
				f()
			})
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~ClientQT", f)
		}
	}
}

func (ptr *ClientQT) DisconnectDestroyClientQT() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~ClientQT")
	}
}

func (ptr *ClientQT) DestroyClientQT() {
	if ptr.Pointer() != nil {
		C.ClientQT_DestroyClientQT(ptr.Pointer())
		ptr.SetPointer(nil)
		runtime.SetFinalizer(ptr, nil)
	}
}

func (ptr *ClientQT) DestroyClientQTDefault() {
	if ptr.Pointer() != nil {
		C.ClientQT_DestroyClientQTDefault(ptr.Pointer())
		ptr.SetPointer(nil)
		runtime.SetFinalizer(ptr, nil)
	}
}

//export callbackClientQT_Event
func callbackClientQT_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt(signal.(func(*std_core.QEvent) bool)(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewClientQTFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *ClientQT) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return C.ClientQT_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e)) != 0
	}
	return false
}

//export callbackClientQT_EventFilter
func callbackClientQT_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt(signal.(func(*std_core.QObject, *std_core.QEvent) bool)(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewClientQTFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *ClientQT) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return C.ClientQT_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event)) != 0
	}
	return false
}

//export callbackClientQT_ChildEvent
func callbackClientQT_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		signal.(func(*std_core.QChildEvent))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewClientQTFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *ClientQT) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackClientQT_ConnectNotify
func callbackClientQT_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		signal.(func(*std_core.QMetaMethod))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewClientQTFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ClientQT) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackClientQT_CustomEvent
func callbackClientQT_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		signal.(func(*std_core.QEvent))(std_core.NewQEventFromPointer(event))
	} else {
		NewClientQTFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *ClientQT) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackClientQT_DeleteLater
func callbackClientQT_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		signal.(func())()
	} else {
		NewClientQTFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *ClientQT) DeleteLaterDefault() {
	if ptr.Pointer() != nil {
		C.ClientQT_DeleteLaterDefault(ptr.Pointer())
		ptr.SetPointer(nil)
		runtime.SetFinalizer(ptr, nil)
	}
}

//export callbackClientQT_Destroyed
func callbackClientQT_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		signal.(func(*std_core.QObject))(std_core.NewQObjectFromPointer(obj))
	}

}

//export callbackClientQT_DisconnectNotify
func callbackClientQT_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		signal.(func(*std_core.QMetaMethod))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewClientQTFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ClientQT) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackClientQT_ObjectNameChanged
func callbackClientQT_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		signal.(func(string))(cGoUnpackString(objectName))
	}

}

//export callbackClientQT_TimerEvent
func callbackClientQT_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		signal.(func(*std_core.QTimerEvent))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewClientQTFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *ClientQT) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ClientQT_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}
