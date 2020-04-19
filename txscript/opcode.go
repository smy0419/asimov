// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txscript

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"golang.org/x/crypto/ripemd160"

	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/crypto"
)

// An opcode defines the information related to a txscript opcode.  opfunc, if
// present, is the function to call to perform the opcode on the script.  The
// current script is passed in as a slice with the first member being the opcode
// itself.
type opcode struct {
	value  byte
	name   string
	length int
	opfunc func(*parsedOpcode, *Engine) error
}

// These constants are the values of the official opcodes used on the btc wiki,
// in bitcoin core and in most if not all other references and software related
// to handling BTC scripts.
const (
	OP_0                   = 0x00 // 0
	OP_FALSE               = 0x00 // 0 - AKA OP_0
	OP_DATA_1              = 0x01 // 1
	OP_DATA_2              = 0x02 // 2
	OP_DATA_3              = 0x03 // 3
	OP_DATA_4              = 0x04 // 4
	OP_DATA_5              = 0x05 // 5
	OP_DATA_6              = 0x06 // 6
	OP_DATA_7              = 0x07 // 7
	OP_DATA_8              = 0x08 // 8
	OP_DATA_9              = 0x09 // 9
	OP_DATA_10             = 0x0a // 10
	OP_DATA_11             = 0x0b // 11
	OP_DATA_12             = 0x0c // 12
	OP_DATA_13             = 0x0d // 13
	OP_DATA_14             = 0x0e // 14
	OP_DATA_15             = 0x0f // 15
	OP_DATA_16             = 0x10 // 16
	OP_DATA_17             = 0x11 // 17
	OP_DATA_18             = 0x12 // 18
	OP_DATA_19             = 0x13 // 19
	OP_DATA_20             = 0x14 // 20
	OP_DATA_21             = 0x15 // 21
	OP_DATA_22             = 0x16 // 22
	OP_DATA_23             = 0x17 // 23
	OP_DATA_24             = 0x18 // 24
	OP_DATA_25             = 0x19 // 25
	OP_DATA_26             = 0x1a // 26
	OP_DATA_27             = 0x1b // 27
	OP_DATA_28             = 0x1c // 28
	OP_DATA_29             = 0x1d // 29
	OP_DATA_30             = 0x1e // 30
	OP_DATA_31             = 0x1f // 31
	OP_DATA_32             = 0x20 // 32
	OP_DATA_33             = 0x21 // 33
	OP_DATA_34             = 0x22 // 34
	OP_DATA_35             = 0x23 // 35
	OP_DATA_36             = 0x24 // 36
	OP_DATA_37             = 0x25 // 37
	OP_DATA_38             = 0x26 // 38
	OP_DATA_39             = 0x27 // 39
	OP_DATA_40             = 0x28 // 40
	OP_DATA_41             = 0x29 // 41
	OP_DATA_42             = 0x2a // 42
	OP_DATA_43             = 0x2b // 43
	OP_DATA_44             = 0x2c // 44
	OP_DATA_45             = 0x2d // 45
	OP_DATA_46             = 0x2e // 46
	OP_DATA_47             = 0x2f // 47
	OP_DATA_48             = 0x30 // 48
	OP_DATA_49             = 0x31 // 49
	OP_DATA_50             = 0x32 // 50
	OP_DATA_51             = 0x33 // 51
	OP_DATA_52             = 0x34 // 52
	OP_DATA_53             = 0x35 // 53
	OP_DATA_54             = 0x36 // 54
	OP_DATA_55             = 0x37 // 55
	OP_DATA_56             = 0x38 // 56
	OP_DATA_57             = 0x39 // 57
	OP_DATA_58             = 0x3a // 58
	OP_DATA_59             = 0x3b // 59
	OP_DATA_60             = 0x3c // 60
	OP_DATA_61             = 0x3d // 61
	OP_DATA_62             = 0x3e // 62
	OP_DATA_63             = 0x3f // 63
	OP_DATA_64             = 0x40 // 64
	OP_DATA_65             = 0x41 // 65
	OP_DATA_66             = 0x42 // 66
	OP_DATA_67             = 0x43 // 67
	OP_DATA_68             = 0x44 // 68
	OP_DATA_69             = 0x45 // 69
	OP_DATA_70             = 0x46 // 70
	OP_DATA_71             = 0x47 // 71
	OP_DATA_72             = 0x48 // 72
	OP_DATA_73             = 0x49 // 73
	OP_DATA_74             = 0x4a // 74
	OP_DATA_75             = 0x4b // 75
	OP_PUSHDATA1           = 0x4c // 76
	OP_PUSHDATA2           = 0x4d // 77
	OP_PUSHDATA4           = 0x4e // 78
	OP_UNKNOWN79           = 0x4f // 79
	OP_UNKNOWN80           = 0x50 // 80
	OP_1                   = 0x51 // 81 - AKA OP_TRUE
	OP_TRUE                = 0x51 // 81
	OP_2                   = 0x52 // 82
	OP_3                   = 0x53 // 83
	OP_4                   = 0x54 // 84
	OP_5                   = 0x55 // 85
	OP_6                   = 0x56 // 86
	OP_7                   = 0x57 // 87
	OP_8                   = 0x58 // 88
	OP_9                   = 0x59 // 89
	OP_10                  = 0x5a // 90
	OP_11                  = 0x5b // 91
	OP_12                  = 0x5c // 92
	OP_13                  = 0x5d // 93
	OP_14                  = 0x5e // 94
	OP_15                  = 0x5f // 95
	OP_16                  = 0x60 // 96
	OP_UNKNOWN97           = 0x61 // 97
	OP_UNKNOWN98           = 0x62 // 98
	OP_UNKNOWN99           = 0x63 // 99
	OP_UNKNOWN100          = 0x64 // 100
	OP_UNKNOWN101          = 0x65 // 101
	OP_UNKNOWN102          = 0x66 // 102
	OP_UNKNOWN103          = 0x67 // 103
	OP_UNKNOWN104          = 0x68 // 104
	OP_VERIFY              = 0x69 // 105
	OP_UNKNOWN106          = 0x6a // 106
	OP_UNKNOWN107          = 0x6b // 107
	OP_UNKNOWN108          = 0x6c // 108
	OP_UNKNOWN109          = 0x6d // 109
	OP_UNKNOWN110          = 0x6e // 110
	OP_UNKNOWN111          = 0x6f // 111
	OP_UNKNOWN112          = 0x70 // 112
	OP_UNKNOWN113          = 0x71 // 113
	OP_UNKNOWN114          = 0x72 // 114
	OP_UNKNOWN115          = 0x73 // 115
	OP_UNKNOWN116          = 0x74 // 116
	OP_UNKNOWN117          = 0x75 // 117
	OP_DUP                 = 0x76 // 118
	OP_UNKNOWN119          = 0x77 // 119
	OP_UNKNOWN120          = 0x78 // 120
	OP_UNKNOWN121          = 0x79 // 121
	OP_UNKNOWN122          = 0x7a // 122
	OP_UNKNOWN123          = 0x7b // 123
	OP_UNKNOWN124          = 0x7c // 124
	OP_UNKNOWN125          = 0x7d // 125
	OP_UNKNOWN126          = 0x7e // 126
	OP_UNKNOWN127          = 0x7f // 127
	OP_UNKNOWN128          = 0x80 // 128
	OP_UNKNOWN129          = 0x81 // 129
	OP_UNKNOWN130          = 0x82 // 130
	OP_UNKNOWN131          = 0x83 // 131
	OP_UNKNOWN132          = 0x84 // 132
	OP_UNKNOWN133          = 0x85 // 133
	OP_UNKNOWN134          = 0x86 // 134
	OP_EQUAL               = 0x87 // 135
	OP_EQUALVERIFY         = 0x88 // 136
	OP_UNKNOWN137          = 0x89 // 137
	OP_UNKNOWN138          = 0x8a // 138
	OP_UNKNOWN139          = 0x8b // 139
	OP_UNKNOWN140          = 0x8c // 140
	OP_UNKNOWN141          = 0x8d // 141
	OP_UNKNOWN142          = 0x8e // 142
	OP_UNKNOWN143          = 0x8f // 143
	OP_UNKNOWN144          = 0x90 // 144
	OP_UNKNOWN145          = 0x91 // 145
	OP_UNKNOWN146          = 0x92 // 146
	OP_UNKNOWN147          = 0x93 // 147
	OP_UNKNOWN148          = 0x94 // 148
	OP_UNKNOWN149          = 0x95 // 149
	OP_UNKNOWN150          = 0x96 // 150
	OP_UNKNOWN151          = 0x97 // 151
	OP_UNKNOWN152          = 0x98 // 152
	OP_UNKNOWN153          = 0x99 // 153
	OP_UNKNOWN154          = 0x9a // 154
	OP_UNKNOWN155          = 0x9b // 155
	OP_UNKNOWN156          = 0x9c // 156
	OP_UNKNOWN157          = 0x9d // 157
	OP_UNKNOWN158          = 0x9e // 158
	OP_UNKNOWN159          = 0x9f // 159
	OP_UNKNOWN160          = 0xa0 // 160
	OP_UNKNOWN161          = 0xa1 // 161
	OP_UNKNOWN162          = 0xa2 // 162
	OP_UNKNOWN163          = 0xa3 // 163
	OP_UNKNOWN164          = 0xa4 // 164
	OP_UNKNOWN165          = 0xa5 // 165
	OP_RIPEMD160           = 0xa6 // 166
	OP_SHA1                = 0xa7 // 167
	OP_SHA256              = 0xa8 // 168
	OP_HASH160             = 0xa9 // 169
	OP_HASH256             = 0xaa // 170
	OP_UNKNOWN171          = 0xab // 171
	OP_CHECKSIG            = 0xac // 172
	OP_CHECKSIGVERIFY      = 0xad // 173
	OP_CHECKMULTISIG       = 0xae // 174
	OP_CHECKMULTISIGVERIFY = 0xaf // 175
	OP_UNKNOWN176          = 0xb0 // 176
	OP_CHECKLOCKTIMEVERIFY = 0xb1 // 177
	OP_CHECKSEQUENCEVERIFY = 0xb2 // 178
	OP_UNKNOWN179          = 0xb3 // 179
	OP_UNKNOWN180          = 0xb4 // 180
	OP_UNKNOWN181          = 0xb5 // 181
	OP_UNKNOWN182          = 0xb6 // 182
	OP_UNKNOWN183          = 0xb7 // 183
	OP_UNKNOWN184          = 0xb8 // 184
	OP_UNKNOWN185          = 0xb9 // 185
	OP_UNKNOWN186          = 0xba // 186
	OP_UNKNOWN187          = 0xbb // 187
	OP_UNKNOWN188          = 0xbc // 188
	OP_UNKNOWN189          = 0xbd // 189
	OP_UNKNOWN190          = 0xbe // 190
	OP_UNKNOWN191          = 0xbf // 191
	OP_TEMPLATE            = 0xc0 // 192
	OP_CREATE              = 0xc1 // 193
	OP_CALL                = 0xc2 // 194
	OP_SPEND               = 0xc3 // 195
	OP_IFLAG_EQUAL         = 0xc4 // 196
	OP_IFLAG_EQUALVERIFY   = 0xc5 // 197
	OP_VOTE                = 0xc6 // 198
	OP_UNKNOWN199          = 0xc7 // 199
	OP_UNKNOWN200          = 0xc8 // 200
	OP_UNKNOWN201          = 0xc9 // 201
	OP_UNKNOWN202          = 0xca // 202
	OP_UNKNOWN203          = 0xcb // 203
	OP_UNKNOWN204          = 0xcc // 204
	OP_UNKNOWN205          = 0xcd // 205
	OP_UNKNOWN206          = 0xce // 206
	OP_UNKNOWN207          = 0xcf // 207
	OP_UNKNOWN208          = 0xd0 // 208
	OP_UNKNOWN209          = 0xd1 // 209
	OP_UNKNOWN210          = 0xd2 // 210
	OP_UNKNOWN211          = 0xd3 // 211
	OP_UNKNOWN212          = 0xd4 // 212
	OP_UNKNOWN213          = 0xd5 // 213
	OP_UNKNOWN214          = 0xd6 // 214
	OP_UNKNOWN215          = 0xd7 // 215
	OP_UNKNOWN216          = 0xd8 // 216
	OP_UNKNOWN217          = 0xd9 // 217
	OP_UNKNOWN218          = 0xda // 218
	OP_UNKNOWN219          = 0xdb // 219
	OP_UNKNOWN220          = 0xdc // 220
	OP_UNKNOWN221          = 0xdd // 221
	OP_UNKNOWN222          = 0xde // 222
	OP_UNKNOWN223          = 0xdf // 223
	OP_UNKNOWN224          = 0xe0 // 224
	OP_UNKNOWN225          = 0xe1 // 225
	OP_UNKNOWN226          = 0xe2 // 226
	OP_UNKNOWN227          = 0xe3 // 227
	OP_UNKNOWN228          = 0xe4 // 228
	OP_UNKNOWN229          = 0xe5 // 229
	OP_UNKNOWN230          = 0xe6 // 230
	OP_UNKNOWN231          = 0xe7 // 231
	OP_UNKNOWN232          = 0xe8 // 232
	OP_UNKNOWN233          = 0xe9 // 233
	OP_UNKNOWN234          = 0xea // 234
	OP_UNKNOWN235          = 0xeb // 235
	OP_UNKNOWN236          = 0xec // 236
	OP_UNKNOWN237          = 0xed // 237
	OP_UNKNOWN238          = 0xee // 238
	OP_UNKNOWN239          = 0xef // 239
	OP_UNKNOWN240          = 0xf0 // 240
	OP_UNKNOWN241          = 0xf1 // 241
	OP_UNKNOWN242          = 0xf2 // 242
	OP_UNKNOWN243          = 0xf3 // 243
	OP_UNKNOWN244          = 0xf4 // 244
	OP_UNKNOWN245          = 0xf5 // 245
	OP_UNKNOWN246          = 0xf6 // 246
	OP_UNKNOWN247          = 0xf7 // 247
	OP_UNKNOWN248          = 0xf8 // 248
	OP_UNKNOWN249          = 0xf9 // 249
	OP_SMALLINTEGER        = 0xfa // 250 - bitcoin core internal
	OP_PUBKEYS             = 0xfb // 251 - bitcoin core internal
	OP_UNKNOWN252          = 0xfc // 252
	OP_PUBKEYHASH          = 0xfd // 253 - bitcoin core internal
	OP_PUBKEY              = 0xfe // 254 - bitcoin core internal
	OP_INVALIDOPCODE       = 0xff // 255 - bitcoin core internal
)

