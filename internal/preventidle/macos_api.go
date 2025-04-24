//go:build darwin
// +build darwin

package preventidle

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreGraphics
#include <ApplicationServices/ApplicationServices.h>
#include <unistd.h>
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"

	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

// CallSendInput 使用 Core Graphics API 模擬鍵盤或滑鼠事件。
// 當 mode 為 "key" 時，以 Shift 鍵 (key code 56) 為示例產生鍵盤按下與釋放事件；
// 當 mode 為 "mouse" 時，取得目前滑鼠位置，並向右平移 1 像素後建立滑鼠移動事件。
func CallSendInput(mode string) error {
	switch mode {
	case "key":
		// 建立事件來源（不能傳 nil，所以要呼叫 CGEventSourceCreate）
		source := C.CGEventSourceCreate(C.kCGEventSourceStateCombinedSessionState)
		if source == (C.CGEventSourceRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create event source for keyboard")
		}
		defer C.CFRelease(C.CFTypeRef(source))

		var key C.CGKeyCode = 56 // 使用 Shift 鍵作為示例

		// 建立鍵盤按下事件
		eventDown := C.CGEventCreateKeyboardEvent(source, key, C.bool(true))
		if eventDown == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create key down event")
		}
		C.CGEventPost(C.kCGHIDEventTap, eventDown)
		C.CFRelease(C.CFTypeRef(eventDown))

		// 延遲 10 毫秒
		C.usleep(10000)

		// 建立鍵盤放開事件
		eventUp := C.CGEventCreateKeyboardEvent(source, key, C.bool(false))
		if eventUp == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create key up event")
		}
		C.CGEventPost(C.kCGHIDEventTap, eventUp)
		C.CFRelease(C.CFTypeRef(eventUp))

		logger.LogInfo("macOS: simulated key press (shift key)")
		return nil

	case "mouse":
		// 建立事件來源
		source := C.CGEventSourceCreate(C.kCGEventSourceStateCombinedSessionState)
		if source == (C.CGEventSourceRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create event source for mouse")
		}
		defer C.CFRelease(C.CFTypeRef(source))

		// 建立一個事件以取得目前滑鼠位置
		event := C.CGEventCreate(source)
		if event == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create event for mouse location")
		}
		location := C.CGEventGetLocation(event)
		C.CFRelease(C.CFTypeRef(event))

		// 向右平移 1 像素
		newX := location.x + 1.0
		newY := location.y
		newPoint := C.CGPointMake(newX, newY)

		// 建立滑鼠移動事件
		mouseEvent := C.CGEventCreateMouseEvent(source, C.kCGEventMouseMoved, newPoint, C.kCGMouseButtonLeft)
		if mouseEvent == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create mouse move event")
		}
		C.CGEventPost(C.kCGHIDEventTap, mouseEvent)
		C.CFRelease(C.CFTypeRef(mouseEvent))

		logger.LogInfo("macOS: simulated mouse move")
		return nil

	default:
		return errors.New("unsupported mode for CallSendInput on macOS")
	}
}

// GetIdleTime 使用 CGEventSourceSecondsSinceLastEventType 取得系統閒置時間（以秒計），並轉換為 time.Duration。
func GetIdleTime() (time.Duration, error) {
	idleSeconds := float64(C.CGEventSourceSecondsSinceLastEventType(C.kCGEventSourceStateCombinedSessionState, C.kCGAnyInputEventType))
	if idleSeconds < 0 {
		return 0, errors.New("invalid idle time returned")
	}
	return time.Duration(idleSeconds * float64(time.Second)), nil
}
