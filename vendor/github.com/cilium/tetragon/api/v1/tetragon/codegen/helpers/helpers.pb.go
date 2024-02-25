// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by protoc-gen-go-tetragon. DO NOT EDIT

package helpers

import (
	fmt "fmt"
	tetragon "github.com/cilium/tetragon/api/v1/tetragon"
)

// ResponseTypeString returns an event's type as a string
func ResponseTypeString(response *tetragon.GetEventsResponse) (string, error) {
	if response == nil {
		return "", fmt.Errorf("Response is nil")
	}

	event := response.Event
	if event == nil {
		return "", fmt.Errorf("Event is nil")
	}

	switch event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		return tetragon.EventType_PROCESS_EXEC.String(), nil
	case *tetragon.GetEventsResponse_ProcessExit:
		return tetragon.EventType_PROCESS_EXIT.String(), nil
	case *tetragon.GetEventsResponse_ProcessKprobe:
		return tetragon.EventType_PROCESS_KPROBE.String(), nil
	case *tetragon.GetEventsResponse_ProcessTracepoint:
		return tetragon.EventType_PROCESS_TRACEPOINT.String(), nil
	case *tetragon.GetEventsResponse_ProcessLoader:
		return tetragon.EventType_PROCESS_LOADER.String(), nil
	case *tetragon.GetEventsResponse_ProcessUprobe:
		return tetragon.EventType_PROCESS_UPROBE.String(), nil
	case *tetragon.GetEventsResponse_ProcessThrottle:
		return tetragon.EventType_PROCESS_THROTTLE.String(), nil
	case *tetragon.GetEventsResponse_Test:
		return tetragon.EventType_TEST.String(), nil
	case *tetragon.GetEventsResponse_RateLimitInfo:
		return tetragon.EventType_RATE_LIMIT_INFO.String(), nil

	}
	return "", fmt.Errorf("Unhandled response type %T", event)
}

// ResponseGetProcess returns a GetEventsResponse's process if it exists
func ResponseGetProcess(response *tetragon.GetEventsResponse) *tetragon.Process {
	if response == nil {
		return nil
	}

	event := response.Event
	if event == nil {
		return nil
	}

	return ResponseInnerGetProcess(event)
}

// ResponseInnerGetProcess returns a GetEventsResponse inner event's process if it exists
func ResponseInnerGetProcess(event tetragon.IsGetEventsResponse_Event) *tetragon.Process {
	switch ev := event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		return ev.ProcessExec.Process
	case *tetragon.GetEventsResponse_ProcessExit:
		return ev.ProcessExit.Process
	case *tetragon.GetEventsResponse_ProcessThrottle:
		return ev.ProcessThrottle.Process
	case *tetragon.GetEventsResponse_ProcessKprobe:
		return ev.ProcessKprobe.Process
	case *tetragon.GetEventsResponse_ProcessTracepoint:
		return ev.ProcessTracepoint.Process
	case *tetragon.GetEventsResponse_ProcessUprobe:
		return ev.ProcessUprobe.Process
	case *tetragon.GetEventsResponse_ProcessLoader:
		return ev.ProcessLoader.Process

	}
	return nil
}

// ResponseGetProcessKprobe returns a GetEventsResponse's process if it exists
func ResponseGetProcessKprobe(response *tetragon.GetEventsResponse) *tetragon.ProcessKprobe {
	if response == nil {
		return nil
	}

	return response.GetProcessKprobe()
}

// ResponseGetParent returns a GetEventsResponse's parent process if it exists
func ResponseGetParent(response *tetragon.GetEventsResponse) *tetragon.Process {
	if response == nil {
		return nil
	}

	event := response.Event
	if event == nil {
		return nil
	}

	return ResponseInnerGetParent(event)
}

// ResponseInnerGetParent returns a GetEventsResponse inner event's parent process if it exists
func ResponseInnerGetParent(event tetragon.IsGetEventsResponse_Event) *tetragon.Process {
	switch ev := event.(type) {
	case *tetragon.GetEventsResponse_ProcessExec:
		return ev.ProcessExec.Parent
	case *tetragon.GetEventsResponse_ProcessExit:
		return ev.ProcessExit.Parent
	case *tetragon.GetEventsResponse_ProcessThrottle:
		return ev.ProcessThrottle.Parent
	case *tetragon.GetEventsResponse_ProcessKprobe:
		return ev.ProcessKprobe.Parent
	case *tetragon.GetEventsResponse_ProcessTracepoint:
		return ev.ProcessTracepoint.Parent
	case *tetragon.GetEventsResponse_ProcessUprobe:
		return ev.ProcessUprobe.Parent

	}
	return nil
}
