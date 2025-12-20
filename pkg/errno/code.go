package errno

var (
	OK                  = New(0, "Success")
	InternalServerError = New(10001, "Internal server error")
	ErrBind             = New(10002, "Error occurred while binding the request body to the struct")
	ErrDatabase         = New(10003, "Invalid token")

	ErrTokenInvalid      = New(20101, "Token invalid")
	ErrPasswordIncorrect = New(20102, "Incorrect password")
	ErrUserBanned        = New(20103, "User not found")

	ErrVerifyCodeInvalid = New(20201, "Verification code invalid")
	ErrUserAlreadyExist  = New(20202, "User already exists")
	ErrRegisterFailed    = New(20203, "Register failed")

	ErrContactNotFound = New(20301, "Contact not found")
	ErrAlreadyFriend   = New(20302, "Already friends")

	ErrGroupNotFound  = New(30401, "Group not found")
	ErrGroupFull      = New(30402, "Group full")
	ErrNotGroupMember = New(30403, "not a member of this group")
)
