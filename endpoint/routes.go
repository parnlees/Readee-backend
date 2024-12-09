package endpoint

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RegisterRoutes(app *fiber.App) {
	// User
	app.Get("/users", GetUsers)
	app.Get("/users/:userId", GetUserSpecific)
	app.Post("/getUserInfoByEmail", GetUserInfoByEmail)
	app.Post("/createUser", CreateUser)
	app.Post("/checkUser", CheckUser)
	app.Patch("/user/resetPassword/:userId", ResetPassword)

	app.Get("/genres", GetGenres)
	app.Get("/genres/:genre_id", GetGenreByID)
	app.Post("/createGenres", CreateGenres)

	app.Patch("/user/edit/:userId", EditUser)

	// Book
	app.Post("/createBook", CreateBook)
	app.Get("/getBook/:BookId", GetBookSpecific)
	app.Get("/getBooks", GetBooks)
	app.Patch("/editBook/:BookId", EditBook)
	app.Delete("/deleteBook/:BookId", DeleteBook)
	app.Get("/getBookByUser/:OwnerId", GetBookByOwnerId)

	app.Get("/books/recommendations/:userId", getBooksForUser)

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

	app.Get("/get_review_rating/:giverId/:receiverId", GetReviewRating)

	//History
	app.Get("/history/:userId", GetHistory)
	app.Get("/tradeCount/:userId", TradeCount)

	//Rating
	app.Get("/getRating/:userId", GetRatingByUserId)
	app.Get("/getAverageRate/:userId", GetAverageRatingByUserId)

	//Room
	app.Post("/createRoom/:senderId/:receiverId", CreateRoom)
	app.Get("/getRoomId/:senderId/:receiverId", GetRoomId)

	//Message
	app.Post("/createMessage", CreateMessage)

	// app.Post("/messages", CreateMessage)
	app.Get("/getAllMessage/:roomId", GetMessagesByRoomId)
	app.Get("/getAllChat/:userId", GetAllChatByUserId)
	app.Get("/chat/:roomId", websocket.New(Chat))
	app.Get("/rooms/:roomId/messages", GetMessagesByRoomId)
	app.Post("/uploadImage", func(c *fiber.Ctx) error {
		// Get the file from the request
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error()) // If no file is uploaded or error occurs
		}

		// Upload the image to Azure and get the URL
		url, err := UploadImage(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error()) // If image upload to Azure fails
		}

		// Return the URL of the uploaded image
		return c.JSON(fiber.Map{"url": url}) // Send the URL as response
	})

	//ads banner
	app.Get("/getALlAds", GetAllAds)
	

	app.Post("/login", Login)
}
