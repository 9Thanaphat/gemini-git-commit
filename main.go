package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/genai"
)

// Define ANSI color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[1;31m"
	ColorGreen  = "\033[1;32m"
	ColorYellow = "\033[1;33m"
	ColorCyan   = "\033[1;36m"
)

type Config struct {
	GeminiAPIKey string `json:"gemini_api_key"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".lazy-commit-config.json"), nil
}

// Function to explicitly save key to the config file
func saveAPIKey(key string) {
	configPath, err := getConfigPath()
	if err != nil {
		log.Fatalf("%sError: Cannot get home directory: %v%s", ColorRed, err, ColorReset)
	}

	config := Config{GeminiAPIKey: key}
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("%sError: Failed to marshal config: %v%s", ColorRed, err, ColorReset)
	}

	err = os.WriteFile(configPath, configData, 0600)
	if err != nil {
		log.Fatalf("%sError: Failed to write config file: %v%s", ColorRed, err, ColorReset)
	}
	fmt.Printf("%s[OK] API Key successfully saved/updated at: %s%s\n", ColorGreen, configPath, ColorReset)
}

// Smart load that reads from Env or Config file
func loadAPIKey() string {
	// Priority 1: Check System Environment Variable first
	if envKey := os.Getenv("GEMINI_API_KEY"); envKey != "" {
		return envKey
	}

	configPath, err := getConfigPath()
	if err != nil {
		return ""
	}

	// Priority 2: Check Config file in Home Directory
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.ReadFile(configPath)
		if err == nil {
			var config Config
			if err := json.Unmarshal(file, &config); err == nil && config.GeminiAPIKey != "" {
				return config.GeminiAPIKey
			}
		}
	}
	return ""
}

func main() {
	// Remove timestamp from log output for cleaner CLI UX
	log.SetFlags(0)

	// 1. Define Command Line Flags
	configFlag := flag.String("config", "", "Set or update your Gemini API Key")
	flag.Parse() // Parse flags to separate them from remaining arguments

	// 2. If user provides -config flag, save it and exit
	if *configFlag != "" {
		saveAPIKey(strings.TrimSpace(*configFlag))
		return // Terminate the program early after configuration
	}

	// 3. Get additional context from remaining arguments (Args that are not flags)
	userContext := strings.Join(flag.Args(), " ")

	// 4. Run 'git diff --cached' to get staged changes
	cmd := exec.Command("git", "diff", "--cached")
	diffBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("%sError: Failed to run 'git diff' (are you in a Git repository?)%s\n%v", ColorRed, ColorReset, err)
	}

	diffStr := string(diffBytes)
	if strings.TrimSpace(diffStr) == "" {
		fmt.Printf("%sNo staged changes found. (Did you forget to run 'git add .'?)\n%s", ColorYellow, ColorReset)
		return
	}

	// 5. Load API Key
	apiKey := loadAPIKey()
	if apiKey == "" {
		log.Fatalf("%sError: GEMINI_API_KEY not found.\nPlease set it using: aic -config \"YOUR_API_KEY\"%s", ColorRed, ColorReset)
	}

	// 6. Connect to Gemini API using the new SDK
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("%sError: Failed to create GenAI client: %v%s", ColorRed, err, ColorReset)
	}

	// 7. Build the prompt
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

	if userContext != "" {
		basePrompt += fmt.Sprintf("\n\nAdditional User Context/Instruction: %s", userContext)
	}

	prompt := fmt.Sprintf("%s\n\nDiff:\n%s", basePrompt, diffStr)

	fmt.Printf("%sAnalyzing code and generating commit message options...%s\n", ColorCyan, ColorReset)

	// 8. Send request to AI
	resp, err := client.Models.GenerateContent(ctx, "gemini-3.5-flash", genai.Text(prompt), nil)
	if err != nil {
		log.Fatalf("\n%sError: AI processing failed: %v%s\n", ColorRed, err, ColorReset)
	}

	// 9. Extract and display the AI response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		commitMsg := resp.Candidates[0].Content.Parts[0].Text

		fmt.Printf("\n%sSuggested Commit Messages:%s\n", ColorCyan, ColorReset)
		fmt.Printf("\n%s%s%s\n\n", ColorGreen, strings.TrimSpace(commitMsg), ColorReset)
	} else {
		fmt.Printf("\n%sNo response received from AI.%s\n", ColorYellow, ColorReset)
	}
}
