package main

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume/delivery"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume/repo"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

// sudo docker rm -f my_container
// sudo docker build -t docker .
// sudo docker run -p 5000:5000 --name my_container -t docker

func main() {
	muxRoute := mux.NewRouter()
	conn := "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable&pool_max_conns=1000"
	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		log.Fatal("No connection to postgres", err)
	}

	fRepo := repo.NewRepoPostgres(pool)
	fUseCase := usecase.NewRepoUseCase(fRepo)
	fHandler := delivery.NewForumHandler(fUseCase)

	forum := muxRoute.PathPrefix("/api").Subrouter()
	{
		forum.HandleFunc("/user/{nickname}/create", fHandler.CreateUser).Methods(http.MethodPost)
		forum.HandleFunc("/user/{nickname}/profile", fHandler.GetUser).Methods(http.MethodGet)
		forum.HandleFunc("/user/{nickname}/profile", fHandler.ChangeUserInfo).Methods(http.MethodPost)

		forum.HandleFunc("/forum/create", fHandler.CreateForum).Methods(http.MethodPost)
		forum.HandleFunc("/forum/{slug}/details", fHandler.GetForumDetails).Methods(http.MethodGet)
		forum.HandleFunc("/forum/{slug}/create", fHandler.CreateThread).Methods(http.MethodPost)
		forum.HandleFunc("/forum/{slug}/users", fHandler.GetUsers).Methods(http.MethodGet)
		forum.HandleFunc("/forum/{slug}/threads", fHandler.GetThreads).Methods(http.MethodGet)

		forum.HandleFunc("/post/{id}/details", fHandler.GetPostDetails).Methods(http.MethodGet)
		forum.HandleFunc("/post/{id}/details", fHandler.ChangePostInfo).Methods(http.MethodPost)

		forum.HandleFunc("/service/status", fHandler.GetStatus).Methods(http.MethodGet)
		forum.HandleFunc("/service/clear", fHandler.Clear).Methods(http.MethodPost)

		forum.HandleFunc("/thread/{slug_or_id}/create", fHandler.CreatePosts).Methods(http.MethodPost)
		forum.HandleFunc("/thread/{slug_or_id}/details", fHandler.GetThreadDetails).Methods(http.MethodGet)
		forum.HandleFunc("/thread/{slug_or_id}/details", fHandler.ChangeThreadInfo).Methods(http.MethodPost)
		//forum.HandleFunc("/thread/{slug_or_id}/posts", fHandler.GetPostOfThread).Methods(http.MethodGet)
		forum.HandleFunc("/thread/{slug_or_id}/vote", fHandler.ChangeVote).Methods(http.MethodPost)
	}

	http.Handle("/", muxRoute)
	log.Print(http.ListenAndServe(":5000", muxRoute))
}
