package main

import (
	"fmt"
	"unsafe"

	"github.com/frei-0xff/win_hotkeys_daemon/winapi"
	"github.com/frei-0xff/win_hotkeys_daemon/wintypes"
)

var (
	keyboardHook  wintypes.HHOOK
	keyCallback   wintypes.HOOKPROC = keyPressCallback
	winKeyPressed bool              = false
)

/*
Handles callbacks for low level keyboard events
Resolves keypress to printable ascii character and appends to log
*/
func keyPressCallback(nCode int, wparam wintypes.WPARAM, lparam wintypes.LPARAM) wintypes.LRESULT {
	// Based on KEYUP events as that should be more reliable for how the resulting text actually looks
	if nCode >= 0 {
		// Resolve struct that holds real event data
		kbd := (*wintypes.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		if kbd.VkCode == wintypes.VK_LWIN || kbd.VkCode == wintypes.VK_RWIN {
			winKeyPressed = wparam == wintypes.WPARAM(wintypes.WM_KEYDOWN)
		}
		if wparam == wintypes.WPARAM(wintypes.WM_KEYDOWN) || wparam == wintypes.WPARAM(wintypes.WM_SYSKEYDOWN) {
			if winKeyPressed {
				fmt.Print("win+")
			}
			fmt.Println(kbd.VkCode)
			if winKeyPressed && kbd.VkCode == wintypes.VK_TAB {
				winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_MENU), 0xb8, 0, 0)                        //Alt Press
				winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_TAB), 0x8f, 0, 0)                         // Tab Press
				winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_TAB), 0x8f, wintypes.KEYEVENTF_KEYUP, 0)  // Tab Release
				winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_MENU), 0xb8, wintypes.KEYEVENTF_KEYUP, 0) // Alt Release
				return 1
			}
		}
	}
	return winapi.CallNextHookEx(keyboardHook, nCode, wparam, lparam)
}

/*
Attaches our initial hooks and runs the message queue
*/
func Start() {
	keyboardHook = winapi.SetWindowsHookEx(
		wintypes.WH_KEYBOARD_LL,
		keyCallback,
		0,
		0,
	)
	defer winapi.UnhookWindowsHookEx(keyboardHook)

	var msg wintypes.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}
}

func main() {
	go Start()
	select {}
}
