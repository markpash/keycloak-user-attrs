package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v11"
	"github.com/joho/godotenv"
)

func newKeycloakClient(ctx context.Context) (gocloak.GoCloak, *gocloak.JWT, error) {
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminRealm := os.Getenv("ADMIN_REALM")
	keycloakServer := os.Getenv("SERVER")

	// Options are for compat with keycloak 17 default paths
	client := gocloak.NewClient(keycloakServer,
		gocloak.SetAuthRealms("realms"),
		gocloak.SetAuthAdminRealms("admin/realms"))

	token, err := client.LoginAdmin(ctx, adminUsername, adminPassword, adminRealm)
	if err != nil {
		return nil, nil, err
	}

	return client, token, nil
}

func main() {
	godotenv.Load()

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "realm or username not provided")
		os.Exit(1)
	}

	realm, username := os.Args[1], os.Args[2]

	getUsersParams := gocloak.GetUsersParams{
		Username: &username,
		Enabled:  gocloak.BoolP(true),
		Exact:    gocloak.BoolP(true),
		Max:      gocloak.IntP(1),
	}

	ctx := context.Background()

	client, token, err := newKeycloakClient(ctx)
	if err != nil {
		panic(err)
	}

	users, err := client.GetUsers(ctx, token.AccessToken, realm, getUsersParams)
	if err != nil {
		panic(err)
	}

	if len(users) == 0 {
		fmt.Fprintln(os.Stderr, "user not found")
		os.Exit(1)
	}

	if users[0].Attributes == nil {
		fmt.Fprintln(os.Stderr, "user has no attributes")
		os.Exit(1)
	}

	jsonUserAttrs, err := json.Marshal(users[0].Attributes)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, string(jsonUserAttrs))
}
