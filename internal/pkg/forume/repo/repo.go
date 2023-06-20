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

func (r *repoPostgres) CheckForumForUniq(ctx context.Context, forum models.Forum) ([]models.Forum, int) {
	const CheckForumForUniq = `SELECT * FROM forum WHERE Slug = $1;`
	rows, err := r.Conn.Query(ctx, CheckForumForUniq, forum.Slug)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer rows.Close()

	var forums []models.Forum
	for rows.Next() {
		var f models.Forum
		err := rows.Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
		if err != nil {
			return nil, http.StatusInternalServerError
		}
		forums = append(forums, f)
	}
	if rows.Err() != nil {
		return nil, http.StatusInternalServerError
	}
	return forums, http.StatusOK
}

func (r *repoPostgres) CreateForum(ctx context.Context, forum models.Forum) ([]models.Forum, int) {
	const CreateForum = `INSERT INTO forum(Title, "user", Slug, Posts, Threads) VALUES ($1, $2, $3, $4, $5);`
	_, err := r.Conn.Exec(ctx, CreateForum, forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads)
	if err != nil {
		return []models.Forum{}, http.StatusInternalServerError
	}
	return []models.Forum{forum}, http.StatusCreated
}

func (r *repoPostgres) GetForumDetails(ctx context.Context, slug string) (models.Forum, error) {
	const GetForumDetails = `SELECT * FROM forum WHERE Slug = $1;`

	var fForum models.Forum
	err := r.Conn.QueryRow(ctx, GetForumDetails, slug).
		Scan(&fForum.Title, &fForum.User, &fForum.Slug, &fForum.Posts, &fForum.Threads)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Forum{}, nil
		}
		return models.Forum{}, err
	}
	return fForum, nil
}

func (r *repoPostgres) CheckThreadForUniq(ctx context.Context, thread models.Thread) ([]models.Thread, int) {
	const CheckThreadForUniq = `SELECT * FROM thread WHERE Slug = $1;`
	rows, err := r.Conn.Query(ctx, CheckThreadForUniq, thread.Slug)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var t models.Thread
		err := rows.Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
		if err != nil {
			return nil, http.StatusInternalServerError
		}
		threads = append(threads, t)
	}
	if rows.Err() != nil {
		return nil, http.StatusInternalServerError
	}
	return threads, http.StatusOK
}

func (r *repoPostgres) CreateThread(ctx context.Context, thread models.Thread) ([]models.Thread, int) {
	const CreateThread = `INSERT INTO thread(Title, Author, Forum, Message, Votes, Slug, Created) VALUES ($1, $2, $3, $4, $5, $6, $7) returning Id;`
	err := r.Conn.QueryRow(ctx, CreateThread, thread.Title, thread.Author, thread.Forum,
		thread.Message, thread.Votes, thread.Slug, thread.Created).Scan(&thread.ID)
	if err != nil {
		return []models.Thread{}, http.StatusInternalServerError
	}
	return []models.Thread{thread}, http.StatusCreated
}

func (r *repoPostgres) GetThreads(ctx context.Context, slug string, params models.RequestParameters) ([]models.Thread, error) {
	var GetThreads = `SELECT * FROM thread WHERE Forum = $1`

	if params.Desc {
		GetThreads = GetThreads + ` AND created <= $2 ORDER BY created DESC, id DESC`
		if params.Since == "" {
			params.Since = "9999-12-31T23:59:59.000Z"
		}
	} else {
		GetThreads = GetThreads + ` AND created >= $2  ORDER BY created, id`
		if params.Since == "" {
			params.Since = "0001-01-01T00:00:00.000Z"
		}
	}
	GetThreads = GetThreads + ` LIMIT $3;`

	rows, err := r.Conn.Query(ctx, GetThreads, slug, params.Since, params.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fThreads []models.Thread
	for rows.Next() {
		var t models.Thread
		err := rows.Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
		if err != nil {
			return nil, err
		}
		fThreads = append(fThreads, t)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return fThreads, nil
}
