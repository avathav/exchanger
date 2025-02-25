# Exchanger
## How to Run with Docker

1. **Build the Docker image:**
   ```bash
   docker build -t exchanger .
   ```

2. **Run the Docker container (exposing port 8080):**
   ```bash
   docker run -p 8080:8080 -e OPENEXCHANGE_API_KEY=YOUR_API_KEY exchanger
   ```

3. **Access the service** at:
   ```
   http://localhost:8080
   ```

Replace `YOUR_API_KEY` with your actual OpenExchangeRates API key.

## Available Endpoints

1. **GET** `/rates?currencies=USD,EUR,GBP`
2. **GET** `/exchange?from=CRYPTO1&to=CRYPTO2&amount=X`

