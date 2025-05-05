//go:build darwin
// +build darwin

package preventidle

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreGraphics
#include <ApplicationServices/ApplicationServices.h>
#include <unistd.h>

#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <IOKit/pwr_mgt/IOPMLib.h>
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

var (
	sysAssertion  C.IOPMAssertionID
	dispAssertion C.IOPMAssertionID
)

// PreventSleep 建立「PreventUserIdleSystemSleep」宣告
func PreventSleep() error {
	cs := C.CString("GoIdleGuard active")
	defer C.free(unsafe.Pointer(cs))
	reason := C.CFStringCreateWithCString(C.kCFAllocatorDefault, cs, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(reason))

	// 1) Prevent system sleep
	if st := C.IOPMAssertionCreateWithName(
		C.kIOPMAssertionTypePreventUserIdleSystemSleep,
		C.kIOPMAssertionLevelOn,
		reason,
		&sysAssertion,
	); st != C.kIOReturnSuccess {
		return fmt.Errorf("system sleep assertion failed: %d", st)
	}

	// 2) Prevent display sleep (screensaver, display dimming)
	if st := C.IOPMAssertionCreateWithName(
		C.kIOPMAssertionTypePreventUserIdleDisplaySleep,
		C.kIOPMAssertionLevelOn,
		reason,
		&dispAssertion,
	); st != C.kIOReturnSuccess {
		// 先解除第一個再回錯誤
		C.IOPMAssertionRelease(sysAssertion)
		return fmt.Errorf("display sleep assertion failed: %d", st)
	}

	return nil
}

// AllowIdle 恢復系統與顯示器閒置行為
func AllowIdle() error {
	if st := C.IOPMAssertionRelease(dispAssertion); st != C.kIOReturnSuccess {
		return fmt.Errorf("release display assertion failed: %d", st)
	}
	if st := C.IOPMAssertionRelease(sysAssertion); st != C.kIOReturnSuccess {
		return fmt.Errorf("release system assertion failed: %d", st)
	}
	return nil
}

// CallSendInput 使用 Core Graphics API 模擬鍵盤或滑鼠事件。
// 當 mode 為 "key" 時，以 Shift 鍵 (key code 56) 為示例產生鍵盤按下與釋放事件；
// 當 mode 為 "mouse" 時，取得目前滑鼠位置，並向右平移 1 像素後建立滑鼠移動事件。
func CallSendInput(mode string) error {
	switch mode {
	case "key":
		source := C.CGEventSourceCreate(C.kCGEventSourceStateCombinedSessionState)
		if source == (C.CGEventSourceRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create event source for keyboard")
		}
		defer C.CFRelease(C.CFTypeRef(source))

		var key C.CGKeyCode = 56 // 空格鍵 (kVK_Space)

		eventDown := C.CGEventCreateKeyboardEvent(source, key, C.bool(true))
		if eventDown == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create key down event")
		}
		C.CGEventPost(C.kCGSessionEventTap, eventDown) // 改用 kCGSessionEventTap
		C.CFRelease(C.CFTypeRef(eventDown))

		C.usleep(10000)

		eventUp := C.CGEventCreateKeyboardEvent(source, key, C.bool(false))
		if eventUp == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create key up event")
		}
		C.CGEventPost(C.kCGSessionEventTap, eventUp)
		C.CFRelease(C.CFTypeRef(eventUp))
		logger.LogInfo("macOS: simulated key press (space key)")
		return nil
	case "mouse":
		source := C.CGEventSourceCreate(C.kCGEventSourceStateCombinedSessionState)
		if source == (C.CGEventSourceRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create event source for mouse")
		}
		defer C.CFRelease(C.CFTypeRef(source))

		event := C.CGEventCreate(source)
		if event == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create event for mouse location")
		}
		location := C.CGEventGetLocation(event)
		C.CFRelease(C.CFTypeRef(event))

		newX := location.x + 1.0
		newY := location.y
		newPoint := C.CGPointMake(newX, newY)

		mouseEvent := C.CGEventCreateMouseEvent(source, C.kCGEventMouseMoved, newPoint, C.kCGMouseButtonLeft)
		if mouseEvent == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create mouse move event")
		}
		C.CGEventPost(C.kCGSessionEventTap, mouseEvent)
		C.CFRelease(C.CFTypeRef(mouseEvent))

		clickEvent := C.CGEventCreateMouseEvent(source, C.kCGEventLeftMouseDown, newPoint, C.kCGMouseButtonLeft)
		if clickEvent == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create mouse down event")
		}
		C.CGEventPost(C.kCGSessionEventTap, clickEvent)
		C.CFRelease(C.CFTypeRef(clickEvent))

		C.usleep(10000)

		clickUpEvent := C.CGEventCreateMouseEvent(source, C.kCGEventLeftMouseUp, newPoint, C.kCGMouseButtonLeft)
		if clickUpEvent == (C.CGEventRef)(unsafe.Pointer(nil)) {
			return errors.New("failed to create mouse up event")
		}
		C.CGEventPost(C.kCGSessionEventTap, clickUpEvent)
		C.CFRelease(C.CFTypeRef(clickUpEvent))

		logger.LogInfo("macOS: simulated mouse move and click")
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
