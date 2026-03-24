package main

import (
    "context"
    "fmt"
    "time"
)

// --- services now return errors and respect context ---

func getName(ctx context.Context, userID string) (string, error) {
    select {
    case <-time.After(100 * time.Millisecond):
        return "Alice", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func getAge(ctx context.Context, userID string) (int, error) {
    select {
    case <-time.After(150 * time.Millisecond):
        return 30, nil
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}

func getEmail(ctx context.Context, userID string) (string, error) {
    select {
    case <-time.After(120 * time.Millisecond):
        return "alice@example.com", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

// --- result wrapper so goroutines can send errors too ---

type nameResult  struct { val string; err error }
type ageResult   struct { val int;    err error }
type emailResult struct { val string; err error }

// --- main function ---

func fetchUserData(ctx context.Context, userID string) UserProfile {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    nameCh  := make(chan nameResult,  1)
    ageCh   := make(chan ageResult,   1)
    emailCh := make(chan emailResult, 1)

    go func() {
        v, err := getName(ctx, userID)
        nameCh <- nameResult{v, err}
    }()
    go func() {
        v, err := getAge(ctx, userID)
        ageCh <- ageResult{v, err}
    }()
    go func() {
        v, err := getEmail(ctx, userID)
        emailCh <- emailResult{v, err}
    }()

    var profile UserProfile

    // collect 3 results, racing against timeout each time
    for i := 0; i < 3; i++ {
        select {
        case r := <-nameCh:
            if r.err != nil {
                fmt.Println("timed out")
                return UserProfile{}
            }
            profile.Name = r.val

        case r := <-ageCh:
            if r.err != nil {
                fmt.Println("timed out")
                return UserProfile{}
            }
            profile.Age = r.val

        case r := <-emailCh:
            if r.err != nil {
                fmt.Println("timed out")
                return UserProfile{}
            }
            profile.Email = r.val

        case <-ctx.Done():
            fmt.Println("timed out")
            return UserProfile{}
        }
    }

    return profile
}

type UserProfile struct {
    Name  string
    Age   int
    Email string
}

func main() {
    profile := fetchUserData(context.Background(), "user-123")
    fmt.Printf("Name: %s, Age: %d, Email: %s\n", profile.Name, profile.Age, profile.Email)
}