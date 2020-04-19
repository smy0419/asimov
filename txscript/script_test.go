// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"bytes"
	"testing"
)

// TestParseOpcode tests for opcode parsing with bad data templates.
func TestParseOpcode(t *testing.T) {
	// Deep copy the array and make one of the opcodes invalid by setting it
	// to the wrong length.
	fakeArray := opcodeArray
	fakeArray[OP_PUSHDATA4] = opcode{value: OP_PUSHDATA4,
		name: "OP_PUSHDATA4", length: -8, opfunc: opcodePushData}

	// This script would be fine if -8 was a valid length.
	_, err := parseScriptTemplate([]byte{OP_PUSHDATA4, 0x1, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00}, &fakeArray)
	if err == nil {
		t.Errorf("no error with dodgy opcode array!")
	}
}

// TestUnparsingInvalidOpcodes tests for errors when unparsing invalid parsed
// opcodes.
func TestUnparsingInvalidOpcodes(t *testing.T) {
	tests := []struct {
		name        string
		pop         *parsedOpcode
		expectedErr error
	}{
		{
			name: "OP_FALSE",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_FALSE],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_FALSE long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_FALSE],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_1 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_1],
				data:   nil,
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_1",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_1],
				data:   make([]byte, 1),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_1 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_1],
				data:   make([]byte, 2),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_2 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_2],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_2",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_2],
				data:   make([]byte, 2),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_2 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_2],
				data:   make([]byte, 3),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_3 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_3],
				data:   make([]byte, 2),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_3",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_3],
				data:   make([]byte, 3),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_3 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_3],
				data:   make([]byte, 4),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_4 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_4],
				data:   make([]byte, 3),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_4",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_4],
				data:   make([]byte, 4),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_4 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_4],
				data:   make([]byte, 5),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_5 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_5],
				data:   make([]byte, 4),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_5",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_5],
				data:   make([]byte, 5),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_5 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_5],
				data:   make([]byte, 6),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_6 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_6],
				data:   make([]byte, 5),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_6",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_6],
				data:   make([]byte, 6),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_6 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_6],
				data:   make([]byte, 7),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_7 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_7],
				data:   make([]byte, 6),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_7",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_7],
				data:   make([]byte, 7),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_7 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_7],
				data:   make([]byte, 8),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_8 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_8],
				data:   make([]byte, 7),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_8",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_8],
				data:   make([]byte, 8),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_8 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_8],
				data:   make([]byte, 9),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_9 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_9],
				data:   make([]byte, 8),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_9",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_9],
				data:   make([]byte, 9),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_9 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_9],
				data:   make([]byte, 10),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_10 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_10],
				data:   make([]byte, 9),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_10",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_10],
				data:   make([]byte, 10),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_10 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_10],
				data:   make([]byte, 11),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_11 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_11],
				data:   make([]byte, 10),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_11",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_11],
				data:   make([]byte, 11),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_11 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_11],
				data:   make([]byte, 12),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_12 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_12],
				data:   make([]byte, 11),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_12",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_12],
				data:   make([]byte, 12),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_12 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_12],
				data:   make([]byte, 13),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_13 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_13],
				data:   make([]byte, 12),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_13",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_13],
				data:   make([]byte, 13),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_13 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_13],
				data:   make([]byte, 14),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_14 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_14],
				data:   make([]byte, 13),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_14",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_14],
				data:   make([]byte, 14),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_14 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_14],
				data:   make([]byte, 15),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_15 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_15],
				data:   make([]byte, 14),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_15",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_15],
				data:   make([]byte, 15),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_15 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_15],
				data:   make([]byte, 16),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_16 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_16],
				data:   make([]byte, 15),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_16",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_16],
				data:   make([]byte, 16),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_16 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_16],
				data:   make([]byte, 17),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_17 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_17],
				data:   make([]byte, 16),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_17",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_17],
				data:   make([]byte, 17),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_17 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_17],
				data:   make([]byte, 18),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_18 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_18],
				data:   make([]byte, 17),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_18",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_18],
				data:   make([]byte, 18),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_18 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_18],
				data:   make([]byte, 19),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_19 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_19],
				data:   make([]byte, 18),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_19",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_19],
				data:   make([]byte, 19),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_19 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_19],
				data:   make([]byte, 20),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_20 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_20],
				data:   make([]byte, 19),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_20",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_20],
				data:   make([]byte, 20),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_20 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_20],
				data:   make([]byte, 21),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_21 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_21],
				data:   make([]byte, 20),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_21",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_21],
				data:   make([]byte, 21),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_21 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_21],
				data:   make([]byte, 22),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_22 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_22],
				data:   make([]byte, 21),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_22",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_22],
				data:   make([]byte, 22),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_22 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_22],
				data:   make([]byte, 23),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_23 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_23],
				data:   make([]byte, 22),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_23",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_23],
				data:   make([]byte, 23),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_23 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_23],
				data:   make([]byte, 24),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_24 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_24],
				data:   make([]byte, 23),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_24",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_24],
				data:   make([]byte, 24),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_24 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_24],
				data:   make([]byte, 25),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_25 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_25],
				data:   make([]byte, 24),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_25",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_25],
				data:   make([]byte, 25),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_25 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_25],
				data:   make([]byte, 26),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_26 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_26],
				data:   make([]byte, 25),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_26",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_26],
				data:   make([]byte, 26),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_26 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_26],
				data:   make([]byte, 27),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_27 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_27],
				data:   make([]byte, 26),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_27",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_27],
				data:   make([]byte, 27),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_27 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_27],
				data:   make([]byte, 28),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_28 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_28],
				data:   make([]byte, 27),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_28",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_28],
				data:   make([]byte, 28),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_28 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_28],
				data:   make([]byte, 29),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_29 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_29],
				data:   make([]byte, 28),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_29",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_29],
				data:   make([]byte, 29),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_29 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_29],
				data:   make([]byte, 30),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_30 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_30],
				data:   make([]byte, 29),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_30",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_30],
				data:   make([]byte, 30),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_30 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_30],
				data:   make([]byte, 31),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_31 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_31],
				data:   make([]byte, 30),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_31",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_31],
				data:   make([]byte, 31),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_31 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_31],
				data:   make([]byte, 32),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_32 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_32],
				data:   make([]byte, 31),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_32",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_32],
				data:   make([]byte, 32),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_32 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_32],
				data:   make([]byte, 33),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_33 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_33],
				data:   make([]byte, 32),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_33",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_33],
				data:   make([]byte, 33),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_33 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_33],
				data:   make([]byte, 34),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_34 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_34],
				data:   make([]byte, 33),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_34",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_34],
				data:   make([]byte, 34),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_34 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_34],
				data:   make([]byte, 35),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_35 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_35],
				data:   make([]byte, 34),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_35",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_35],
				data:   make([]byte, 35),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_35 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_35],
				data:   make([]byte, 36),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_36 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_36],
				data:   make([]byte, 35),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_36",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_36],
				data:   make([]byte, 36),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_36 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_36],
				data:   make([]byte, 37),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_37 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_37],
				data:   make([]byte, 36),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_37",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_37],
				data:   make([]byte, 37),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_37 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_37],
				data:   make([]byte, 38),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_38 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_38],
				data:   make([]byte, 37),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_38",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_38],
				data:   make([]byte, 38),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_38 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_38],
				data:   make([]byte, 39),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_39 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_39],
				data:   make([]byte, 38),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_39",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_39],
				data:   make([]byte, 39),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_39 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_39],
				data:   make([]byte, 40),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_40 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_40],
				data:   make([]byte, 39),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_40",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_40],
				data:   make([]byte, 40),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_40 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_40],
				data:   make([]byte, 41),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_41 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_41],
				data:   make([]byte, 40),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_41",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_41],
				data:   make([]byte, 41),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_41 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_41],
				data:   make([]byte, 42),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_42 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_42],
				data:   make([]byte, 41),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_42",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_42],
				data:   make([]byte, 42),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_42 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_42],
				data:   make([]byte, 43),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_43 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_43],
				data:   make([]byte, 42),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_43",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_43],
				data:   make([]byte, 43),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_43 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_43],
				data:   make([]byte, 44),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_44 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_44],
				data:   make([]byte, 43),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_44",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_44],
				data:   make([]byte, 44),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_44 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_44],
				data:   make([]byte, 45),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_45 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_45],
				data:   make([]byte, 44),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_45",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_45],
				data:   make([]byte, 45),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_45 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_45],
				data:   make([]byte, 46),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_46 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_46],
				data:   make([]byte, 45),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_46",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_46],
				data:   make([]byte, 46),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_46 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_46],
				data:   make([]byte, 47),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_47 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_47],
				data:   make([]byte, 46),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_47",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_47],
				data:   make([]byte, 47),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_47 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_47],
				data:   make([]byte, 48),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_48 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_48],
				data:   make([]byte, 47),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_48",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_48],
				data:   make([]byte, 48),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_48 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_48],
				data:   make([]byte, 49),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_49 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_49],
				data:   make([]byte, 48),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_49",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_49],
				data:   make([]byte, 49),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_49 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_49],
				data:   make([]byte, 50),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_50 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_50],
				data:   make([]byte, 49),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_50",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_50],
				data:   make([]byte, 50),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_50 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_50],
				data:   make([]byte, 51),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_51 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_51],
				data:   make([]byte, 50),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_51",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_51],
				data:   make([]byte, 51),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_51 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_51],
				data:   make([]byte, 52),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_52 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_52],
				data:   make([]byte, 51),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_52",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_52],
				data:   make([]byte, 52),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_52 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_52],
				data:   make([]byte, 53),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_53 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_53],
				data:   make([]byte, 52),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_53",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_53],
				data:   make([]byte, 53),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_53 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_53],
				data:   make([]byte, 54),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_54 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_54],
				data:   make([]byte, 53),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_54",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_54],
				data:   make([]byte, 54),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_54 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_54],
				data:   make([]byte, 55),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_55 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_55],
				data:   make([]byte, 54),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_55",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_55],
				data:   make([]byte, 55),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_55 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_55],
				data:   make([]byte, 56),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_56 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_56],
				data:   make([]byte, 55),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_56",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_56],
				data:   make([]byte, 56),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_56 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_56],
				data:   make([]byte, 57),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_57 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_57],
				data:   make([]byte, 56),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_57",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_57],
				data:   make([]byte, 57),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_57 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_57],
				data:   make([]byte, 58),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_58 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_58],
				data:   make([]byte, 57),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_58",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_58],
				data:   make([]byte, 58),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_58 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_58],
				data:   make([]byte, 59),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_59 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_59],
				data:   make([]byte, 58),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_59",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_59],
				data:   make([]byte, 59),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_59 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_59],
				data:   make([]byte, 60),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_60 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_60],
				data:   make([]byte, 59),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_60",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_60],
				data:   make([]byte, 60),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_60 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_60],
				data:   make([]byte, 61),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_61 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_61],
				data:   make([]byte, 60),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_61",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_61],
				data:   make([]byte, 61),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_61 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_61],
				data:   make([]byte, 62),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_62 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_62],
				data:   make([]byte, 61),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_62",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_62],
				data:   make([]byte, 62),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_62 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_62],
				data:   make([]byte, 63),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_63 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_63],
				data:   make([]byte, 62),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_63",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_63],
				data:   make([]byte, 63),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_63 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_63],
				data:   make([]byte, 64),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_64 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_64],
				data:   make([]byte, 63),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_64",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_64],
				data:   make([]byte, 64),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_64 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_64],
				data:   make([]byte, 65),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_65 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_65],
				data:   make([]byte, 64),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_65",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_65],
				data:   make([]byte, 65),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_65 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_65],
				data:   make([]byte, 66),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_66 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_66],
				data:   make([]byte, 65),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_66",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_66],
				data:   make([]byte, 66),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_66 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_66],
				data:   make([]byte, 67),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_67 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_67],
				data:   make([]byte, 66),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_67",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_67],
				data:   make([]byte, 67),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_67 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_67],
				data:   make([]byte, 68),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_68 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_68],
				data:   make([]byte, 67),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_68",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_68],
				data:   make([]byte, 68),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_68 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_68],
				data:   make([]byte, 69),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_69 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_69],
				data:   make([]byte, 68),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_69",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_69],
				data:   make([]byte, 69),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_69 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_69],
				data:   make([]byte, 70),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_70 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_70],
				data:   make([]byte, 69),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_70",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_70],
				data:   make([]byte, 70),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_70 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_70],
				data:   make([]byte, 71),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_71 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_71],
				data:   make([]byte, 70),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_71",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_71],
				data:   make([]byte, 71),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_71 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_71],
				data:   make([]byte, 72),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_72 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_72],
				data:   make([]byte, 71),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_72",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_72],
				data:   make([]byte, 72),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_72 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_72],
				data:   make([]byte, 73),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_73 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_73],
				data:   make([]byte, 72),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_73",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_73],
				data:   make([]byte, 73),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_73 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_73],
				data:   make([]byte, 74),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_74 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_74],
				data:   make([]byte, 73),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_74",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_74],
				data:   make([]byte, 74),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_74 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_74],
				data:   make([]byte, 75),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_75 short",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_75],
				data:   make([]byte, 74),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DATA_75",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_75],
				data:   make([]byte, 75),
			},
			expectedErr: nil,
		},
		{
			name: "OP_DATA_75 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DATA_75],
				data:   make([]byte, 76),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_PUSHDATA1",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUSHDATA1],
				data:   []byte{0, 1, 2, 3, 4},
			},
			expectedErr: nil,
		},
		{
			name: "OP_PUSHDATA2",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUSHDATA2],
				data:   []byte{0, 1, 2, 3, 4},
			},
			expectedErr: nil,
		},
		{
			name: "OP_PUSHDATA4",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUSHDATA1],
				data:   []byte{0, 1, 2, 3, 4},
			},
			expectedErr: nil,
		},
		{
			name: "OP_TRUE",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_TRUE],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_TRUE long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_TRUE],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_2",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_2],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_2 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_2],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_2",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_2],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_2 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_2],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_3",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_3],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_3 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_3],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_4",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_4],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_4 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_4],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_5",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_5],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_5 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_5],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_6",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_6],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_6 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_6],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_7",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_7],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_7 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_7],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_8",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_8],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_8 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_8],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_9",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_9],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_9 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_9],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_10",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_10],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_10 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_10],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_11",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_11],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_11 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_11],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_12",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_12],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_12 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_12],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_13",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_13],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_13 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_13],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_14",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_14],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_14 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_14],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_15",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_15],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_15 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_15],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_16",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_16],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_16 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_16],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_VERIFY",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_VERIFY],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_VERIFY long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_VERIFY],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_DUP",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DUP],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_DUP long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_DUP],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_EQUAL",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_EQUAL],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_EQUAL long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_EQUAL],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_EQUALVERIFY",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_EQUALVERIFY],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_EQUALVERIFY long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_EQUALVERIFY],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_RIPEMD160",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_RIPEMD160],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_RIPEMD160 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_RIPEMD160],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_SHA1",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_SHA1],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_SHA1 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_SHA1],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_SHA256",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_SHA256],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_SHA256 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_SHA256],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_HASH160",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_HASH160],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_HASH160 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_HASH160],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_HASH256",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_HASH256],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_HASH256 long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_HASH256],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_CHECKSIG",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKSIG],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_CHECKSIG long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKSIG],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_CHECKSIGVERIFY",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKSIGVERIFY],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_CHECKSIGVERIFY long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKSIGVERIFY],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_CHECKMULTISIG",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKMULTISIG],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_CHECKMULTISIG long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKMULTISIG],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_CHECKMULTISIGVERIFY",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKMULTISIGVERIFY],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_CHECKMULTISIGVERIFY long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_CHECKMULTISIGVERIFY],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_PUBKEYHASH",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUBKEYHASH],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_PUBKEYHASH long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUBKEYHASH],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_PUBKEY",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUBKEY],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_PUBKEY long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_PUBKEY],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
		{
			name: "OP_INVALIDOPCODE",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_INVALIDOPCODE],
				data:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "OP_INVALIDOPCODE long",
			pop: &parsedOpcode{
				opcode: &opcodeArray[OP_INVALIDOPCODE],
				data:   make([]byte, 1),
			},
			expectedErr: scriptError(ErrInternal, ""),
		},
	}

	for _, test := range tests {
		_, err := test.pop.bytes()
		if e := tstCheckScriptError(err, test.expectedErr); e != nil {
			t.Errorf("Parsed opcode test '%s': %v", test.name, e)
			continue
		}
	}
}

