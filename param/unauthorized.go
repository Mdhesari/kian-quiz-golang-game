package param

type UnAuthorizedResponse struct {
	Message string `json:"message"`
}

func GetDefaultUnAuthorizedResponse() UnAuthorizedResponse {
	return UnAuthorizedResponse{Message: "UnAuthorized."}
}
