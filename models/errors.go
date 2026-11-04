package models

import "errors"

// Erros de domínio. As camadas de baixo (repository/service) os devolvem e
// o handler é o único lugar que decide qual status HTTP cada um vira.
var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
)

// ValidationError indica dados de entrada inválidos. A mensagem é segura
// para ser devolvida ao cliente.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }
