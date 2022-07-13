package main

import (
	"fmt"
	"unsafe"

	"github.com/frei-0xff/win_hotkeys_daemon/winapi"
	"github.com/frei-0xff/win_hotkeys_daemon/wintypes"
	"github.com/tevino/abool/v2"
)

var (
	keyboardHook    wintypes.HHOOK
	keyCallback     wintypes.HOOKPROC = keyPressCallback
	winKeyPressed   *abool.AtomicBool
	altTabEmulating *abool.AtomicBool
)

/*
Handles callbacks for low level keyboard events
*/
func keyPressCallback(nCode int, wparam wintypes.WPARAM, lparam wintypes.LPARAM) wintypes.LRESULT {
	if nCode >= 0 {
		// Resolve struct that holds real event data
		kbd := (*wintypes.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		if kbd.VkCode == wintypes.VK_LWIN || kbd.VkCode == wintypes.VK_RWIN {
			if wparam == wintypes.WPARAM(wintypes.WM_KEYDOWN) {
				winKeyPressed.Set()
			} else {
				if altTabEmulating.IsSet() {
					winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_MENU), 0xb8, wintypes.KEYEVENTF_KEYUP, 0) // Alt Release
					altTabEmulating.UnSet()
				}
				winKeyPressed.UnSet()
			}
		}
		if wparam == wintypes.WPARAM(wintypes.WM_KEYDOWN) || wparam == wintypes.WPARAM(wintypes.WM_SYSKEYDOWN) {
			if winKeyPressed.IsSet() {
				fmt.Print("win+")
			}
			fmt.Println(kbd.VkCode)
			if winKeyPressed.IsSet() && kbd.VkCode == wintypes.VK_TAB {
				if altTabEmulating.IsNotSet() {
					altTabEmulating.Set()
					winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_MENU), 0xb8, 0, 0) //Alt Press
				}
				if altTabEmulating.IsSet() {
					winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_TAB), 0x8f, 0, 0)                        // Tab Press
					winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_TAB), 0x8f, wintypes.KEYEVENTF_KEYUP, 0) // Tab Release
					return 1
				}
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
	winKeyPressed = abool.New()
	altTabEmulating = abool.New()
	go Start()
	select {}
}
