package repo

import (
	"context"
	"fmt"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	mapset "github.com/deckarep/golang-set"
	"github.com/jackc/pgtype"
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
	query += ` RETURNING id, path, created;`

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

	var pathf pgtype.Int4Array
	var ftime time.Time
	for i, _ := range posts {
		if rows.Next() {
			err = rows.Scan(&posts[i].ID, &pathf, &ftime)
			posts[i].Created = ftime
			fmt.Println(posts[i].ID, pathf, ftime)
			if err != nil {
				fmt.Println("err", err)
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

func (r *repoPostgres) GetUsers(ctx context.Context, slug string, params models.RequestParameters) ([]models.User, error) {
	var GetUsers = `SELECT Nickname, FullName, About, Email FROM users_forum WHERE Slug = $1`
	var rows pgx.Rows
	var err error

	if params.Since == "" {
		if params.Desc {
			GetUsers = GetUsers + ` ORDER BY Nickname DESC`
		} else {
			GetUsers = GetUsers + ` ORDER BY Nickname`
		}
		GetUsers = GetUsers + ` LIMIT $2;`
	} else {
		if params.Desc {
			GetUsers = GetUsers + ` AND Nickname < $2 ORDER BY Nickname DESC`
		} else {
			GetUsers = GetUsers + ` AND Nickname > $2  ORDER BY Nickname`
		}
		GetUsers = GetUsers + ` LIMIT $3;`
	}

	if params.Since == "" {
		rows, err = r.Conn.Query(ctx, GetUsers, slug, params.Limit)
	} else {
		rows, err = r.Conn.Query(ctx, GetUsers, slug, params.Since, params.Limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fUsers []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.NickName, &u.FullName, &u.About, &u.Email)
		if err != nil {
			return nil, err
		}
		fUsers = append(fUsers, u)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return fUsers, nil
}

func (r *repoPostgres) GetPostDetails(ctx context.Context, id int, related []string) (models.PostDetailed, error) {
	var (
		flagUser   bool
		flagForum  bool
		flagThread bool
	)
	var errScan error
	var fPost models.PostDetailed
	var fAuthor models.User
	var fForum models.Forum
	var fThread models.Thread
	fPost.Author = &fAuthor
	fPost.Forum = &fForum
	fPost.Thread = &fThread

	for _, param := range related {
		if param == "user" {
			flagUser = true
		}
		if param == "forum" {
			flagForum = true
		}
		if param == "thread" {
			flagThread = true
		}
	}

	var GetPostDetails = "SELECT post.Id, post.Author, post.Created, post.Forum, post.isEdited, " +
		"post.Message, post.Parent, post.Thread"

	if flagUser {
		GetPostDetails += ", users.Nickname, users.FullName, users.About, users.Email"
	}
	if flagForum {
		GetPostDetails += ", forum.Title, forum.\"user\", forum.Slug, forum.Posts, forum.Threads"
	}
	if flagThread {
		GetPostDetails += ", thread.Id, thread.Title, thread.Author, thread.Forum, " +
			"thread.Message, thread.Votes, thread.Slug, thread.Created"
	}
	GetPostDetails += " FROM post"

	if flagUser {
		GetPostDetails += " JOIN users ON post.Author = users.Nickname"
	}
	if flagForum {
		GetPostDetails += " JOIN forum ON post.Forum = forum.Slug"
	}
	if flagThread {
		GetPostDetails += " JOIN thread ON post.Thread = thread.Id"
	}
	GetPostDetails += ` WHERE post.Id = $1;`

	if flagUser && flagForum && flagThread {
		errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
			&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
			&fAuthor.NickName, &fAuthor.FullName, &fAuthor.About, &fAuthor.Email,
			&fForum.Title, &fForum.User, &fForum.Slug, &fForum.Posts, &fForum.Threads,
			&fThread.ID, &fThread.Title, &fThread.Author, &fThread.Forum,
			&fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
	} else {
		if flagUser && flagForum {
			errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
				&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
				&fAuthor.NickName, &fAuthor.FullName, &fAuthor.About, &fAuthor.Email,
				&fForum.Title, &fForum.User, &fForum.Slug, &fForum.Posts, &fForum.Threads)
		} else {
			if flagUser && flagThread {
				errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
					&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
					&fAuthor.NickName, &fAuthor.FullName, &fAuthor.About, &fAuthor.Email,
					&fThread.ID, &fThread.Title, &fThread.Author, &fThread.Forum,
					&fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
			} else {
				if flagForum && flagThread {
					errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
						&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
						&fForum.Title, &fForum.User, &fForum.Slug, &fForum.Posts, &fForum.Threads,
						&fThread.ID, &fThread.Title, &fThread.Author, &fThread.Forum,
						&fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
				} else {
					if flagUser {
						errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
							&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
							&fAuthor.NickName, &fAuthor.FullName, &fAuthor.About, &fAuthor.Email)
					} else {
						if flagForum {
							errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
								&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
								&fForum.Title, &fForum.User, &fForum.Slug, &fForum.Posts, &fForum.Threads)
						} else {
							if flagThread {
								errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
									&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread,
									&fThread.ID, &fThread.Title, &fThread.Author, &fThread.Forum,
									&fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
							} else {
								errScan = r.Conn.QueryRow(ctx, GetPostDetails, id).Scan(&fPost.Post.ID, &fPost.Post.Author, &fPost.Post.Created,
									&fPost.Post.Forum, &fPost.Post.IsEdited, &fPost.Post.Message, &fPost.Post.Parent, &fPost.Post.Thread)
							}
						}
					}
				}
			}
		}
	}
	if errScan != nil {
		if errScan == pgx.ErrNoRows {
			return models.PostDetailed{}, nil
		}
		return models.PostDetailed{}, errScan
	}
	return fPost, nil
	// SELECT post.*, author.*
	//FROM (SELECT * FROM post WHERE id = $1) AS post
	//JOIN author ON post.author = author.id;
}

func (r *repoPostgres) ChangePostInfo(ctx context.Context, post models.Post) (models.Post, int) {
	const ChangePostInfo = `UPDATE post SET Message = $1, IsEdited = true WHERE Id = $2;`
	_, err := r.Conn.Exec(ctx, ChangePostInfo, post.Message, post.ID)
	if err == nil {
		return post, http.StatusOK
	}
	return models.Post{}, http.StatusInternalServerError
}

func (r *repoPostgres) GetStatus(ctx context.Context) (models.Info, int) {
	countUsers := "SELECT count(*) FROM users"
	countForums := "SELECT count(*) FROM forum"
	countThreads := "SELECT count(*) FROM thread"
	countPosts := "SELECT count(*) FROM post"

	var info models.Info

	err := r.Conn.QueryRow(ctx, countUsers).Scan(&info.Users)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Info{}, http.StatusNotFound
		}
		return models.Info{}, http.StatusInternalServerError
	}
	err = r.Conn.QueryRow(ctx, countForums).Scan(&info.Forums)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Info{}, http.StatusNotFound
		}
		return models.Info{}, http.StatusInternalServerError
	}
	err = r.Conn.QueryRow(ctx, countThreads).Scan(&info.Threads)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Info{}, http.StatusNotFound
		}
		return models.Info{}, http.StatusInternalServerError
	}
	err = r.Conn.QueryRow(ctx, countPosts).Scan(&info.Posts)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Info{}, http.StatusNotFound
		}
		return models.Info{}, http.StatusInternalServerError
	}
	return info, http.StatusOK
}

