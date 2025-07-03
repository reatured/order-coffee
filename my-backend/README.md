# My Backend (Go)

This is a simple Go backend server that responds with "Hello, World!" at the root endpoint.

## Project Plan: Coffee Ordering App

This project will be a simple coffee ordering app that allows your friends to order coffee online.

### Features
- **Frontend:** Built with React, providing a user-friendly interface for selecting and ordering coffee.
- **Backend:** Built with Go, serving API endpoints for coffee types and order submission.
- **Order Notification:** When a user submits an order, the backend will send the order details to your email address using SMTP.
- **Data Storage:** Coffee types and orders will be stored in a local JSON file (no external database required).
- **Deployment:** The backend will be deployed on [Railway](https://railway.app/).

### Planned API Endpoints
- `GET /coffees` — List all available coffee types
- `POST /order` — Submit a new coffee order (triggers an email to you)
- (Optional) `GET /orders` — List all orders (for admin/friends to see)

### Email Sending Requirements
- An email account with SMTP access (e.g., Gmail)
- (Optional) An App Password if using Gmail with 2FA
- Email credentials should be stored securely using environment variables
- The backend will use Go's `net/smtp` package or a third-party library to send emails

### Order Flow
1. User submits the order form on the frontend
2. Frontend sends a POST request to the backend `/order` endpoint
3. Backend receives the order, saves it to the local JSON file (optional), and sends an email to you with the order details
4. Backend responds to the frontend with a success message

### Next Steps
1. Set up the backend API endpoints in Go
2. Create a local JSON file to store coffee types and orders
3. Implement email sending functionality in the backend
4. Build the React frontend to interact with the backend
5. Deploy the backend to Railway
6. Connect the frontend to the deployed backend

## How to Run

1. Make sure you have Go installed (https://go.dev/dl/).
2. In this directory, run:

```sh
go run main.go
```

3. Visit [http://localhost:8080](http://localhost:8080) in your browser.

---

Feel free to modify and expand this project! 