package defs



type ErrResponse struct {
	HttpSC      int
	ErrorDetail ErrDetail
}

type ErrDetail struct {
	ErrorMsg string `json:"error"`
	ErrorCode string `json:"error_code"`
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{
		HttpSC: 400,
		ErrorDetail: ErrDetail{
			ErrorMsg: "请求body不正确",
			ErrorCode: "001",
		},
	}

	ErrorNotAuthUser = ErrResponse{
		HttpSC: 401,
		ErrorDetail: ErrDetail{
			ErrorMsg: "User authentication failed.",
			ErrorCode: "002",
		},
	}

	ErrorDBError = ErrResponse{
		HttpSC: 500,
		ErrorDetail: ErrDetail{
			ErrorMsg: "DB ops failed",
			ErrorCode: "003",
		},
	}

	ErrorInternalFaults = ErrResponse{
		HttpSC: 500,
		ErrorDetail: ErrDetail{
			ErrorMsg: "Internal service error",
			ErrorCode: "004",
		},
	}
)