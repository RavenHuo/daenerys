/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package errors

// Codes ecode error interface which has a code & message.
type Codes interface {

	// Error return Code in string form
	Error() string
	// Code get error code.
	Code() int

	// Message get code message.
	Message() string
	// Equal for compatible.
	Equal(error) bool
}
