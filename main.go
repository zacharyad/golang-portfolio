package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Booking struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	UUID      string `json:"uuid"`
	StartTime string `json:"start_time"`
	RoomName  string `json:"room_name"`
	GroupSize byte   `json:"group_size"`
}

var dummyBookings = []Booking{
	{
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
		GroupSize: 5,
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
		GroupSize: 3,
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
		GroupSize: 8,
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
		GroupSize: 4,
	},
	{
		Name:  "Charlie Wilson",
		Email: "charlie.wilson@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
		GroupSize: 2,
	},
	{
		Name:  "Billy Griller",
		Email: "Bill.grill@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
		GroupSize: 8,
	},
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	fmt.Println("env var for api key: ", os.Getenv("FH_API_KEY"))

	// Create a new engine for HTML templates
	engine := html.NewFileSystem(http.Dir("./views"), ".html")

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Route for the root path
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	// API route for fetching bookings
	app.Get("/api/bookings", handleBookings)

	// Serve static files for CSS and JavaScript
	app.Static("/static", "./static")

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

func handleBookings(c *fiber.Ctx) error {
	return c.JSON(dummyBookings)
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/template/html/v2"
// )

// type Booking struct {
// 	Name      string `json:"name"`
// 	Email     string `json:"email"`
// 	Phone     string `json:"phone"`
// 	UUID      string `json:"uuid"`
// 	StartTime string `json:"start_time"`
// 	RoomName  string `json:"room_name"`
// 	GroupSize byte   `json:"group_size"`
// }

// var Bookings []Booking

// func main() {
// 	// Create a new engine for HTML templates
// 	engine := html.NewFileSystem(http.Dir("./views"), ".html")

// 	// Create a new Fiber instance
// 	app := fiber.New(fiber.Config{
// 		Views: engine,
// 	})

// 	// Route for the root path
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.Render("index", fiber.Map{})
// 	})

// 	// API route for AJAX searches
// 	app.Get("/api/bookings", handleBookings)

// 	// Serve static files for CSS and JavaScript
// 	app.Static("/static", "./static")

// 	// Start the server
// 	log.Fatal(app.Listen(":3000"))
// }
// func handleBookings(c *fiber.Ctx) error {
// 	// This is a mock function. Replace it with actual API call and data processing
// 	bookings := fetchBookings()
// 	return c.JSON(bookings)
// }

// func fetchBookings() []Booking {
// 	baseAPIURL := "https://fareharbor.com/api/external/v1/companies/lockedmanhattan"
// 	client := resty.New()
// 	availabilityPk := []string{"1234", "4321", "1423"}
// 	var bookings []Booking

// 	for _, availPK := range availabilityPk {
// 		apiAllBookingsReqString := fmt.Sprintf("%s/availabilies/%s/bookings", baseAPIURL, availPK)
// 		// call out to fh api
// 		resp, err := client.R().Get(apiAllBookingsReqString)
// 		if err != nil {
// 			log.Printf("Error calling out to get all bookings for pk: %v \n ", availPK)
// 			return nil
// 		}

// 		err = json.Unmarshal(resp.Body(), &bookings)
// 		if err != nil {
// 			log.Fatalf("Failed to parse API response: %v", err)
// 		}

// 		// add to bookings slice
// 		bookings = append(bookings, Booking{})
// 	}

// 	for _, booking := range bookings {
// 		apiDetailedBookingReqString := fmt.Sprintf("%s/bookings/%s/", baseAPIURL, booking.UUID)
// 		// call out to fh api
// 		resp, err := client.R().Get(apiDetailedBookingReqString)

// 		if err != nil {
// 			return nil
// 		}

// 		fmt.Println(resp.Body())

// 		detailedBooking := Booking{}

// 		err = json.Unmarshal(resp.Body(), &detailedBooking)
// 		if err != nil {
// 			log.Fatalf("Failed to parse API response: %v", err)
// 		}

// 		// add completed bookings to the bookings slice
// 		bookings = append(bookings, detailedBooking)
// 	}

// 	processedBookings := make([]Booking, len(bookings))
// 	for i, booking := range bookings {
// 		processedBookings[i] = Booking{
// 			Name:      booking.Name,
// 			Email:     booking.Email,
// 			Phone:     booking.Phone,
// 			UUID:      booking.UUID,
// 			StartTime: booking.StartTime,
// 			RoomName:  booking.RoomName,
// 			GroupSize: booking.GroupSize,
// 		}
// 	}

// 	return processedBookings
// }
