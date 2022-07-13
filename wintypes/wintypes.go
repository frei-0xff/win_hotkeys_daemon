// Based on: https://github.com/vgo0/gologger/blob/master/wintypes/wintypes.go
package wintypes

const (
	//https://docs.microsoft.com/en-us/windows/win32/winauto/event-constants
	EVENT_OBJECT_FOCUS DWORD = 0x8005
	//https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-keydown
	WM_KEYDOWN DWORD = 0x100
	//https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-keyup
	WM_KEYUP DWORD = 0x101
	//https://docs.microsoft.com/en-us/windows/win32/inputdev/wm-syskeydown
	WM_SYSKEYDOWN DWORD = 0x104

	WINEVENT_OUTOFCONTEXT   DWORD = 4
	WINEVENT_SKIPOWNPROCESS DWORD = 2

	//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexa
	WH_KEYBOARD_LL = 13
)

// https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
type (
	//typedef int BOOL;
	BOOL int32
	//typedef unsigned char BYTE;
	BYTE byte
	//typedef unsigned long DWORD;
	DWORD uint32
	//typedef PVOID HANDLE;
	HANDLE uintptr
	//typedef HANDLE HHOOK;
	HHOOK HANDLE
	//typedef HANDLE HINSTANCE;
	HINSTANCE HANDLE
	//typedef HINSTANCE HMODULE;
	HMODULE HANDLE
	//typedef HANDLE HWND;
	HWND HANDLE
	//typedef long LONG;
	LONG int32
	/*
		#if defined(_WIN64)
		typedef __int64 LONG_PTR;
		#else
		typedef long LONG_PTR;
		#endif
	*/
	LONG_PTR uintptr
	//typedef LONG_PTR LPARAM;
	LPARAM LONG_PTR
	//typedef LONG_PTR LRESULT;
	LRESULT LONG_PTR
	//typedef UINT_PTR WPARAM;
	WPARAM uintptr
	//https://docs.microsoft.com/en-us/windows/win32/winauto/hwineventhook
	//typedef HANDLE HWINEVENTHOOK;
	HWINEVENTHOOK HANDLE
	//typedef BYTE *PBYTE;
	PBYTE []BYTE
	/*
		https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-hookproc
		LRESULT Hookproc(
			int code,
			[in] WPARAM wParam,
			[in] LPARAM lParam
		)
	*/
	HOOKPROC func(int, WPARAM, LPARAM) LRESULT
	/*
		https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-wineventproc
		void Wineventproc(
			HWINEVENTHOOK hWinEventHook,
			DWORD event,
			HWND hwnd,
			LONG idObject,
			LONG idChild,
			DWORD idEventThread,
			DWORD dwmsEventTime
		)
	*/
	WINEVENTPROC func(HWINEVENTHOOK, DWORD, HWND, LONG, LONG, DWORD, DWORD) uintptr
)

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
typedef struct tagMSG {
  HWND   hwnd;
  UINT   message;
  WPARAM wParam;
  LPARAM lParam;
  DWORD  time;
  POINT  pt;
  DWORD  lPrivate;
} MSG, *PMSG, *NPMSG, *LPMSG;
*/
type MSG struct {
	Hwnd     HWND
	Message  uint32
	WParam   WPARAM
	LParam   LPARAM
	Time     DWORD
	Pt       POINT
	LPrivate DWORD
}

/*
https://docs.microsoft.com/en-us/previous-versions/dd162805(v=vs.85)
typedef struct tagPOINT {
  LONG x;
  LONG y;
} POINT, *PPOINT;
*/
type POINT struct {
	X, Y LONG
}

