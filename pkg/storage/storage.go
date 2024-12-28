package storage
import "github.com/krishkumar84/bdcoe-golang-portal/pkg/types"


type Storage interface {
	// CreateUser(name string, email string, age int)(int64, error)
	// GetUserById(id int64)(types.User, error) 
	// GetAllUsers()([]types.User, error)
	CreateUser(name string, email string, password string, studentId string) (string, error)
	GetUserByEmail(email string) (*types.User, error)
}
