
# ✨ Car Inventory Web Application ✨

This is a Go-based web application for managing and viewing a collection of cars with their details such as manufacturer, category, specifications, etc. The application supports search, filtering, comparison, and recently viewed cars.

## 🚀 Features
- Display cars with their details.
- Search and filter cars by name, category, manufacturer, and year.
- Compare multiple cars
- Recently viewed cars tracking using cookies.
- Manufacturer preference tracking using cookies.
- Graceful error handling for non-existent cars.


## 📦 Dependencies
- `encoding/json`
- `html/template`
- `net/http`
- `os`
- `strconv`
- `strings`
- `fmt`

## 🔨 Installation API
1. Install node.js:
```sh
    npm install
```
2. Run Car API:
```sh
    npm start
```

## 🔨 Installation
1. Clone the repository:
```sh
    git clone <repository-url>
```
2. Navigate to the project directory:
```sh
    cd cars
```
3. Install dependencies (if any).

4. Ensure `data.json` exists in the `api/` directory.

5. Run the application:
```sh
    go run .
```

6. Open your browser and go to:
```
    http://localhost:8080
```

### 🛑 Stop the Server

-   Press `Ctrl + C` if running the server in VS Code or a terminal.
    
-   If using a separate terminal window, simply close it.
    

----------

## 📄 Available Endpoints
- `/` - Home page (listing cars with filtering options).
- `/car/{id}` - Display details for a specific car.
- `/compare?carIDs={id1},{id2}` - Compare multiple cars.

## 📌 Future Improvements
- Implement database support (e.g., SQLite, PostgreSQL).
- Add user authentication.
- Enhance front-end design with frameworks like TailwindCSS or Bootstrap.
- Implement API endpoints for CRUD operations.

✨ This project is designed to be **simple, efficient, and easy to use**. Contributions and feedback are always welcome! 🚀