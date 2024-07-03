package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
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
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Pet Project",
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
}

func main() {
	app_INIT()

	APPKEY := GetEnvVal("APPKEY")

	fmt.Println(APPKEY)

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
	//app.Get("/api/rooms", handleItems)

	// Serve static files for CSS and JavaScript
	app.Static("/static", "./static")

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

func app_INIT() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return
	}
}

func GetEnvVal(envkey string) string {
	BASE_API_ENV := os.Getenv("BASE_API_ENV")
	return os.Getenv(BASE_API_ENV + envkey)
}

func handleBookings(c *fiber.Ctx) error {
	allItem, err := GetAllItems()

	if err != nil {
		fmt.Println("issue getting all items")
		return err
	}

	var allItemsAvails []string

	for _, item := range allItem {
		itemsAvails, err := getAllAvailabilitiesForToday(strconv.Itoa(item.Pk))
		if err != nil {
			fmt.Println("issue getting all avails for: ", item.Name)
			return err
		}

		allItemsAvails = append(allItemsAvails, itemsAvails...)
	}

	fmt.Println("all item's Avails: ", allItemsAvails)

	return c.JSON(dummyBookings)
}

func todaysDate() string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())

	if len(day) == 1 {
		day = "0" + day
	}

	if len(month) == 1 {
		month = "0" + month
	}

	return year + "-" + month + "-" + day
}

func getAllAvailabilitiesForToday(itemPk string) ([]string, error) {
	url := "https://fareharbor.com/api/external/v1/companies/" + GetEnvVal("SHORTNAME") + "/items/" + itemPk + "/availabilities/date/" + todaysDate() + "/"
	APPKEY_VAL := os.Getenv("FH_API_" + "APPKEY")
	USERKEY_VAL := os.Getenv("FH_API_" + "USERKEY")

	fmt.Println("URL", url)

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.Header.Add("X-FareHarbor-API-User", USERKEY_VAL)
	req.Header.Add("X-FareHarbor-API-App", APPKEY_VAL)
	req.SetRequestURI(url)

	if err := agent.Parse(); err != nil {
		return nil, err
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("%v %v", errs, code)
	}

	var response struct {
		Avails AllAvails `json:"availabilities"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	fmt.Println("response from all avails api req ", response, code)

	var allAvails []string

	for _, v := range response.Avails {
		allAvails = append(allAvails, strconv.Itoa(v.Pk))

	}
	return allAvails, nil
}

func GetAllItems() (AllItems, error) {
	url := "https://fareharbor.com/api/external/v1/companies/" + GetEnvVal("SHORTNAME") + "/items/"
	APPKEY_VAL := os.Getenv("FH_API_" + "APPKEY")
	USERKEY_VAL := os.Getenv("FH_API_" + "USERKEY")

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.Header.Add("X-FareHarbor-API-User", USERKEY_VAL)
	req.Header.Add("X-FareHarbor-API-App", APPKEY_VAL)
	req.SetRequestURI(url)

	if err := agent.Parse(); err != nil {
		return nil, err
	}

	code, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("%v %v", errs, code)
	}

	var response struct {
		Items AllItems `json:"items"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var allItems AllItems

	for _, v := range response.Items {
		if v.Name == "Gift Card" || v.Name == "Locked Shirts" || v.Name == "Gift Certificate" {
			break
		}

		allItems = append(allItems, v)
	}

	return allItems, nil
}

type Item struct {
	Pk   int    `json:"pk"`
	Name string `json:"name"`
}
type AllItems []Item

type Avail struct {
	Pk int `json:"pk"`
}

type AllAvails []Avail
