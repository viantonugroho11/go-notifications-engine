package config

import (
	"context"
	"errors"

	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type FCM struct {
	ProjectID string `json:"project_id"`
}

func ConnectToFirebase(projectID string) (*firebase.App, error) {
	opt := option.WithCredentialsFile("path/to/refreshToken.json")
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: projectID,
		AuthOverride: &map[string]interface{}{
			"refreshToken": "path/to/refreshToken.json",
		},
		DatabaseURL: "https://<DATABASE_NAME>.firebaseio.com",
		StorageBucket: "gs://<BUCKET_NAME>.appspot.com",
		ServiceAccountID: "serviceAccount:${PROJECT_ID}@appspot.gserviceaccount.com",
	}, opt)
	if err != nil {
		return nil, errors.New("error initializing app")
	}
	return app, nil
}

