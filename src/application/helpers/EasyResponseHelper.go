package helper

import internals "auth-server-proxy/src/application/adapters/internals"

func EasyErrorRespond(cErrorCode string, cErrorMessage string) internals.ResponseAdapter {
	return internals.ResponseAdapter{
		StatusCode: 500,
		Errors: []internals.FieldErrorAdapter{
			{
				Code:    cErrorCode,
				Message: cErrorMessage,
				Field:   "",
			},
		},
	}
}

func EasyListErrorRespond(errorList []internals.FieldErrorAdapter, statusCode int) internals.ResponseAdapter {
	return internals.ResponseAdapter{
		StatusCode: statusCode,
		Errors:     errorList,
	}
}

func EasyEmptyRespond(statusCode int) internals.ResponseAdapter {
	return internals.ResponseAdapter{
		StatusCode: statusCode,
	}
}

func EasySuccessRespond(dataResponse interface{}, statusCode int) internals.ResponseAdapter {
	return internals.ResponseAdapter{
		StatusCode: statusCode,
		Response:   dataResponse,
	}
}
