//go:build darwin

package main

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreGraphics

#include <ApplicationServices/ApplicationServices.h>
#include <CoreGraphics/CoreGraphics.h>
#include <stdlib.h>

// Returns the focused window title as a malloc'd C string (caller must free),
// or NULL on failure (e.g. Accessibility permission not granted).
char* activeWindowTitle() {
	AXUIElementRef syswide = AXUIElementCreateSystemWide();

	AXUIElementRef focusedApp = NULL;
	AXUIElementCopyAttributeValue(syswide, kAXFocusedApplicationAttribute, (CFTypeRef*)&focusedApp);
	CFRelease(syswide);
	if (focusedApp == NULL) {
		return NULL;
	}

	AXUIElementRef focusedWindow = NULL;
	AXUIElementCopyAttributeValue(focusedApp, kAXFocusedWindowAttribute, (CFTypeRef*)&focusedWindow);
	CFRelease(focusedApp);
	if (focusedWindow == NULL) {
		return NULL;
	}

	CFStringRef title = NULL;
	AXUIElementCopyAttributeValue(focusedWindow, kAXTitleAttribute, (CFTypeRef*)&title);
	CFRelease(focusedWindow);
	if (title == NULL) {
		return NULL;
	}

	CFIndex bufLen = CFStringGetMaximumSizeForEncoding(CFStringGetLength(title), kCFStringEncodingUTF8) + 1;
	char* buf = (char*)malloc(bufLen);
	char* result = NULL;
	if (CFStringGetCString(title, buf, bufLen, kCFStringEncodingUTF8)) {
		result = buf;
	} else {
		free(buf);
	}
	CFRelease(title);
	return result;
}

long long idleTimeMs() {
	double secs = CGEventSourceSecondsSinceLastEventType(
		kCGEventSourceStateHIDSystemState,
		kCGAnyInputEventType
	);
	return (long long)(secs * 1000.0);
}
*/
import "C"
import "unsafe"

func getActiveWindowTitle() string {
	cs := C.activeWindowTitle()
	if cs == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

func getIdleTimeMs() int64 {
	return int64(C.idleTimeMs())
}
