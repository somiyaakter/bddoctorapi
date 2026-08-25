🇧🇩 BD Doctor Scraper API

A production-ready REST API for collecting, storing, and serving doctor
information from Bangladesh-based healthcare sources.

The project is built with Go, PostgreSQL, and a web scraping pipeline. It
provides structured doctor and chamber information through a paginated REST API.

🚀 Live API

Production API

BD Doctor Scraper API

Example Request GET /api/v1/doctors?page=1&page_size=20

Example:

https://bddoctorscraper.vercel.app/api/v1/doctors?page=1&page_size=20 ✨
Features Doctor information scraping Doctor profile information BMDC
registration number Degrees and qualifications Experience information
Specialties Designation Workplace information Doctor profile image Doctor
profile URL Multiple chamber information Appointment phone numbers Visiting
hours Pagination PostgreSQL persistence RESTful API Production deployment with
Vercel Environment-based configuration 🏗️ Architecture ┌──────────────────────┐
│ External Website │ │ Doctor Profiles │ └──────────┬───────────┘ │ │ Scraping ▼
┌──────────────────────┐ │ Go Scraper │ └──────────┬───────────┘ │ │ Store ▼
┌──────────────────────┐ │ PostgreSQL │ │ Database │ └──────────┬───────────┘ │
│ Query ▼ ┌──────────────────────┐ │ Go REST API │ └──────────┬───────────┘ │ ▼
┌──────────────────────┐ │ Clients │ │ Web / Mobile / Apps │
└──────────────────────┘ 🛠️ Tech Stack Technology Purpose Go Backend API &
scraper PostgreSQL Database database/sql Database access HTML Parser / Scraper
Data extraction REST API Data access Vercel Production deployment Git & GitHub
Version control 📁 Project Structure BDDoctorScraper/ │ ├── cmd/ │ └── api/ │
└── main.go │ ├── internal/ │ ├── config/ │ │ └── config.go │ │ │ ├── doctor/ │
│ ├── model.go │ │ ├── handler.go │ │ ├── service.go │ │ └── repository.go │ │ │
├── scraper/ │ │ ├── scraper.go │ │ └── parser.go │ │ │ └── database/ │ └──
database.go │ ├── migrations/ │ └── ... │ ├── .env ├── .gitignore ├── go.mod ├──
go.sum └── README.md

The exact structure may vary depending on the current implementation.

⚙️ Environment Variables

The application requires a PostgreSQL database connection.

Create a .env file locally:

DATABASE_URL=postgresql://username:password@host:5432/database Available
Environment Variables Variable Required Description DATABASE_URL Yes PostgreSQL
connection string PORT No HTTP server port; defaults to 8080

The application loads environment variables using os.Getenv().

For local development, godotenv can load variables from .env.

💻 Local Development

1. Clone the repository git clone <your-repository-url> cd BDDoctorScraper
2. Install dependencies go mod download
3. Configure environment variables

Create .env:

DATABASE_URL=your_postgresql_connection_string 4. Run the application go run
./cmd/api

The API will normally be available at:

http://localhost:8080 🗄️ Database

The project uses PostgreSQL for persistent storage.

The main entities are:

Doctor │ ├── Doctor Information │ └── Chambers ├── Name ├── Address ├── Visiting
Hours └── Appointment Phone

A doctor can have multiple chambers.

👨‍⚕️ Doctor Data Model

A doctor record contains information such as:

{ "id": 1, "name": "Dr. Example", "bmdc_reg_no": "A-123456", "degrees": "MBBS,
FCPS", "experience_years": 10, "specialties": "Medicine", "designation":
"Consultant", "workplace": "Example Hospital", "image_url":
"https://example.com/image.jpg", "profile_url":
"https://example.com/doctor/example" } 🏥 Chamber Data

A doctor can have multiple chamber locations.

Example:

{ "id": 1, "doctor_id": 1, "name": "Example Diagnostic Center", "address":
"Dhaka, Bangladesh", "visiting_hour": "5:00 PM - 8:00 PM", "appointment_phone":
"01XXXXXXXXX" } 🔌 API Documentation Get Doctors GET /api/v1/doctors

Returns a paginated list of doctors.

Query Parameters Parameter Type Default Description page integer 1 Page number
page_size integer 20 Number of records per page Example GET
/api/v1/doctors?page=1&page_size=20

