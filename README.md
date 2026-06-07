# gemini-git-commit

A lightweight, fast Git CLI tool written in Go that automatically generates 3 distinct conventional commit message options using the power of Gemini 3.5 Flash.

---

## Features

* **AI-Powered:** Analyzes your staged changes (`git diff --cached`) using the advanced **Gemini 3.5 Flash** model.
* **Conventional Commits Standard:** Ensures all generated options strictly follow the `<type>: <description>` convention (e.g., `feat:`, `fix:`, `docs:`, `refactor:`).
* **Smart Configuration:** Securely saves your API key locally. Set it once and forget it, no need to export environment variables for every new terminal session.
* **Colorized UI:** Clean and readable terminal output with color-coded feedback for a better developer experience.
* **Context Aware:** Allows you to append additional context or custom instructions from the command line to guide the AI.
* **Pure Go:** Compiles into a single lightweight binary, making it extremely fast and easy to distribute.

---

### Usage
Stage your changes using `git add` and then fire up the tool:
```bash
git add .
aic

<img width="683" height="192" alt="image" src="https://github.com/user-attachments/assets/550c59f6-743f-44ab-bf3b-b1f7f91d814f" />

---

## Installation

### Prerequisites
* [Go](https://go.dev/doc/install) (version 1.21 or higher recommended)
* [Git](https://git-scm.com/) installed and configured on your machine

### Setup Steps

1.  **Clone the repository:**
    ```bash
    git clone git@github.com:9Thanaphat/gemini-git-commit.git
    cd gemini-git-commit
    ```

2.  **Download and tidy up Go modules:**
    ```bash
    go mod tidy
    ```

3.  **Build the executable binary:**
    ```bash
    go build -o aic main.go
    ```

4.  **Move the binary to your system PATH:**
    Moving it to `/usr/local/bin/` allows you to run the tool globally from any directory on your system.
    ```bash
    sudo mv aic /usr/local/bin/
    ```

---

## Configuration (API Key)

Before using the tool, you need a Gemini API key from [Google AI Studio](https://aistudio.google.com/). You no longer need to manually export environment variables every time. You can set up your key using one of the following methods:

### Method 1: The Config Flag (Recommended)
You can set or update your API key globally at any time using the `-config` flag. The tool will save it securely in your home directory (`~/.lazy-commit-config.json`).
```bash
aic -config "your_actual_gemini_api_key_here"
