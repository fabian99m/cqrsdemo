package errors

var FileNotFound = Status{
	Code:    01,
	Message: "file not found",
}

var InvalidParams = Status{
	Code:    01,
	Message: "%v",
}

var FileSizeInvalid = Status{
	Code:    02,
	Message: "file sizeMb %f no allowed",
}

var MissingCommandName = Status{
	Code:    03,
	Message: "invalid command name",
}

var CommandNotRegistered = Status{
	Code:    04,
	Message: "command %s not registered",
}

var GenericError = Status{
	Code:    99,
	Message: "%v",
}
