package handlers

import (
	"cars-viewer/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var carData models.CarData

func LoadCarsData() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 3) // Channel to receive errors from goroutines

	// Fetch Car Models
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchCarModels(); err != nil {
			errChan <- err
		}
	}()

	// Fetch Manufacturers
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchManufacturers(); err != nil {
			errChan <- err
		}
	}()

	// Fetch Categories
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchCategories(); err != nil {
			errChan <- err
		}
	}()

	// Wait for all API calls to finish
	wg.Wait()

	// Close error channel and check for errors
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Populate Manufacturer and Category data for each Car (no need to go async for this)
	for i, car := range carData.CarModels {
		for _, manufacturer := range carData.Manufacturers {
			if car.ManufacturerID == manufacturer.ID {
				carData.CarModels[i].Manufacturer = manufacturer
				break
			}
		}
		for _, category := range carData.Categories {
			if car.CategoryID == category.ID {
				carData.CarModels[i].Category = category
				break
			}
		}
	}
	return nil
}

func fetchCarModels() error {
	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		return fmt.Errorf("failed to fetch car models: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&carData.CarModels); err != nil {
		return fmt.Errorf("failed to decode car models: %v", err)
	}
	return nil
}

func fetchManufacturers() error {
	resp, err := http.Get("http://localhost:3000/api/manufacturers")
	if err != nil {
		return fmt.Errorf("failed to fetch manufacturers: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&carData.Manufacturers); err != nil {
		return fmt.Errorf("failed to decode manufacturers: %v", err)
	}
	return nil
}

func fetchCategories() error {
	resp, err := http.Get("http://localhost:3000/api/categories")
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&carData.Categories); err != nil {
		return fmt.Errorf("failed to decode categories: %v", err)
	}
	return nil
}

func generateYearRange(minYear, maxYear int) []int {
	var years []int
	for year := minYear; year <= maxYear; year++ {
		years = append(years, year)
	}
	return years
}

func CarDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/car/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		//http.Error(w, "Invalid car ID", http.StatusBadRequest)
		renderCarNotFoundPage(w)
		return
	}

	// Fetch car data from the API
	apiURL := fmt.Sprintf("http://localhost:3000/api/models/%d", id)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		renderCarNotFoundPage(w)
		return
	}
	defer resp.Body.Close()

	var selectedCar models.Car
	if err := json.NewDecoder(resp.Body).Decode(&selectedCar); err != nil {
		HandleServerError(w, err)
		return
	}

	// 🔥 Link the Manufacturer data to the selectedCar
	for _, manufacturer := range carData.Manufacturers {
		if selectedCar.ManufacturerID == manufacturer.ID {
			selectedCar.Manufacturer = manufacturer
			break
		}
	}

	// 🔥 Link the Category data to the selectedCar
	for _, category := range carData.Categories {
		if selectedCar.CategoryID == category.ID {
			selectedCar.Category = category
			break
		}
	}

	// Update "recently viewed cars" cookie
	const maxViewedCars = 7
	var viewedCars []string

	if cookie, err := r.Cookie("viewedCars"); err == nil {
		viewedCars = strings.Split(cookie.Value, ",")
	}

	carIDStr := strconv.Itoa(selectedCar.ID)
	// Remove the car if it's already in the list to prevent duplication
	for i, id := range viewedCars {
		if id == carIDStr {
			viewedCars = append(viewedCars[:i], viewedCars[i+1:]...)
			break
		}
	}

	// Add the car to the front of the list
	viewedCars = append([]string{carIDStr}, viewedCars...)

	// Limit the list to maxViewedCars
	if len(viewedCars) > maxViewedCars {
		viewedCars = viewedCars[:maxViewedCars]
	}

	// Save the updated list to a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "viewedCars",
		Value: strings.Join(viewedCars, ","),
		Path:  "/",
	})

	tmpl, err := template.ParseFiles("templates/car_details.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, selectedCar)
}

