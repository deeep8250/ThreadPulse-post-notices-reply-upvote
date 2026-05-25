package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	replieshandler "github.com/threadpulse/internal/replies/handlers"
	repliesrepo "github.com/threadpulse/internal/replies/repository"
	repliesservice "github.com/threadpulse/internal/replies/services"
	threadhandler "github.com/threadpulse/internal/threads/handlers"
	threadrepo "github.com/threadpulse/internal/threads/repository"
	threadservice "github.com/threadpulse/internal/threads/services"

	"github.com/gin-gonic/gin"
	authhandler "github.com/threadpulse/internal/auth/handlers"
	authrepo "github.com/threadpulse/internal/auth/repository"
	authservice "github.com/threadpulse/internal/auth/services"
	"github.com/threadpulse/internal/middleware"
	upvotehandler "github.com/threadpulse/internal/upvotes/handlers"
	upvoterepo "github.com/threadpulse/internal/upvotes/repositories"
	upvoteservice "github.com/threadpulse/internal/upvotes/services"
)

func Routes() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.Use(middleware.ErrorHandler())

	AuthRepo := authrepo.NewAuthRepo()
	AuthService := authservice.NewAuthService(AuthRepo)
	AuthHandler := authhandler.NewAuthHandler(AuthService)

	authHandler := r.Group("/auth")
	{
		authHandler.POST("/register", AuthHandler.RegisterHandler)
		authHandler.POST("/login", AuthHandler.Login)
	}

	RepliesRepo := repliesrepo.NewRepliesRepo()
	repliesService := repliesservice.NewRepliesService(RepliesRepo)
	repliesHandler := replieshandler.NewRepliesHandler(repliesService)

	//threads
	ThreadRepo := threadrepo.NewThreadRepo()
	ThreadService := threadservice.NewThreadsService(ThreadRepo)
	ThreadHandler := threadhandler.NewThreadHandler(ThreadService)

	// upvotes
	UpvoteRepo := upvoterepo.NewUpvotesRepository()
	UpvoteWorker := upvoterepo.NewUpvoteWorker(UpvoteRepo)
	UpvoteWorker.Start()
	UpvoteService := upvoteservice.NewUpvoteService(UpvoteRepo, UpvoteWorker)
	UpvoteHandler := upvotehandler.NewUpvoteHandler(UpvoteService)

	Protected := r.Group("/private", middleware.Miiddleware())
	{
		Protected.POST("/thread", middleware.RateLimiter(5, time.Minute), ThreadHandler.CreateThreadHandler)
		Protected.PATCH("/thread/:id", middleware.RateLimiter(5, time.Minute), ThreadHandler.UpdateThreadHandler)
		Protected.DELETE("/thread/:id", ThreadHandler.DeleteThreadHandler)

		//replies
		Protected.POST("/thread/:id/reply", repliesHandler.CreateRepliesHandler)
		Protected.PATCH("/thread/reply/:id", repliesHandler.UpdateRepliesHandler)
		Protected.DELETE("/thread/reply/:id", repliesHandler.DeleteReplyHandler)

		//upvotes
		Protected.POST("/thread/:id/upvote", UpvoteHandler.Upvote)

	}
	Public := r.Group("/public")
	{
		Public.GET("/threads", ThreadHandler.GetAllThreadHandler)
		Public.GET("/thread/:id", ThreadHandler.GetThreadByIdHandler)
		Public.GET("/thread/:id/replies", repliesHandler.GetAllRepliesHandler)
		Public.GET("/thread/:id/upvotes", UpvoteHandler.GetAllUpvotes)
		Public.GET("/thread/hot", ThreadHandler.GetHotThreads)
	}

	r.Run(":8080")

}
