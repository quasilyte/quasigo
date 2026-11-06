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
	slice := ctx.SliceArg[[]byte](0)
	value := getslot(ctx.slotptr, 1).Byte()
	ctx.SetSliceResult(append(slice, value))
}

func nativeAppend64(ctx NativeCallContext) {
	slice := ctx.SliceArg[[]uint64](0)
	value := getslot(ctx.slotptr, 1).Scalar
	ctx.SetSliceResult(append(slice, value))
}

func nativeBytesToString(ctx NativeCallContext) {
	slice := ctx.SliceArg[[]byte](0)
	ctx.SetStringResult(string(slice))
}