// TestRemoveOpcodeByData ensures that removing data carrying opcodes based on
// the data they contain works as expected.
func TestRemoveOpcodeByData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []byte
		remove []byte
		err    error
		after  []byte
	}{
		{
			name:   "simple case",
			before: []byte{OP_DATA_4, 1, 2, 3, 4},
			remove: []byte{1, 2, 3, 4},
			after:  nil,
		},
		{
			name:   "simple case (miss)",
			before: []byte{OP_DATA_4, 1, 2, 3, 4},
			remove: []byte{1, 2, 3, 5},
			after:  []byte{OP_DATA_4, 1, 2, 3, 4},
		},
		{
			// padded to keep it canonical.
			name: "simple case (pushdata1)",
			before: append(append([]byte{OP_PUSHDATA1, 76},
				bytes.Repeat([]byte{0}, 72)...),
				[]byte{1, 2, 3, 4}...),
			remove: []byte{1, 2, 3, 4},
			after:  nil,
		},
		{
			name: "simple case (pushdata1 miss)",
			before: append(append([]byte{OP_PUSHDATA1, 76},
				bytes.Repeat([]byte{0}, 72)...),
				[]byte{1, 2, 3, 4}...),
			remove: []byte{1, 2, 3, 5},
			after: append(append([]byte{OP_PUSHDATA1, 76},
				bytes.Repeat([]byte{0}, 72)...),
				[]byte{1, 2, 3, 4}...),
		},
		{
			name:   "simple case (pushdata1 miss noncanonical)",
			before: []byte{OP_PUSHDATA1, 4, 1, 2, 3, 4},
			remove: []byte{1, 2, 3, 4},
			after:  []byte{OP_PUSHDATA1, 4, 1, 2, 3, 4},
		},
		{
			name: "simple case (pushdata2)",
			before: append(append([]byte{OP_PUSHDATA2, 0, 1},
				bytes.Repeat([]byte{0}, 252)...),
				[]byte{1, 2, 3, 4}...),
			remove: []byte{1, 2, 3, 4},
			after:  nil,
		},
		{
			name: "simple case (pushdata2 miss)",
			before: append(append([]byte{OP_PUSHDATA2, 0, 1},
				bytes.Repeat([]byte{0}, 252)...),
				[]byte{1, 2, 3, 4}...),
			remove: []byte{1, 2, 3, 4, 5},
			after: append(append([]byte{OP_PUSHDATA2, 0, 1},
				bytes.Repeat([]byte{0}, 252)...),
				[]byte{1, 2, 3, 4}...),
		},
		{
			name:   "simple case (pushdata2 miss noncanonical)",
			before: []byte{OP_PUSHDATA2, 4, 0, 1, 2, 3, 4},
			remove: []byte{1, 2, 3, 4},
			after:  []byte{OP_PUSHDATA2, 4, 0, 1, 2, 3, 4},
		},
		{
			// This is padded to make the push canonical.
			name: "simple case (pushdata4)",
			before: append(append([]byte{OP_PUSHDATA4, 0, 0, 1, 0},
				bytes.Repeat([]byte{0}, 65532)...),
				[]byte{1, 2, 3, 4}...),
			remove: []byte{1, 2, 3, 4},
			after:  nil,
		},
		{
			name:   "simple case (pushdata4 miss noncanonical)",
			before: []byte{OP_PUSHDATA4, 4, 0, 0, 0, 1, 2, 3, 4},
			remove: []byte{1, 2, 3, 4},
			after:  []byte{OP_PUSHDATA4, 4, 0, 0, 0, 1, 2, 3, 4},
		},
		{
			// This is padded to make the push canonical.
			name: "simple case (pushdata4 miss)",
			before: append(append([]byte{OP_PUSHDATA4, 0, 0, 1, 0},
				bytes.Repeat([]byte{0}, 65532)...), []byte{1, 2, 3, 4}...),
			remove: []byte{1, 2, 3, 4, 5},
			after: append(append([]byte{OP_PUSHDATA4, 0, 0, 1, 0},
				bytes.Repeat([]byte{0}, 65532)...), []byte{1, 2, 3, 4}...),
		},
		{
			name:   "invalid opcode ",
			before: []byte{OP_UNKNOWN187},
			remove: []byte{1, 2, 3, 4},
			after:  []byte{OP_UNKNOWN187},
		},
		{
			name:   "invalid length (instruction)",
			before: []byte{OP_PUSHDATA1},
			remove: []byte{1, 2, 3, 4},
			err:    scriptError(ErrMalformedPush, ""),
		},
		{
			name:   "invalid length (data)",
			before: []byte{OP_PUSHDATA1, 255, 254},
			remove: []byte{1, 2, 3, 4},
			err:    scriptError(ErrMalformedPush, ""),
		},
	}

	// tstRemoveOpcodeByData is a convenience function to parse the provided
	// raw script, remove the passed data, then unparse the result back
	// into a raw script.
	tstRemoveOpcodeByData := func(script []byte, data []byte) ([]byte, error) {
		pops, err := parseScript(script)
		if err != nil {
			return nil, err
		}
		pops = removeOpcodeByData(pops, data)
		return unparseScript(pops)
	}

	for _, test := range tests {
		result, err := tstRemoveOpcodeByData(test.before, test.remove)
		if e := tstCheckScriptError(err, test.err); e != nil {
			t.Errorf("%s: %v", test.name, e)
			continue
		}

		if !bytes.Equal(test.after, result) {
			t.Errorf("%s: value does not equal expected: exp: %q"+
				" got: %q", test.name, test.after, result)
		}
	}
}

// TestIsUnspendable ensures the IsUnspendable function returns the expected
// results.
func TestIsUnspendable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pkScript []byte
		expected bool
	}{
		{
			// Unspendable
			pkScript: []byte{0x6a, 0x04, 0x74, 0x65, 0x73, 0x74},
			expected: true,
		},
		{
			// Spendable
			pkScript: []byte{0x76, 0xa9, 0x15, 0x66, 0x29, 0x95, 0xa0,
				0xfe, 0x68, 0x43, 0xfa, 0x9b, 0x95, 0x45,
				0x97, 0xf0, 0xdc, 0xa7, 0xa4, 0x4d, 0xf6,
				0xfa, 0x0b, 0x5c, 0xc5, 0xac},
			expected: false,
		},
	}

	for i, test := range tests {
		res := IsUnspendable(test.pkScript)
		if res != test.expected {
			t.Errorf("TestIsUnspendable #%d failed: got %v want %v",
				i, res, test.expected)
			continue
		}
	}
}
