package qruntime

type NativeCallContext struct {
	env     *EvalEnv
	slotptr *Slot
}

func (ncc NativeCallContext) BoolArg(index int) bool {
	return getslot(ncc.slotptr, byte(index)).Bool()
}

func (ncc NativeCallContext) ByteArg(index int) byte {
	return getslot(ncc.slotptr, byte(index)).Byte()
}

func (ncc NativeCallContext) SliceArg[T any](index int) []T {
	return *getslot(ncc.slotptr, byte(index)).Slice[T]()
}

func (ncc NativeCallContext) PtrArg[T any](index int) *T {
	return getslot(ncc.slotptr, byte(index)).Ptr[T]()
}

func (ncc NativeCallContext) IntArg(index int) int {
	return getslot(ncc.slotptr, byte(index)).Int()
}

func (ncc NativeCallContext) FloatArg(index int) float64 {
	return getslot(ncc.slotptr, byte(index)).Float()
}

func (ncc NativeCallContext) StringArg(index int) string {
	return getslot(ncc.slotptr, byte(index)).String()
}

func (ncc NativeCallContext) InterfaceArg(index int) any {
	return getslot(ncc.slotptr, byte(index)).Interface()
}

func (ncc NativeCallContext) VariadicArg() []any {
	return ncc.env.vararg
}

func (ncc NativeCallContext) SetIntResult(v int)  { ncc.env.result.SetInt(v) }
func (ncc NativeCallContext) SetIntResult2(v int) { ncc.env.result2.SetInt(v) }

func (ncc NativeCallContext) SetFloatResult(v float64)  { ncc.env.result.SetFloat(v) }
func (ncc NativeCallContext) SetFloatResult2(v float64) { ncc.env.result2.SetFloat(v) }

func (ncc NativeCallContext) SetBoolResult(v bool)  { ncc.env.result.SetBool(v) }
func (ncc NativeCallContext) SetBoolResult2(v bool) { ncc.env.result2.SetBool(v) }

func (ncc NativeCallContext) SetStringResult(v string)  { ncc.env.result.SetString(v) }
func (ncc NativeCallContext) SetStringResult2(v string) { ncc.env.result2.SetString(v) }

func (ncc NativeCallContext) SetInterfaceResult(v any)  { ncc.env.result.SetInterface(v) }
func (ncc NativeCallContext) SetInterfaceResult2(v any) { ncc.env.result2.SetInterface(v) }

func (ncc NativeCallContext) SetSliceResult[T any](v []T)  { *ncc.env.result.Slice[T]() = v }
func (ncc NativeCallContext) SetSliceResult2[T any](v []T) { *ncc.env.result2.Slice[T]() = v }

func (ncc NativeCallContext) SetPtrResult[T any](v *T)  { ncc.env.result.SetPtr(v) }
func (ncc NativeCallContext) SetPtrResult2[T any](v *T) { ncc.env.result2.SetPtr(v) }
