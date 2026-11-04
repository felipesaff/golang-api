package models

import "errors"

// Erros de domínio. As camadas de baixo (repository/service) os devolvem e
// o handler é o único lugar que decide qual status HTTP cada um vira.
var (
	// Recurso não encontrado (404).
	ErrOrderNotFound   = errors.New("order not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")

	// Conflito com o estado atual dos dados (409).
	ErrUserAlreadyExists = errors.New("email or document already registered")
	ErrUserHasOrders     = errors.New("user has orders and cannot be deleted")
	ErrProductInUse      = errors.New("product is used in orders and cannot be deleted")

	// A requisição aponta para algo que não existe (422). Ex.: criar um pedido
	// com um produto inexistente. Costuma vir combinado com um ErrXNotFound:
	//   fmt.Errorf("%w: %w", ErrInvalidReference, ErrProductNotFound)
	ErrInvalidReference = errors.New("invalid reference")
)

// ValidationError indica dados de entrada inválidos (400). A mensagem é segura
// para ser devolvida ao cliente.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }
