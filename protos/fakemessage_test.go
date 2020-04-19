// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import "io"

// fakeMessage implements the Message interface and is used to force encode
// errors in messages.
type fakeMessage struct {
	command        string
	payload        []byte
	forceEncodeErr bool
	forceLenErr    bool
}

// VVSDecode doesn't do anything.  It just satisfies the protos.Message
// interface.
func (msg *fakeMessage) VVSDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	return nil
}

// VVSEncode writes the payload field of the fake message or forces an error
// if the forceEncodeErr flag of the fake message is set.  It also satisfies the
// protos.Message interface.
func (msg *fakeMessage) VVSEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	if msg.forceEncodeErr {
		err := &MessageError{
			Func:        "fakeMessage.VVSEncode",
			Description: "intentional error",
		}
		return err
	}

	_, err := w.Write(msg.payload)
	return err
}

// Command returns the command field of the fake message and satisfies the
// Message interface.
func (msg *fakeMessage) Command() string {
	return msg.command
}

// MaxPayloadLength returns the length of the payload field of fake message
// or a smaller value if the forceLenErr flag of the fake message is set.  It
// satisfies the Message interface.
func (msg *fakeMessage) MaxPayloadLength(pver uint32) uint32 {
	lenp := uint32(len(msg.payload))
	if msg.forceLenErr {
		return lenp - 1
	}

	return lenp
}
