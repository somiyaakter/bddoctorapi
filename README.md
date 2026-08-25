# 🇧🇩 BD Doctor Scraper API

A production-ready REST API for scraping, storing, and serving doctor information from Bangladesh-based healthcare sources.

Built with **Go** and **PostgreSQL**.

## 🚀 Live API

**Production API:**

https://bddoctorscraper.vercel.app/api/v1/doctors

Example:

~~~http
GET /api/v1/doctors?page=1&page_size=20
~~~

---

## ✨ Features

- Doctor information scraping
- Doctor profile information
- BMDC registration number
- Doctor specialties
- Designation and workplace
- Doctor image and profile URL
- Multiple chamber information
- Visiting hours
- Appointment phone numbers
- PostgreSQL database
- REST API
- Pagination
- Vercel deployment

---

## 🛠️ Tech Stack

- **Go** — Backend API and scraper
- **PostgreSQL** — Database
- **REST API** — Data access
- **HTML Web Scraping** — Data collection
- **Git & GitHub** — Version control
- **Vercel** — Deployment

---

## ⚙️ Environment Variables

Create a `.env` file in the project root:

~~~env
DATABASE_URL=your_postgresql_connection_string
PORT=8080
~~~

### Available Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | - | PostgreSQL database connection string |
| `PORT` | No | `8080` | HTTP server port |

> **Important:** Never commit your `.env` file to Git.

---

## 💻 Getting Started

### Clone the Repository

~~~bash
git clone https://github.com/somiyaakter/bddoctorapi.git
cd bddoctorapi
~~~

### Install Dependencies

~~~bash
go mod download
~~~

### Configure Environment Variables

Create a `.env` file:

~~~env
DATABASE_URL=your_postgresql_connection_string
PORT=8080
~~~

### Run the Application

~~~bash
go run .
~~~

The API will be available at:

~~~text
http://localhost:8080
~~~

---

## 🔌 API Documentation

### Get Doctors

Returns a paginated list of doctors.

~~~http
GET /api/v1/doctors
~~~

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | `1` | Page number |
| `page_size` | integer | `20` | Number of doctors per page |

### Example Request

~~~http
GET /api/v1/doctors?page=1&page_size=20
~~~

### Production Request

~~~http
GET https://bddoctorscraper.vercel.app/api/v1/doctors?page=1&page_size=20
~~~

---

## 📦 Example Response

~~~json
{
  "data": [
    {
      "id": 1,
      "name": "Dr. Example",
      "bmdc_reg_no": "A-123456",
      "degrees": "MBBS, FCPS",
      "experience_years": 10,
      "specialties": "Medicine",
      "designation": "Consultant",
      "workplace": "Example Hospital",
      "image_url": "https://example.com/image.jpg",
      "profile_url": "https://example.com/doctor/example",
      "chambers": [
        {
          "id": 1,
          "doctor_id": 1,
          "name": "Example Diagnostic Center",
          "address": "Dhaka, Bangladesh",
          "visiting_hour": "5:00 PM - 8:00 PM",
          "appointment_phone": "01XXXXXXXXX"
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 7413,
    "total_pages": 371
  }
}
~~~

---

## 📄 Pagination

The API supports pagination for efficient data retrieval.

### First Page

~~~http
GET /api/v1/doctors?page=1&page_size=20
~~~

### Second Page

~~~http
GET /api/v1/doctors?page=2&page_size=20
~~~

### Custom Page Size

~~~http
GET /api/v1/doctors?page=1&page_size=100
~~~

---

## 🕷️ Scraping Workflow

The scraper follows this workflow:

~~~text
Source Website
      ↓
Fetch Doctor Pages
      ↓
Parse HTML
      ↓
Extract Doctor Information
      ↓
Extract Chamber Information
      ↓
Store Data in PostgreSQL
      ↓
Expose Data through REST API
~~~

---

## ☁️ Deployment

The API is deployed on **Vercel**.

### Required Environment Variable

Configure the following environment variable in the Vercel project:

~~~env
DATABASE_URL=your_postgresql_connection_string
~~~

`PORT` is optional and defaults to `8080`.

### Deployment Flow

~~~text
GitHub Repository
       ↓
Vercel
       ↓
Build & Deploy
       ↓
Production API
~~~

---

## 🔐 Security

Never commit sensitive environment variables to Git.

Do not commit:

- `.env`
- `.env.local`
- `.env.*.local`
- Database credentials
- API keys
- Private keys

Production secrets should be configured through Vercel Environment Variables.

### Recommended `.gitignore`

~~~gitignore
.env
.env.local
.env.*.local
~~~

---

## ⚠️ Responsible Scraping

This project should be used responsibly.

- Respect the target website's policies.
- Avoid excessive requests.
- Use appropriate request delays.
- Respect applicable rate limits.
- Do not collect unnecessary personal information.
- Use scraped data for legitimate purposes.
- Follow applicable laws and regulations.

---

## 🧭 Future Improvements

- Doctor search
- Doctor filtering
- Location-based filtering
- Doctor detail endpoint
- Chamber detail endpoint
- API authentication
- API key management
- Per-client rate limiting
- Response caching
- Swagger / OpenAPI documentation
- Automated tests
- Structured logging
- Scheduled scraping

---

## 🤝 Contributing

Contributions are welcome.

### Create a Feature Branch

~~~bash
git checkout -b feature/your-feature
~~~

### Commit Changes

~~~bash
git add .
git commit -m "add: your feature"
~~~

### Push Your Branch

~~~bash
git push origin feature/your-feature
~~~

Then open a Pull Request.

---

## 📝 License

This project is currently under development.

Add your preferred license here if you plan to distribute the project publicly.

---

## 👨‍💻 Author

**Somiya Aakter**

Computer Science & Engineering Student

---

## ⭐ Project Status

**Active Development**

The API is currently deployed and serving doctor data through the production endpoint.

### Production API

https://bddoctorscraper.vercel.app/api/v1/doctors