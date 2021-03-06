// Copyright (c) 2020 Red Hat, Inc.
//
// The Universal Permissive License (UPL), Version 1.0
//
// Subject to the condition set forth below, permission is hereby granted to any
// person obtaining a copy of this software, associated documentation and/or data
// (collectively the "Software"), free of charge and under any and all copyright
// rights in the Software, and any and all patent rights owned or freely
// licensable by each licensor hereunder covering either (i) the unmodified
// Software as contributed to or provided by such licensor, or (ii) the Larger
// Works (as defined below), to deal in both
//
// (a) the Software, and
// (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
// one is included with the Software (each a "Larger Work" to which the Software
// is contributed by such licensors),
//
// without restriction, including without limitation the rights to copy, create
// derivative works of, display, perform, and distribute the Software and make,
// use, sell, offer for sale, import, export, have made, and have sold the
// Software and the Larger Work(s), and to sublicense the foregoing rights on
// either these or other terms.
//
// This license is subject to the following condition:
// The above copyright notice and either this complete permission notice or at
// a minimum a reference to the UPL must be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package test

import (
	"net/http"
	"strconv"

	"github.com/onsi/gomega/ghttp"
	rhjmcv1alpha2 "github.com/rh-jmc-team/container-jfr-operator/pkg/apis/rhjmc/v1alpha2"
	jfrclient "github.com/rh-jmc-team/container-jfr-operator/pkg/client"
)

func NewDumpHandler() http.HandlerFunc {
	return createRecordingHandler(30, true)
}

func NewDumpFailHandler() http.HandlerFunc {
	return createRecordingHandler(30, false)
}

func NewStartHandler() http.HandlerFunc {
	return createRecordingHandler(0, true)
}

func NewStartFailHandler() http.HandlerFunc {
	return createRecordingHandler(0, false)
}

func createRecordingHandler(duration int64, succeed bool) http.HandlerFunc {
	desc := NewRecordingDescriptors("CREATED", duration)[0]
	handlers := []http.HandlerFunc{
		ghttp.VerifyRequest(http.MethodPost, "/api/v1/targets/1.2.3.4:8001/recordings"),
		ghttp.VerifyContentType("application/x-www-form-urlencoded"),
		ghttp.VerifyFormKV("recordingName", "test-recording"),
		ghttp.VerifyFormKV("events", "jdk.socketRead:enabled=true,jdk.socketWrite:enabled=true"),
		verifyToken(),
	}
	if duration > 0 {
		handlers = append(handlers, ghttp.VerifyFormKV("duration", strconv.Itoa(int(duration))))
	}
	if succeed {
		handlers = append(handlers, ghttp.RespondWithJSONEncoded(http.StatusOK, desc))
	} else {
		handlers = append(handlers, ghttp.RespondWith(http.StatusBadRequest,
			"Recording with name \"test-recording\" already exists"))
	}
	return ghttp.CombineHandlers(handlers...)
}

func NewStopHandler() http.HandlerFunc {
	return stopHandler(true)
}

func NewStopFailHandler() http.HandlerFunc {
	return stopHandler(false)
}

func stopHandler(succeed bool) http.HandlerFunc {
	handlers := []http.HandlerFunc{
		ghttp.VerifyRequest(http.MethodPatch, "/api/v1/targets/1.2.3.4:8001/recordings/test-recording"),
		ghttp.VerifyContentType("text/plain"),
		ghttp.VerifyBody([]byte("stop")),
		verifyToken(),
	}
	if succeed {
		handlers = append(handlers, ghttp.RespondWith(http.StatusOK, nil))
	} else {
		handlers = append(handlers, ghttp.RespondWith(http.StatusNotFound,
			"Recording with name \"test-recording\" not found"))
	}
	return ghttp.CombineHandlers(handlers...)
}

func NewSaveHandler() http.HandlerFunc {
	return saveHandler(true)
}

func NewSaveFailHandler() http.HandlerFunc {
	return saveHandler(false)
}