// Conditional execution constants.
const (
	OpCondTrue  = 1
)

// opcodeArray holds details about all possible opcodes such as how many bytes
// the opcode and any associated data should take, its human-readable name, and
// the handler function.
var opcodeArray = [256]opcode{
	// Data push opcodes.
	OP_FALSE:     {OP_FALSE, "OP_0", 1, opcodeFalse},
	OP_DATA_1:    {OP_DATA_1, "OP_DATA_1", 2, opcodePushData},
	OP_DATA_2:    {OP_DATA_2, "OP_DATA_2", 3, opcodePushData},
	OP_DATA_3:    {OP_DATA_3, "OP_DATA_3", 4, opcodePushData},
	OP_DATA_4:    {OP_DATA_4, "OP_DATA_4", 5, opcodePushData},
	OP_DATA_5:    {OP_DATA_5, "OP_DATA_5", 6, opcodePushData},
	OP_DATA_6:    {OP_DATA_6, "OP_DATA_6", 7, opcodePushData},
	OP_DATA_7:    {OP_DATA_7, "OP_DATA_7", 8, opcodePushData},
	OP_DATA_8:    {OP_DATA_8, "OP_DATA_8", 9, opcodePushData},
	OP_DATA_9:    {OP_DATA_9, "OP_DATA_9", 10, opcodePushData},
	OP_DATA_10:   {OP_DATA_10, "OP_DATA_10", 11, opcodePushData},
	OP_DATA_11:   {OP_DATA_11, "OP_DATA_11", 12, opcodePushData},
	OP_DATA_12:   {OP_DATA_12, "OP_DATA_12", 13, opcodePushData},
	OP_DATA_13:   {OP_DATA_13, "OP_DATA_13", 14, opcodePushData},
	OP_DATA_14:   {OP_DATA_14, "OP_DATA_14", 15, opcodePushData},
	OP_DATA_15:   {OP_DATA_15, "OP_DATA_15", 16, opcodePushData},
	OP_DATA_16:   {OP_DATA_16, "OP_DATA_16", 17, opcodePushData},
	OP_DATA_17:   {OP_DATA_17, "OP_DATA_17", 18, opcodePushData},
	OP_DATA_18:   {OP_DATA_18, "OP_DATA_18", 19, opcodePushData},
	OP_DATA_19:   {OP_DATA_19, "OP_DATA_19", 20, opcodePushData},
	OP_DATA_20:   {OP_DATA_20, "OP_DATA_20", 21, opcodePushData},
	OP_DATA_21:   {OP_DATA_21, "OP_DATA_21", 22, opcodePushData},
	OP_DATA_22:   {OP_DATA_22, "OP_DATA_22", 23, opcodePushData},
	OP_DATA_23:   {OP_DATA_23, "OP_DATA_23", 24, opcodePushData},
	OP_DATA_24:   {OP_DATA_24, "OP_DATA_24", 25, opcodePushData},
	OP_DATA_25:   {OP_DATA_25, "OP_DATA_25", 26, opcodePushData},
	OP_DATA_26:   {OP_DATA_26, "OP_DATA_26", 27, opcodePushData},
	OP_DATA_27:   {OP_DATA_27, "OP_DATA_27", 28, opcodePushData},
	OP_DATA_28:   {OP_DATA_28, "OP_DATA_28", 29, opcodePushData},
	OP_DATA_29:   {OP_DATA_29, "OP_DATA_29", 30, opcodePushData},
	OP_DATA_30:   {OP_DATA_30, "OP_DATA_30", 31, opcodePushData},
	OP_DATA_31:   {OP_DATA_31, "OP_DATA_31", 32, opcodePushData},
	OP_DATA_32:   {OP_DATA_32, "OP_DATA_32", 33, opcodePushData},
	OP_DATA_33:   {OP_DATA_33, "OP_DATA_33", 34, opcodePushData},
	OP_DATA_34:   {OP_DATA_34, "OP_DATA_34", 35, opcodePushData},
	OP_DATA_35:   {OP_DATA_35, "OP_DATA_35", 36, opcodePushData},
	OP_DATA_36:   {OP_DATA_36, "OP_DATA_36", 37, opcodePushData},
	OP_DATA_37:   {OP_DATA_37, "OP_DATA_37", 38, opcodePushData},
	OP_DATA_38:   {OP_DATA_38, "OP_DATA_38", 39, opcodePushData},
	OP_DATA_39:   {OP_DATA_39, "OP_DATA_39", 40, opcodePushData},
	OP_DATA_40:   {OP_DATA_40, "OP_DATA_40", 41, opcodePushData},
	OP_DATA_41:   {OP_DATA_41, "OP_DATA_41", 42, opcodePushData},
	OP_DATA_42:   {OP_DATA_42, "OP_DATA_42", 43, opcodePushData},
	OP_DATA_43:   {OP_DATA_43, "OP_DATA_43", 44, opcodePushData},
	OP_DATA_44:   {OP_DATA_44, "OP_DATA_44", 45, opcodePushData},
	OP_DATA_45:   {OP_DATA_45, "OP_DATA_45", 46, opcodePushData},
	OP_DATA_46:   {OP_DATA_46, "OP_DATA_46", 47, opcodePushData},
	OP_DATA_47:   {OP_DATA_47, "OP_DATA_47", 48, opcodePushData},
	OP_DATA_48:   {OP_DATA_48, "OP_DATA_48", 49, opcodePushData},
	OP_DATA_49:   {OP_DATA_49, "OP_DATA_49", 50, opcodePushData},
	OP_DATA_50:   {OP_DATA_50, "OP_DATA_50", 51, opcodePushData},
	OP_DATA_51:   {OP_DATA_51, "OP_DATA_51", 52, opcodePushData},
	OP_DATA_52:   {OP_DATA_52, "OP_DATA_52", 53, opcodePushData},
	OP_DATA_53:   {OP_DATA_53, "OP_DATA_53", 54, opcodePushData},
	OP_DATA_54:   {OP_DATA_54, "OP_DATA_54", 55, opcodePushData},
	OP_DATA_55:   {OP_DATA_55, "OP_DATA_55", 56, opcodePushData},
	OP_DATA_56:   {OP_DATA_56, "OP_DATA_56", 57, opcodePushData},
	OP_DATA_57:   {OP_DATA_57, "OP_DATA_57", 58, opcodePushData},
	OP_DATA_58:   {OP_DATA_58, "OP_DATA_58", 59, opcodePushData},
	OP_DATA_59:   {OP_DATA_59, "OP_DATA_59", 60, opcodePushData},
	OP_DATA_60:   {OP_DATA_60, "OP_DATA_60", 61, opcodePushData},
	OP_DATA_61:   {OP_DATA_61, "OP_DATA_61", 62, opcodePushData},
	OP_DATA_62:   {OP_DATA_62, "OP_DATA_62", 63, opcodePushData},
	OP_DATA_63:   {OP_DATA_63, "OP_DATA_63", 64, opcodePushData},
	OP_DATA_64:   {OP_DATA_64, "OP_DATA_64", 65, opcodePushData},
	OP_DATA_65:   {OP_DATA_65, "OP_DATA_65", 66, opcodePushData},
	OP_DATA_66:   {OP_DATA_66, "OP_DATA_66", 67, opcodePushData},
	OP_DATA_67:   {OP_DATA_67, "OP_DATA_67", 68, opcodePushData},
	OP_DATA_68:   {OP_DATA_68, "OP_DATA_68", 69, opcodePushData},
	OP_DATA_69:   {OP_DATA_69, "OP_DATA_69", 70, opcodePushData},
	OP_DATA_70:   {OP_DATA_70, "OP_DATA_70", 71, opcodePushData},
	OP_DATA_71:   {OP_DATA_71, "OP_DATA_71", 72, opcodePushData},
	OP_DATA_72:   {OP_DATA_72, "OP_DATA_72", 73, opcodePushData},
	OP_DATA_73:   {OP_DATA_73, "OP_DATA_73", 74, opcodePushData},
	OP_DATA_74:   {OP_DATA_74, "OP_DATA_74", 75, opcodePushData},
	OP_DATA_75:   {OP_DATA_75, "OP_DATA_75", 76, opcodePushData},
	OP_PUSHDATA1: {OP_PUSHDATA1, "OP_PUSHDATA1", -1, opcodePushData},
	OP_PUSHDATA2: {OP_PUSHDATA2, "OP_PUSHDATA2", -2, opcodePushData},
	OP_PUSHDATA4: {OP_PUSHDATA4, "OP_PUSHDATA4", -4, opcodePushData},
	OP_UNKNOWN79: {OP_UNKNOWN79, "OP_UNKNOWN79", 1, opcodeInvalid},
	OP_UNKNOWN80: {OP_UNKNOWN80, "OP_UNKNOWN80", 1, opcodeInvalid},
	OP_TRUE:      {OP_TRUE, "OP_1", 1, opcodeN},
	OP_2:         {OP_2, "OP_2", 1, opcodeN},
	OP_3:         {OP_3, "OP_3", 1, opcodeN},
	OP_4:         {OP_4, "OP_4", 1, opcodeN},
	OP_5:         {OP_5, "OP_5", 1, opcodeN},
	OP_6:         {OP_6, "OP_6", 1, opcodeN},
	OP_7:         {OP_7, "OP_7", 1, opcodeN},
	OP_8:         {OP_8, "OP_8", 1, opcodeN},
	OP_9:         {OP_9, "OP_9", 1, opcodeN},
	OP_10:        {OP_10, "OP_10", 1, opcodeN},
	OP_11:        {OP_11, "OP_11", 1, opcodeN},
	OP_12:        {OP_12, "OP_12", 1, opcodeN},
	OP_13:        {OP_13, "OP_13", 1, opcodeN},
	OP_14:        {OP_14, "OP_14", 1, opcodeN},
	OP_15:        {OP_15, "OP_15", 1, opcodeN},
	OP_16:        {OP_16, "OP_16", 1, opcodeN},

	// Control opcodes.
	OP_UNKNOWN97:           {OP_UNKNOWN97, "OP_UNKNOWN97", 1, opcodeInvalid},
	OP_UNKNOWN98:           {OP_UNKNOWN98, "OP_UNKNOWN98", 1, opcodeInvalid},
	OP_UNKNOWN99:           {OP_UNKNOWN99, "OP_UNKNOWN99", 1, opcodeInvalid},
	OP_UNKNOWN100:          {OP_UNKNOWN100, "OP_UNKNOWN100", 1, opcodeInvalid},
	OP_UNKNOWN101:          {OP_UNKNOWN101, "OP_UNKNOWN101", 1, opcodeInvalid},
	OP_UNKNOWN102:          {OP_UNKNOWN102, "OP_UNKNOWN102", 1, opcodeInvalid},
	OP_UNKNOWN103:          {OP_UNKNOWN103, "OP_UNKNOWN103", 1, opcodeInvalid},
	OP_UNKNOWN104:          {OP_UNKNOWN104, "OP_UNKNOWN104", 1, opcodeInvalid},
	OP_VERIFY:              {OP_VERIFY, "OP_VERIFY", 1, opcodeVerify},
	OP_UNKNOWN106:          {OP_UNKNOWN106, "OP_UNKNOWN106", 1, opcodeInvalid},
	OP_CHECKLOCKTIMEVERIFY: {OP_CHECKLOCKTIMEVERIFY, "OP_CHECKLOCKTIMEVERIFY", 1, opcodeInvalid},
	OP_CHECKSEQUENCEVERIFY: {OP_CHECKSEQUENCEVERIFY, "OP_CHECKSEQUENCEVERIFY", 1, opcodeInvalid},

	// Stack opcodes.
	OP_UNKNOWN107: {OP_UNKNOWN107, "OP_UNKNOWN107", 1, opcodeInvalid},
	OP_UNKNOWN108: {OP_UNKNOWN108, "OP_UNKNOWN108", 1, opcodeInvalid},
	OP_UNKNOWN109: {OP_UNKNOWN109, "OP_UNKNOWN109", 1, opcodeInvalid},
	OP_UNKNOWN110: {OP_UNKNOWN110, "OP_UNKNOWN110", 1, opcodeInvalid},
	OP_UNKNOWN111: {OP_UNKNOWN111, "OP_UNKNOWN111", 1, opcodeInvalid},
	OP_UNKNOWN112: {OP_UNKNOWN112, "OP_UNKNOWN112", 1, opcodeInvalid},
	OP_UNKNOWN113: {OP_UNKNOWN113, "OP_UNKNOWN113", 1, opcodeInvalid},
	OP_UNKNOWN114: {OP_UNKNOWN114, "OP_UNKNOWN114", 1, opcodeInvalid},
	OP_UNKNOWN115: {OP_UNKNOWN115, "OP_UNKNOWN115", 1, opcodeInvalid},
	OP_UNKNOWN116: {OP_UNKNOWN116, "OP_UNKNOWN116", 1, opcodeInvalid},
	OP_UNKNOWN117: {OP_UNKNOWN117, "OP_UNKNOWN117", 1, opcodeInvalid},
	OP_DUP:        {OP_DUP, "OP_DUP", 1, opcodeDup},
	OP_UNKNOWN119: {OP_UNKNOWN119, "OP_UNKNOWN119", 1, opcodeInvalid},
	OP_UNKNOWN120: {OP_UNKNOWN120, "OP_UNKNOWN120", 1, opcodeInvalid},
	OP_UNKNOWN121: {OP_UNKNOWN121, "OP_UNKNOWN121", 1, opcodeInvalid},
	OP_UNKNOWN122: {OP_UNKNOWN122, "OP_UNKNOWN122", 1, opcodeInvalid},
	OP_UNKNOWN123: {OP_UNKNOWN123, "OP_UNKNOWN123", 1, opcodeInvalid},
	OP_UNKNOWN124: {OP_UNKNOWN124, "OP_UNKNOWN124", 1, opcodeInvalid},
	OP_UNKNOWN125: {OP_UNKNOWN125, "OP_UNKNOWN125", 1, opcodeInvalid},
	OP_UNKNOWN126: {OP_UNKNOWN126, "OP_UNKNOWN126", 1, opcodeInvalid},
	OP_UNKNOWN127: {OP_UNKNOWN127, "OP_UNKNOWN127", 1, opcodeInvalid},
	OP_UNKNOWN128: {OP_UNKNOWN128, "OP_UNKNOWN128", 1, opcodeInvalid},
	OP_UNKNOWN129: {OP_UNKNOWN129, "OP_UNKNOWN129", 1, opcodeInvalid},
	OP_UNKNOWN130: {OP_UNKNOWN130, "OP_UNKNOWN130", 1, opcodeInvalid},
	OP_UNKNOWN131: {OP_UNKNOWN131, "OP_UNKNOWN131", 1, opcodeInvalid},
	OP_UNKNOWN132: {OP_UNKNOWN132, "OP_UNKNOWN132", 1, opcodeInvalid},
	OP_UNKNOWN133: {OP_UNKNOWN133, "OP_UNKNOWN133", 1, opcodeInvalid},
	OP_UNKNOWN134: {OP_UNKNOWN134, "OP_UNKNOWN134", 1, opcodeInvalid},

	// Bitwise logic opcodes.
	OP_EQUAL:       {OP_EQUAL, "OP_EQUAL", 1, opcodeEqual},
	OP_EQUALVERIFY: {OP_EQUALVERIFY, "OP_EQUALVERIFY", 1, opcodeEqualVerify},

	OP_UNKNOWN137: {OP_UNKNOWN137, "OP_UNKNOWN137", 1, opcodeInvalid},
	OP_UNKNOWN138: {OP_UNKNOWN138, "OP_UNKNOWN138", 1, opcodeInvalid},
	OP_UNKNOWN139: {OP_UNKNOWN139, "OP_UNKNOWN139", 1, opcodeInvalid},
	OP_UNKNOWN140: {OP_UNKNOWN140, "OP_UNKNOWN140", 1, opcodeInvalid},
	OP_UNKNOWN141: {OP_UNKNOWN141, "OP_UNKNOWN141", 1, opcodeInvalid},
	OP_UNKNOWN142: {OP_UNKNOWN142, "OP_UNKNOWN142", 1, opcodeInvalid},
	OP_UNKNOWN143: {OP_UNKNOWN143, "OP_UNKNOWN143", 1, opcodeInvalid},
	OP_UNKNOWN144: {OP_UNKNOWN144, "OP_UNKNOWN144", 1, opcodeInvalid},
	OP_UNKNOWN145: {OP_UNKNOWN145, "OP_UNKNOWN145", 1, opcodeInvalid},
	OP_UNKNOWN146: {OP_UNKNOWN146, "OP_UNKNOWN146", 1, opcodeInvalid},
	OP_UNKNOWN147: {OP_UNKNOWN147, "OP_UNKNOWN147", 1, opcodeInvalid},
	OP_UNKNOWN148: {OP_UNKNOWN148, "OP_UNKNOWN148", 1, opcodeInvalid},
	OP_UNKNOWN149: {OP_UNKNOWN149, "OP_UNKNOWN149", 1, opcodeInvalid},
	OP_UNKNOWN150: {OP_UNKNOWN150, "OP_UNKNOWN150", 1, opcodeInvalid},
	OP_UNKNOWN151: {OP_UNKNOWN151, "OP_UNKNOWN151", 1, opcodeInvalid},
	OP_UNKNOWN152: {OP_UNKNOWN152, "OP_UNKNOWN152", 1, opcodeInvalid},
	OP_UNKNOWN153: {OP_UNKNOWN153, "OP_UNKNOWN153", 1, opcodeInvalid},
	OP_UNKNOWN154: {OP_UNKNOWN154, "OP_UNKNOWN154", 1, opcodeInvalid},
	OP_UNKNOWN155: {OP_UNKNOWN155, "OP_UNKNOWN155", 1, opcodeInvalid},
	OP_UNKNOWN156: {OP_UNKNOWN156, "OP_UNKNOWN156", 1, opcodeInvalid},
	OP_UNKNOWN157: {OP_UNKNOWN157, "OP_UNKNOWN157", 1, opcodeInvalid},
	OP_UNKNOWN158: {OP_UNKNOWN158, "OP_UNKNOWN158", 1, opcodeInvalid},
	OP_UNKNOWN159: {OP_UNKNOWN159, "OP_UNKNOWN159", 1, opcodeInvalid},
	OP_UNKNOWN160: {OP_UNKNOWN160, "OP_UNKNOWN160", 1, opcodeInvalid},
	OP_UNKNOWN161: {OP_UNKNOWN161, "OP_UNKNOWN161", 1, opcodeInvalid},
	OP_UNKNOWN162: {OP_UNKNOWN162, "OP_UNKNOWN162", 1, opcodeInvalid},
	OP_UNKNOWN163: {OP_UNKNOWN163, "OP_UNKNOWN163", 1, opcodeInvalid},
	OP_UNKNOWN164: {OP_UNKNOWN164, "OP_UNKNOWN164", 1, opcodeInvalid},
	OP_UNKNOWN165: {OP_UNKNOWN165, "OP_UNKNOWN165", 1, opcodeInvalid},

	// Crypto opcodes.
	OP_RIPEMD160:           {OP_RIPEMD160, "OP_RIPEMD160", 1, opcodeRipemd160},
	OP_SHA1:                {OP_SHA1, "OP_SHA1", 1, opcodeSha1},
	OP_SHA256:              {OP_SHA256, "OP_SHA256", 1, opcodeSha256},
	OP_HASH160:             {OP_HASH160, "OP_HASH160", 1, opcodeHash160},
	OP_HASH256:             {OP_HASH256, "OP_HASH256", 1, opcodeHash256},
	OP_UNKNOWN171:          {OP_UNKNOWN171, "OP_UNKNOWN171", 1, opcodeInvalid},
	OP_CHECKSIG:            {OP_CHECKSIG, "OP_CHECKSIG", 1, opcodeCheckSig},
	OP_CHECKSIGVERIFY:      {OP_CHECKSIGVERIFY, "OP_CHECKSIGVERIFY", 1, opcodeCheckSigVerify},
	OP_CHECKMULTISIG:       {OP_CHECKMULTISIG, "OP_CHECKMULTISIG", 1, opcodeCheckMultiSig},
	OP_CHECKMULTISIGVERIFY: {OP_CHECKMULTISIGVERIFY, "OP_CHECKMULTISIGVERIFY", 1, opcodeCheckMultiSigVerify},

	OP_UNKNOWN176: {OP_UNKNOWN176, "OP_UNKNOWN176", 1, opcodeInvalid},
	OP_UNKNOWN179: {OP_UNKNOWN179, "OP_UNKNOWN179", 1, opcodeInvalid},
	OP_UNKNOWN180: {OP_UNKNOWN180, "OP_UNKNOWN180", 1, opcodeInvalid},
	OP_UNKNOWN181: {OP_UNKNOWN181, "OP_UNKNOWN181", 1, opcodeInvalid},
	OP_UNKNOWN182: {OP_UNKNOWN182, "OP_UNKNOWN182", 1, opcodeInvalid},
	OP_UNKNOWN183: {OP_UNKNOWN183, "OP_UNKNOWN183", 1, opcodeInvalid},
	OP_UNKNOWN184: {OP_UNKNOWN184, "OP_UNKNOWN184", 1, opcodeInvalid},
	OP_UNKNOWN185: {OP_UNKNOWN185, "OP_UNKNOWN185", 1, opcodeInvalid},
	OP_UNKNOWN186: {OP_UNKNOWN186, "OP_UNKNOWN186", 1, opcodeInvalid},
	OP_UNKNOWN187: {OP_UNKNOWN187, "OP_UNKNOWN187", 1, opcodeInvalid},
	OP_UNKNOWN188: {OP_UNKNOWN188, "OP_UNKNOWN188", 1, opcodeInvalid},
	OP_UNKNOWN189: {OP_UNKNOWN189, "OP_UNKNOWN189", 1, opcodeInvalid},
	OP_UNKNOWN190: {OP_UNKNOWN190, "OP_UNKNOWN190", 1, opcodeInvalid},
	OP_UNKNOWN191: {OP_UNKNOWN191, "OP_UNKNOWN191", 1, opcodeInvalid},
	OP_TEMPLATE:          {OP_TEMPLATE, "OP_TEMPLATE", 1, opcodeInvalid},
	OP_CREATE:            {OP_CREATE, "OP_CREATE", 1, opcodeInvalid},
	OP_CALL:              {OP_CALL, "OP_CALL", 1, opcodeInvalid},
	OP_SPEND:             {OP_SPEND, "OP_SPEND", 1, opcodeInvalid},
	OP_IFLAG_EQUAL:       {OP_IFLAG_EQUAL, "OP_IFLAG_EQUAL", 1, opcodeIgnoreFlagEqual},
	OP_IFLAG_EQUALVERIFY: {OP_IFLAG_EQUALVERIFY, "OP_IFLAG_EQUALVERIFY", 1, opcodeIgnoreFlagEqualVerify},
	OP_VOTE:              {OP_VOTE, "OP_VOTE", 1, opcodeInvalid},
	OP_UNKNOWN199:        {OP_UNKNOWN199, "OP_UNKNOWN199", 1, opcodeInvalid},
	OP_UNKNOWN200:        {OP_UNKNOWN200, "OP_UNKNOWN200", 1, opcodeInvalid},
	OP_UNKNOWN201:        {OP_UNKNOWN201, "OP_UNKNOWN201", 1, opcodeInvalid},
	OP_UNKNOWN202:        {OP_UNKNOWN202, "OP_UNKNOWN202", 1, opcodeInvalid},
	OP_UNKNOWN203:        {OP_UNKNOWN203, "OP_UNKNOWN203", 1, opcodeInvalid},
	OP_UNKNOWN204:        {OP_UNKNOWN204, "OP_UNKNOWN204", 1, opcodeInvalid},
	OP_UNKNOWN205:        {OP_UNKNOWN205, "OP_UNKNOWN205", 1, opcodeInvalid},
	OP_UNKNOWN206:        {OP_UNKNOWN206, "OP_UNKNOWN206", 1, opcodeInvalid},
	OP_UNKNOWN207:        {OP_UNKNOWN207, "OP_UNKNOWN207", 1, opcodeInvalid},
	OP_UNKNOWN208:        {OP_UNKNOWN208, "OP_UNKNOWN208", 1, opcodeInvalid},
	OP_UNKNOWN209:        {OP_UNKNOWN209, "OP_UNKNOWN209", 1, opcodeInvalid},
	OP_UNKNOWN210:        {OP_UNKNOWN210, "OP_UNKNOWN210", 1, opcodeInvalid},
	OP_UNKNOWN211:        {OP_UNKNOWN211, "OP_UNKNOWN211", 1, opcodeInvalid},
	OP_UNKNOWN212:        {OP_UNKNOWN212, "OP_UNKNOWN212", 1, opcodeInvalid},
	OP_UNKNOWN213:        {OP_UNKNOWN213, "OP_UNKNOWN213", 1, opcodeInvalid},
	OP_UNKNOWN214:        {OP_UNKNOWN214, "OP_UNKNOWN214", 1, opcodeInvalid},
	OP_UNKNOWN215:        {OP_UNKNOWN215, "OP_UNKNOWN215", 1, opcodeInvalid},
	OP_UNKNOWN216:        {OP_UNKNOWN216, "OP_UNKNOWN216", 1, opcodeInvalid},
	OP_UNKNOWN217:        {OP_UNKNOWN217, "OP_UNKNOWN217", 1, opcodeInvalid},
	OP_UNKNOWN218:        {OP_UNKNOWN218, "OP_UNKNOWN218", 1, opcodeInvalid},
	OP_UNKNOWN219:        {OP_UNKNOWN219, "OP_UNKNOWN219", 1, opcodeInvalid},
	OP_UNKNOWN220:        {OP_UNKNOWN220, "OP_UNKNOWN220", 1, opcodeInvalid},
	OP_UNKNOWN221:        {OP_UNKNOWN221, "OP_UNKNOWN221", 1, opcodeInvalid},
	OP_UNKNOWN222:        {OP_UNKNOWN222, "OP_UNKNOWN222", 1, opcodeInvalid},
	OP_UNKNOWN223:        {OP_UNKNOWN223, "OP_UNKNOWN223", 1, opcodeInvalid},
	OP_UNKNOWN224:        {OP_UNKNOWN224, "OP_UNKNOWN224", 1, opcodeInvalid},
	OP_UNKNOWN225:        {OP_UNKNOWN225, "OP_UNKNOWN225", 1, opcodeInvalid},
	OP_UNKNOWN226:        {OP_UNKNOWN226, "OP_UNKNOWN226", 1, opcodeInvalid},
	OP_UNKNOWN227:        {OP_UNKNOWN227, "OP_UNKNOWN227", 1, opcodeInvalid},
	OP_UNKNOWN228:        {OP_UNKNOWN228, "OP_UNKNOWN228", 1, opcodeInvalid},
	OP_UNKNOWN229:        {OP_UNKNOWN229, "OP_UNKNOWN229", 1, opcodeInvalid},
	OP_UNKNOWN230:        {OP_UNKNOWN230, "OP_UNKNOWN230", 1, opcodeInvalid},
	OP_UNKNOWN231:        {OP_UNKNOWN231, "OP_UNKNOWN231", 1, opcodeInvalid},
	OP_UNKNOWN232:        {OP_UNKNOWN232, "OP_UNKNOWN232", 1, opcodeInvalid},
	OP_UNKNOWN233:        {OP_UNKNOWN233, "OP_UNKNOWN233", 1, opcodeInvalid},
	OP_UNKNOWN234:        {OP_UNKNOWN234, "OP_UNKNOWN234", 1, opcodeInvalid},
	OP_UNKNOWN235:        {OP_UNKNOWN235, "OP_UNKNOWN235", 1, opcodeInvalid},
	OP_UNKNOWN236:        {OP_UNKNOWN236, "OP_UNKNOWN236", 1, opcodeInvalid},
	OP_UNKNOWN237:        {OP_UNKNOWN237, "OP_UNKNOWN237", 1, opcodeInvalid},
	OP_UNKNOWN238:        {OP_UNKNOWN238, "OP_UNKNOWN238", 1, opcodeInvalid},
	OP_UNKNOWN239:        {OP_UNKNOWN239, "OP_UNKNOWN239", 1, opcodeInvalid},
	OP_UNKNOWN240:        {OP_UNKNOWN240, "OP_UNKNOWN240", 1, opcodeInvalid},
	OP_UNKNOWN241:        {OP_UNKNOWN241, "OP_UNKNOWN241", 1, opcodeInvalid},
	OP_UNKNOWN242:        {OP_UNKNOWN242, "OP_UNKNOWN242", 1, opcodeInvalid},
	OP_UNKNOWN243:        {OP_UNKNOWN243, "OP_UNKNOWN243", 1, opcodeInvalid},
	OP_UNKNOWN244:        {OP_UNKNOWN244, "OP_UNKNOWN244", 1, opcodeInvalid},
	OP_UNKNOWN245:        {OP_UNKNOWN245, "OP_UNKNOWN245", 1, opcodeInvalid},
	OP_UNKNOWN246:        {OP_UNKNOWN246, "OP_UNKNOWN246", 1, opcodeInvalid},
	OP_UNKNOWN247:        {OP_UNKNOWN247, "OP_UNKNOWN247", 1, opcodeInvalid},
	OP_UNKNOWN248:        {OP_UNKNOWN248, "OP_UNKNOWN248", 1, opcodeInvalid},
	OP_UNKNOWN249:        {OP_UNKNOWN249, "OP_UNKNOWN249", 1, opcodeInvalid},

	// Bitcoin Core internal use opcode.  Defined here for completeness.
	OP_SMALLINTEGER: {OP_SMALLINTEGER, "OP_SMALLINTEGER", 1, opcodeInvalid},
	OP_PUBKEYS:      {OP_PUBKEYS, "OP_PUBKEYS", 1, opcodeInvalid},
	OP_UNKNOWN252:   {OP_UNKNOWN252, "OP_UNKNOWN252", 1, opcodeInvalid},
	OP_PUBKEYHASH:   {OP_PUBKEYHASH, "OP_PUBKEYHASH", 1, opcodeInvalid},
	OP_PUBKEY:       {OP_PUBKEY, "OP_PUBKEY", 1, opcodeInvalid},

	OP_INVALIDOPCODE: {OP_INVALIDOPCODE, "OP_INVALIDOPCODE", 1, opcodeInvalid},
}