/*
https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
typedef struct tagKBDLLHOOKSTRUCT {
  DWORD     vkCode;
  DWORD     scanCode;
  DWORD     flags;
  DWORD     time;
  ULONG_PTR dwExtraInfo;
} KBDLLHOOKSTRUCT, *LPKBDLLHOOKSTRUCT, *PKBDLLHOOKSTRUCT;
*/
type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-keybd_event
const (
	KEYEVENTF_EXTENDEDKEY DWORD = 0x01
	KEYEVENTF_KEYUP       DWORD = 0x02
)

//https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
const (
	VK_LBUTTON             DWORD = 0x01
	VK_RBUTTON             DWORD = 0x02
	VK_CANCEL              DWORD = 0x03
	VK_MBUTTON             DWORD = 0x04
	VK_XBUTTON1            DWORD = 0x05
	VK_XBUTTON2            DWORD = 0x06
	VK_BACK                DWORD = 0x08
	VK_TAB                 DWORD = 0x09
	VK_CLEAR               DWORD = 0x0C
	VK_RETURN              DWORD = 0x0D
	VK_SHIFT               DWORD = 0x10
	VK_CONTROL             DWORD = 0x11
	VK_MENU                DWORD = 0x12
	VK_PAUSE               DWORD = 0x13
	VK_CAPITAL             DWORD = 0x14
	VK_KANA                DWORD = 0x15
	VK_HANGUEL             DWORD = 0x15
	VK_HANGUL              DWORD = 0x15
	VK_IME_ON              DWORD = 0x16
	VK_JUNJA               DWORD = 0x17
	VK_FINAL               DWORD = 0x18
	VK_HANJA               DWORD = 0x19
	VK_KANJI               DWORD = 0x19
	VK_IME_OFF             DWORD = 0x1A
	VK_ESCAPE              DWORD = 0x1B
	VK_CONVERT             DWORD = 0x1C
	VK_NONCONVERT          DWORD = 0x1D
	VK_ACCEPT              DWORD = 0x1E
	VK_MODECHANGE          DWORD = 0x1F
	VK_SPACE               DWORD = 0x20
	VK_PRIOR               DWORD = 0x21
	VK_NEXT                DWORD = 0x22
	VK_END                 DWORD = 0x23
	VK_HOME                DWORD = 0x24
	VK_LEFT                DWORD = 0x25
	VK_UP                  DWORD = 0x26
	VK_RIGHT               DWORD = 0x27
	VK_DOWN                DWORD = 0x28
	VK_SELECT              DWORD = 0x29
	VK_PRINT               DWORD = 0x2A
	VK_EXECUTE             DWORD = 0x2B
	VK_SNAPSHOT            DWORD = 0x2C
	VK_INSERT              DWORD = 0x2D
	VK_DELETE              DWORD = 0x2E
	VK_HELP                DWORD = 0x2F
	VK_LWIN                DWORD = 0x5B
	VK_RWIN                DWORD = 0x5C
	VK_APPS                DWORD = 0x5D
	VK_SLEEP               DWORD = 0x5F
	VK_NUMPAD0             DWORD = 0x60
	VK_NUMPAD1             DWORD = 0x61
	VK_NUMPAD2             DWORD = 0x62
	VK_NUMPAD3             DWORD = 0x63
	VK_NUMPAD4             DWORD = 0x64
	VK_NUMPAD5             DWORD = 0x65
	VK_NUMPAD6             DWORD = 0x66
	VK_NUMPAD7             DWORD = 0x67
	VK_NUMPAD8             DWORD = 0x68
	VK_NUMPAD9             DWORD = 0x69
	VK_MULTIPLY            DWORD = 0x6A
	VK_ADD                 DWORD = 0x6B
	VK_SEPARATOR           DWORD = 0x6C
	VK_SUBTRACT            DWORD = 0x6D
	VK_DECIMAL             DWORD = 0x6E
	VK_DIVIDE              DWORD = 0x6F
	VK_F1                  DWORD = 0x70
	VK_F2                  DWORD = 0x71
	VK_F3                  DWORD = 0x72
	VK_F4                  DWORD = 0x73
	VK_F5                  DWORD = 0x74
	VK_F6                  DWORD = 0x75
	VK_F7                  DWORD = 0x76
	VK_F8                  DWORD = 0x77
	VK_F9                  DWORD = 0x78
	VK_F10                 DWORD = 0x79
	VK_F11                 DWORD = 0x7A
	VK_F12                 DWORD = 0x7B
	VK_F13                 DWORD = 0x7C
	VK_F14                 DWORD = 0x7D
	VK_F15                 DWORD = 0x7E
	VK_F16                 DWORD = 0x7F
	VK_F17                 DWORD = 0x80
	VK_F18                 DWORD = 0x81
	VK_F19                 DWORD = 0x82
	VK_F20                 DWORD = 0x83
	VK_F21                 DWORD = 0x84
	VK_F22                 DWORD = 0x85
	VK_F23                 DWORD = 0x86
	VK_F24                 DWORD = 0x87
	VK_NUMLOCK             DWORD = 0x90
	VK_SCROLL              DWORD = 0x91
	VK_LSHIFT              DWORD = 0xA0
	VK_RSHIFT              DWORD = 0xA1
	VK_LCONTROL            DWORD = 0xA2
	VK_RCONTROL            DWORD = 0xA3
	VK_LMENU               DWORD = 0xA4
	VK_RMENU               DWORD = 0xA5
	VK_BROWSER_BACK        DWORD = 0xA6
	VK_BROWSER_FORWARD     DWORD = 0xA7
	VK_BROWSER_REFRESH     DWORD = 0xA8
	VK_BROWSER_STOP        DWORD = 0xA9
	VK_BROWSER_SEARCH      DWORD = 0xAA
	VK_BROWSER_FAVORITES   DWORD = 0xAB
	VK_BROWSER_HOME        DWORD = 0xAC
	VK_VOLUME_MUTE         DWORD = 0xAD
	VK_VOLUME_DOWN         DWORD = 0xAE
	VK_VOLUME_UP           DWORD = 0xAF
	VK_MEDIA_NEXT_TRACK    DWORD = 0xB0
	VK_MEDIA_PREV_TRACK    DWORD = 0xB1
	VK_MEDIA_STOP          DWORD = 0xB2
	VK_MEDIA_PLAY_PAUSE    DWORD = 0xB3
	VK_LAUNCH_MAIL         DWORD = 0xB4
	VK_LAUNCH_MEDIA_SELECT DWORD = 0xB5
	VK_LAUNCH_APP1         DWORD = 0xB6
	VK_LAUNCH_APP2         DWORD = 0xB7
	VK_OEM_1               DWORD = 0xBA
	VK_OEM_PLUS            DWORD = 0xBB
	VK_OEM_COMMA           DWORD = 0xBC
	VK_OEM_MINUS           DWORD = 0xBD
	VK_OEM_PERIOD          DWORD = 0xBE
	VK_OEM_2               DWORD = 0xBF
	VK_OEM_3               DWORD = 0xC0
	VK_OEM_4               DWORD = 0xDB
	VK_OEM_5               DWORD = 0xDC
	VK_OEM_6               DWORD = 0xDD
	VK_OEM_7               DWORD = 0xDE
	VK_OEM_8               DWORD = 0xDF
	VK_OEM_102             DWORD = 0xE2
	VK_PROCESSKEY          DWORD = 0xE5
	VK_PACKET              DWORD = 0xE7
	VK_ATTN                DWORD = 0xF6
	VK_CRSEL               DWORD = 0xF7
	VK_EXSEL               DWORD = 0xF8
	VK_EREOF               DWORD = 0xF9
	VK_PLAY                DWORD = 0xFA
	VK_ZOOM                DWORD = 0xFB
	VK_NONAME              DWORD = 0xFC
	VK_PA1                 DWORD = 0xFD
	VK_OEM_CLEAR           DWORD = 0xFE
)
