package ifaceops

//test:disasm_both
// ifaceops.testIfaceNil1 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   IsNilInterface temp0 = err
//   ReturnScalar temp0
func testIfaceNil1(err error) bool {
	return err == nil
}

//test:disasm_both
// ifaceops.testIfaceNil2 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   IsNilInterface temp0 = err
//   ReturnScalar temp0
func testIfaceNil2(err error) bool {
	return nil == err
}

//test:disasm_both
// ifaceops.testIfaceNotNil1 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   IsNotNilInterface temp0 = err
//   ReturnScalar temp0
func testIfaceNotNil1(err error) bool {
	return err != nil
}

//test:disasm_both
// ifaceops.testIfaceNotNil2 code=5 frame=48 (2 slots: 1 params, 1 locals)
//   IsNotNilInterface temp0 = err
//   ReturnScalar temp0
func testIfaceNotNil2(err error) bool {
	return nil != err
}
