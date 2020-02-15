// Code generated by "stringer -type=BPlusErrorCode"; DO NOT EDIT.

package err

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ServiceNotFound-1000]
	_ = x[OperationNotFound-1001]
	_ = x[DecodingError-1002]
	_ = x[CannotGenerateHTTPRequest-1003]
	_ = x[CannotGenerateHTTPRequest1-1004]
	_ = x[CannotGenerateHTTPRequestForPayload-1005]
	_ = x[ResponseUnmarshalException-1006]
	_ = x[ParamsNotExpected-1007]
	_ = x[HTTPCallFailed-1008]
	_ = x[CannotReadResponseBody-1009]
	_ = x[CannotMakeStateEntity-1010]
	_ = x[ErrorInDecoding-1011]
	_ = x[ErrorInAutoState-1012]
	_ = x[AutoStateNotConfigured-1013]
	_ = x[InvalidState-1014]
	_ = x[InvalidEvent-1015]
	_ = x[CannotReadFile-1016]
	_ = x[EventNotFoundInRequest-1017]
	_ = x[ParameterMissingInRequest-1018]
	_ = x[ErrorInObtainingSTM-1019]
}

const _BPlusErrorCode_name = "ServiceNotFoundOperationNotFoundDecodingErrorCannotGenerateHTTPRequestCannotGenerateHTTPRequest1CannotGenerateHTTPRequestForPayloadResponseUnmarshalExceptionParamsNotExpectedHTTPCallFailedCannotReadResponseBodyCannotMakeStateEntityErrorInDecodingErrorInAutoStateAutoStateNotConfiguredInvalidStateInvalidEventCannotReadFileEventNotFoundInRequestParameterMissingInRequestErrorInObtainingSTM"

var _BPlusErrorCode_index = [...]uint16{0, 15, 32, 45, 70, 96, 131, 157, 174, 188, 210, 231, 246, 262, 284, 296, 308, 322, 344, 369, 388}

func (i BPlusErrorCode) String() string {
	i -= 1000
	if i < 0 || i >= BPlusErrorCode(len(_BPlusErrorCode_index)-1) {
		return "BPlusErrorCode(" + strconv.FormatInt(int64(i+1000), 10) + ")"
	}
	return _BPlusErrorCode_name[_BPlusErrorCode_index[i]:_BPlusErrorCode_index[i+1]]
}
