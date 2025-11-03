## 📛 *CheckGuard – Fraud Phone Blocklist System**

A Go-based system to help retail stores verify and block fraudulent phone numbers used in suspicious check-cashing activities.

Built to prevent fraud by maintaining a centralized blocklist with full audit details.

---

## 🚀 **Features**

| Category            | Details                                                  |
| ------------------- | -------------------------------------------------------- |
| Block Fraud Numbers | Add & block suspicious phone numbers                     |
| Search              | Search by phone number, store location                   |
| Validation          | Strict E.164 phone format validation using DB CHECK rule |
| Audit               | Tracks incident date, notes, timestamps                  |
| UI                  | Simple HTML-based front-end for store clerks             |
| Persistence         | PostgreSQL storage                                       |
| REST APIs           | JSON CRUD APIs                                           |
| Testing             | Unit test coverage for handlers & DB layer               |
| Folder Structure    | Clean MVC-style architecture                             |

---

## 🧠 **Tech Stack**

| Layer        | Technology                             |
| ------------ | -------------------------------------- |
| Backend      | Go (net/http)                          |
| Database     | PostgreSQL                             |
| DB Migration | Custom schema loader (`db/schema.sql`) |
| Tests        | Go testing package (`testing`)         |
| Frontend     | HTML, CSS, JS (basic)                  |

---

## 📁 **Project Structure**

```
blocklist_app/
├── main.go
├── go.mod
├── routes/
│   └── routes.go
├── handlers/
│   └── handlers.go
├── db/
│   ├── db.go
│   └── schema.sql
├── models/
│   └── blocked_number.go
├── static/
│   ├── index.html
│   └── style.css
└── tests/
    └── handlers_test.go
```

---

## 🗃️ **Database Schema**

Schema file location: `db/schema.sql`

Key rules included:

✔ `CHECK (phone_number ~ '^\+[1-9][0-9]{9,14}$')`
✔ Prevents invalid numbers
✔ Prevents numbers starting with +0
✔ Indexes for fast search

Table: `blocked_numbers`

| Column                  | Description               |
| ----------------------- | ------------------------- |
| id                      | Primary key               |
| phone_number            | E.164 format (unique)     |
| reason                  | Fraud reason              |
| store_location          | City/store ID             |
| incident_date           | When fraud occurred       |
| check_amount            | Amount on presented check |
| notes                   | Additional notes          |
| created_at / updated_at | Audit                     |

---

## 🏗️ **Architecture Flow**

```
User → UI Form → HTTP Request → Router → Handler → DB Layer → PostgreSQL
```

---

## ⚙️ **Setup & Installation**

### ✅ Prerequisite

* Go 1.20+
* PostgreSQL 14+
* Git

### ✅ Clone Repo

```sh
git clone https://github.com/yourname/blocklist_app.git
cd blocklist_app
```

### ✅ DB Setup

Create DB:

```sql
CREATE DATABASE checkguard;
```

### ✅ Configure DB Env

Edit in `db/db.go`:

```
postgres://postgres:YOUR_PASSWORD@localhost:5432/checkguard?sslmode=disable
```

### ✅ Install Dependencies

```sh
go mod tidy
```

### ✅ Run App

```sh
go run main.go
```

Visit UI:

```
http://localhost:8080
```

---

## 📡 **API Endpoints**

### ➕ Add Blocked Number

`POST /block`

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

## 🧩 **Static File Handling**

Served via Go:

```go
fs := http.FileServer(http.Dir("static"))
http.Handle("/static/", http.StripPrefix("/static/", fs))
```

Maps `/static/*` → `static/` folder.

---

## ✅ **Key Security & Validation Rules**

* DB-level phone validation
* Prevents invalid entries
* Logs & audit timestamps
* Unique index on phone numbers
* Server-side input validation

---

## 📌 **Future Enhancements**

| Feature               | Status    |
| --------------------- | --------- |
| JWT Auth              | ⏳ Planned |
| Admin Dashboard       | ⏳         |
| Redis Cache           | ⏳         |
| Cloud-ready migration | ⏳         |

---

## 🧑‍💻 **Author**

**Muzaffar Khan**

---

## ⭐ **Contribute**

Pull Requests welcome. Open issues for suggestions.

---
