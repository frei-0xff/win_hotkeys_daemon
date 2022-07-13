package main

import (
	"fmt"
	"unsafe"

	"github.com/frei-0xff/win_hotkeys_daemon/winapi"
	"github.com/frei-0xff/win_hotkeys_daemon/wintypes"
)

var (
	keyboardHook    wintypes.HHOOK
	keyCallback     wintypes.HOOKPROC = keyPressCallback
	altTabEmulating bool              = false
)

/*
Handles callbacks for low level keyboard events
*/
func keyPressCallback(nCode int, wparam wintypes.WPARAM, lparam wintypes.LPARAM) wintypes.LRESULT {
	if nCode >= 0 {
		// Resolve struct that holds real event data
		kbd := (*wintypes.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		if kbd.ScanCode != 0xff {
			winKeyPressed := winapi.GetKeyState(int(wintypes.VK_LWIN)) != 0 || winapi.GetKeyState(int(wintypes.VK_RWIN)) != 0
			if (kbd.VkCode == wintypes.VK_LWIN || kbd.VkCode == wintypes.VK_RWIN) &&
				wparam != wintypes.WPARAM(wintypes.WM_KEYDOWN) && altTabEmulating {
				altTabEmulating = false
				winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_MENU), 0xff, wintypes.KEYEVENTF_KEYUP, 0) // Alt Release
			}
			if wparam == wintypes.WPARAM(wintypes.WM_KEYDOWN) || wparam == wintypes.WPARAM(wintypes.WM_SYSKEYDOWN) {
				fmt.Print(winapi.GetKeyState(int(wintypes.VK_LWIN)), " ")
				fmt.Print(winapi.GetKeyState(int(wintypes.VK_RWIN)), " ")
				if winKeyPressed {
					fmt.Print("win+")
				}
				fmt.Println(kbd.VkCode, " ", kbd.ScanCode)
				if winKeyPressed && kbd.VkCode == wintypes.VK_TAB {
					if !altTabEmulating {
						altTabEmulating = true
						winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_MENU), 0xff, 0, 0) //Alt Press
					}
					winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_TAB), 0xff, 0, 0)                        // Tab Press
					winapi.KeybdEvent(wintypes.BYTE(wintypes.VK_TAB), 0xff, wintypes.KEYEVENTF_KEYUP, 0) // Tab Release
					return 1
				}
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
