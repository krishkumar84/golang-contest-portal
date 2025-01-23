# GOLANG Contest Portal Backend

A high-performance, scalable online judge system built with **Go**, **MongoDB**, and **Docker**. This system supports programming contests, user authentication, problem management, and test case validation.
A contest platform..
---

## 🚀 Features

- **High Performance**: Built with Go for maximum efficiency and concurrent request handling.
- **Scalable Architecture**: Microservices-based design with Docker containerization.
- **Secure Authentication**: JWT-based authentication with role-based access control.
- **Database Optimization**: MongoDB aggregation pipelines for efficient data retrieval.
- **CI/CD Pipeline**: Automated deployment using GitHub Actions.
- **SSL/TLS Support**: Secure communication with Certbot integration.
- **Load Balancing**: Nginx reverse proxy for better request distribution.

---

## 🏗 Architecture

### Technology Stack
- **Backend**: Go 1.23
- **Database**: MongoDB
- **Containerization**: Docker & Docker Compose
- **Web Server**: Nginx
- **SSL**: Certbot
- **CI/CD**: GitHub Actions

---

## 🔒 Authentication & Authorization

- **JWT-based authentication** for secure user sessions.
- **Role-based access control** for user and admin privileges.
- **Secure cookie management**.
- Protected routes with middleware.

---

## 📡 API Endpoints

### **Authentication**
- `POST /api/signup` - Register a new user.
- `POST /api/login` - User login.

### **Contests**
- `POST /api/contest` - Create a new contest.
- `GET /api/contest` - Retrieve all contests.
- `GET /api/contest/{id}` - Retrieve contest details by ID.
- `PUT /api/contest/{id}` - Update contest information.
- `DELETE /api/contest/{id}` - Delete a contest.

### **Questions**
- `POST /api/question` - Create a new question.
- `GET /api/question/{id}` - Retrieve question details by ID.
- `PUT /api/question/{id}` - Update question details.
- `POST /api/contest/{id}/question` - Add a question to a contest.
- `DELETE /api/contest/{contestId}/question/{questionId}` - Remove a question from a contest.

### **Test Cases**
- `POST /api/testcase` - Create a new test case.
- `PUT /api/testcase/{id}` - Update an existing test case.
- `POST /api/question/{id}/testcase` - Add a test case to a question.
- `DELETE /api/question/{questionId}/testcase/{testCaseId}` - Remove a test case from a question.

---

## 🚀 Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.23+ (for local development)
- MongoDB (for local development)

### Running Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/krishkumar84/bdcoe-golang-portal.git
   cd bdcoe-golang-portal
   cp config/example.yaml config/local.yaml
   docker-compose up -d
2. Config.yaml:
```bash   
env: "development"
DatabaseURL: "mongodb://localhost:27017"
DatabaseName: "bdcoe_portal"
JwtSecret: "your-secret-key"
```

💪 Performance & Scalability

### MongoDB Optimization
- Efficient aggregation pipelines for complex queries.
- Indexed collections for faster lookups.
- Proper document structure for optimal data retrieval.

### Go Performance Features
- Goroutines for concurrent request handling.
- Efficient connection pooling.
- Context-based timeout management.
- Structured logging with `slog`.

### Docker & Infrastructure
- Multi-stage builds for smaller and efficient images.
- Container orchestration with Docker Compose.
- Nginx load balancing for optimized request distribution.
- Automated SSL certificate renewal with Certbot.

## 🔐 Security Features
- JWT token-based authentication.
- Password hashing for secure credentials.
- HTTPS enforcement with SSL/TLS.
- Secure cookie configuration.
- Role-based access control.
- Request validation and input sanitization.

## 📦 Deployment
This project uses **GitHub Actions** for CI/CD. Automated deployment targets an EC2 instance with Docker-based container orchestration.

## 📝 Contributing
1. Fork the repository.
2. Create a feature branch:
    ```bash
    git checkout -b feature/AmazingFeature
    ```
3. Commit your changes:
    ```bash
    git commit -m "Add some AmazingFeature"
    ```
4. Push to the branch:
    ```bash
    git push origin feature/AmazingFeature
    ```
5. Open a Pull Request.

## 📄 License
This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for more details.



http_server:
  address: ":8000"