func HandleServerError(w http.ResponseWriter, err error) {
	log.Println("Internal Server Error:", err) // Log the error for debugging
	w.WriteHeader(http.StatusInternalServerError)

	tmpl, tmplErr := template.ParseFiles("templates/error.html")
	if tmplErr != nil {
		log.Println("Error loading error template:", tmplErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				HandleServerError(w, fmt.Errorf("%v", err))
			}
		}()
		next(w, r)
	}
}

func renderCarNotFoundPage(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("templates/car_not_found.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing 'car_not_found.html' template: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func GetManufacturerNameByID(id int) string {
	for _, manufacturer := range carData.Manufacturers {
		if manufacturer.ID == id {
			return manufacturer.Name
		}
	}
	return ""
}

func filterCars(r *http.Request) []models.Car {
	searchQuery := strings.ToLower(r.URL.Query().Get("search"))
	selectedCategory := strings.ToLower(r.URL.Query().Get("category"))
	selectedManufacturer := strings.ToLower(r.URL.Query().Get("manufacturer"))
	minYear, _ := strconv.Atoi(r.URL.Query().Get("minYear"))
	maxYear, _ := strconv.Atoi(r.URL.Query().Get("maxYear"))

	// Fetch all car models from API
	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		log.Println("Error fetching cars from API:", err)
		return nil
	}
	defer resp.Body.Close()

	var allCars []models.Car
	if err := json.NewDecoder(resp.Body).Decode(&allCars); err != nil {
		log.Println("Error decoding car data:", err)
		return nil
	}

	var filteredCars []models.Car
	for _, car := range carData.CarModels {
		// Fetch full Manufacturer data
		for _, manufacturer := range carData.Manufacturers {
			if manufacturer.ID == car.ManufacturerID {
				car.Manufacturer = manufacturer
				break
			}
		}
		// Fetch full Category data
		for _, category := range carData.Categories {
			if category.ID == car.CategoryID {
				car.Category = category
				break
			}
		}
		// Apply filters
		if searchQuery != "" && !strings.Contains(strings.ToLower(car.Name), searchQuery) {
			continue
		}
		if selectedCategory != "" && strings.ToLower(car.Category.Name) != selectedCategory {
			continue
		}
		if selectedManufacturer != "" && strings.ToLower(car.Manufacturer.Name) != selectedManufacturer {
			continue
		}
		if (minYear != 0 && car.Year < minYear) || (maxYear != 0 && car.Year > maxYear) {
			continue
		}
		filteredCars = append(filteredCars, car)
	}
	return filteredCars
}

func HomeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modelYears := generateYearRange(1990, 2025)
		filteredCars := filterCars(r)
		preferredManufacturer := ""
		var recentlyViewedCars []models.Car

		if cookie, err := r.Cookie("viewedCars"); err == nil {
			viewedIDs := strings.Split(cookie.Value, ",")
			for _, idStr := range viewedIDs {
				id, _ := strconv.Atoi(idStr)
				for _, car := range carData.CarModels {
					if car.ID == id {
						recentlyViewedCars = append(recentlyViewedCars, car)
						break
					}
				}
			}
		}

		//panic("Simulated server error")

		if cookie, err := r.Cookie("preferredManufacturer"); err == nil {
			preferredManufacturer = cookie.Value
		}

		data := models.TemplateData{
			CarData: models.CarData{
				Manufacturers: carData.Manufacturers,
				Categories:    carData.Categories,
				CarModels:     filteredCars,
			},
			ModelYears:            modelYears,
			PreferredManufacturer: preferredManufacturer,
		}
		err := tmpl.Execute(w, struct {
			models.TemplateData
			RecentlyViewedCars []models.Car
		}{
			TemplateData:       data,
			RecentlyViewedCars: recentlyViewedCars,
		})
		if err != nil {
			HandleServerError(w, err)
			//http.Error(w, fmt.Sprintf("Error executing template: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	carIDs := r.URL.Query()["carIDs"]
	var selectedCars []models.Car
	for _, idStr := range carIDs {
		id, _ := strconv.Atoi(idStr)
		for _, car := range carData.CarModels {
			if car.ID == id {
				selectedCars = append(selectedCars, car)
				break
			}
		}
	}

	tmpl, err := template.ParseFiles("templates/compare.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, selectedCars)
}
