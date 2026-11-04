package user

type CreateUserRequestDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
}
