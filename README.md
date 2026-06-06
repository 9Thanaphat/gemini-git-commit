# gemini-git-commit

A lightweight, fast Git CLI tool written in Go that automatically generates 3 distinct conventional commit message options using the power of Gemini 3.5 Flash.

---

## 🌟 Features

* **AI-Powered:** Analyzes your staged changes (`git diff --cached`) using the advanced **Gemini 3.5 Flash** model.
* **Conventional Commits Standard:** Ensures all generated options strictly follow the `<type>: <description>` convention (e.g., `feat:`, `fix:`, `docs:`, `refactor:`).
* **Context Aware:** Allows you to append additional context or custom instructions from the command line to guide the AI.
* **Pure Go:** Compiles into a single lightweight binary, making it extremely fast and easy to distribute.

---

## 🛠️ Installation

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

## 🔑 Environment Setup

Before using the tool, you need to acquire a Gemini API key from Google AI Studio and expose it as an environment variable in your terminal.

```bash
export GEMINI_API_KEY="your_actual_gemini_api_key_here"
