# 🇧🇩 BD Doctor Scraper API

A production-ready REST API for scraping, storing, and serving
doctor information from Bangladesh-based healthcare sources.

Built with Go and PostgreSQL.
## 🚀 Live API

https://bddoctorscraper.vercel.app/api/v1/doctors

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

## 🛠️ Tech Stack

- Go
- PostgreSQL
- REST API
- HTML Web Scraping
- Git & GitHub
- Vercel

## ⚙️ Environment Variables

Create a `.env` file in the project root:

DATABASE_URL=your_postgresql_connection_string

### Available Variables

| Variable | Required | Default |
|---|---|---|
| DATABASE_URL | Yes | - |
| PORT | No | 8080 |

## 💻 Getting Started

### Clone the repository

```bash
git clone ...
cd bddoctorscraper
```md

## 🔌 API Documentation

### Get Doctors

GET /api/v1/doctors
### Query Parameters

| Parameter | Type | Default | Description |
|---|---|---|---|
| page | integer | 1 | Page number |
## 🕷️ Scraping Workflow

The scraper follows this workflow:

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

## ☁️ Deployment

The API is deployed on Vercel.

### Environment Variables

The following environment variable must be configured in Vercel:

DATABASE_URL=your_postgresql_connection_string

`PORT` is optional and defaults to `8080`.
## 🔐 Security

Never commit sensitive environment variables to Git.

Do not commit:

- `.env`
- Database credentials
- API keys
- Private keys

Production secrets should be configured through Vercel Environment Variables.

## ⚠️ Responsible Scraping

This project should be used responsibly.

- Respect the target website's policies.
- Avoid excessive requests.
- Use appropriate request delays.
- Respect applicable rate limits.
- Do not collect unnecessary personal information.
- Use scraped data for legitimate purposes.