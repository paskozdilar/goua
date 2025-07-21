package goua

/*
#include "open62541.h"
*/
import "C"

type NodeIDType C.int

const (
	NodeIDType_Numeric    NodeIDType = C.UA_NODEIDTYPE_NUMERIC
	NodeIDType_String     NodeIDType = C.UA_NODEIDTYPE_STRING
	NodeIDType_Guid       NodeIDType = C.UA_NODEIDTYPE_GUID
	NodeIDType_ByteString NodeIDType = C.UA_NODEIDTYPE_BYTESTRING
)
