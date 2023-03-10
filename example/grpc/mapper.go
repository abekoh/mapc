package grpc

import "github.com/abekoh/mapc/example/domain"

// MapUserToUser maps User to User.
// This function is generated by mapc.
// DO NOT EDIT this function.
func MapUserToUser(x domain.User) User {
	return User{
		Name: x.Name,
		// state:
		// sizeCache:
		// unknownFields:
		// Id:
	}
}

// MapTaskToTask maps Task to Task.
// This function is generated by mapc.
// DO NOT EDIT this function.
func MapTaskToTask(x domain.Task) Task {
	return Task{
		Title:       x.Title,
		Description: x.Description,
		// state:
		// sizeCache:
		// unknownFields:
		// Id:
		// StoryPoint:
		// RegisteredAt:
		// User:
		// Subtasks:
	}
}

// MapSubTaskToSubTask maps SubTask to SubTask.
// This function is generated by mapc.
// DO NOT EDIT this function.
func MapSubTaskToSubTask(x domain.SubTask) SubTask {
	return SubTask{
		Title:       x.Title,
		Description: x.Description,
		// state:
		// sizeCache:
		// unknownFields:
		// Id:
		// RegisteredAt:
		// User:
	}
}
