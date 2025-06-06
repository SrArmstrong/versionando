package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func ConnectFirestore() {
	ctx := context.Background()

	// Para producción con archivo de credenciales
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		log.Fatalf("Error inicializando Firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error inicializando Firestore: %v", err)
	}

	FirestoreClient = client
	log.Println("Conectado a Firestore")
}

func CloseFirestore() {
	if FirestoreClient != nil {
		FirestoreClient.Close()
	}
}
