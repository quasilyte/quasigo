package qruntime

import (
	"unsafe"
)

func nativeMakeSlice(ctx NativeCallContext) {
	elemSize := ctx.IntArg(0)
	sliceLen := ctx.IntArg(1)
	sliceCap := ctx.IntArg(2)
	mem := make([]byte, sliceLen*elemSize, sliceCap*elemSize)
	raw := *(*goSlice)(unsafe.Pointer(&mem))
	ctx.env.result.setRawSlice(goSlice{
		data: raw.data,
		len:  sliceLen,
		cap:  sliceCap,
	})
}

func nativeAppend8(ctx NativeCallContext) {
	slice := getslot(ctx.slotptr, 0).ByteSlice()
	value := getslot(ctx.slotptr, 1).Byte()
	ctx.env.result.SetByteSlice(append(slice, value))
}

func nativeAppend64(ctx NativeCallContext) {
	slice := getslot(ctx.slotptr, 0).slice64()
	value := getslot(ctx.slotptr, 1).Scalar
	ctx.env.result.setSlice64(append(slice, value))
}

func nativeBytesToString(ctx NativeCallContext) {
	slice := ctx.ByteSliceArg(0)
	ctx.SetStringResult(string(slice))
}