func (r *repoPostgres) Clear(ctx context.Context) int {
	const ClearAll = `TRUNCATE TABLE users, forum, thread, post, vote, users_forum  CASCADE;`
	_, err := r.Conn.Exec(ctx, ClearAll)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (r *repoPostgres) GetPosts(ctx context.Context, idPost int, params models.RequestParameters) ([]models.Post, error) {
	return []models.Post{}, nil
}

func (r *repoPostgres) GetPostsFlat(ctx context.Context, params models.RequestParameters, threadID int) ([]models.Post, error) {
	var rows pgx.Rows
	var err error
	GetPosts := `SELECT Id, Author, Created, Forum, isEdited, Message, Parent, Thread FROM post WHERE Thread = $1`

	if params.SinceInt == 0 {
		if params.Desc {
			GetPosts = GetPosts + ` ORDER BY Id DESC`
		} else {
			GetPosts = GetPosts + ` ORDER BY Id`
		}
		GetPosts = GetPosts + ` LIMIT $2;`
		rows, err = r.Conn.Query(ctx, GetPosts, threadID, params.Limit)
	} else {
		if params.Desc {
			GetPosts = GetPosts + ` AND Id < $2 ORDER BY Id DESC`
		} else {
			GetPosts = GetPosts + ` AND Id > $2  ORDER BY Id`
		}
		GetPosts = GetPosts + ` LIMIT $3;`
		rows, err = r.Conn.Query(ctx, GetPosts, threadID, params.SinceInt, params.Limit)

	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		p := models.Post{}
		err := rows.Scan(&p.ID, &p.Author, &p.Created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
		if err != nil {
			return posts, models.InternalError
		}
		fmt.Println(p.ID, " ", p.Parent, " ", p.Path.Elements)
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *repoPostgres) GetPostsTree(ctx context.Context, params models.RequestParameters, thread int) ([]models.Post, error) {
	var rows pgx.Rows
	var errQuery error
	selectPosts := `SELECT post.Id, post.Author, post.Created, post.Forum, post.IsEdited, post.Message, post.Parent, post.Thread, post.Path
                  FROM post`

	if params.Limit == 100 {
		if params.SinceInt != 0 && params.Desc {
			selectPosts += ` JOIN post last ON last.id = $2 
                       WHERE post.path < last.path AND post.thread = $1 
                       ORDER BY post.path DESC, post.id DESC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, params.SinceInt)
		}
		if params.SinceInt == 0 && params.Desc {
			selectPosts += ` WHERE post.Thread = $1 
                        ORDER BY post.Path DESC, post.Id DESC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread)
		}
		if params.SinceInt != 0 && !params.Desc {
			selectPosts += ` JOIN post last ON last.id = $2 
                       WHERE post.path > last.path AND post.thread = $1 
                       ORDER BY post.path ASC, post.id ASC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, params.SinceInt)
		}
		if params.SinceInt == 0 && !params.Desc {
			selectPosts += ` WHERE post.Thread = $1 
                       ORDER BY post.Path ASC, post.Id ASC`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread)
		}
	} else {
		if params.SinceInt != 0 && params.Desc {
			selectPosts += ` JOIN post last ON last.id = $2 
                       WHERE post.path < last.path AND post.thread = $1 
                       ORDER BY post.path DESC, post.id DESC LIMIT $3`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, params.SinceInt, params.Limit)
		}
		if params.SinceInt == 0 && params.Desc {
			selectPosts += ` WHERE post.Thread = $1 
                       ORDER BY post.Path DESC, post.Id DESC LIMIT $2`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, params.Limit)
		}
		if params.SinceInt != 0 && !params.Desc {
			selectPosts += ` JOIN post last ON last.id = $2 
                       WHERE post.path > last.path AND post.thread = $1 
                       ORDER BY post.path ASC, post.id ASC LIMIT $3`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, params.SinceInt, params.Limit)
		}
		if params.SinceInt == 0 && !params.Desc {
			selectPosts += ` WHERE post.Thread = $1 
                       ORDER BY post.Path ASC, post.Id ASC LIMIT $2`
			rows, errQuery = r.Conn.Query(ctx, selectPosts, thread, params.Limit)
		}
	}
	if errQuery != nil {
		return []models.Post{}, models.InternalError
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		postOne := models.Post{}
		err := rows.Scan(&postOne.ID, &postOne.Author, &postOne.Created, &postOne.Forum, &postOne.IsEdited, &postOne.Message, &postOne.Parent, &postOne.Thread, &postOne.Path)

		if err != nil {
			return []models.Post{}, models.InternalError
		}
		fmt.Println(postOne.ID, " ", postOne.Parent, " ", postOne.Path.Elements)
		posts = append(posts, postOne)
	}
	return posts, nil
}

//func (r *repoPostgres) GetPostsTree(ctx context.Context, params models.RequestParameters, threadID int) ([]models.Post, error) {
//	//var rows pgx.Rows
//	//var err error
//	//GetPosts := `SELECT Id, Author, Created, Forum, IsEdited, Message, Parent, Thread FROM post WHERE Thread = $1`
//	//
//	//if params.SinceInt == 0 {
//	//	if params.Desc {
//	//		if params.Limit != 100 {
//	//			GetPosts = GetPosts + ` ORDER BY Path DESC, Id DESC`
//	//			GetPosts = GetPosts + ` LIMIT $2;`
//	//		} else {
//	//			GetPosts = GetPosts + ` ORDER BY Path, Id DESC`
//	//		}
//	//	} else {
//	//		GetPosts = GetPosts + ` ORDER BY Path, Id`
//	//		GetPosts = GetPosts + ` LIMIT $2;`
//	//	}
//	//	fmt.Println("repo start ", GetPosts)
//	//	rows, err = r.Conn.Query(ctx, GetPosts, threadID, params.Limit)
//	//	if err != nil {
//	//		return []models.Post{}, err
//	//	}
//	//} else {
//	//	if params.Desc {
//	//		GetPosts = `SELECT post.Id, post.Author, post.Created, post.Forum,
//	//  			post.IsEdited, post.Message, post.Parent, post.Thread
//	//			FROM post JOIN post parent ON parent.Id = $2 WHERE post.Path < parent.Path AND  post.Thread = $1
//	//			ORDER BY post.Path DESC, post.Id DESC`
//	//	} else {
//	//		GetPosts = `SELECT post.Id, post.Author, post.Created,
//	//			post.Forum, post.IsEdited, post.Message, post.Parent, post.Thread
//	//			FROM post JOIN post parent ON parent.Id = $2 WHERE post.Path > parent.Path AND  post.Thread = $1
//	//			ORDER BY post.Path, post.Id`
//	//	}
//	//	GetPosts = GetPosts + ` LIMIT $3;`
//	//	rows, err = r.Conn.Query(ctx, GetPosts, threadID, params.SinceInt, params.Limit)
//	//	if err != nil {
//	//		return []models.Post{}, err
//	//	}
//	//}
//
//	var rows pgx.Rows
//
//	query := `SELECT id, author, created, forum, isedited, message, parent, thread
//			  FROM post
//			  WHERE thread = $1 `
//
//	if params.Limit == 100 && params.SinceInt == 0 {
//		if params.Desc {
//			query += `ORDER BY path, id DESC`
//		} else {
//			query += `ORDER BY path, id ASC`
//		}
//		rows, _ = r.Conn.Query(ctx, query, threadID)
//	} else {
//		if params.Limit != 100 && params.SinceInt == 0 {
//			if params.Desc {
//				query += ` ORDER BY path DESC, id DESC LIMIT $2`
//			} else {
//				query += ` ORDER BY path, id ASC LIMIT $2`
//			}
//			rows, _ = r.Conn.Query(ctx, query, threadID, params.Limit)
//		}
//
//		if params.Limit != 100 && params.SinceInt != 0 {
//			if params.Desc {
//				query = `SELECT post.id, post.author, post.created,
//				post.forum, post.isedited, post.message, post.parent, post.thread
//				FROM post JOIN post parent ON parent.id = $2 WHERE post.path < parent.path AND  post.thread = $1
//				ORDER BY post.path DESC, post.id DESC LIMIT $3`
//			} else {
//				query = `SELECT post.id, post.author, post.created,
//				post.forum, post.isedited, post.message, post.parent, post.thread
//				FROM post JOIN post parent ON parent.id = $2 WHERE post.path > parent.path AND  post.thread = $1
//				ORDER BY post.path ASC, post.id ASC LIMIT $3`
//			}
//			rows, _ = r.Conn.Query(ctx, query, threadID, params.Since, params.Limit)
//		}
//
//		if params.Limit == 100 && params.SinceInt != 0 {
//			if params.Desc {
//				query = `SELECT post.id, post.author, post.created,
//				post.forum, post.isedited, post.message, post.parent, post.thread
//				FROM post JOIN post parent ON parent.id = $2 WHERE post.path < parent.path AND  post.thread = $1
//				ORDER BY post.path DESC, post.id DESC`
//			} else {
//				query = `SELECT post.id, post.author, post.created,
//				post.forum, post.isedited, post.message, post.parent, post.thread
//				FROM post JOIN post parent ON parent.id = $2 WHERE post.path > parent.path AND  post.thread = $1
//				ORDER BY post.path ASC, post.id ASC`
//			}
//			rows, _ = r.Conn.Query(ctx, query, threadID, params.Since)
//		}
//	}
//
//	posts := make([]models.Post, 0)
//	defer rows.Close()
//	for rows.Next() {
//		p := models.Post{}
//		err := rows.Scan(&p.ID, &p.Author, &p.Created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
//		if err != nil {
//			fmt.Println(err)
//			return posts, models.InternalError
//		}
//		posts = append(posts, p)
//	}
//
//	return posts, nil
//}

func (r *repoPostgres) GetPostsParent(ctx context.Context, params models.RequestParameters, thread int) ([]models.Post, error) {
	selectPostParents := fmt.Sprintf(`SELECT Id FROM post WHERE Thread = %d AND Parent = 0`, thread)

	if params.SinceInt == 0 {
		if params.Desc {
			selectPostParents += ` ORDER BY Id DESC `
		} else {
			selectPostParents += ` ORDER BY Id ASC `
		}
	} else {
		if params.Desc {
			selectPostParents += fmt.Sprintf(` AND Path[1] < (SELECT Path[1] FROM post WHERE Id = %d) ORDER BY Id DESC `, params.SinceInt)
		} else {
			selectPostParents += fmt.Sprintf(` AND Path[1] > (SELECT Path[1] FROM post WHERE Id = %d) ORDER BY Id ASC `, params.SinceInt)
		}
	}

	if params.Limit != 100 {
		selectPostParents += fmt.Sprintf(" LIMIT %d", params.Limit)
	}

	selectPosts := fmt.Sprintf(`SELECT Id, Author, Created, Forum, IsEdited, Message, Parent, Thread FROM post WHERE Path[1] = ANY (%s) `, selectPostParents)

	if params.Desc {
		selectPosts += ` ORDER BY Path[1] DESC, Path, Id `
	} else {
		selectPosts += ` ORDER BY Path[1] ASC, Path, Id `
	}

	rows, _ := r.Conn.Query(ctx, selectPosts)
	defer rows.Close()
	posts := make([]models.Post, 0)
	for rows.Next() {
		onePost := models.Post{}
		err := rows.Scan(&onePost.ID, &onePost.Author, &onePost.Created, &onePost.Forum, &onePost.IsEdited, &onePost.Message, &onePost.Parent, &onePost.Thread)
		if err != nil {
			return posts, models.InternalError
		}
		posts = append(posts, onePost)
	}

	return posts, nil
}

//func (r *repoPostgres) GetPostsParent(ctx context.Context, params models.RequestParameters, threadID int) ([]models.Post, error) {
//	var rows pgx.Rows
//
//	parents := fmt.Sprintf(`SELECT id FROM post WHERE thread = %d AND parent = 0 `, threadID)
//
//	if params.SinceInt != 0 {
//		if params.Desc {
//			parents += ` AND path[1] < ` + fmt.Sprintf(`(SELECT path[1] FROM post WHERE id = %s) `, params.Since)
//		} else {
//			parents += ` AND path[1] > ` + fmt.Sprintf(`(SELECT path[1] FROM post WHERE id = %s) `, params.Since)
//		}
//	}
//
//	if params.Desc {
//		parents += ` ORDER BY id DESC `
//	} else {
//		parents += ` ORDER BY id `
//	}
//
//	parents += fmt.Sprintf(` LIMIT %d`, params.Limit)
//
//	query := fmt.Sprintf(
//		`SELECT Id, Author, Created, Forum, IsEdited, Message, Parent, Thread FROM post WHERE Path[1] = ANY (%s) `, parents)
//
//	if params.Desc {
//		query += ` ORDER BY path[1] DESC, path,  id `
//	} else {
//		query += ` ORDER BY path[1], path,  id `
//	}
//
//	rows, _ = r.Conn.Query(ctx, query)
//	posts := make([]models.Post, 0)
//	defer rows.Close()
//	for rows.Next() {
//		p := models.Post{}
//		err := rows.Scan(&p.ID, &p.Author, &p.Created, &p.Forum, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
//		if err != nil {
//			return posts, models.InternalError
//		}
//		posts = append(posts, p)
//	}
//	return posts, nil
//}
