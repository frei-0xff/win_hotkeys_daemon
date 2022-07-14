// Base on: https://github.com/vgo0/gologger/blob/master/winapi/winapi.go
package main

import (
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	setWindowsHookEx    = user32.NewProc("SetWindowsHookExA")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	unhookWinEvent      = user32.NewProc("UnhookWinEvent")
	getMessage          = user32.NewProc("GetMessageW")
	toAscii             = user32.NewProc("ToAscii")
	getKeyboardState    = user32.NewProc("GetKeyboardState")
	attachThreadInput   = user32.NewProc("AttachThreadInput")
	setWinEventHook     = user32.NewProc("SetWinEventHook")
	getWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	getWindowText       = user32.NewProc("GetWindowTextW")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	translateMessage    = user32.NewProc("TranslateMessage")
	dispatchMessage     = user32.NewProc("DispatchMessage")
	keybdEvent          = user32.NewProc("keybd_event")
	getAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
	postMessage         = user32.NewProc("PostMessageA")
	openClipboard       = user32.NewProc("OpenClipboard")
	getClipboardData    = user32.NewProc("GetClipboardData")
	closeClipboard      = user32.NewProc("CloseClipboard")
	setCursorPos        = user32.NewProc("SetCursorPos")

	kernel32           = windows.NewLazySystemDLL("kernel32.dll")
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	getThreadId        = kernel32.NewProc("GetThreadId")
	winExec            = kernel32.NewProc("WinExec")
	globalLock         = kernel32.NewProc("GlobalLock")
	globalUnlock       = kernel32.NewProc("GlobalUnlock")
)

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
HWND GetForegroundWindow();
Get currently active window if needed
*/
func GetForegroundWindow() HWND {
	ret, _, _ := getForegroundWindow.Call()

	return HWND(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextlengthw
int GetWindowTextLengthW(
  [in] HWND hWnd
);
Get length of window title to fetch via GetWindowText
*/
func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := getWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextw
int GetWindowTextW(
  [in]  HWND   hWnd,
  [out] LPWSTR lpString,
  [in]  int    nMaxCount
);

Get window title (calls GetWindowTextLength and converts to string for you)
*/
func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	getWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	title := syscall.UTF16ToString(buf)
	if title == "" {
		title = "Unknown"
	}
	return title
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getthreadid
DWORD GetThreadId(
  [in] HANDLE Thread
);
*/
func GetThreadId(Thread HANDLE) DWORD {
	ret, _, _ := getThreadId.Call(uintptr(Thread))

	return DWORD(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwineventhook
HWINEVENTHOOK SetWinEventHook(
  [in] DWORD        eventMin,
  [in] DWORD        eventMax,
  [in] HMODULE      hmodWinEventProc,
  [in] WINEVENTPROC pfnWinEventProc,
  [in] DWORD        idProcess,
  [in] DWORD        idThread,
  [in] DWORD        dwFlags
);

Used by us to detect changes in selected / foreground window via EVENT_OBJECT_FOCUS
*/
func SetWinEventHook(eventMin DWORD, eventMax DWORD, hmodWinEventProc HMODULE, pfnWinEventProc WINEVENTPROC, idProcess DWORD, idThread DWORD, dwFlags DWORD) HWINEVENTHOOK {
	ret, _, _ := setWinEventHook.Call(
		uintptr(eventMin),
		uintptr(eventMax),
		uintptr(hmodWinEventProc),
		uintptr(syscall.NewCallback(pfnWinEventProc)),
		uintptr(idProcess),
		uintptr(idThread),
		uintptr(dwFlags),
	)

	return HWINEVENTHOOK(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-attachthreadinput
BOOL AttachThreadInput(
  [in] DWORD idAttach,
  [in] DWORD idAttachTo,
  [in] BOOL  fAttach
);

Used to allow our thread receiving low level keyboard events to get an accurate keyboard state for use in ToAscii
*/
func AttachThreadInput(idAttach DWORD, idAttachTo DWORD, fAttach BOOL) BOOL {
	ret, _, _ := attachThreadInput.Call(
		uintptr(idAttach),
		uintptr(idAttachTo),
		uintptr(fAttach),
	)

	return BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentthreadid
DWORD GetCurrentThreadId();
*/
func GetCurrentThreadId() DWORD {
	ret, _, _ := getCurrentThreadId.Call()

	return DWORD(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-toascii
int ToAscii(
  [in]           UINT       uVirtKey,
  [in]           UINT       uScanCode,
  [in, optional] const BYTE *lpKeyState,
  [out]          LPWORD     lpChar,
  [in]           UINT       uFlags
);

This is used to take a current keyboard state and the VkCode from a low level keyboard event and turn it into the real text
This approach takes care of things like taking into account if caps lock or shift keys are active
*/
func ToAscii(uVirtKey DWORD, uScanCode DWORD, lpKeyState *[256]byte, lpChar *uint16, uFlags DWORD) int {
	ret, _, _ := toAscii.Call(
		uintptr(uVirtKey),
		uintptr(uScanCode),
		uintptr(unsafe.Pointer(&(*lpKeyState)[0])),
		uintptr(unsafe.Pointer(lpChar)),
		uintptr(uFlags),
	)

	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardstate
BOOL GetKeyboardState(
  [out] PBYTE lpKeyState
);
*/
func GetKeyboardState(lpKeyState *[256]byte) BOOL {
	ret, _, _ := getKeyboardState.Call(
		uintptr(unsafe.Pointer(&(*lpKeyState)[0])),
	)
	return BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexa
HHOOK SetWindowsHookExA(
  [in] int       idHook,
  [in] HOOKPROC  lpfn,
  [in] HINSTANCE hmod,
  [in] DWORD     dwThreadId
);
*/
func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := setWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
BOOL UnhookWindowsHookEx(
  [in] HHOOK hhk
);
*/
func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := unhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwinevent
BOOL UnhookWinEvent(
  [in] HWINEVENTHOOK hWinEventHook
);
*/
func UnhookWinEvent(hWinEventHook HWINEVENTHOOK) BOOL {
	ret, _, _ := unhookWinEvent.Call(uintptr(hWinEventHook))

	return BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callnexthookex
LRESULT CallNextHookEx(
  [in, optional] HHOOK  hhk,
  [in]           int    nCode,
  [in]           WPARAM wParam,
  [in]           LPARAM lParam
);
*/
func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := callNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-dispatchmessage
LRESULT DispatchMessage(
  [in] const MSG *lpMsg
);
*/
func DispatchMessage(msg *MSG) LRESULT {
	ret, _, _ := dispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)),
	)

	return LRESULT(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-translatemessage
BOOL TranslateMessage(
  [in] const MSG *lpMsg
);
*/
func TranslateMessage(msg *MSG) BOOL {
	ret, _, _ := translateMessage.Call(
		uintptr(unsafe.Pointer(msg)),
	)

	return BOOL(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessage
BOOL GetMessage(
  [out]          LPMSG lpMsg,
  [in, optional] HWND  hWnd,
  [in]           UINT  wMsgFilterMin,
  [in]           UINT  wMsgFilterMax
);
*/
func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := getMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-keybd_event
void keybd_event(
  [in] BYTE      bVk,
  [in] BYTE      bScan,
  [in] DWORD     dwFlags,
  [in] ULONG_PTR dwExtraInfo
);
*/
func KeybdEvent(bVk, bScan BYTE, dwFlags DWORD, dwExtraInfo DWORD) {
	keybdEvent.Call(
		uintptr(bVk),
		uintptr(bScan),
		uintptr(dwFlags),
		uintptr(dwExtraInfo),
	)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeystate
SHORT GetAsyncKeyState(
  [in] int nVirtKey
);
*/
func GetAsyncKeyState(nVirtKey int) int {
	ret, _, _ := getAsyncKeyState.Call(
		uintptr(nVirtKey))
	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-winexec
UINT WinExec(
  [in] LPCSTR lpCmdLine,
  [in] UINT   uCmdShow
);
*/
func WinExec(lpCmdLine string, uCmdShow DWORD) int {
	buf := []byte(lpCmdLine)
	ret, _, _ := winExec.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(uCmdShow))
	return int(ret)
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postmessagea
BOOL PostMessageA(
  [in, optional] HWND   hWnd,
  [in]           UINT   Msg,
  [in]           WPARAM wParam,
  [in]           LPARAM lParam
);
*/
func PostMessage(hwnd HWND, Msg DWORD, wParam WPARAM, lParam LPARAM) BOOL {
	ret, _, _ := postMessage.Call(
		uintptr(hwnd),
		uintptr(Msg),
		uintptr(wParam),
		uintptr(lParam),
	)
	return BOOL(ret)
}

/*
https://docs.microsoft.com/ru-ru/windows/win32/api/winuser/nf-winuser-getclipboarddata?redirectedfrom=MSDN
HANDLE GetClipboardData(
  [in] UINT uFormat
);
*/
func GetClipboardData() (text string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return
	}
	h, _, _ := getClipboardData.Call(CF_TEXT)
	if h == 0 {
		closeClipboard.Call()
		return
	}
	l, _, _ := globalLock.Call(h)
	if l == 0 {
		closeClipboard.Call()
		return
	}
	s := (*[1 << 20]byte)(unsafe.Pointer(l))[:]
	for i, v := range s {
		if v == 0 {
			text = string(s[0:i])
			break
		}
	}
	l, _, _ = globalUnlock.Call(h)
	if l == 0 {
		closeClipboard.Call()
		return
	}
	ret, _, _ = closeClipboard.Call()
	if ret == 0 {
		return
	}
	return
}

/*
https://docs.microsoft.com/ru-ru/windows/win32/api/winuser/nf-winuser-setcursorpos
BOOL SetCursorPos(
  [in] int X,
  [in] int Y
);
*/
func SetCursorPos(X, Y int) BOOL {
	ret, _, _ := setCursorPos.Call(
		uintptr(X),
		uintptr(Y),
	)
	return BOOL(ret)
}
