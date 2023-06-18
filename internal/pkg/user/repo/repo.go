package repo

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) user.Repository {
	return &repoPostgres{Conn: Conn}
}

func (r *repoPostgres) CreateUser(ctx context.Context, user models.User) error {
	const CreateUser = `INSERT INTO users(Nickname, FullName, About, Email) VALUES ($1, $2, $3, $4);`
	_, err := r.Conn.Exec(ctx, CreateUser, user.NickName, user.FullName, user.About, user.Email)
	if err != nil {
		return models.InternalError
	}
	return nil
}

func (r *repoPostgres) CheckUserForUniq(ctx context.Context, user models.User) ([]models.User, error) {
	const CheckUserForUniq = `SELECT * FROM users WHERE Nickname = $1 OR Email = $2;`
	return nil, nil
	//users, err := r.Conn.Exec(ctx, CheckUserForUniq, user.NickName, user.Email)
	//if err != nil {
	//	return nil, err
	//}
	////return []models.User{users}, nil
	//return nil, nil
}
