package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
	"github.com/rojanmagar2001/gotodo/internal/infrastructure/jsonstore"
)

func runSeedCommand(args []string) error {
	fs := flag.NewFlagSet("seed", flag.ContinueOnError)

	var (
		n         = fs.Int("n", 1000, "number of todos to generate")
		seed      = fs.Int64("seed", 42, "random seed (deterministic datasets)")
		file      = fs.String("file", "", "path to todos.json (default ~/.gotodo/todos.json)")
		overwrite = fs.Bool("overwrite", false, "overwrite existing file (DANGEROUS)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	dbPath, err := defaultDBPath(*file)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return err
	}

	if *overwrite {
		_ = os.Remove(dbPath)
	}

	repo := jsonstore.NewRepository(dbPath)
	rng := rand.New(rand.NewSource(*seed))
	now := time.Now().UTC()

	// Generate and store in batches (fast + avoids huge memory usage)
	const batchSize = 500
	created := 0
	ctx := context.Background()

	for created < *n {
		batch := min(batchSize, *n-created)

		for i := 0; i < batch; i++ {
			td, err := genTodo(rng, now, created+i)
			if err != nil {
				return err
			}
			if err := repo.Create(ctx, td); err != nil {
				return err
			}
		}

		created += batch
		fmt.Fprintf(os.Stderr, "\rSeeded %d/%d", created, *n)
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
	return nil
}

func defaultDBPath(p string) (string, error) {
	if p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gotodo", "todos.json"), nil
}

func genTodo(rng *rand.Rand, now time.Time, i int) (todo.Todo, error) {
	// Deterministic ID: t000001, t000002...
	id := todo.TodoID(fmt.Sprintf("t%06d", i+1))

	titleStr := randomTitle(rng)
	title, err := todo.NewTitle(titleStr)
	if err != nil {
		return todo.Todo{}, err
	}

	priority := randomPriority(rng)

	tags := todo.NewTags(randomTags(rng))

	var due *todo.DueDate
	if rng.Float64() < 0.65 { // 65% have due dates
		// due between -10 and +30 days from now
		days := rng.Intn(41) - 10
		iso := now.AddDate(0, 0, days).Format("2006-01-02")
		dd, err := todo.ParseDueDate(iso)
		if err != nil {
			return todo.Todo{}, err
		}
		due = &dd
	}

	createdAt := now.Add(-time.Duration(rng.Intn(60*24)) * time.Hour) // up to ~60 days old

	td, events, err := todo.NewTodo(todo.NewTodoParams{
		ID:       id,
		Title:    title,
		Priority: priority,
		Tags:     tags,
		DueDate:  due,
		Now:      createdAt,
	})
	_ = events // seeding doesn’t need publishing
	if err != nil {
		return todo.Todo{}, err
	}

	// Random status distribution:
	// ~70% active, ~20% done, ~8% archived, ~2% deleted
	p := rng.Float64()
	switch {
	case p < 0.20:
		// done
		td2, _, err := td.Complete(createdAt.Add(2 * time.Hour))
		if err == nil {
			td = td2
		}
	case p < 0.28:
		// archived (must be done first according to your rules)
		td2, _, err := td.Complete(createdAt.Add(2 * time.Hour))
		if err == nil {
			td3, _, err := td2.Archive(createdAt.Add(3 * time.Hour))
			if err == nil {
				td = td3
			}
		}
	case p < 0.30:
		// deleted (soft delete)
		td2, _, err := td.SoftDelete(createdAt.Add(4 * time.Hour))
		if err == nil {
			td = td2
		}
	default:
		// active
	}

	// Ensure UpdatedAt looks realistic
	if td.UpdatedAt.Before(td.CreatedAt) {
		td.UpdatedAt = td.CreatedAt
	}
	return td, nil
}

func randomPriority(rng *rand.Rand) todo.Priority {
	x := rng.Float64()
	switch {
	case x < 0.20:
		return todo.PriorityHigh
	case x < 0.60:
		return todo.PriorityMedium
	default:
		return todo.PriorityLow
	}
}

func randomTags(rng *rand.Rand) []string {
	pool := []string{
		"work", "home", "study", "go", "fitness", "health", "finance",
		"errands", "reading", "project", "chores", "shopping", "dev",
	}
	// 0..4 tags
	k := rng.Intn(5)
	out := make([]string, 0, k)
	for i := 0; i < k; i++ {
		out = append(out, pool[rng.Intn(len(pool))])
	}
	return out
}

func randomTitle(rng *rand.Rand) string {
	verbs := []string{"Write", "Review", "Fix", "Plan", "Refactor", "Learn", "Organize", "Ship", "Draft", "Test"}
	nouns := []string{"report", "feature", "module", "notes", "budget", "workout", "PR", "meeting", "docs", "todo app"}
	extras := []string{"today", "this week", "ASAP", "with tests", "cleanly", "v2", "for release", "before lunch", "tonight", ""}

	v := verbs[rng.Intn(len(verbs))]
	n := nouns[rng.Intn(len(nouns))]
	e := extras[rng.Intn(len(extras))]

	if e == "" {
		return fmt.Sprintf("%s %s", v, n)
	}
	return fmt.Sprintf("%s %s %s", v, n, e)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// (optional) so your seed command can reuse ListSpec defaults later
var _ = ports.ListSpec{}
