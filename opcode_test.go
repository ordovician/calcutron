package calcutron

import "testing"

func testOpcodeEqual(t *testing.T, s string, expect Opcode) {
	got := ParseOpcode(s)
	if got != expect {
		t.Errorf("expect %v from parsing %s but but got %v", expect, s, got)
	}
}

func TestParseOpcode(t *testing.T) {
	testOpcodeEqual(t, "HLT", HLT)
	testOpcodeEqual(t, "ADD", ADD)
	testOpcodeEqual(t, "SuB", SUB)
	testOpcodeEqual(t, "subi", SUBI)
}

func TestParseOpcodeAll(t *testing.T) {
	for _, op := range AllOpcodes {
		testOpcodeEqual(t, op.String(), op)
	}
}
