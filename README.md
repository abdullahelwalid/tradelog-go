# TradeLog Backend

## Overview
TradeLog is a trading journal application designed to help traders efficiently log, track, and analyze their trades. The system provides features such as trade entry management, performance analytics, and insights to enhance trading strategies. The backend is built using **Golang**, providing a fast, efficient, and scalable solution with a well-structured architecture.

## Tech Stack
- **Language:** Golang
- **Framework:** Standard Go modules
- **Database:** PostgreSQL (or any other configured database)
- **API Protocol:** RESTful APIs
- **Authentication:** JWT-based authentication
- **Containerization:** Docker
- **Deployment:** AWS ECS
- **Task Runner:** `air` for live reloading in development

## Project Structure
The backend follows a modular structure for maintainability and scalability.

```
/tradelog-backend
│── cmd/main/       # Main application entry point
│── pkg/            # Contains core business logic modules
│── bin/            # Compiled binaries
│── .github/workflows # GitHub Actions for CI/CD
│── .air.toml       # Configuration for `air` live reload
│── Dockerfile      # Docker container setup
│── docker-compose.yml # Docker Compose configuration
│── go.mod          # Go module dependencies
│── go.sum          # Dependency checksums
│── task_definition.json # AWS ECS task definition
```

## Backend Installation & Setup
### Prerequisites
- Go 1.22+
- Docker & Docker Compose
- PostgreSQL database
- `.env` file with necessary environment variables (API keys, DB credentials, etc.)

### Running Locally
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/tradelog-backend.git
   cd tradelog-backend
   ```
2. Install dependencies:
   ```sh
   go mod download
   ```
3. Start the backend with Air for live reloading:
   ```sh
   air
   ```
   or manually with:
   ```sh
   go run cmd/main/main.go
   ```
4. Run with Docker:
   ```sh
   docker-compose up --build
   ```

## API Documentation
The backend exposes RESTful endpoints to interact with TradeLog. Below are key API routes:

### Authentication
#### **User Signup**
- **Endpoint:** `POST /api/auth/signup`
- **Payload:**
  ```json
  {
    "email": "user@example.com",
    "password": "SecurePass123!",
    "firstName": "John",
    "lastName": "Doe"
  }
  ```
- **Response:**
  ```json
  {
    "message": "User registered successfully. Please verify your email."
  }
  ```

#### **User Login**
- **Endpoint:** `POST /api/auth/login`
- **Payload:**
  ```json
  {
    "email": "user@example.com",
    "password": "SecurePass123!"
  }
  ```
- **Response:**
  ```json
  {
    "token": "jwt_token_here"
  }
  ```

### Trade Management
#### **Create Trade**
- **Endpoint:** `POST /api/trades`
- **Headers:** `Authorization: Bearer <token>`
- **Payload:**
  ```json
  {
    "symbol": "AAPL",
    "entryPrice": 150.25,
    "exitPrice": 155.00,
    "quantity": 10,
    "tradeType": "long",
    "date": "2025-03-04"
  }
  ```
- **Response:**
  ```json
  {
    "id": "trade123",
    "message": "Trade logged successfully."
  }
  ```

### Deployment (AWS ECS)
1. Build and push Docker image:
   ```sh
   docker build -t tradelog-backend .
   docker tag tradelog-backend:latest <aws_account_id>.dkr.ecr.<region>.amazonaws.com/tradelog-backend:latest
   docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/tradelog-backend:latest
   ```
2. Deploy to ECS using task definition:
   ```sh
   aws ecs update-service --cluster tradelog-cluster --service tradelog-service --force-new-deployment
   ```

## Contributing
- Fork the repository
- Create a feature branch
- Submit a pull request

## License
This project is licensed under the MIT License.
