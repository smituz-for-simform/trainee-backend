# 📘 Trainee Backend (Go + PostgreSQL)

A simple CRUD Phonebook backend built using **Go (Gin)** and **PostgreSQL**. This project is designed to help you understand backend development along with DevOps concepts like environment variables and database setup.

---

# 🚀 What This Backend Does

This backend provides APIs to manage contacts:

- **GET /contacts** → Fetch all contacts
- **POST /contacts** → Add a new contact
- **PUT /contacts** → Update an existing contact
- **DELETE /contacts/:id** → Delete a contact

---

# 🛠️ Prerequisites

Make sure you have the following installed:

- Go (latest version)
- PostgreSQL (installed and running)
- Git
- Postman or curl (for testing APIs)

---

# ⚙️ Step 1: Clone the Repository

```bash
git clone <your-backend-repo-url>
cd trainee_backend
```

---

# ⚙️ Step 2: Setup PostgreSQL Database

## 1. Open PostgreSQL shell

```bash
sudo -i -u postgres
psql
```

## 2. Create database

```sql
CREATE DATABASE phonebook;
```

## 3. Create user

```sql
CREATE USER appuser WITH PASSWORD 'strongpassword';
```

## 4. Grant access

```sql
GRANT ALL PRIVILEGES ON DATABASE phonebook TO appuser;
```

## 5. Grant schema permissions (important)

```sql
\c phonebook
GRANT ALL ON SCHEMA public TO appuser;
```

## 6. Create table

```sql
CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    name TEXT,
    phone TEXT
);
```

## 7. Exit

```sql
\q
exit
```

---

# 🔐 Step 3: Create `.env` File

Create a file named `.env` in the root of the project:

```bash
touch .env
```

Add the following:

```env
DB_URL=postgres://appuser:strongpassword@localhost:5432/phonebook
```

---

# 🧠 What is DB_URL?

This is a connection string used by your Go app to connect to PostgreSQL.

Format:

```
postgres://username:password@host:port/database
```

Example used here:

- username → appuser
- password → strongpassword
- host → localhost (your machine)
- port → 5432 (default PostgreSQL port)
- database → phonebook

---

# ⚠️ Important

Add `.env` to `.gitignore` so it is not pushed to GitHub:

```bash
echo ".env" >> .gitignore
```

---

# ▶️ Step 4: Run the Backend

```bash
go mod tidy
go run main.go
```

Server will start on:

```
http://localhost:8080
```

---

# 🧪 Step 5: Test APIs

## ➤ Create Contact (POST)

- URL: `http://localhost:8080/add_contact`
- Method: POST
- Body (JSON):

```json
{
  "name": "John Doe",
  "phone": "9876543210"
}
```

---

## ➤ Get Contacts (GET)

- URL: `http://localhost:8080/get_contacts`

---

## ➤ Update Contact (PUT)

- URL: `http://localhost:8080/update_contact`
- Body:

```json
{
  "id": 1,
  "name": "Updated Name",
  "phone": "9999999999"
}
```

---

## ➤ Delete Contact (DELETE)

- URL: `http://localhost:8080/del_contact/1`

---

# 🔍 Debugging Tips

- Check DB data:
  ```sql
  SELECT * FROM contacts;
  ```
- If connection fails → verify DB_URL
- If table not found → ensure table is created in `phonebook`
- If permission denied → re-run schema grant step

---

# 💡 DevOps Notes

This project is structured to help you later with:

- Dockerizing backend and database
- Using environment variables securely
- Deploying on Kubernetes

---

# ✅ You’re Good To Go

If everything is set correctly, you should be able to:

- Run backend
- Connect to database
- Perform all CRUD operations successfully

---

Happy Learning 🚀
