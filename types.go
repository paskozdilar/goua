package goua

/*
#include "open62541.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Builtin types:

type (
	Boolean C.UA_Boolean
	SByte   C.UA_SByte
	Byte    C.UA_Byte
	Int16   C.UA_Int16
	UInt16  C.UA_UInt16
	Int32   C.UA_Int32
	UInt32  C.UA_UInt32
	Int64   C.UA_Int64
	UInt64  C.UA_UInt64
	Float   C.UA_Float
	Double  C.UA_Double
)

type String C.UA_String

func (s *String) String() string {
	if s == nil || s.data == nil {
		return ""
	}
	return C.GoStringN((*C.char)(unsafe.Pointer(s.data)), C.int(s.length))
}

type ByteString = String

// TODO: DateTime

type GUID C.UA_Guid

func (g *GUID) String() string {
	if g == nil {
		return ""
	}
	return C.GoStringN((*C.char)(unsafe.Pointer(&g.data1)), 16)
}

func (g *GUID) Parse(s String) error {
	if s.length != 16 {
		return fmt.Errorf("GUID must be 16 bytes long, got %d bytes", s.length)
	}
	var guid C.UA_Guid
	retval := StatusCode(C.UA_Guid_parse(&guid, C.UA_String(s)))
	if retval != StatusCode_Good {
		return fmt.Errorf("failed to parse GUID: %s", retval)
	}
	return nil
}