func saveHandler(succeed bool) http.HandlerFunc {
	handlers := []http.HandlerFunc{
		ghttp.VerifyRequest(http.MethodPatch, "/api/v1/targets/1.2.3.4:8001/recordings/test-recording"),
		ghttp.VerifyContentType("text/plain"),
		ghttp.VerifyBody([]byte("save")),
		verifyToken(),
	}
	if succeed {
		handlers = append(handlers, ghttp.RespondWith(http.StatusOK, "saved-test-recording.jfr"))
	} else {
		handlers = append(handlers, ghttp.RespondWith(http.StatusNotFound,
			"Recording with name \"test-recording\" not found"))
	}
	return ghttp.CombineHandlers(handlers...)
}

func NewListHandler(descriptors []jfrclient.RecordingDescriptor) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodGet, "/api/v1/targets/1.2.3.4:8001/recordings"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusOK, descriptors),
	)
}

func NewListFailHandler(descriptors []jfrclient.RecordingDescriptor) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodGet, "/api/v1/targets/1.2.3.4:8001/recordings"),
		verifyToken(),
		ghttp.RespondWith(http.StatusInternalServerError, "test message"),
	)
}

func NewRecordingDescriptors(state string, duration int64) []jfrclient.RecordingDescriptor {
	return []jfrclient.RecordingDescriptor{
		{
			Name:        "test-recording",
			State:       state,
			StartTime:   1597090030341,
			Duration:    duration,
			DownloadURL: "http://path/to/test-recording.jfr",
		},
	}
}

func NewListSavedHandler(saved []jfrclient.SavedRecording) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodGet, "/api/v1/recordings"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusOK, saved),
	)
}

func NewListSavedFailHandler(saved []jfrclient.SavedRecording) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodGet, "/api/v1/recordings"),
		verifyToken(),
		ghttp.RespondWith(http.StatusNotImplemented, "Archive path /bad/dir does not exist"),
	)
}

func NewSavedRecordings() []jfrclient.SavedRecording {
	return []jfrclient.SavedRecording{
		{
			Name:        "saved-test-recording.jfr",
			DownloadURL: "http://path/to/saved-test-recording.jfr",
			ReportURL:   "http://path/to/saved-test-recording.html",
		},
	}
}

func NewDeleteHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodDelete, "/api/v1/targets/1.2.3.4:8001/recordings/test-recording"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
	)
}

func NewDeleteFailHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodDelete, "/api/v1/targets/1.2.3.4:8001/recordings/test-recording"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusNotFound,
			"No recording with name \"test-recording\" found"),
	)
}

func NewDeleteSavedHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodDelete, "/api/v1/recordings/saved-test-recording.jfr"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
	)
}

func NewDeleteSavedFailHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodDelete, "/api/v1/recordings/saved-test-recording.jfr"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusNotFound, "saved-test-recording.jfr"),
	)
}

func NewListEventTypesHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodGet, "/api/v1/targets/1.2.3.4:8001/events"),
		verifyToken(),
		ghttp.RespondWithJSONEncoded(http.StatusOK, NewEventTypes()),
	)
}

func NewListEventTypesFailHandler() http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest(http.MethodGet, "/api/v1/targets/1.2.3.4:8001/events"),
		verifyToken(),
		ghttp.RespondWith(http.StatusUnauthorized, nil),
	)
}

func NewEventTypes() []rhjmcv1alpha2.EventInfo {
	return []rhjmcv1alpha2.EventInfo{
		{
			TypeID:      "jdk.socketRead",
			Name:        "Socket Read",
			Description: "Reading data from a socket",
			Category:    []string{"Java Application"},
			Options: map[string]rhjmcv1alpha2.OptionDescriptor{
				"enabled": {
					Name:         "Enabled",
					Description:  "Record event",
					DefaultValue: "false",
				},
				"stackTrace": {
					Name:         "Stack Trace",
					Description:  "Record stack traces",
					DefaultValue: "false",
				},
				"threshold": {
					Name:         "Threshold",
					Description:  "Record event with duration above or equal to threshold",
					DefaultValue: "0ns[ns]",
				},
			},
		},
	}
}

func verifyToken() http.HandlerFunc {
	return ghttp.VerifyHeaderKV("Authorization", "Bearer myToken")
}
