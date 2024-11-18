package dtos

type Response struct {
	Body struct {
		Message string `json:"message" example:"User created successfully!" doc:"Confirmation message"`
	}
}
