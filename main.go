package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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
}

type APIResponse struct {
	Bookings []struct {
		Contact struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"contact"`
		UUID      string `json:"uuid"`
		CreatedAt string `json:"created_at"`
	} `json:"bookings"`
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
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
	},
	{
		Name:  "Charlie Wilson",
		Email: "charlie.wilson@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Billy Griller",
		Email: "Bill.grill@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
	},
	{
		Name:  "Charlie Wilson",
		Email: "charlie.wilson@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Billy Griller",
		Email: "Bill.grill@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
	},
	{
		Name:  "Charlie Wilson",
		Email: "charlie.wilson@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Billy Griller",
		Email: "Bill.grill@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
	},
	{
		Name:  "Charlie Wilson",
		Email: "charlie.wilson@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Billy Griller",
		Email: "Bill.grill@example.com",

		UUID:      "o5p4q-r3s2t-13579",
		StartTime: "2024-07-19T11:15:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Enter Sequence",
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
	},
	{
		Name:  "John Doe",
		Email: "johnD@example.com",

		UUID:      "b1f5e-d3c2a-98765",
		StartTime: "2024-07-15T14:00:00Z",
		RoomName:  "Pet Project",
	},
	{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",

		UUID:      "a2b3c-e4f5d-12345",
		StartTime: "2024-07-16T10:30:00Z",
		RoomName:  "The Witching Hour",
	},
	{
		Name:  "Bob Johnson",
		Email: "bob.johnson@example.com",

		UUID:      "c6d7e-f8g9h-56789",
		StartTime: "2024-07-17T09:00:00Z",
		RoomName:  "The Dinner Party",
	},
	{
		Name:      "Alice Brown",
		Email:     "alice.brown@example.com",
		UUID:      "i9j8k-l7m6n-24680",
		StartTime: "2024-07-18T13:45:00Z",
		RoomName:  "Kingdom Quest",
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

	log.Fatal(app.Listen(":8080"))
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
			log.Println("issue getting all items")
			return err
		}
	}

	allItemsAvails, err := getAllAvailabilitiesConcurrently(cached_items)
	if err != nil {
		log.Println("issue getting all availabilities")
		return err
	}

	// now use GetConcurrentBookings(context.Background(), allItemsAvails) to get all the bookings based on availability PKs
	// add to the return of thie request
	return c.JSON(fiber.Map{
		"bookings":       dummyBookings,
		"availabilities": allItemsAvails,
	})
}

func getAllAvailabilitiesConcurrently(items AllItems) (map[string][]string, error) {
	log.Println("Starting getAllAvailabilitiesConcurrently")
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

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

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("X-FareHarbor-API-User", USERKEY_VAL)
	req.Header.Add("X-FareHarbor-API-App", APPKEY_VAL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response struct {
		Items AllItems `json:"items"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
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

func GetConcurrentBookings(ctx context.Context, availabilities map[string][]string) ([]Booking, error) {
	var (
		bookings    []Booking
		errorList   []error
		mu          sync.Mutex
		wg          sync.WaitGroup
		rateLimiter = time.NewTicker(time.Second / 30)
		errorChan   = make(chan error, len(availabilities))
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	APPKEY_VAL := os.Getenv("FH_API_APPKEY")
	USERKEY_VAL := os.Getenv("FH_API_USERKEY")

	for roomName, availabilityPKs := range availabilities {
		for _, pk := range availabilityPKs {
			wg.Add(1)
			go func(room, availabilityPK string) {
				defer wg.Done()

				select {
				case <-rateLimiter.C:
				case <-ctx.Done():
					errorChan <- ctx.Err()
					return
				}

				url := fmt.Sprintf("https://fareharbor.com/api/external/v1/companies/%s/availabilities/%s/bookings/", GetEnvVal("SHORTNAME"), availabilityPK)

				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					errorChan <- fmt.Errorf("failed to create request for %s: %v", availabilityPK, err)
					return
				}

				req.Header.Add("X-FareHarbor-API-User", USERKEY_VAL)
				req.Header.Add("X-FareHarbor-API-App", APPKEY_VAL)

				resp, err := client.Do(req)
				if err != nil {
					errorChan <- fmt.Errorf("failed to send request for %s: %v", availabilityPK, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					errorChan <- fmt.Errorf("unexpected status code for %s: %d", availabilityPK, resp.StatusCode)
					return
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					errorChan <- fmt.Errorf("failed to read response body for %s: %v", availabilityPK, err)
					return
				}

				var apiResp APIResponse
				if err := json.Unmarshal(body, &apiResp); err != nil {
					errorChan <- fmt.Errorf("failed to unmarshal response for %s: %v", availabilityPK, err)
					return
				}

				mu.Lock()
				for _, b := range apiResp.Bookings {
					bookings = append(bookings, Booking{
						Name:      b.Contact.Name,
						Email:     b.Contact.Email,
						UUID:      b.UUID,
						StartTime: b.CreatedAt, // Using CreatedAt as StartTime for now
						RoomName:  room,
					})
				}
				mu.Unlock()

				log.Printf("Successfully processed bookings for %s, %s", room, availabilityPK)
			}(roomName, pk)
		}
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		errorList = append(errorList, err)
	}

	rateLimiter.Stop()

	var finalError error
	if len(errorList) > 0 {
		finalError = fmt.Errorf("encountered %d errors while fetching bookings", len(errorList))
		for _, err := range errorList {
			log.Printf("Error: %v", err)
		}
	}

	return bookings, finalError
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
