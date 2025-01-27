package enums

type InternalHttpMessage string

const (
	Success       InternalHttpMessage = "La petición se completó exitosamente"
	SuccessEmpty  InternalHttpMessage = "La solicitud se completó exitosamente, pero su respuesta no tiene ningún contenido"
	InternalError InternalHttpMessage = "La solicitud no pudo ser procesada debido a un problema interno"
	RequestError  InternalHttpMessage = "La petición no pudo completarse debido a que la solicitud no fue válida"
)
