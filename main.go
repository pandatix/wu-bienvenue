package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/ctfer-io/go-ctfd/api"
)

const (
	username = "PandatiX"
	password = "CTFer.io au BreizhCTF 2025 ?"

	url = "https://ctf.bzh"

	flagRegex = "BZHCTF{.*}"
)

func main() {
	ctx := context.Background()

	fmt.Printf("[+] Getting nonce and session\n")
	nonce, session, err := api.GetNonceAndSession(url)
	if err != nil {
		fatal(err)
	}

	client := api.NewClient(url, nonce, session, "")

	fmt.Printf("[+] Logging in with user %s\n", username)
	if err := client.Login(&api.LoginParams{
		Name:     username,
		Password: password,
	}, api.WithContext(ctx)); err != nil {
		fatal(err)
	}

	fmt.Printf("[+] Getting challenges\n")
	challs, err := client.GetChallenges(&api.GetChallengesParams{}, api.WithContext(ctx))
	if err != nil {
		fatal(err)
	}
	if len(challs) != 1 {
		log.Fatal("challenges are either not published yet or the \"bienvenue\" challenge is not the only one visible by dependency relationships")
	}

	fmt.Printf("[+] Extracting flag\n")
	re := regexp.MustCompile(flagRegex)
	flag := re.FindString(challs[0].Description)
	fmt.Printf("    Flag: %s\n", flag)

	fmt.Printf("[+] Submitting flag to challenge %d", challs[0].ID)
	if _, err := client.PostChallengesAttempt(&api.PostChallengesAttemptParams{
		ChallengeID: challs[0].ID,
		Submission:  flag,
	}, api.WithContext(ctx)); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Printf("[-] Fatal error: %s", err)
	os.Exit(1)
}
