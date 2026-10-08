package feedback

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const MaxContentLength = 5000

var (
	ErrInvalidCourseID = errors.New("courseId must not be empty and must be at most 120 characters")
	ErrInvalidContent  = errors.New("content must not be empty and must be at most 5000 characters")
)

type CreateInput struct {
	CourseID string `json:"courseId"`
	Content  string `json:"content"`
}

type Feedback struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"courseId" bson:"courseId"`
	Content   string    `json:"content" bson:"content"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type Store interface {
	Create(context.Context, CreateInput) (Feedback, error)
}

func Validate(input CreateInput) error {
	courseID := strings.TrimSpace(input.CourseID)
	if courseID == "" || utf8.RuneCountInString(courseID) > 120 {
		return ErrInvalidCourseID
	}

	content := strings.TrimSpace(input.Content)
	if content == "" || utf8.RuneCountInString(content) > MaxContentLength {
		return ErrInvalidContent
	}

	return nil
}

func Normalize(input CreateInput) CreateInput {
	input.CourseID = strings.TrimSpace(input.CourseID)
	input.Content = strings.TrimSpace(input.Content)
	return input
}
