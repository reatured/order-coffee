package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "net/smtp"
    "github.com/joho/godotenv"
)

type Coffee struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Order struct {
    Name     string `json:"name"`
    CoffeeID int    `json:"coffeeId"`
    Notes    string `json:"notes"`
    Email    string `json:"email,omitempty"`
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        h(w, r)
    }
}

func getCoffeesHandler(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadFile("coffees.json")
    if err != nil {
        http.Error(w, "Could not read coffees", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
    var order Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        http.Error(w, "Invalid order", http.StatusBadRequest)
        return
    }

    // Optionally save order to orders.json
    orders := []Order{}
    ordersFile := "orders.json"
    if _, err := os.Stat(ordersFile); err == nil {
        data, _ := ioutil.ReadFile(ordersFile)
        json.Unmarshal(data, &orders)
    }
    orders = append(orders, order)
    ordersData, _ := json.MarshalIndent(orders, "", "  ")
    ioutil.WriteFile(ordersFile, ordersData, 0644)

    // Send email to admin
    emailSent, emailErr := sendOrderEmail(order)

    // Send confirmation email to customer if email is provided
    customerEmailSent := false
    customerEmailErr := error(nil)
    if order.Email != "" {
        customerEmailSent, customerEmailErr = sendConfirmationEmail(order)
    }

    w.Header().Set("Content-Type", "application/json")
    resp := map[string]interface{}{
        "status":    "ok",
        "adminEmailSent": emailSent,
        "customerEmailSent": customerEmailSent,
    }
    if emailErr != nil {
        resp["adminError"] = emailErr.Error()
    }
    if customerEmailErr != nil {
        resp["customerError"] = customerEmailErr.Error()
    }
    json.NewEncoder(w).Encode(resp)
}

func sendOrderEmail(order Order) (bool, error) {
    // Load coffee name
    data, _ := ioutil.ReadFile("coffees.json")
    var coffees []Coffee
    json.Unmarshal(data, &coffees)
    coffeeName := "Unknown"
    for _, c := range coffees {
        if c.ID == order.CoffeeID {
            coffeeName = c.Name
            break
        }
    }

    // Email config (use environment variables in production)
    from := os.Getenv("SMTP_USER")
    pass := os.Getenv("SMTP_PASSWORD")
    to := os.Getenv("CONTACT_EMAIL")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")

    subject := "New Coffee Order"
    body := fmt.Sprintf("Name: %s\nCoffee: %s\nNotes: %s", order.Name, coffeeName, order.Notes)
    if order.Email != "" {
        body += fmt.Sprintf("\nCustomer Email: %s", order.Email)
    }
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, pass, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        log.Println("Failed to send email:", err)
        return false, err
    }
    return true, nil
}

func sendConfirmationEmail(order Order) (bool, error) {
    // Load coffee name
    data, _ := ioutil.ReadFile("coffees.json")
    var coffees []Coffee
    json.Unmarshal(data, &coffees)
    coffeeName := "Unknown"
    for _, c := range coffees {
        if c.ID == order.CoffeeID {
            coffeeName = c.Name
            break
        }
    }

    from := os.Getenv("SMTP_USER")
    pass := os.Getenv("SMTP_PASSWORD")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")

    to := order.Email
    subject := "Your Coffee Order Confirmation"
    body := fmt.Sprintf("Hi %s,\n\nThank you for your order!\n\nOrder Details:\nCoffee: %s\nNotes: %s\n\nWe will process your order soon.\n\nBest regards,\nCoffee Shop", order.Name, coffeeName, order.Notes)
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, pass, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        log.Println("Failed to send confirmation email to customer:", err)
        return false, err
    }
    return true, nil
}

func main() {
    godotenv.Load(".env")
    http.HandleFunc("/coffees", withCORS(getCoffeesHandler))
    http.HandleFunc("/order", withCORS(orderHandler))
    fmt.Println("Server running at http://localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", nil))
}