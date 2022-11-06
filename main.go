package main

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

const LLKHF_UP = 0x80

var (
	keyboardHook      HHOOK
	altTabEmulating   bool = false
	isWinKeyPressed   bool
	isShiftKeyPressed bool
)

func pressKey(key DWORD) {
	KeybdEvent(BYTE(key), 0xff, 0, 0)
}

func releaseKey(key DWORD) {
	KeybdEvent(BYTE(key), 0xff, KEYEVENTF_KEYUP, 0)
}

func runProgram(path string, flags DWORD) {
	go func() {
		WinExec(path, flags)
	}()
}

/*
Handles callbacks for low level keyboard events
*/
func keyPressCallback(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
	if nCode >= 0 {
		kbd := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		if kbd.ScanCode != 0xff {
			if kbd.Flags&LLKHF_UP != 0 { // Key UP
				if kbd.VkCode == VK_LWIN || kbd.VkCode == VK_RWIN || kbd.VkCode == VK_APPS {
					isWinKeyPressed = false
					if altTabEmulating {
						altTabEmulating = false
						releaseKey(VK_MENU)
					}
				}
				if kbd.VkCode == VK_LSHIFT || kbd.VkCode == VK_RSHIFT || kbd.VkCode == VK_SHIFT {
					isShiftKeyPressed = false
				}
			} else { // Key DOWN
				if kbd.VkCode == VK_LWIN || kbd.VkCode == VK_RWIN || kbd.VkCode == VK_APPS {
					isWinKeyPressed = true
				}
				if kbd.VkCode == VK_LSHIFT || kbd.VkCode == VK_RSHIFT || kbd.VkCode == VK_SHIFT {
					isShiftKeyPressed = true
				}

				if isWinKeyPressed && kbd.VkCode == VK_TAB {
					if !altTabEmulating {
						altTabEmulating = true
						pressKey(VK_MENU)
					}
				}
				if isWinKeyPressed && kbd.VkCode == VK_W {
					runProgram(`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`, SW_MAXIMIZE)
					return 1
				}
				if isWinKeyPressed && kbd.VkCode == VK_RETURN {
					runProgram(`C:\Program Files\WezTerm\wezterm-gui.exe`, SW_MAXIMIZE)
					return 1
				}
				if isWinKeyPressed && kbd.VkCode == VK_Q {
					PostMessage(GetForegroundWindow(), WM_CLOSE, 0, 0)
					return 1
				}
				// if isWinKeyPressed && kbd.VkCode == VK_J {
				//     pressKey(VK_MENU)
				//     pressKey(VK_ESCAPE)
				//     releaseKey(VK_ESCAPE)
				//     releaseKey(VK_MENU)
				//     return 1
				// }
				// if isWinKeyPressed && kbd.VkCode == VK_K {
				//     pressKey(VK_MENU)
				//     pressKey(VK_LSHIFT)
				//     pressKey(VK_ESCAPE)
				//     releaseKey(VK_ESCAPE)
				//     releaseKey(VK_LSHIFT)
				//     releaseKey(VK_MENU)
				//     return 1
				// }
				if isWinKeyPressed && kbd.VkCode == VK_M {
					if isShiftKeyPressed {
						clpText := strings.ReplaceAll(GetClipboardData(), `"`, ``)
						runProgram(fmt.Sprintf(`C:/ProgramData/chocolatey/lib/mpv.install/tools/mpv.com "%s"`, clpText), SW_MINIMIZE)
						return 1
					}
				}
				if isWinKeyPressed && kbd.VkCode == VK_RIGHT && !isShiftKeyPressed {
					SetCursorPos(2560, 1592)
					SetCursorPos(2560, 1592)
					return 1
				}
				if isWinKeyPressed && kbd.VkCode == VK_LEFT && !isShiftKeyPressed {
					SetCursorPos(960, 540)
					SetCursorPos(960, 540)
					return 1
				}
			}
		}
	}
	return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
}

func windowChangeCallback(hWinEventHook HWINEVENTHOOK, event DWORD, hwnd HWND,
	idObject LONG, idChild LONG, idEventThread DWORD,
	dwmsEventTime DWORD) uintptr {

	if GetWindowText(hwnd) == "Представление задач" {
		PostMessage(hwnd, WM_CLOSE, 0, 0)
	}
	return uintptr(0)
}

func Start() {
	keyboardHook = SetWindowsHookEx(
		WH_KEYBOARD_LL,
		keyPressCallback,
		0,
		0,
	)
	defer UnhookWindowsHookEx(keyboardHook)
	windowSwitchHook := SetWinEventHook(
		EVENT_OBJECT_FOCUS,
		EVENT_OBJECT_FOCUS,
		0,
		windowChangeCallback,
		0,
		0,
		0|2,
	)
	defer UnhookWinEvent(windowSwitchHook)

	var msg MSG
	for GetMessage(&msg, 0, 0, 0) != 0 {
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	go Start()
	select {}
}
