package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Booking struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	UUID      string `json:"uuid"`
	StartTime string `json:"start_time"`
	RoomName  string `json:"room_name"`
	GroupSize byte   `json:"group_size"`
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

var dummyBookings = []Booking{
	{
		Name:      "John Doe",
		Email:     "johnD@example.com",
		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
		GroupSize: 5,
	},
}

var cached_items AllItems
var today int

func main() {
	app_INIT()

	APPKEY := GetEnvVal("APPKEY")
	fmt.Println(APPKEY)

	engine := html.NewFileSystem(http.Dir("./views"), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Get("/api/bookings", handleBookings)
	app.Static("/static", "./static")

	log.Fatal(app.Listen(":3000"))
}

func app_INIT() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func GetEnvVal(envkey string) string {
	BASE_API_ENV := os.Getenv("BASE_API_ENV")
	return os.Getenv(BASE_API_ENV + envkey)
}

func handleBookings(c *fiber.Ctx) error {

	var err error

	if len(cached_items) == 0 || today != time.Now().Day() {
		today = time.Now().Day()
		cached_items, err = GetAllItems()
		if err != nil {
			fmt.Println("issue getting all items")
			return err
		}
	}

	allItemsAvails, err := getAllAvailabilitiesConcurrently(cached_items)
	if err != nil {
		fmt.Println("issue getting all availabilities")
		return err
	}

	return c.JSON(fiber.Map{
		"bookings":       dummyBookings,
		"availabilities": allItemsAvails,
	})
}

func getAllAvailabilitiesConcurrently(items AllItems) (map[string][]string, error) {
	log.Println("Starting getAllAvailabilitiesConcurrently")

	var wg sync.WaitGroup
	var mu sync.Mutex
	allItemsAvails := make(map[string][]string)
	errors := make(chan error, len(items))

	for _, item := range items {
		wg.Add(1)
		go func(item Item) {
			defer wg.Done()
			log.Printf("Starting goroutine for item: %s (PK: %d)", item.Name, item.Pk)

			itemAvails, err := getAllAvailabilitiesForToday(strconv.Itoa(item.Pk))
			if err != nil {
				log.Printf("Error getting avails for item %s (PK: %d): %v", item.Name, item.Pk, err)
				errors <- fmt.Errorf("error getting avails for item %s (PK: %d): %v", item.Name, item.Pk, err)
				return
			}

			mu.Lock()
			allItemsAvails[item.Name] = itemAvails

			mu.Unlock()
			log.Printf("Completed processing for item: %s (PK: %d)", item.Name, item.Pk)
		}(item)
	}

	wg.Wait()
	close(errors)

	var errStrings []string
	for err := range errors {
		errStrings = append(errStrings, err.Error())
	}

	if len(errStrings) > 0 {
		return allItemsAvails, fmt.Errorf("errors occurred: %s", strings.Join(errStrings, "; "))
	}

	log.Println("Successfully completed getAllAvailabilitiesConcurrently")
	return allItemsAvails, nil
}
func getAllAvailabilitiesForToday(itemPk string) ([]string, error) {
	log.Printf("Starting getAllAvailabilitiesForToday for item PK: %s", itemPk)

	url := fmt.Sprintf("https://fareharbor.com/api/external/v1/companies/%s/items/%s/availabilities/date/%s/",
		GetEnvVal("SHORTNAME"), itemPk, todaysDate())

	APPKEY_VAL := os.Getenv("FH_API_APPKEY")
	USERKEY_VAL := os.Getenv("FH_API_USERKEY")

	if APPKEY_VAL == "" || USERKEY_VAL == "" {
		return nil, fmt.Errorf("API keys not set properly")
	}

	log.Printf("Preparing to make request to URL: %s", url)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("X-FareHarbor-API-User", USERKEY_VAL)
	req.Header.Add("X-FareHarbor-API-App", APPKEY_VAL)

	log.Printf("Sending request for item PK: %s", itemPk)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Received response with status code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Avails AllAvails `json:"availabilities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var allAvails []string
	for _, v := range response.Avails {
		allAvails = append(allAvails, strconv.Itoa(v.Pk))
	}

	log.Printf("Successfully processed availabilities for item PK: %s", itemPk)
	return allAvails, nil
}

func GetAllItems() (AllItems, error) {
	url := "https://fareharbor.com/api/external/v1/companies/" + GetEnvVal("SHORTNAME") + "/items/"
	APPKEY_VAL := os.Getenv("FH_API_APPKEY")
	USERKEY_VAL := os.Getenv("FH_API_USERKEY")

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
			continue
		}
		allItems = append(allItems, v)
	}

	return allItems, nil
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
