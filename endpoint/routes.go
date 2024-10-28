package endpoint

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// User
	app.Get("/users", GetUsers)
	app.Get("/users/:userId", GetUserSpecific)
	app.Post("/createUser", CreateUser)

	app.Get("/genres", GetGenres)
	app.Get("/genres/:genre_id", GetGenreByID)
	app.Post("/createGenres", CreateGenres)

	// Book
	app.Post("/createBook", CreateBook)
	app.Get("/getBook/:BookId", GetBookSpecific)
	app.Get("/getBooks", GetBooks)
	app.Patch("/editBook/:BookId", EditBook)
	app.Delete("/deleteBook/:BookId", DeleteBook)
	app.Get("/getBookByUser/:OwnerId", GetBookByOwnerId)

	// Log
	app.Post("/books/:bookId/like/:userId", LikeBook)
	app.Post("/books/:bookId/unlike/:userId", UnLikeBook)
	app.Get("/getLogs/:liker_id", GetLogsByUserID)
	app.Post("/unlikeLogs/:bookLikeId/:likerId", UnlikeLogs)

	// Genres
	app.Get("/userGenres", GetUserGenres)
	app.Get("/userGenres/:user_user_id", GetUserGenreByUserID)
	app.Post("/createUserGenres", CreateUserGenres)
	app.Put("/userGenre/edit", EditGenre)

	//Match
	app.Get("/getMatches/:userId", GetMatchBook)
	app.Get("/getAllMatches/:matchId", GetMatchById)
	app.Delete("/deleteMatch/:matchId", DeleteMatch)

	//Trade
	app.Post("/trades/:matchId/send-request/:initiatorId", SendTradeRequest)
	app.Post("/trades/:matchId/cancel-request", CancelTradeRequest)
	app.Post("/trades/:matchId/accept", AcceptTradeRequest)
	app.Post("/trades/:matchId/reject", RejectTradeRequest)

	//Review & Rating
	app.Post("/review_rating", SubmitRatingAndReview)
	app.Get("/avgRating/:userId", GetAverageRating)
	app.Get("/reviews/received/:userId", GetReceivedReviewsAndRatings)
	app.Get("/reviews/given/:userId", GetGivenReviewsAndRatingsWithTradedBooks)

	//History
	app.Get("/history/:userId", GetHistory)

	//Rating
	app.Get("/getRating/:userId", GetRatingByUserId)

	//Room
	app.Post("/createRoom/:senderId/:receiverId", CreateRoom)
	app.Get("/getRoomId/:senderId/:receiverId", GetRoomId)

	//Message
	app.Post("/createMessage", CreateMessage)
	app.Get("/getAllMessage/:roomId", GetMessagesByRoomId)
	app.Get("/getAllChat/:userId", GetAllChatByUserId)

	app.Post("/login", Login)
}
