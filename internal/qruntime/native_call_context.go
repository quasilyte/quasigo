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

func (ncc NativeCallContext) ByteSliceArg(index int) []byte {
	return getslot(ncc.slotptr, byte(index)).ByteSlice()
}

func (ncc NativeCallContext) IntSliceArg(index int) []int {
	return getslot(ncc.slotptr, byte(index)).IntSlice()
}

func (ncc NativeCallContext) FloatSliceArg(index int) []float64 {
	return getslot(ncc.slotptr, byte(index)).FloatSlice()
}

func (ncc NativeCallContext) StringSliceArg(index int) []string {
	return getslot(ncc.slotptr, byte(index)).StringSlice()
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

func (ncc NativeCallContext) SetByteSliceResult(v []byte)  { ncc.env.result.SetByteSlice(v) }
func (ncc NativeCallContext) SetByteSliceResult2(v []byte) { ncc.env.result2.SetByteSlice(v) }

func (ncc NativeCallContext) SetIntSliceResult(v []int)  { ncc.env.result.SetIntSlice(v) }
func (ncc NativeCallContext) SetIntSliceResult2(v []int) { ncc.env.result2.SetIntSlice(v) }

func (ncc NativeCallContext) SetFloatSliceResult(v []float64)  { ncc.env.result.SetFloatSlice(v) }
func (ncc NativeCallContext) SetFloatSliceResult2(v []float64) { ncc.env.result2.SetFloatSlice(v) }

func (ncc NativeCallContext) SetStringSliceResult(v []string)  { ncc.env.result.SetStringSlice(v) }
func (ncc NativeCallContext) SetStringSliceResult2(v []string) { ncc.env.result2.SetStringSlice(v) }

func (ncc NativeCallContext) SetInterfaceResult(v any)  { ncc.env.result.SetInterface(v) }
func (ncc NativeCallContext) SetInterfaceResult2(v any) { ncc.env.result2.SetInterface(v) }
