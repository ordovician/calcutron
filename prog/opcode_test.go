package prog

import "testing"

func testOpcodeEqual(t *testing.T, s string, expect Opcode) {
	got, _ := ParseOpcode(s)
	if got != expect {
		t.Errorf("expect %v from parsing %s but but got %v", expect, s, got)
	}
}

func TestParseOpcode(t *testing.T) {
	testOpcodeEqual(t, "rJumP", RJUMP)
	testOpcodeEqual(t, "ADD", ADD)
	testOpcodeEqual(t, "AddI", ADDI)
	testOpcodeEqual(t, "sub", SUB)
}

func TestParseOpcodeAll(t *testing.T) {
	for _, op := range AllOpcodes {
		testOpcodeEqual(t, op.String(), op)
	}
}
