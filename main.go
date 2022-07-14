package main

import (
	"fmt"
	"strings"
	"unsafe"
)

var (
	keyboardHook    HHOOK
	keyCallback     HOOKPROC = keyPressCallback
	altTabEmulating bool     = false
)

func winKeyState() bool {
	return GetAsyncKeyState(int(VK_LWIN)) > 1 ||
		GetAsyncKeyState(int(VK_RWIN)) > 1 ||
		GetAsyncKeyState(int(VK_APPS)) > 1
}

func shiftKeyState() bool {
	return GetAsyncKeyState(int(VK_SHIFT)) > 1
}

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
			if altTabEmulating && wparam != WPARAM(WM_KEYDOWN) &&
				(kbd.VkCode == VK_LWIN || kbd.VkCode == VK_RWIN || kbd.VkCode == VK_APPS) {
				altTabEmulating = false
				releaseKey(VK_MENU)
			}
			if wparam == WPARAM(WM_KEYDOWN) || wparam == WPARAM(WM_SYSKEYDOWN) {
				isWinKeyPressed := winKeyState()
				if isWinKeyPressed && kbd.VkCode == VK_TAB {
					if !altTabEmulating {
						altTabEmulating = true
						pressKey(VK_MENU)
					}
					pressKey(VK_TAB)
					releaseKey(VK_TAB)
					return 1
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
					if shiftKeyState() {
						clpText := strings.ReplaceAll(GetClipboardData(), `"`, ``)
						runProgram(fmt.Sprintf(`C:/ProgramData/chocolatey/lib/mpv.install/tools/mpv.com "%s"`, clpText), SW_MINIMIZE)
						return 1
					}
				}
				if isWinKeyPressed && kbd.VkCode == VK_RIGHT && !shiftKeyState() {
					SetCursorPos(2560, 1592)
					SetCursorPos(2560, 1592)
					return 1
				}
				if isWinKeyPressed && kbd.VkCode == VK_LEFT && !shiftKeyState() {
					SetCursorPos(960, 540)
					SetCursorPos(960, 540)
					return 1
				}
			}
		}
	}
	return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
}

/*
Attaches our initial hooks and runs the message queue
*/
func Start() {
	keyboardHook = SetWindowsHookEx(
		WH_KEYBOARD_LL,
		keyCallback,
		0,
		0,
	)
	defer UnhookWindowsHookEx(keyboardHook)

	var msg MSG
	for GetMessage(&msg, 0, 0, 0) != 0 {
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}

func main() {
	go Start()
	select {}
}
