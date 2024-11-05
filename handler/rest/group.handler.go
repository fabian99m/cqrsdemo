package handler

type GroupHandler struct {
	CommandRestHandler CommandRestHandler
	FileRestHandler    FileHandler
}

func NewGruopHandler(commandRestHandler *CommandRestHandler, fileHandler *FileHandler) *GroupHandler {
	return &GroupHandler{
		*commandRestHandler,
		*fileHandler,
	}
}
