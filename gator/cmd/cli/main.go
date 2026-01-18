package main

import (
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/config"
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/repl"
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/database"
	"log"
	"os"
	"database/sql"
	"context"
)


type State struct {
	cfg *config.Config
	db *database.Queries
}


func main() {
	// Init cfg
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Init db
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	
	// Init State
	state := &State{
		cfg: cfg,
db: dbQueries,
	}
	
	// Init cmd
	repl := repl.NewREPL(state, "gator")
	registerCommands(repl)
	
	// Handle
	args := os.Args[1:]
	if err := repl.HandleCommand(args); err != nil {
		log.Fatal(err)
	}
}

func registerCommands(r *repl.REPL[*State]) {
	r.Register(repl.Command[*State]{
		Name: "login",
		Desc: "Login user",
		MinArgs: 1,
		Exec: handlerLogin,
	})
	r.Register(repl.Command[*State]{
		Name: "register",
		Desc: "Register new user",
		MinArgs: 1,
		Exec: handlerRegister,
	})
	r.Register(repl.Command[*State]{
		Name: "reset",
		Desc: "Resets all data",
		MinArgs: 0,
		Exec: handlerReset,
	})
	r.Register(repl.Command[*State]{
		Name: "users",
		Desc: "List all users",
		MinArgs: 0,
		Exec: handlerUsers,
	})
	r.Register(repl.Command[*State]{
		Name: "agg",
		Desc: "Run feed loop",
		MinArgs: 1,
		Exec: handlerAgg,
	})
	r.Register(repl.Command[*State]{
		Name: "addfeed",
		Desc: "Add a feed and follow",
		MinArgs: 2,
		Exec: middlewareLoggedIn(handlerAddFeed),
	})
	r.Register(repl.Command[*State]{
		Name: "feeds",
		Desc: "List feeds",
		MinArgs: 0,
		Exec: handlerFeeds,
	})
	r.Register(repl.Command[*State]{
		Name: "follow",
		Desc: "Follow a feed",
		MinArgs: 1,
		Exec: middlewareLoggedIn(handlerFollow),
	})
	r.Register(repl.Command[*State]{
		Name: "following",
		Desc: "List feeds followed",
		MinArgs: 0,
		Exec: middlewareLoggedIn(handlerFollowing),
	})
	r.Register(repl.Command[*State]{
		Name: "unfollow",
		Desc: "Unfollow feed",
		MinArgs: 1,
		Exec: middlewareLoggedIn(handlerUnfollow),
	})
	r.Register(repl.Command[*State]{
		Name: "browse",
		Desc: "Browse Posts",
		MinArgs: 0,
		Exec: middlewareLoggedIn(handlerBrowse),
	})

}

func middlewareLoggedIn(handler func(s *State, args[]string, user database.User) error) func(*State, []string) error {
	return func(s *State, args []string) error {
		userName := s.cfg.CurUser
		user, err := s.db.GetUser(
			context.Background(),
			userName)
		if err != nil {
			return err
		}

		err = handler(s, args, user)

		return err
	}
}
