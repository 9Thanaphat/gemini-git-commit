package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/genai"
)

func main() {

	// 1. Get additional context from arguments (if any)
	userContext := ""
	if len(os.Args) > 1 {
		// Join all arguments after 'aic' into a single string
		userContext = strings.Join(os.Args[1:], " ")
	}

	// 2. Run 'git diff --cached' to get staged changes
	cmd := exec.Command("git", "diff", "--cached")
	diffBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error: Failed to run 'git diff' (are you in a Git repository?)\n%v", err)
	}

	diffStr := string(diffBytes)
	if strings.TrimSpace(diffStr) == "" {
		fmt.Println("No staged changes found. (Did you forget to run 'git add .'?)\n")
		return
	}

	// 3. Check for API Key in Environment Variables
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: GEMINI_API_KEY not found in environment.\nPlease set it using: export GEMINI_API_KEY='your_api_key'")
	}

	// 4. Connect to Gemini API using the new SDK
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Error: Failed to create GenAI client: %v", err)
	}

	// 5. Build the prompt to require 3 options and append user context (if any)
	basePrompt := `You are an expert developer. Generate 3 different conventional commit message options in English based on the following git diff.

Format MUST be exactly:
1. <type>: <description>
2. <type>: <description>
3. <type>: <description>

Allowed Types:
- feat: (new feature)
- fix: (bug fix)
- docs: (documentation, README)
- style: (formatting, no logic change)
- refactor: (code restructuring, no bug fix or new feature)
- test: (add/edit tests)
- chore: (updating tasks, dependencies, etc.)

Strict Rules:
- Output ONLY the 3 options as a numbered list.
- The descriptions MUST be in English.
- Do NOT include markdown backticks or any explanations.`

	// Append additional context if the user provided arguments
	if userContext != "" {
		basePrompt += fmt.Sprintf("\n\nAdditional User Context/Instruction: %s", userContext)
	}

	prompt := fmt.Sprintf("%s\n\nDiff:\n%s", basePrompt, diffStr)

	fmt.Println("Analyzing code and generating commit message options...")

	// 6. Send request to AI
	resp, err := client.Models.GenerateContent(ctx, "gemini-3.5-flash", genai.Text(prompt), nil)
	if err != nil {
		log.Fatalf("\nError: AI processing failed: %v", err)
	}

	// 7. Extract and display the AI response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		commitMsg := resp.Candidates[0].Content.Parts[0].Text

		fmt.Printf("\nSuggested Commit Messages:\n")
		// \033[1;32m sets the terminal text color to green
		fmt.Printf("\n\033[1;32m%s\033[0m\n\n", strings.TrimSpace(commitMsg))
	} else {
		fmt.Println("\nNo response received from AI.")
	}
}
