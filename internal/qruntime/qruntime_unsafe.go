package qruntime

import (
	"unsafe"
)

const SizeofSlot = unsafe.Sizeof(Slot{})

// Works for both empty and non-empty interfaces.
type goIface struct {
	typeinfo uint64
	data     unsafe.Pointer
}

type goString struct {
	data unsafe.Pointer
	len  int
}

type Slot struct {
	// TODO: make fields unexported, adjust the user code as needed.
	Ptr     unsafe.Pointer
	Scalar  uint64
	Scalar2 uint64
}

func (slot Slot) IsNil() bool {
	return slot.Ptr == nil
}

func (slot Slot) IsNilInterface() bool {
	return slot.Ptr == nil && slot.Scalar == 0
}

func (slot *Slot) SetBool(v bool) {
	if v {
		slot.Scalar = 1
	} else {
		slot.Scalar = 0
	}
}

func (slot Slot) Bool() bool {
	return *(*bool)(unsafe.Pointer(&slot.Scalar))
}

func (slot *Slot) MoveInterface(src *Slot) {
	slot.Ptr = src.Ptr
	slot.Scalar = src.Scalar
}

func (slot *Slot) SetInterface(v interface{}) {
	raw := (*goIface)(unsafe.Pointer(&v))
	slot.Ptr = raw.data
	slot.Scalar = raw.typeinfo
}

func (slot Slot) Interface() interface{} {
	v := goIface{
		typeinfo: slot.Scalar,
		data:     slot.Ptr,
	}
	return *(*interface{})(unsafe.Pointer(&v))
}

func (slot *Slot) SetString(v string) {
	raw := *(*goString)(unsafe.Pointer(&v))
	slot.Ptr = raw.data
	slot.Scalar = uint64(raw.len)
}

func (slot Slot) String() string {
	return *(*string)(unsafe.Pointer(&slot))
}

func (slot *Slot) SetInt(v int) {
	slot.Scalar = uint64(v)
}

func (slot Slot) Int() int { return int(slot.Scalar) }