// opcodeOnelineRepls defines opcode names which are replaced when doing a
// one-line disassembly.  This is done to match the output of the reference
// implementation while not changing the opcode names in the nicer full
// disassembly.
var opcodeOnelineRepls = map[string]string{
	"OP_UNKNOWN79": "-1",
	"OP_0":       "0",
	"OP_1":       "1",
	"OP_2":       "2",
	"OP_3":       "3",
	"OP_4":       "4",
	"OP_5":       "5",
	"OP_6":       "6",
	"OP_7":       "7",
	"OP_8":       "8",
	"OP_9":       "9",
	"OP_10":      "10",
	"OP_11":      "11",
	"OP_12":      "12",
	"OP_13":      "13",
	"OP_14":      "14",
	"OP_15":      "15",
	"OP_16":      "16",
}

// parsedOpcode represents an opcode that has been parsed and includes any
// potential data associated with it.
type parsedOpcode struct {
	opcode *opcode
	data   []byte
}

// checkMinimalDataPush returns whether or not the current data push uses the
// smallest possible opcode to represent it.  For example, the value 15 could
// be pushed with OP_DATA_1 15 (among other variations); however, OP_15 is a
// single opcode that represents the same value and is only a single byte versus
// two bytes.
func (pop *parsedOpcode) checkMinimalDataPush() error {
	data := pop.data
	dataLen := len(data)
	opcode := pop.opcode.value

	if dataLen == 0 && opcode != OP_0 {
		str := fmt.Sprintf("zero length data push is encoded with "+
			"opcode %s instead of OP_0", pop.opcode.name)
		return scriptError(ErrMinimalData, str)
	} else if dataLen == 1 && data[0] >= 1 && data[0] <= 16 {
		if opcode != OP_1+data[0]-1 {
			// Should have used OP_1 .. OP_16
			str := fmt.Sprintf("data push of the value %d encoded "+
				"with opcode %s instead of OP_%d", data[0],
				pop.opcode.name, data[0])
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen == 1 && data[0] == 0x81 {
		if opcode != OP_UNKNOWN79 {
			str := fmt.Sprintf("data push of the value -1 encoded "+
				"with opcode %s instead of OP_UNKNOWN79",
				pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 75 {
		if int(opcode) != dataLen {
			// Should have used a direct push
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OP_DATA_%d", dataLen,
				pop.opcode.name, dataLen)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 255 {
		if opcode != OP_PUSHDATA1 {
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OP_PUSHDATA1",
				dataLen, pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	} else if dataLen <= 65535 {
		if opcode != OP_PUSHDATA2 {
			str := fmt.Sprintf("data push of %d bytes encoded "+
				"with opcode %s instead of OP_PUSHDATA2",
				dataLen, pop.opcode.name)
			return scriptError(ErrMinimalData, str)
		}
	}
	return nil
}

// print returns a human-readable string representation of the opcode for use
// in script disassembly.
func (pop *parsedOpcode) print(oneline bool) string {
	// The reference implementation one-line disassembly replaces opcodes
	// which represent values (e.g. OP_0 through OP_16 and OP_UNKNOWN79)
	// with the raw value.  However, when not doing a one-line dissassembly,
	// we prefer to show the actual opcode names.  Thus, only replace the
	// opcodes in question when the oneline flag is set.
	opcodeName := pop.opcode.name
	if oneline {
		if replName, ok := opcodeOnelineRepls[opcodeName]; ok {
			opcodeName = replName
		}

		// Nothing more to do for non-data push opcodes.
		if pop.opcode.length == 1 {
			return opcodeName
		}

		return fmt.Sprintf("%x", pop.data)
	}

	// Nothing more to do for non-data push opcodes.
	if pop.opcode.length == 1 {
		return opcodeName
	}

	// Add length for the OP_PUSHDATA# opcodes.
	retString := opcodeName
	switch pop.opcode.length {
	case -1:
		retString += fmt.Sprintf(" 0x%02x", len(pop.data))
	case -2:
		retString += fmt.Sprintf(" 0x%04x", len(pop.data))
	case -4:
		retString += fmt.Sprintf(" 0x%08x", len(pop.data))
	}

	return fmt.Sprintf("%s 0x%02x", retString, pop.data)
}

// bytes returns any data associated with the opcode encoded as it would be in
// a script.  This is used for unparsing scripts from parsed opcodes.
func (pop *parsedOpcode) bytes() ([]byte, error) {
	var retbytes []byte
	if pop.opcode.length > 0 {
		retbytes = make([]byte, 1, pop.opcode.length)
	} else {
		retbytes = make([]byte, 1, 1+len(pop.data)-
			pop.opcode.length)
	}

	retbytes[0] = pop.opcode.value
	if pop.opcode.length == 1 {
		if len(pop.data) != 0 {
			str := fmt.Sprintf("internal consistency error - "+
				"parsed opcode %s has data length %d when %d "+
				"was expected", pop.opcode.name, len(pop.data),
				0)
			return nil, scriptError(ErrInternal, str)
		}
		return retbytes, nil
	}
	nbytes := pop.opcode.length
	if pop.opcode.length < 0 {
		l := len(pop.data)
		// tempting just to hardcode to avoid the complexity here.
		switch pop.opcode.length {
		case -1:
			retbytes = append(retbytes, byte(l))
			nbytes = int(retbytes[1]) + len(retbytes)
		case -2:
			retbytes = append(retbytes, byte(l&0xff),
				byte(l>>8&0xff))
			nbytes = int(binary.LittleEndian.Uint16(retbytes[1:])) +
				len(retbytes)
		case -4:
			retbytes = append(retbytes, byte(l&0xff),
				byte((l>>8)&0xff), byte((l>>16)&0xff),
				byte((l>>24)&0xff))
			nbytes = int(binary.LittleEndian.Uint32(retbytes[1:])) +
				len(retbytes)
		}
	}

	retbytes = append(retbytes, pop.data...)

	if len(retbytes) != nbytes {
		str := fmt.Sprintf("internal consistency error - "+
			"parsed opcode %s has data length %d when %d was "+
			"expected", pop.opcode.name, len(retbytes), nbytes)
		return nil, scriptError(ErrInternal, str)
	}

	return retbytes, nil
}

// *******************************************
// Opcode implementation functions start here.
// *******************************************

// opcodeInvalid is a common handler for disabled opcodes.  It returns an
// appropriate error indicating the opcode is disabled.  While it would
// ordinarily make more sense to detect if the script contains any disabled
// opcodes before executing in an initial parse step, the consensus rules
// dictate the script doesn't fail until the program counter passes over a
// disabled opcode (even when they appear in a branch that is not executed).
func opcodeInvalid(op *parsedOpcode, vm *Engine) error {
	str := fmt.Sprintf("attempt to execute disabled opcode %s",
		op.opcode.name)
	return scriptError(ErrDisabledOpcode, str)
}

// opcodeFalse pushes an empty array to the data stack to represent false.  Note
// that 0, when encoded as a number according to the numeric encoding consensus
// rules, is an empty array.
func opcodeFalse(op *parsedOpcode, vm *Engine) error {
	vm.dstack.PushByteArray(nil)
	return nil
}

// opcodePushData is a common handler for the vast majority of opcodes that push
// raw data (bytes) to the data stack.
func opcodePushData(op *parsedOpcode, vm *Engine) error {
	vm.dstack.PushByteArray(op.data)
	return nil
}

// opcodeN is a common handler for the small integer data push opcodes.  It
// pushes the numeric value the opcode represents (which will be from 1 to 16)
// onto the data stack.
func opcodeN(op *parsedOpcode, vm *Engine) error {
	// The opcodes are all defined consecutively, so the numeric value is
	// the difference.
	vm.dstack.PushInt(scriptNum(op.opcode.value - (OP_1 - 1)))
	return nil
}

// abstractVerify examines the top item on the data stack as a boolean value and
// verifies it evaluates to true.  An error is returned either when there is no
// item on the stack or when that item evaluates to false.  In the latter case
// where the verification fails specifically due to the top item evaluating
// to false, the returned error will use the passed error code.
func abstractVerify(op *parsedOpcode, vm *Engine, c ErrorCode) error {
	verified, err := vm.dstack.PopBool()
	if err != nil {
		return err
	}

	if !verified {
		str := fmt.Sprintf("%s failed", op.opcode.name)
		return scriptError(c, str)
	}
	return nil
}

// opcodeVerify examines the top item on the data stack as a boolean value and
// verifies it evaluates to true.  An error is returned if it does not.
func opcodeVerify(op *parsedOpcode, vm *Engine) error {
	return abstractVerify(op, vm, ErrVerify)
}

// opcodeDup duplicates the top item on the data stack.
//
// Stack transformation: [... x1 x2 x3] -> [... x1 x2 x3 x3]
func opcodeDup(op *parsedOpcode, vm *Engine) error {
	return vm.dstack.DupN(1)
}

// opcodeEqual removes the top 2 items of the data stack, compares them as raw
// bytes, and pushes the result, encoded as a boolean, back to the stack.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeEqual(op *parsedOpcode, vm *Engine) error {
	a, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	b, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	vm.dstack.PushBool(bytes.Equal(a, b))
	return nil
}

// opcodeEqualVerify is a combination of opcodeEqual and opcodeVerify.
// Specifically, it removes the top 2 items of the data stack, compares them,
// and pushes the result, encoded as a boolean, back to the stack.  Then, it
// examines the top item on the data stack as a boolean value and verifies it
// evaluates to true.  An error is returned if it does not.
//
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeEqualVerify(op *parsedOpcode, vm *Engine) error {
	err := opcodeEqual(op, vm)
	if err == nil {
		err = abstractVerify(op, vm, ErrEqualVerify)
	}
	return err
}

// opcodeIgnoreEqual removes the top 2 items of the data stack, compares them as raw
// bytes, and pushes the result, encoded as a boolean, back to the stack.
//
// Stack transformation: [... x1 x2] -> [... bool]
func opcodeIgnoreFlagEqual(op *parsedOpcode, vm *Engine) error {
	a, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	if len(a) == 0 {
		return scriptError(ErrInvalidFlags, "address miss flag")
	}
	// ignore the first byte
	a = a[1:]
	b, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	vm.dstack.PushBool(bytes.Equal(a, b))
	return nil
}

// opcodeEqualVerify is a combination of opcodeEqual and opcodeVerify.
// Specifically, it removes the top 2 items of the data stack, compares them,
// and pushes the result, encoded as a boolean, back to the stack.  Then, it
// examines the top item on the data stack as a boolean value and verifies it
// evaluates to true.  An error is returned if it does not.
//
// Stack transformation: [... x1 x2] -> [... bool] -> [...]
func opcodeIgnoreFlagEqualVerify(op *parsedOpcode, vm *Engine) error {
	err := opcodeIgnoreFlagEqual(op, vm)
	if err == nil {
		err = abstractVerify(op, vm, ErrEqualVerify)
	}
	return err
}

// calcHash calculates the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

// opcodeRipemd160 treats the top item of the data stack as raw bytes and
// replaces it with ripemd160(data).
//
// Stack transformation: [... x1] -> [... ripemd160(x1)]
func opcodeRipemd160(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	vm.dstack.PushByteArray(calcHash(buf, ripemd160.New()))
	return nil
}

// opcodeSha1 treats the top item of the data stack as raw bytes and replaces it
// with sha1(data).
//
// Stack transformation: [... x1] -> [... sha1(x1)]
func opcodeSha1(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	hash := sha1.Sum(buf)
	vm.dstack.PushByteArray(hash[:])
	return nil
}

// opcodeSha256 treats the top item of the data stack as raw bytes and replaces
// it with sha256(data).
//
// Stack transformation: [... x1] -> [... sha256(x1)]
func opcodeSha256(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(buf)
	vm.dstack.PushByteArray(hash[:])
	return nil
}

// opcodeHash160 treats the top item of the data stack as raw bytes and replaces
// it with ripemd160(sha256(data)).
//
// Stack transformation: [... x1] -> [... ripemd160(sha256(x1))]
func opcodeHash160(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	hash := sha256.Sum256(buf)
	vm.dstack.PushByteArray(calcHash(hash[:], ripemd160.New()))
	return nil
}

// opcodeHash256 treats the top item of the data stack as raw bytes and replaces
// it with sha256(sha256(data)).
//
// Stack transformation: [... x1] -> [... sha256(sha256(x1))]
func opcodeHash256(op *parsedOpcode, vm *Engine) error {
	buf, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	vm.dstack.PushByteArray(common.DoubleHashB(buf))
	return nil
}

// opcodeCheckSig treats the top 2 items on the stack as a public key and a
// signature and replaces them with a bool which indicates if the signature was
// successfully verified.
//
// The process of verifying a signature requires calculating a signature hash in
// the same way the transaction signer did.  It involves hashing portions of the
// transaction based on the hash type byte (which is the final byte of the
// signature) and the portion of the script starting from the most recent
// OP_UNKNOWN171 (or the beginning of the script if there are none) to the
// end of the script (with any other OP_UNKNOWN171s removed).  Once this
// "script hash" is calculated, the signature is checked using standard
// cryptographic methods against the provided public key.
//
// Stack transformation: [... signature pubkey] -> [... bool]
func opcodeCheckSig(op *parsedOpcode, vm *Engine) error {
	pkBytes, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	fullSigBytes, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	// The signature actually needs needs to be longer than this, but at
	// least 1 byte is needed for the hash type below.  The full length is
	// checked depending on the script flags and upon parsing the signature.
	if len(fullSigBytes) < 1 {
		vm.dstack.PushBool(false)
		return nil
	}

	// Trim off hashtype from the signature string and check if the
	// signature and pubkey conform to the strict encoding requirements
	// depending on the flags.
	//
	// NOTE: When the strict encoding flags are set, any errors in the
	// signature or public encoding here result in an immediate script error
	// (and thus no result bool is pushed to the data stack).  This differs
	// from the logic below where any errors in parsing the signature is
	// treated as the signature failure resulting in false being pushed to
	// the data stack.  This is required because the more general script
	// validation consensus rules do not have the new strict encoding
	// requirements enabled by the flags.
	hashType := SigHashType(fullSigBytes[len(fullSigBytes)-1])
	sigBytes := fullSigBytes[:len(fullSigBytes)-1]
	if err := vm.checkHashTypeEncoding(hashType); err != nil {
		return err
	}
	if err := vm.checkSignatureEncoding(sigBytes); err != nil {
		return err
	}
	if err := vm.checkPubKeyEncoding(pkBytes); err != nil {
		return err
	}

	// Get script starting from the most recent OP_UNKNOWN171.
	subScript := vm.subScript()

	// Generate the signature hash based on the signature hash type.
	var hash []byte
	// Remove the signature since there is no way for a signature
	// to sign itself.
	subScript = removeOpcodeByData(subScript, fullSigBytes)

	hash = calcSignatureHash(subScript, hashType, &vm.tx, vm.txIdx)

	pubKey, err := crypto.ParsePubKey(pkBytes, crypto.S256())
	if err != nil {
		vm.dstack.PushBool(false)
		return nil
	}

	// We take DER Signature as a standard
	signature, err := crypto.ParseDERSignature(sigBytes, crypto.S256())

	if err != nil {
		vm.dstack.PushBool(false)
		return nil
	}

	valid := signature.Verify(hash, pubKey)

	if !valid && vm.hasFlag(ScriptVerifyNullFail) && len(sigBytes) > 0 {
		str := "signature not empty on failed checksig"
		return scriptError(ErrNullFail, str)
	}

	vm.dstack.PushBool(valid)
	return nil
}

// opcodeCheckSigVerify is a combination of opcodeCheckSig and opcodeVerify.
// The opcodeCheckSig function is invoked followed by opcodeVerify.  See the
// documentation for each of those opcodes for more details.
//
// Stack transformation: signature pubkey] -> [... bool] -> [...]
func opcodeCheckSigVerify(op *parsedOpcode, vm *Engine) error {
	err := opcodeCheckSig(op, vm)
	if err == nil {
		err = abstractVerify(op, vm, ErrCheckSigVerify)
	}
	return err
}

// parsedSigInfo houses a raw signature along with its parsed form and a flag
// for whether or not it has already been parsed.  It is used to prevent parsing
// the same signature multiple times when verifying a multisig.
type parsedSigInfo struct {
	signature       []byte
	parsedSignature *crypto.Signature
	parsed          bool
}

func parseMultiPubkeys(vm *Engine) ([][]byte, int, error) {
	numKeys, err := vm.dstack.PopInt()
	if err != nil {
		return nil, 0, err
	}

	numPubKeys := int(numKeys.Int32())
	if numPubKeys < 0 {
		str := fmt.Sprintf("number of pubkeys %d is negative",
			numPubKeys)
		return nil, 0, scriptError(ErrInvalidPubKeyCount, str)
	}
	if numPubKeys > MaxPubKeysPerMultiSig {
		str := fmt.Sprintf("too many pubkeys: %d > %d",
			numPubKeys, MaxPubKeysPerMultiSig)
		return nil, 0, scriptError(ErrInvalidPubKeyCount, str)
	}
	vm.numOps += numPubKeys
	if vm.numOps > MaxOpsPerScript {
		str := fmt.Sprintf("exceeded max operation limit of %d",
			MaxOpsPerScript)
		return nil, 0, scriptError(ErrTooManyOperations, str)
	}

	pubKeys := make([][]byte, 0, numPubKeys)
	for i := 0; i < numPubKeys; i++ {
		pubKey, err := vm.dstack.PopByteArray()
		if err != nil {
			return nil, 0, err
		}

		if err := vm.checkPubKeyEncoding(pubKey); err != nil {
			return nil, 0, err
		}
		pubKeys = append(pubKeys, pubKey)
	}

	numSigs, err := vm.dstack.PopInt()
	if err != nil {
		return nil, 0, err
	}
	numSignatures := int(numSigs.Int32())
	if numSignatures < 0 {
		str := fmt.Sprintf("number of signatures %d is negative",
			numSignatures)
		return nil, 0, scriptError(ErrInvalidSignatureCount, str)

	}
	if numSignatures > numPubKeys {
		str := fmt.Sprintf("more signatures than pubkeys: %d > %d",
			numSignatures, numPubKeys)
		return nil, 0, scriptError(ErrInvalidSignatureCount, str)
	}
	return pubKeys, numSignatures, nil
}

func checkMultiSig(op *parsedOpcode, vm *Engine, numSignatures, numPubKeys int,
	signatures []*parsedSigInfo, pubKeys [][]byte) (bool, error) {
	// Get script starting from the most recent OP_UNKNOWN171.
	script := vm.subScript()

	success := true
	numPubKeys++
	pubKeyIdx := -1
	signatureIdx := 0
	for numSignatures > 0 {
		// When there are more signatures than public keys remaining,
		// there is no way to succeed since too many signatures are
		// invalid, so exit early.
		pubKeyIdx++
		numPubKeys--
		if numSignatures > numPubKeys {
			success = false
			break
		}

		sigInfo := signatures[signatureIdx]
		pubKey := pubKeys[pubKeyIdx]

		// The order of the signature and public key evaluation is
		// important here since it can be distinguished by an
		// OP_CHECKMULTISIG NOT when the strict encoding flag is set.

		rawSig := sigInfo.signature
		if len(rawSig) == 0 {
			// Skip to the next pubkey if signature is empty.
			continue
		}

		// Split the signature into hash type and signature components.
		hashType := SigHashType(rawSig[len(rawSig)-1])
		signature := rawSig[:len(rawSig)-1]

		// Only parse and check the signature encoding once.
		var parsedSig *crypto.Signature
		var err error
		if !sigInfo.parsed {
			if err = vm.checkHashTypeEncoding(hashType); err != nil {
				return false, err
			}
			if err = vm.checkSignatureEncoding(signature); err != nil {
				return false, err
			}

			// Parse the signature.
			// We take DER Signature as a standard
			parsedSig, err = crypto.ParseDERSignature(signature,
				crypto.S256())
			sigInfo.parsed = true
			if err != nil {
				continue
			}
			sigInfo.parsedSignature = parsedSig
		} else {
			// Skip to the next pubkey if the signature is invalid.
			if sigInfo.parsedSignature == nil {
				continue
			}

			// Use the already parsed signature.
			parsedSig = sigInfo.parsedSignature
		}

		// Parse the pubkey.
		parsedPubKey, err := crypto.ParsePubKey(pubKey, crypto.S256())
		if err != nil {
			continue
		}

		// Generate the signature hash based on the signature hash type.
		var sHash = calcSignatureHash(script, hashType, &vm.tx, vm.txIdx)
		valid := parsedSig.Verify(sHash, parsedPubKey)
		if valid {
			// PubKey verified, move on to the next signature.
			signatureIdx++
			numSignatures--
		}
	}
	return success, nil
}

// opcodeCheckMultiSig treats the top item on the stack as an integer number of
// public keys, followed by that many entries as raw data representing the public
// keys, followed by the integer number of signatures, followed by that many
// entries as raw data representing the signatures.
//
// All of the aforementioned stack items are replaced with a bool which
// indicates if the requisite number of signatures were successfully verified.
//
// See the opcodeCheckSigVerify documentation for more details about the process
// for verifying each signature.
//
// Stack transformation:
// [... [sig ...] numsigs [pubkey ...] numpubkeys] -> [... bool]
func opcodeCheckMultiSig(op *parsedOpcode, vm *Engine) error {
	pubKeys, numSignatures, err := parseMultiPubkeys(vm)
	if err != nil {
		return err
	}
	numPubKeys := len(pubKeys)

	signatures := make([]*parsedSigInfo, 0, numSignatures)
	for i := 0; i < numSignatures; i++ {
		signature, err := vm.dstack.PopByteArray()
		if err != nil {
			return err
		}
		sigInfo := &parsedSigInfo{signature: signature}
		signatures = append(signatures, sigInfo)
	}

	success, err := checkMultiSig(op, vm, numSignatures, numPubKeys, signatures, pubKeys)
	if err != nil {
		return err
	}

	if !success && vm.hasFlag(ScriptVerifyNullFail) {
		for _, sig := range signatures {
			if len(sig.signature) > 0 {
				str := "not all signatures empty on failed checkmultisig"
				return scriptError(ErrNullFail, str)
			}
		}
	}

	vm.dstack.PushBool(success)
	return nil
}

// opcodeCheckMultiSigVerify is a combination of opcodeCheckMultiSig and
// opcodeVerify.  The opcodeCheckMultiSig is invoked followed by opcodeVerify.
// See the documentation for each of those opcodes for more details.
//
// Stack transformation:
// [... [sig ...] numsigs [pubkey ...] numpubkeys] -> [... bool] -> [...]
func opcodeCheckMultiSigVerify(op *parsedOpcode, vm *Engine) error {
	err := opcodeCheckMultiSig(op, vm)
	if err == nil {
		err = abstractVerify(op, vm, ErrCheckMultiSigVerify)
	}
	return err
}

// OpcodeByName is a map that can be used to lookup an opcode by its
// human-readable name (OP_CHECKMULTISIG, OP_CHECKSIG, etc).
var OpcodeByName = make(map[string]byte)

func init() {
	// Initialize the opcode name to value map using the contents of the
	// opcode array.  Also add entries for "OP_FALSE", and "OP_TRUE"
	// since they are aliases for "OP_0", and "OP_1" respectively.
	for _, op := range opcodeArray {
		OpcodeByName[op.name] = op.value
	}
	OpcodeByName["OP_FALSE"] = OP_FALSE
	OpcodeByName["OP_TRUE"] = OP_TRUE
}
