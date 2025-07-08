package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// User represents a user in the system
type User struct {
	ID       int
	Name     string
	Email    string
	Password string // This should be hashed in production
}

// UserService handles user operations
type UserService struct {
	users []User
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		users: make([]User, 0),
	}
}

// AddUser adds a new user to the service
func (us *UserService) AddUser(name, email, password string) {
	// TODO: Add validation
	user := User{
		ID:       len(us.users) + 1,
		Name:     name,
		Email:    email,
		Password: password, // Should hash this
	}
	us.users = append(us.users, user)
}

// GetUser retrieves a user by ID
func (us *UserService) GetUser(id int) *User {
	for i := 0; i < len(us.users); i++ {
		if us.users[i].ID == id {
			return &us.users[i]
		}
	}
	return nil
}

// GetUserByEmail finds a user by email
func (us *UserService) GetUserByEmail(email string) *User {
	for _, user := range us.users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

// DeleteUser removes a user by ID
func (us *UserService) DeleteUser(id int) bool {
	for i, user := range us.users {
		if user.ID == id {
			us.users = append(us.users[:i], us.users[i+1:]...)
			return true
		}
	}
	return false
}

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	// Simple email validation - not comprehensive
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ProcessUserInput handles user input from command line
func ProcessUserInput(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: program <command> [args...]")
		return
	}

	command := args[1]
	userService := NewUserService()

	switch command {
	case "add":
		if len(args) < 5 {
			fmt.Println("Usage: add <name> <email> <password>")
			return
		}
		name := args[2]
		email := args[3]
		password := args[4]

		if !ValidateEmail(email) {
			fmt.Println("Invalid email format")
			return
		}

		userService.AddUser(name, email, password)
		fmt.Println("User added successfully")

	case "get":
		if len(args) < 3 {
			fmt.Println("Usage: get <id>")
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid ID format")
			return
		}

		user := userService.GetUser(id)
		if user == nil {
			fmt.Println("User not found")
			return
		}

		fmt.Printf("User: %+v\n", *user)

	case "delete":
		if len(args) < 3 {
			fmt.Println("Usage: delete <id>")
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid ID format")
			return
		}

		if userService.DeleteUser(id) {
			fmt.Println("User deleted successfully")
		} else {
			fmt.Println("User not found")
		}

	default:
		fmt.Println("Unknown command:", command)
	}
}

func main() {
	ProcessUserInput(os.Args)
}
