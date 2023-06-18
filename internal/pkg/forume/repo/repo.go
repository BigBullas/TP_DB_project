package repo

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type repoPostgres struct {
	Conn *pgxpool.Pool
}

func NewRepoPostgres(Conn *pgxpool.Pool) forume.Repository {
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
	rows, err := r.Conn.Query(ctx, CheckUserForUniq, user.NickName, user.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.NickName, &u.FullName, &u.About, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return users, nil
}

func (r *repoPostgres) GetUser(ctx context.Context, nickname string) (models.User, error) {
	const GetUser = `SELECT * FROM users WHERE Nickname = $1;`

	var fUser models.User
	err := r.Conn.QueryRow(ctx, GetUser, nickname).Scan(&fUser.NickName, &fUser.FullName, &fUser.About, &fUser.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return fUser, nil
}

func (r *repoPostgres) ChangeUserInfo(ctx context.Context, user models.User) (models.User, int) {
	const ChangeUserInfo = `UPDATE users SET FullName = $1, About = $2, Email = $3 WHERE Nickname = $4;`
	_, err := r.Conn.Exec(ctx, ChangeUserInfo, user.FullName, user.About, user.Email, user.NickName)
	if err == nil {
		return user, http.StatusOK
	}
	return models.User{}, http.StatusInternalServerError

}
