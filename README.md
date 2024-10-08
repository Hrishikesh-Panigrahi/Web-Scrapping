﻿# Web-Scrapping
A Web Scrapping built with Go, and the Gin framework.

## Technologies Used
- Backend: Go
- Frontend: HTML, Tailwind
- Containerization: Docker

## Requirements

- Go 1.22+
- A SQL database (e.g., PostgreSQL, MySQL, SQLite)
- Git

## Installation
1. **Clone the repository:**
```sh
   git clone github.com/Hrishikesh-Panigrahi/Web-Scrapping
   cd GoCMS
```
2. **Install dependencies:**
```sh
   go mod tidy
   ```
And you are all set

## Docker
1. **Pull the Docker Image:**
```sh
docker pull hrishikeshpanigrahi025/web_scrapping 
```
2. **Run the Docker container:**
```sh
docker run -p 8000:8000 hrishikeshpanigrahi025/web_scrapping
```
And you are all set

3. **Open your browser and navigate to http://localhost:8000.**


## Run Locally
To run the project locally, you have 3 options:

### 1. Launch Debugger
   - Open your project in Visual Studio Code.
   - Set breakpoints as needed.
   - Launch the debugger by pressing F5 or by selecting Run > Start Debugging from the menu.
   - For Delve Debugger, visit the [delve GitHub repository](https://github.com/go-delve/delve)

### 2. Run Air
   - Ensure you have Air installed for live reloading.
   - Start Air by running the following command in your terminal:
```sh
   air
```
   - For more information, visit the [Air GitHub repository](https://github.com/air-verse/air).

### 3. Run go run main.go Command
   - Open your terminal.
   - Navigate to the project directory.
   - Run the following command to start the application:
```sh
   go run main.go
```
