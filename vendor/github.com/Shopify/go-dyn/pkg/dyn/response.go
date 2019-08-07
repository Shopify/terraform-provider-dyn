package dyn

import "fmt"

// values for responseHeader.Status
const (
	responseSuccess    = "success"
	responseFailure    = "failure"
	responseIncomplete = "incomplete"
)

// values for responseMessage.Level
const (
	responseMessageFatal = "FATAL"
	responseMessageError = "ERROR"
	responseMessageWarn  = "WARN"
	responseMessageInfo  = "INFO"
)

// values for responseMessage.ErrorCode
const (
	responseMessageDeprecatedRequest  = "DEPRECATED_REQUEST"  // The requested command is deprecated
	responseMessageIllegalOperation   = "ILLEGAL_OPERATION"   // The operation is not allowed with this data set
	responseMessageInternalError      = "INTERNAL_ERROR"      // An error occurred that cannot be classified.
	responseMessageInvalidData        = "INVALID_DATA"        // A field contained data that was invalid
	responseMessageInvalidRequest     = "INVALID_REQUEST"     // The request was not recognized as a valid command
	responseMessageInvalidVersion     = "INVALID_VERSION"     // The version number passed in was invalid
	responseMessageMissingData        = "MISSING_DATA"        // A required field was not provided
	responseMessageNotFound           = "NOT_FOUND"           // No results were found
	responseMessageOperationFailed    = "OPERATION_FAILED"    // The operation failed to complete successfully
	responseMessagePermissionDenied   = "PERMISSION_DENIED"   // This user does not have permission to perform this action
	responseMessageServiceUnavailable = "SERVICE_UNAVAILABLE" // The requested service is currently unavailable.
	responseMessageTargetExists       = "TARGET_EXISTS"       // Attempted to add a duplicate resource
	responseMessageUnknownError       = "UNKNOWN_ERROR"       // An error occurred that cannot be classified
)

// common header for API responses
type responseHeader struct {
	JobID    int               `json:"job_id,omitempty"`
	Status   string            `json:"status"`
	Messages []responseMessage `json:"msgs,omitempty"`
}

// messages within API response header
type responseMessage struct {
	Source    string `json:"SOURCE"`
	Level     string `json:"LVL"`
	Info      string `json:"INFO"`
	ErrorCode string `json:"ERR_CD"`
}

// Error implements the error interface for the responseMessage type.
func (m responseMessage) Error() string {
	return fmt.Sprintf("%v: %v", m.ErrorCode, m.Info)
}
