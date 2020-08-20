package errno

var (
	/*
			code = 0 说明是正确返回；code > 0 说明是错误返回
			10002 首位表示服务级别错误，1为系统级错误；2为普通错误，通常是由用户非法操作引起的；
		 	第二第三位00表示服务模块代码；最后两位02表示具体错误代码
	*/

	OK                  = &Errno{Code: 0, Msg: "OK"}
	InternalServerError = &Errno{Code: 10001, Msg: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Msg: "Error occurred while binding the request body to the struct."}

	ErrValidation   = &Errno{Code: 20001, Msg: "Validation failed."}
	ErrDatabase     = &Errno{Code: 20002, Msg: "Database error."}
	ErrToken        = &Errno{Code: 20003, Msg: "Error occurred while signing the JSON web token."}
	ErrTokenInvalid = &Errno{Code: 20004, Msg: "The token was invalid."}
)
