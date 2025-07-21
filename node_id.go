package goua

/*
#include "open62541.h"
*/
import "C"
import "unsafe"

/*
 * Usage:
 *
 *      var cNodeID *C.UA_NodeId
 *      var n *NodeID
 *      n.Parse(cNodeID)
 *      switch n.IdentifierType {
 *      case NodeIDType_Numeric:
 *          fmt.Printf("Numeric NodeID: %d\n", n.NodeIDNumeric())
 *      case NodeIDType_String:
 *          fmt.Printf("String NodeID: %s\n", n.NodeID.(String))
 *      case NodeIDType_Guid:
 *          fmt.Printf("GUID NodeID: %s\n", n.NodeIDGUID())
 *      case NodeIDType_ByteString:
 *          fmt.Printf("ByteString NodeID: %s\n", n.NodeIDByteString())
 *      }
 *
 *
 */

type NodeID struct {
	NamespaceIndex UInt16
	IdentifierType NodeIDType
	nodeID         any // oneof: UInt32, String, Guid, ByteString
}

func (n *NodeID) Parse(cNodeID *C.UA_NodeId) {
	n.NamespaceIndex = UInt16(cNodeID.namespaceIndex)
	n.IdentifierType = NodeIDType(cNodeID.identifierType)

	pIdent := unsafe.Pointer(&cNodeID.identifier[0])

	switch n.IdentifierType {
	case NodeIDType_Numeric:
		n.nodeID = *(*UInt32)(pIdent)
	case NodeIDType_String:
		n.nodeID = *(*String)(pIdent)
	case NodeIDType_Guid:
		n.nodeID = *(*GUID)(pIdent)
	case NodeIDType_ByteString:
		n.nodeID = *(*ByteString)(pIdent)
	}
}

func (n *NodeID) NodeIDNumeric() UInt32 {
	return n.nodeID.(UInt32)
}

func (n *NodeID) NodeIDString() String {
	return n.nodeID.(String)
}

func (n *NodeID) NodeIDGuid() GUID {
	return n.nodeID.(GUID)
}

func (n *NodeID) NodeIDByteString() ByteString {
	return n.nodeID.(ByteString)
}
