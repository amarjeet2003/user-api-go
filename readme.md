# User Registration HTTP Server Application

This application provides REST APIs for user registration, updating user details, and searching for users. It is developed in Golang with MySQL as the database.

## Functional Requirements

1. **Create User**: Allows creating users with fields such as first name, last name, unique username, and date of birth.
2. **Update User**: Users can update all fields, including the username if it doesn't violate integrity constraints.
3. **Search User**: Provides a list of users, allowing searching by name with the default sorting order by age.

## Tech Stack

- Golang
- MySQL

## Project Structure

The project follows a basic MVC pattern:

- `controllers/`: Contains controller logic for handling HTTP requests.
- `models/`: Defines the data models used in the application.
- `repository/`: Implements the repository layer for database interactions.
- `routes/`: Defines HTTP routes and their corresponding handlers.
- `main.go`: Initiates the HTTP server and sets up routes.

## Installation and Setup

1. Clone the repository:

git clone https://github.com/amarjeet2003/user-api-go.git

2. Make sure you have a database named "users_db" in your MySQL.

3. Create a `.env` file with the following environment variables:

DB_USERNAME=<your_database_username>
DB_PASSWORD=<your_database_password>
DB_HOST=<your_database_host>
DB_PORT=<your_database_port>
DB_NAME=<your_database_name>


4. Run the application:

go run .

## Endpoints

- `POST /users/create`: Create a new user.
- `PUT /users/update/{id}`: Update user details by ID.
- `GET /users/search?name={search_query}`: Search for users by name.