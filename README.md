# 📛 CheckGuard – Fraud Phone Blocklist System

A Go-based system for retail stores to verify and block fraudulent phone numbers used in suspicious check-cashing activities.
Prevents fraud by maintaining a centralized blocklist with full audit logging.

---

## 🚀 Features

| Category            | Details                                           |
| ------------------- | ------------------------------------------------- |
| Block Fraud Numbers | Add & block suspicious phone numbers              |
| Search              | Search by phone number or store location          |
| Validation          | Strict E.164 phone format validation (`DB CHECK`) |
| Audit               | Tracks incident date, notes, timestamps           |
| UI                  | Simple HTML front-end for store clerks            |
| Persistence         | PostgreSQL storage                                |
| REST APIs           | JSON CRUD APIs                                    |
| Testing             | Unit tests for handlers & utils                   |
| Structure           | Clean MVC-style architecture                      |

---

## 🧠 Tech Stack

| Layer    | Technology                             |
| -------- | -------------------------------------- |
| Backend  | Go 1.20+ (net/http)                    |
| Database | PostgreSQL 14+                         |
| Frontend | HTML, CSS, JS (basic)                  |
| Testing  | Go `testing` package                   |
| DB Setup | Custom schema loader (`db/schema.sql`) |

---

## 📁 Project Structure

```
BLOCKLIST_APP/
├── main.go
├── go.mod
├── go.sum
├── internal/
│   ├── api/
│   │   ├── handlers/handlers.go
│   │   └── routes/routes.go
│   ├── config/config.go
│   ├── db/
│   │   ├── db.go
│   │   └── schema.sql
│   ├── models/model.go
│   └── utils/
│       ├── utils_test.go
│       └── validators.go
├── static/
│   ├── css/
│   ├── js/
│   └── templates/
├── package.json   # Optional, frontend deps
└── README.md
```

---

## 🗃️ Database Schema

**Table: blocked_numbers**

| Column                  | Description           |
| ----------------------- | --------------------- |
| id                      | Primary key           |
| phone_number            | E.164 format (unique) |
| reason                  | Fraud reason          |
| store_location          | City/store ID         |
| incident_date           | Date of incident      |
| check_amount            | Amount on check       |
| notes                   | Additional notes      |
| created_at / updated_at | Audit timestamps      |

**Key rules**:

* `CHECK (phone_number ~ '^\+[1-9][0-9]{9,14}$')`
* Prevents invalid numbers and numbers starting with +0
* Indexed for fast search

---

## ⚙️ Setup & Installation

### Prerequisites

* Go 1.20+
* PostgreSQL 14+
* Git

### Clone Repo

```bash
git clone https://github.com/Mkhan2217/blocklist_app
cd blocklist_app
```

### Database Setup

```sql
CREATE DATABASE blocklistdb;
```

Update connection in `internal/db/db.go`:

```go
postgres://postgres:YOUR_PASSWORD@localhost:5432/blocklistdb?sslmode=disable
```

### Install Dependencies

```bash
go mod tidy
```

### Run App

```bash
go run main.go
```

Visit: `http://localhost:8080`

---

## 📡 API Endpoints

### ➕ Add Blocked Number

`POST /block` (JSON, Content-Type: application/json)

```json
{
  "phone_number": "+18005551234",
  "reason": "Suspicious check",
  "store_location": "Walmart NY",
  "check_amount": 450.00,
  "notes": "Forgery attempt"
}
```

### 🔍 Search

`GET /search?phone=+18005551234`

### ❌ Unblock Number

`DELETE /unblock?phone=+18005551234`

---

## 🧩 Static Files

```go
fs := http.FileServer(http.Dir("static"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
```

Maps `/static/*` → `static/` folder.

---

## ✅ Security & Validation

* DB-level phone validation & unique index
* Audit timestamps
* Server-side input validation

---

## 🧑‍💻 Author

**Muzaffar Khan**