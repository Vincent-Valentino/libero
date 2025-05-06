# Libero ML - Data Scraping Microservice

This microservice is responsible for scraping sports data (standings, results, player stats, etc.) based on user preferences stored in the main Libero application.

## Setup

1.  **Create a Virtual Environment:**
    ```bash
    python -m venv venv
    ```
    *   On Windows (cmd/powershell):
        ```powershell
        .\venv\Scripts\activate
        ```
    *   On macOS/Linux (bash/zsh):
        ```bash
        source venv/bin/activate
        ```

2.  **Install Dependencies:**
    ```bash
    pip install -r requirements.txt
    ```

## Running the Service (Development)

Use `uvicorn` to run the FastAPI application. The `--reload` flag enables auto-reloading when code changes are detected.

```bash
uvicorn main:app --reload --port 8001
```

The API will be available at `http://127.0.0.1:8001`. You can access the interactive API documentation (Swagger UI) at `http://127.0.0.1:8001/docs`.

## Endpoints

*   `GET /status`: Returns the current status of the scraper.
*   `POST /scrape`: Triggers a background scraping task.

## TODO

*   Implement database connection to read preferences from the main application's DB.
*   Implement actual web scraping logic.
*   Implement data storage for scraped results (likely a separate DB).
*   Implement scheduling (e.g., using APScheduler).
*   Add proper configuration management (e.g., environment variables).
*   Add authentication/security for endpoints if necessary.