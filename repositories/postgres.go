package repositories

import (
	"errors"

	"github.com/lib/pq"
)

// Códigos de erro do Postgres que traduzimos para erros de domínio.
// Lista completa: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
)

func hasPgCode(err error, code string) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == pq.ErrorCode(code)
}

// rowScanner é satisfeito tanto por *sql.Row quanto por *sql.Rows,
// permitindo uma única função de scan por entidade.
type rowScanner interface {
	Scan(dest ...any) error
}
