// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"syscall"
	"unsafe"
)

var (
	modntdll = syscall.NewLazyDLL("ntdll.dll")

	procAlpcGetMessageAttribute          = modntdll.NewProc("AlpcGetMessageAttribute")
	procNtAlpcAcceptConnectPort          = modntdll.NewProc("NtAlpcAcceptConnectPort")
	procNtAlpcCancelMessage              = modntdll.NewProc("NtAlpcCancelMessage")
	procNtAlpcConnectPort                = modntdll.NewProc("NtAlpcConnectPort")
	procNtAlpcCreatePort                 = modntdll.NewProc("NtAlpcCreatePort")
	procNtAlpcDisconnectPort             = modntdll.NewProc("NtAlpcDisconnectPort")
	procNtAlpcSendWaitReceivePort        = modntdll.NewProc("NtAlpcSendWaitReceivePort")
	procRtlCreateUnicodeStringFromAsciiz = modntdll.NewProc("RtlCreateUnicodeStringFromAsciiz")
)

// NTSTATUS
// NtAlpcCreatePort(
//   __out PHANDLE PortHandle,
//   __in POBJECT_ATTRIBUTES ObjectAttributes,
//   __in_opt PALPC_PORT_ATTRIBUTES PortAttributes
//   );
func NtAlpcCreatePort(pObjectAttributes *OBJECT_ATTRIBUTES, pPortAttributes *ALPC_PORT_ATTRIBUTES) (hPort HANDLE, e error) {

	ret, _, _ := procNtAlpcCreatePort.Call(
		uintptr(unsafe.Pointer(&hPort)),
		uintptr(unsafe.Pointer(pObjectAttributes)),
		uintptr(unsafe.Pointer(pPortAttributes)),
	)

	if ret != ERROR_SUCCESS {
		return hPort, fmt.Errorf("0x%x", ret)
	}

	return
}

// NTSTATUS
// NtAlpcAcceptConnectPort(
//     __out PHANDLE PortHandle,
//     __in HANDLE ConnectionPortHandle,
//     __in ULONG Flags,
//     __in POBJECT_ATTRIBUTES ObjectAttributes,
//     __in PALPC_PORT_ATTRIBUTES PortAttributes,
//     __in_opt PVOID PortContext,
//     __in PPORT_MESSAGE ConnectionRequest,
//     __inout_opt PALPC_MESSAGE_ATTRIBUTES ConnectionMessageAttributes,
//     __in BOOLEAN AcceptConnection
//     );
func NtAlpcAcceptConnectPort(
	hSrvConnPort HANDLE,
	flags uint32,
	pObjAttr *OBJECT_ATTRIBUTES,
	pPortAttr *ALPC_PORT_ATTRIBUTES,
	pContext *AlpcPortContext,
	pConnReq *AlpcShortMessage,
	pConnMsgAttrs *ALPC_MESSAGE_ATTRIBUTES,
	accept uintptr,
) (hPort HANDLE, e error) {

	ret, _, _ := procNtAlpcAcceptConnectPort.Call(
		uintptr(unsafe.Pointer(&hPort)),
		uintptr(hSrvConnPort),
		uintptr(flags),
		uintptr(unsafe.Pointer(pObjAttr)),
		uintptr(unsafe.Pointer(pPortAttr)),
		uintptr(unsafe.Pointer(pContext)),
		uintptr(unsafe.Pointer(pConnReq)),
		uintptr(unsafe.Pointer(pConnMsgAttrs)),
		accept,
	)

	if ret != ERROR_SUCCESS {
		e = fmt.Errorf("0x%x", ret)
	}
	return
}

// NTSTATUS
// NtAlpcSendWaitReceivePort(
//     __in HANDLE PortHandle,
//     __in ULONG Flags,
//     __in_opt PPORT_MESSAGE SendMessage,
//     __in_opt PALPC_MESSAGE_ATTRIBUTES SendMessageAttributes,
//     __inout_opt PPORT_MESSAGE ReceiveMessage,
//     __inout_opt PULONG BufferLength,
//     __inout_opt PALPC_MESSAGE_ATTRIBUTES ReceiveMessageAttributes,
//     __in_opt PLARGE_INTEGER Timeout
//     );
func NtAlpcSendWaitReceivePort(
	hPort HANDLE,
	flags uint32,
	sendMsg *AlpcShortMessage, // Should actually point to PORT_MESSAGE + payload
	sendMsgAttrs *ALPC_MESSAGE_ATTRIBUTES,
	recvMsg *AlpcShortMessage,
	recvBufLen *uint32,
	recvMsgAttrs *ALPC_MESSAGE_ATTRIBUTES,
	timeout *int64, // use native int64
) (e error) {

	ret, _, _ := procNtAlpcSendWaitReceivePort.Call(
		uintptr(hPort),
		uintptr(flags),
		uintptr(unsafe.Pointer(sendMsg)),
		uintptr(unsafe.Pointer(sendMsgAttrs)),
		uintptr(unsafe.Pointer(recvMsg)),
		uintptr(unsafe.Pointer(recvBufLen)),
		uintptr(unsafe.Pointer(recvMsgAttrs)),
		uintptr(unsafe.Pointer(timeout)),
	)

	if ret != ERROR_SUCCESS {
		e = fmt.Errorf("0x%x", ret)
	}
	return
}

// NTSYSAPI
// PVOID
// NTAPI
// AlpcGetMessageAttribute(
//     __in PALPC_MESSAGE_ATTRIBUTES Buffer,
//     __in ULONG AttributeFlag
//     );

// This basically returns a pointer to the correct struct for whichever
// message attribute you asked for. In Go terms, it returns unsafe.Pointer
// which you should then cast. Example:

// ptr := AlpcGetMessageAttribute(&recvMsgAttrs, ALPC_MESSAGE_CONTEXT_ATTRIBUTE)
// if ptr != nil {
//     context := (*ALPC_CONTEXT_ATTR)(ptr)
// }
func AlpcGetMessageAttribute(buf *ALPC_MESSAGE_ATTRIBUTES, attr uint32) unsafe.Pointer {

	ret, _, _ := procAlpcGetMessageAttribute.Call(
		uintptr(unsafe.Pointer(buf)),
		uintptr(attr),
	)
	return unsafe.Pointer(ret)
}

// NTSYSCALLAPI
// NTSTATUS
// NTAPI
// NtAlpcCancelMessage(
//     __in HANDLE PortHandle,
//     __in ULONG Flags,
//     __in PALPC_CONTEXT_ATTR MessageContext
//     );
func NtAlpcCancelMessage(hPort HANDLE, flags uint32, pMsgContext *ALPC_CONTEXT_ATTR) (e error) {

	ret, _, _ := procNtAlpcCancelMessage.Call(
		uintptr(hPort),
		uintptr(flags),
		uintptr(unsafe.Pointer(pMsgContext)),
	)

	if ret != ERROR_SUCCESS {
		e = fmt.Errorf("0x%x", ret)
	}
	return
}

// NTSYSCALLAPI
// NTSTATUS
// NTAPI
// NtAlpcDisconnectPort(
//     __in HANDLE PortHandle,
//     __in ULONG Flags
//     );
func NtAlpcDisconnectPort(hPort HANDLE, flags uint32) (e error) {

	ret, _, _ := procNtAlpcDisconnectPort.Call(
		uintptr(hPort),
		uintptr(flags),
	)

	if ret != ERROR_SUCCESS {
		e = fmt.Errorf("0x%x", ret)
	}
	return
}
