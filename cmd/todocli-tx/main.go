package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/achiku/gotodocli/txmodel"
)

// App application global
type App struct {
	DB     txmodel.DBer
	Logger *log.Logger
}

// NewApp creates new app
func NewApp(v bool) (*App, error) {
	constr := os.Getenv("TODOCLI_CONSTR")
	db, err := txmodel.NewDB(constr)
	if err != nil {
		return nil, err
	}
	var logger *log.Logger
	if v {
		logger = log.New(os.Stdout, "[app]", log.Lmicroseconds)
	} else {
		logger = log.New(ioutil.Discard, "[app]", log.Lmicroseconds)
	}
	app := &App{
		DB:     db,
		Logger: logger,
	}
	return app, nil
}

var (
	verbose = flag.Bool("v", false, "verbose")
	command = flag.String("command", "", "sub command")
)

func main() {
	flag.Parse()

	app, err := NewApp(*verbose)
	if err != nil {
		log.Fatal(err)
	}
	switch *command {
	case "":
		app.Logger.Fatal("command has to be specified")
	case "add":
		app.Logger.Print("add is specified")
	case "start":
		app.Logger.Print("start is specified")
	case "stop":
		app.Logger.Print("stop is specified")
	case "done":
		app.Logger.Print("done is specified")
	}
}
