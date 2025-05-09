# PR Reviewer

PR Reviewer is a tool designed to assist developers in reviewing pull requests efficiently. It provides features to analyze code changes, suggest improvements, and ensure code quality.

## 🔧 Features

- ✅ Modular provider structure (Bitbucket implemented, GitHub-ready)
- ✅ Admin-authenticated API key system
- ✅ Secure `POST /analyze` endpoint - to pull pr, analyze it using AI and post comments
- ✅ Asynchronous background processing
- ✅ Pulls full PR diff, maps to line numbers, sends context to AI
- ✅ Receives structured AI response and posts comments to Bitbucket
- ✅ Built-in token hashing, permission handling, and session support

## 💡 Technologies Used

- **Golang** – web backend + modular architecture
- **Ollama** – locally hosted AI (swap in any LLM endpoint easily)
- **Bitbucket API** – PR diff fetching and commenting
- **PostgreSQL** – session & key storage
- **bcrypt + SHA256** – secure token hashing

---

## Installation

1. Clone the repository:
    ```
    git clone https://github.com/yourusername/pr-reviewer.git
    ```
2. Navigate to the project directory:
    ```
    cd pr-reviewer
    ```

## Usage

1. Install and Run Ollama - [Ollama website](https://ollama.com/download) 
2. Pull preferred model (tested with `qwen2.5-coder:7b`)
3. Update Env variable `OLLAMA_MODEL` in `docker-compose.yaml` file  

4. Run the tool, by using script (runs docker-compose):
    ```
    bin/up
    ```
5. Follow Swagger docs [http://localhost:8080/swagger](http://localhost:8080/swagger) for Creating Admin user and making call to `/analyze` endpoint


## ⚙️ Configurable + Extensible
- Add new providers (e.g. GitHub) by implementing the `Provider` interface
- Add new AI Client by implementing the `AIClient` interface
- API access secured by hashed tokens (stored via SHA256 or bcrypt)

## 🧪 Current Limitations
- Line number mapping from AI responses may be imprecise (improvements underway!!)
- Currently supports only Bitbucket repos - By using Username/App Password
- Error handling/logging can be expanded

## 🛠️ Roadmap
- [ ] Improve AI line-number accuracy  
- [ ] Block multiple processing on the same PR (unless PR was changed)
- [ ] Add more verbose error handling
- [ ] Track analyzing process in DB (started/done/error)
 
## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
