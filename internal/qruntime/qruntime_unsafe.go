package qruntime

import (
	"math"
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

type goSlice struct {
	data unsafe.Pointer
	len  int
	cap  int
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

func (slot Slot) slice64() []uint64 {
	return *(*[]uint64)(unsafe.Pointer(&slot))
}

func (slot *Slot) setRawSlice(s goSlice) {
	slot.Ptr = s.data
	slot.Scalar = uint64(s.len)
	slot.Scalar2 = uint64(s.cap)
}

func (slot *Slot) setSlice64(v []uint64) {
	slot.setRawSlice(*(*goSlice)(unsafe.Pointer(&v)))
}

func (slot *Slot) SetByteSlice(v []byte) {
	slot.setRawSlice(*(*goSlice)(unsafe.Pointer(&v)))
}

func (slot Slot) ByteSlice() []byte {
	return *(*[]byte)(unsafe.Pointer(&slot))
}

func (slot *Slot) SetInt(v int) {
	slot.Scalar = uint64(v)
}

func (slot *Slot) SetFloat(v float64) {
	slot.Scalar = math.Float64bits(v)
}

func (slot *Slot) SetByte(v byte) {
	slot.Scalar = uint64(v)
}

func (slot Slot) Int() int { return int(slot.Scalar) }

func (slot Slot) Float() float64 { return math.Float64frombits(slot.Scalar) }

func (slot Slot) Byte() byte { return byte(slot.Scalar) }

// With go 1.19 this function is less efficient than specialized
// functions, but it can be fixed with GOEXPERIMENT=unified.
// With go 1.20 (and above?) no extra knobs are needed.
func ptradd[T any](p *T, n int) *T {
	return (*T)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(n)))
}

//go:nosplit
func unpack16(code *byte, offset int) int {
	ptr := ptradd(code, offset)
	return int(int16(uint16(*ptr) | uint16(*ptradd(ptr, 1))<<8))
}

//go:nosplit
func unpack8(code *byte, offset int) byte {
	return *ptradd(code, offset)
}

//go:nosplit
func unpack8x2(code *byte, offset int) (byte, byte) {
	ptr := ptradd(code, offset)
	return *ptr, *(ptradd(ptr, 1))
}

//go:nosplit
func unpack8x3(code *byte, offset int) (byte, byte, byte) {
	ptr := ptradd(code, offset)
	return *ptr, *(ptradd(ptr, 1)), *(ptradd(ptr, 2))
}

//go:nosplit
func unpack8x4(code *byte, offset int) (byte, byte, byte, byte) {
	ptr := ptradd(code, offset)
	return *ptr, *(ptradd(ptr, 1)), *(ptradd(ptr, 2)), *(ptradd(ptr, 3))
}

//go:nosplit
func getslot(slotptr *Slot, index byte) *Slot {
	return ptradd(slotptr, int(SizeofSlot)*int(index))
}

//go:nosplit
func canAllocFrame(current, end *Slot, frameSize int) bool {
	return uintptr(unsafe.Pointer(current))+uintptr(frameSize) < uintptr(unsafe.Pointer(end))
}

//go:nosplit
func nextFrameSlot(current *Slot, frameSize int) *Slot {
	return ptradd(current, frameSize)
}