Production:

https://bddoctorscraper.vercel.app/api/v1/doctors?page=1&page_size=20 📄 API
Response

Example response structure:

{ "data": [ { "id": 1, "name": "Dr. Example", "bmdc_reg_no": "A-123456",
"degrees": "MBBS, FCPS", "experience_years": 10, "specialties": "Medicine",
"designation": "Consultant", "workplace": "Example Hospital", "image_url":
"https://example.com/image.jpg", "profile_url":
"https://example.com/doctor/example", "chambers": [ { "id": 1, "doctor_id": 1,
"name": "Example Diagnostic Center", "address": "Dhaka, Bangladesh",
"visiting_hour": "5:00 PM - 8:00 PM", "appointment_phone": "01XXXXXXXXX" } ] }
], "pagination": { "page": 1, "page_size": 20, "total_items": 7413,
"total_pages": 371 } } 📚 Pagination

The API supports pagination to avoid returning thousands of records in a single
request.

First page GET /api/v1/doctors?page=1&page_size=20 Second page GET
/api/v1/doctors?page=2&page_size=20 Larger page size GET
/api/v1/doctors?page=1&page_size=100

The response contains:

{ "pagination": { "page": 1, "page_size": 20, "total_items": 7413,
"total_pages": 371 } } 🕷️ Scraping Pipeline

The scraper collects doctor information from the configured source website and
converts it into structured data.

The general workflow is:

Source Website ↓ Fetch Doctor Pages ↓ Parse HTML ↓ Extract Doctor Information ↓
Extract Chamber Information ↓ Validate / Normalize Data ↓ Store in PostgreSQL ↓
Serve through REST API 🔄 Data Flow Doctor Website │ ▼ Scraper │ ▼ HTML Parser │
▼ Doctor + Chamber Data │ ▼ PostgreSQL │ ▼ REST API │ ▼ Frontend / Client
Application ☁️ Deployment

The application can be deployed to Vercel.

Required Environment Variable

In the Vercel project:

Settings ↓ Environment Variables ↓ DATABASE_URL

Set:

DATABASE_URL=your_postgresql_connection_string

PORT is optional because the application defaults to:

8080 🔐 Security

Sensitive configuration should never be committed to Git.

Do not commit:

.env .env.local database credentials API secrets private keys

Make sure .gitignore contains:

.env .env.local .env.\*.local

Production secrets should be stored in the deployment platform's
environment-variable system.

🧪 Testing the API

Using curl:

curl "https://bddoctorscraper.vercel.app/api/v1/doctors?page=1&page_size=20"

Using a browser:

https://bddoctorscraper.vercel.app/api/v1/doctors?page=1&page_size=20 📊 Current
Dataset

The current API contains approximately:

Doctors : 7,413 Pagination : 20 records/page Total Pages : 371

These numbers may change as the scraper collects or updates additional records.

🧭 Future Improvements

Potential future improvements include:

Doctor search

Filter by specialty

Filter by workplace

Filter by location

Filter by BMDC registration number

Doctor detail endpoint

Chamber detail endpoint

Sorting

API authentication

API key management

Per-client rate limiting

Response caching

Background scraping jobs

Scheduled scraping

Scraping monitoring

Retry mechanism

Structured logging

Automated tests

OpenAPI / Swagger documentation

⚠️ Responsible Scraping

This project should be used responsibly.

When scraping external websites:

Respect the website's terms of service. Respect applicable robots.txt rules
where appropriate. Avoid excessive request rates. Implement reasonable delays
and rate limits. Do not collect unnecessary personal information. Use collected
data only for legitimate purposes. Keep the source website's policies and
applicable laws in mind. 🤝 Contributing

Contributions are welcome.

Development workflow git checkout -b feature/your-feature

Make your changes, then:

git add . git commit -m "add: your feature" git push origin feature/your-feature

Open a pull request after pushing the branch.

📝 License

Add your preferred license here.

For example:

MIT License

If this project contains scraped data, make sure the project's data usage and
redistribution terms are compatible with the source website and applicable laws.

👨‍💻 Author

Somiya Aakter

Computer Science & Engineering Student

⭐ Project Status

Status: Active Development

The API is currently deployed and serving doctor data through the production
endpoint.

Production API:

https://bddoctorscraper.vercel.app/api/v1/doctors
