# Uptime Tracker

A full-stack application designed to monitor, record, and visualize the uptime of configured websites.

## Getting Started

Follow these steps to get the development environment running on your local machine.

### 1. Clone the Repository
```bash
git clone [https://github.com/DrewCoding/Uptime-Tracker.git](https://github.com/DrewCoding/Uptime-Tracker.git)
cd Uptime-Tracker
```

### 2. Start Infrastructure Services
Initialize the database using the docker-compose.yml file
```bash
docker compose up -d
```

### 3. Launch the Backend API
Open a new terminal window and start the Go server.
```bash
cd backend/cmd/api
go run main.go
```

### 4. Launch the Frontend UI
Install the dependencies, and start the development server.
```bash
cd frontend
npm install
npm run dev
```

### Running the Uptime Tracker
The tracker polls the target websites and records their uptime status. To execute a health check run:
```bash
cd backend/cmd/sentinel
go run main.go
```
