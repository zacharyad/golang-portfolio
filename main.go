package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type ApiResponse struct {
	Data []Item `json:"data"`
}

type Item struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
	UUID  int    `json:"userid"`
}

var items []Item

func main() {
	// Create a new engine for HTML templates
	engine := html.NewFileSystem(http.Dir("./views"), ".html")

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Fetch items from external API on startup
	fetchItems()

	// Route for the root path
	app.Get("/", handleSearch)

	// API route for AJAX searches
	app.Get("/api/search", handleAPISearch)

	// Serve static files for CSS and JavaScript
	app.Static("/static", "./static")

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

func fetchItems() {
	client := resty.New()

	// Query the external API
	resp, err := client.R().Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		log.Fatalf("Failed to query external API: %v", err)
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(resp.Body(), &apiResponse.Data)

	if err != nil {
		log.Fatalf("Failed to parse API response: %v", err)
	}

	items = apiResponse.Data
	fmt.Printf("Fetched %d items from external API\n", len(items))
}

func handleSearch(c *fiber.Ctx) error {
	query := c.Query("search")

	var results []Item
	if query != "" {
		results = performSearch(query)
	} else {
		results = items // Show all items if no search query
	}

	return c.Render("index", fiber.Map{
		"Query": query,
		"Items": results,
	})
}

func handleAPISearch(c *fiber.Ctx) error {
	query := c.Query("search")
	results := performSearch(query)

	fmt.Println(results)
	return c.JSON(results)
}

func performSearch(query string) []Item {
	query = strings.ToLower(query)
	var results []Item

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Title), query) ||
			fmt.Sprintf("%d", item.ID) == query ||
			fmt.Sprintf("%d", item.UUID) == query {
			results = append(results, item)
		}
	}

	return results
}

type Booking struct {
	name       string // contact.name on bookings array
	email      string // contact.email on bookings array
	phone      string // contact.email on bookings array
	uuid       string // uuid on bookings array
	start_time string // availability.start_time on booking req
	room_name  string // availability.item.name on booking req
	group_size byte   // availability.capacity on booking req
}

func fetchBookingsForToday() ([]Booking, error) {
	// get all bookings
	// range over to create Booking for each item found and firing off another api request to get specific booking info
	// bookin info for each will include: start_time, room_name, and group_size.

	user1 := Booking{
		name:       "zach Droge",
		email:      "zachEmail@email.com",
		phone:      "7853417421",
		uuid:       "1234uuid-8932-78976932-674326732",
		start_time: "2023-08-18T03:30:00-0600",
		room_name:  "The Witching Hour",
		group_size: 8,
	}

	returnArray := []Booking{}

	returnArray = append(returnArray, user1)

	if len(returnArray) == 0 {
		return nil, fmt.Errorf("error in constructing Today's Users array")
	}

	return returnArray, nil
}

// baseURL := "https://fareharbor.com/api/external/v1/companies/<shortname>"

// firstReq := baseURL + "/availabilities/<availability.pk>/bookings/"

//baseURL + "/bookings/<Booking.uuid>/"
