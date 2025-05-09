package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"pr-reviewer/internal/api"
	"pr-reviewer/internal/database"

	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
)

func main() {

	fmt.Println(`
	 _______________________________________________________________________________
	|  ____________________  __________            .__                              |
	|  \______   \______   \ \______   \ _______  _|__| ______  _  __ ___________   |
	|	|     ___/|       _/  |       _// __ \  \/ /  |/ __ \ \/ \/ // __ \_  __ \  | 
	|	|    |    |    |   \  |    |   \  ___/\   /|  \  ___/\     /\  ___/|  | \/  |
	|	|____|    |____|_  /  |____|_  /\___  >\_/ |__|\___  >\/\_/  \___  >__|     |
	|	    			 \/          \/     \/             \/            \/         |
	|_______________________________________________________________________________|				    
	`)
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all CPU cores

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
		},
	})

	logrus.WithField("version", os.Getenv("BITBUCKET_COMMIT_SHORT")).Info("Starting up.")

	db, err := database.New()
	if err != nil {
		logrus.WithError(err).Fatal("Error verifying database.")
	}

	logrus.Info("Database is ready to use.")

	//Creating new router
	router, err := api.NewRouter(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	const addr = "0.0.0.0:8080"
	server := http.Server{
		Handler:           router,
		Addr:              addr,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	logrus.WithField("address", addr).Info("Server is ready to listen.")
	//Starting server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed.")
	}
}
