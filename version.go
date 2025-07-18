package goua

/*
#include "open62541.h"
*/
import "C"

const (
    VersionMajor = C.UA_OPEN62541_VER_MAJOR
    VersionMinor = C.UA_OPEN62541_VER_MINOR
    VersionPatch = C.UA_OPEN62541_VER_PATCH
    VersionLabel = C.UA_OPEN62541_VER_LABEL
    VersionCommit = C.UA_OPEN62541_VER_COMMIT
    Version = C.UA_OPEN62541_VERSION
)
