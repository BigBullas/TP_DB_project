package repo

import (
	"context"
	"fmt"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	mapset "github.com/deckarep/golang-set"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"time"
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

func (r *repoPostgres) GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error) {
	const GetThreadBySlug = `SELECT * FROM thread WHERE Slug = $1;`
	var fThread models.Thread
	err := r.Conn.QueryRow(ctx, GetThreadBySlug, slug).
		Scan(&fThread.ID, &fThread.Title, &fThread.Author, &fThread.Forum,
			&fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Thread{}, nil
		}
		return models.Thread{}, err
	}
	return fThread, nil
}

func (r *repoPostgres) GetThreadById(ctx context.Context, id int) (models.Thread, error) {
	const GetThreadBySlug = `SELECT * FROM thread WHERE Id = $1;`
	var fThread models.Thread
	err := r.Conn.QueryRow(ctx, GetThreadBySlug, id).
		Scan(&fThread.ID, &fThread.Title, &fThread.Author, &fThread.Forum,
			&fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Thread{}, nil
		}
		return models.Thread{}, err
	}
	return fThread, nil
}

func (r *repoPostgres) CreatePosts(ctx context.Context, posts []models.Post, thread models.Thread) ([]models.Post, int) {
	created := time.Now()
	uniqueParents := mapset.NewSet()
	var uniqueParentsFlag bool
	uniqueAuthors := mapset.NewSet()
	values := make([]interface{}, 0)
	query := "INSERT INTO post (Author, Created, Forum, IsEdited, Message, Parent, Thread) VALUES"

	for k := range posts {
		post := &posts[k]
		post.Forum = thread.Forum
		post.Thread = thread.ID
		post.Created = created

		if post.Parent != 0 {
			uniqueParentsFlag = true
		}
		uniqueParents.Add(post.Parent)

		if post.Author == "" {
			return nil, http.StatusBadRequest
		}
		uniqueAuthors.Add(post.Author)

		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", k*7+1, k*7+2, k*7+3, k*7+4, k*7+5, k*7+6, k*7+7)
		values = append(values, post.Author, post.Created, post.Forum, post.IsEdited, post.Message, post.Parent, post.Thread)
	}
	query = query[:len(query)-1]
	query += ` RETURNING id;`

	queryCheckAuthors := "SELECT EXISTS (SELECT 1 FROM users WHERE Nickname IN ("
	it := 0
	argsNickname := make([]interface{}, uniqueAuthors.Cardinality())
	for author := range uniqueAuthors.Iter() {
		queryCheckAuthors += fmt.Sprintf(" $%d,", it+1)
		argsNickname[it] = author
		it++
	}
	queryCheckAuthors = queryCheckAuthors[:len(queryCheckAuthors)-1] + "));"
	var existsAuthor bool
	errAuthor := r.Conn.QueryRow(ctx, queryCheckAuthors, argsNickname...).Scan(&existsAuthor)
	if errAuthor != nil {
		return nil, http.StatusInternalServerError
	}
	if !existsAuthor {
		return nil, http.StatusNotFound
	}

	if uniqueParentsFlag {
		args := make([]interface{}, uniqueParents.Cardinality())
		i := 0
		queryCheckParents := "SELECT EXISTS (SELECT 1 FROM post WHERE Id IN ("
		queryCheckParentsPart2 := ")) as exists, thread FROM post WHERE Id IN ("
		for id := range uniqueParents.Iter() {
			queryCheckParents += fmt.Sprintf(" $%d,", i+1)
			queryCheckParentsPart2 += fmt.Sprintf(" $%d,", i+1)
			args[i] = id
			i++
		}
		queryCheckParentsPart2 = queryCheckParentsPart2[:len(queryCheckParentsPart2)-1] + ");"
		queryCheckParents = queryCheckParents[:len(queryCheckParents)-1] + queryCheckParentsPart2

		var exists bool
		var fThread int
		rows, err := r.Conn.Query(ctx, queryCheckParents, args...)
		if err != nil {
			return []models.Post{}, http.StatusInternalServerError
		}
		defer rows.Close()

		for j := 0; j < uniqueParents.Cardinality(); j++ {
			if rows.Next() {
				err = rows.Scan(&exists, &fThread)
				if err != nil {
					return nil, http.StatusInternalServerError
				}
				if fThread != thread.ID {
					return nil, http.StatusConflict
				}
			}
		}
		if rows.Err() != nil {
			return nil, http.StatusInternalServerError
		}

		if !exists {
			return []models.Post{}, http.StatusConflict
		}
	}

	rows, err := r.Conn.Query(ctx, query, values...)
	if err != nil {
		return []models.Post{}, http.StatusInternalServerError
	}
	defer rows.Close()

	for i, _ := range posts {
		if rows.Next() {
			err = rows.Scan(&posts[i].ID)
			if err != nil {
				return nil, http.StatusInternalServerError
			}
		}
	}
	if rows.Err() != nil {
		return nil, http.StatusInternalServerError
	}
	return posts, http.StatusCreated
}

func (r *repoPostgres) ChangeVote(ctx context.Context, vote models.Vote, thread models.Thread) (models.Thread, error) {
	const GetVote = `SELECT Author, Voice, Thread FROM vote WHERE Author = $1 AND Thread = $2;`
	var fVote models.Vote
	err := r.Conn.QueryRow(ctx, GetVote, vote.Nickname, vote.Thread).Scan(&fVote.Nickname, &fVote.Voice, &fVote.Thread)
	if err != nil {
		if err == pgx.ErrNoRows {
			const CreateVote = `INSERT INTO vote(Author, Voice, Thread) VALUES ($1, $2, $3);`
			_, err := r.Conn.Exec(ctx, CreateVote, vote.Nickname, vote.Voice, vote.Thread)
			if err != nil {
				return models.Thread{}, err
			}
			return models.Thread{}, nil
		} else {
			return models.Thread{}, err
		}
	}

	if fVote.Voice == vote.Voice {
		return thread, nil
	}

	const UpdateVote = `UPDATE vote SET Voice=$1 WHERE Author=$2 AND Thread=$3;`
	_, err = r.Conn.Exec(ctx, UpdateVote, vote.Voice, vote.Nickname, vote.Thread)
	if err != nil {
		return models.Thread{}, err
	}
	return models.Thread{}, nil
}

func (r *repoPostgres) ChangeThreadInfo(ctx context.Context, thread models.Thread) (models.Thread, int) {
	const ChangeThreadInfo = `UPDATE thread SET Title = $1, Message = $2 WHERE Id = $3;`
	_, err := r.Conn.Exec(ctx, ChangeThreadInfo, thread.Title, thread.Message, thread.ID)
	if err == nil {
		return thread, http.StatusOK
	}
	return models.Thread{}, http.StatusInternalServerError
}
