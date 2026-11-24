package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	gapi "tunelport-svc-totems/internal/infrastructure/grpc"
	"tunelport-svc-totems/internal/infrastructure/persistence/postgres"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const databaseDSNSecretID = "totems-database-dsn"

func main() {
	db, err := connectDatabase()
	if err != nil {
		log.Fatalf("database connectivity check failed: %v", err)
	}

	repo := postgres.NewPostgresTotemRepository(db)
	totemServer := gapi.NewServer(repo)

	address := resolveAddress()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()
	gapi.RegisterTotemServiceServer(grpcServer, totemServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("totem.TotemService", healthpb.HealthCheckResponse_SERVING)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gRPC server...")
		healthServer.Shutdown()
		grpcServer.GracefulStop()
	}()

	log.Printf("Database connectivity OK. Starting gRPC service on %s...", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

func connectDatabase() (*gorm.DB, error) {
	dsn, err := resolveDSN()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return db, nil
}

func resolveDSN() (string, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = os.Getenv("GCP_PROJECT")
	}

	var secretValue string
	if projectID != "" {
		value, err := readDatabaseDSNSecret(projectID)
		if err != nil {
			return "", err
		}
		secretValue = value
	}

	return resolveDatabaseDSN(secretValue, nil, os.Getenv("DATABASE_DSN"))
}

func readDatabaseDSNSecret(projectID string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("create Secret Manager client: %w", err)
	}
	defer client.Close()

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, databaseDSNSecretID),
	})
	if status.Code(err) == codes.NotFound {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("access Secret Manager secret %q: %w", databaseDSNSecretID, err)
	}
	if result == nil || result.GetPayload() == nil {
		return "", nil
	}

	return string(result.GetPayload().GetData()), nil
}

func resolveDatabaseDSN(secretValue string, secretErr error, environmentValue string) (string, error) {
	if secretErr != nil {
		return "", secretErr
	}
	if strings.TrimSpace(secretValue) != "" {
		return secretValue, nil
	}
	if environmentValue != "" {
		return environmentValue, nil
	}

	return "", errors.New("DATABASE_DSN must be set in Secret Manager secret totems-database-dsn or in the DATABASE_DSN environment variable")
}

func resolveAddress() string {
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		log.Fatal("SERVER_ADDRESS must be set")
	}

	return address
}
