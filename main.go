package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v33/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	appID         string
	privateKeyPath string
	repoURL       string
	cloneDir      string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tokenclone",
		Short: "A small Golang utility to clone a GitHub repository using Github app credentials.",
		Run: func(cmd *cobra.Command, args []string) {
			if appID == "" || privateKeyPath == "" || repoURL == "" || cloneDir == "" {
				log.Fatalf("All flags --app_id, --pem_path, --repo_url, and --clone_dir are required")
			}

			privateKey, err := os.ReadFile(privateKeyPath)
			if err != nil {
				log.Fatalf("Error reading private key: %v", err)
			}

			signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
			if err != nil {
				log.Fatalf("Error parsing private key: %v", err)
			}

			// Generate JWT
			jwtToken, err := generateJWT(appID, signKey)
			if err != nil {
				log.Fatalf("Error generating JWT: %v", err)
			}

			// Get GitHub client using JWT
			client := getGitHubClient(jwtToken)

			// Get installation ID
			installationID, err := getInstallationID(client)
			if err != nil {
				log.Fatalf("Error fetching installation ID: %v", err)
			}

			// Generate installation token
			installationToken, err := getInstallationToken(client, installationID)
			if err != nil {
				log.Fatalf("Error generating installation token: %v", err)
			}

			// Check repository access
			repo, err := checkRepoAccess(client, installationToken, repoURL)
			if err != nil {
				log.Fatalf("Error accessing repository: %v", err)
			}

			printRepoDetails(repo)

			// Clone the repository
			err = cloneRepo(repoURL, cloneDir, installationToken)
			if err != nil {
				log.Fatalf("Error cloning repository: %v", err)
			}

			fmt.Println("Repository cloned successfully")
		},
	}

	rootCmd.Flags().StringVar(&appID, "app_id", "", "GitHub App ID")
	rootCmd.Flags().StringVar(&privateKeyPath, "pem_path", "", "Path to the GitHub App private key PEM file")
	rootCmd.Flags().StringVar(&repoURL, "repo_url", "", "URL of the repository to clone")
	rootCmd.Flags().StringVar(&cloneDir, "clone_dir", "", "Directory to clone the repository into")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateJWT(appID string, key *rsa.PrivateKey) (string, error) {
	now := time.Now()
	// Create the JWT claims, which includes the registered claims
	claims := jwt.StandardClaims{
		Issuer:    appID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute * 10).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with our private key
	jwtToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func getGitHubClient(jwtToken string) *github.Client {
	// Create a new OAuth2 token source using the JWT token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: jwtToken},
	)

	// Create a new HTTP client using the OAuth2 token source
	tc := oauth2.NewClient(context.Background(), ts)

	// Return a new GitHub client using the HTTP client
	return github.NewClient(tc)
}

func getInstallationID(client *github.Client) (int64, error) {
	ctx := context.Background()
	installations, _, err := client.Apps.ListInstallations(ctx, nil)
	if err != nil {
		return 0, err
	}

	if len(installations) == 0 {
		return 0, fmt.Errorf("no installations found")
	}

	return installations[0].GetID(), nil
}

func getInstallationToken(client *github.Client, installationID int64) (string, error) {
	// Create a new context
	ctx := context.Background()

	// Generate a new installation token
	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", err
	}

	return token.GetToken(), nil
}

func checkRepoAccess(client *github.Client, token, repoURL string) (*github.Repository, error) {
	ctx := context.Background()

	// Extract owner and repo name from URL
	repoParts := strings.Split(strings.TrimPrefix(repoURL, "https://github.com/"), "/")
	if len(repoParts) != 2 {
		return nil, fmt.Errorf("invalid repository URL format: %s", repoURL)
	}
	owner, repo := repoParts[0], strings.TrimSuffix(repoParts[1], ".git")

	// Create a new OAuth2 token source using the installation token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// Create a new HTTP client using the OAuth2 token source
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	// Check if the repository is accessible
	repository, resp, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("Repository %s/%s not found", owner, repo)
		} else {
			log.Printf("Error accessing repository: %v", err)
		}
		return nil, err
	}

	return repository, nil
}

func printRepoDetails(repo *github.Repository) {
	fmt.Printf("Repository Details:\n")
	fmt.Printf("Name: %s\n", repo.GetName())
	fmt.Printf("Full Name: %s\n", repo.GetFullName())
	fmt.Printf("Clone URL: %s\n", repo.GetCloneURL())
}

func cloneRepo(repoURL, cloneDir, token string) error {
	_, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: "x-access-token", // Username is ignored, but must be non-empty
			Password: token,
		},
	})
	return err
}
