package handler

type GruopHandler struct {
	CommandRestHandler CommandRestHandler
	FileRestHandler FileHandler
}

func NewGruopHandler(commandRestHandler *CommandRestHandler, fileHandler *FileHandler) *GruopHandler {
	return &GruopHandler{
		*commandRestHandler,
		*fileHandler,
	}
}