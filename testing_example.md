
# E-Library Microservices Project

## gRPC Endpoint Testing Report

**Team Members:**
Bekbolat — User Service
Nurassyl — Book Service
Nurkhan — Borrow Service

**Project:** E-Library Microservices System
**Testing Tool:** Postman gRPC / grpcurl
**Backend:** Go + PostgreSQL + gRPC + Docker
**Total tested methods:** 36

---

# 1. Purpose

The purpose of this report is to test all gRPC methods of the E-Library system.

System architecture:

```txt
Frontend
     ↓
API Gateway
     ↓
gRPC
 ┌──────────────┬──────────────┬─────────────┐
 │User Service  │Book Service │Borrow Service
 └──────────────┴──────────────┴─────────────┘
```

---

# 2. Environment

OS:

```txt
Windows 11
```

Language:

```txt
Go
```

Database:

```txt
PostgreSQL
```

Testing:

```txt
Postman gRPC
```

Ports:

```txt
User Service: 50051
Book Service: 50052
Borrow Service: 50053
```

---

# 3. User Service gRPC Testing

### RegisterUser

```json
{
  "name":"Test User",
  "email":"test@test.com",
  "password":"123456",
  "role":"student"
}
```

Expected:

```txt
User created successfully
```

---

### LoginUser

```json
{
  "email":"test@test.com",
  "password":"123456"
}
```

---

### GetUser

```json
{
   "id":"1"
}
```

---

### GetAllUsers

```json
{}
```

---

### UpdateUser

```json
{
 "id":"1",
 "name":"Updated User"
}
```

---

### UpdateRole

```json
{
 "id":"1",
 "role":"admin"
}
```

---

### UpdatePassword

```json
{
 "id":"1",
 "password":"newpassword123"
}
```

---

### DeleteUser

```json
{
 "id":"1"
}
```

---

### GetUsersByRole

```json
{
 "role":"student"
}
```

---

### CountUsers

```json
{}
```

---

### CheckUserExists

```json
{
 "id":"1"
}
```

---

# 4. Book Service gRPC Testing

### CreateBook

```json
{
 "title":"Clean Code",
 "author":"Robert C Martin",
 "category":"Programming",
 "available":"Yes"
}
```

---

### GetBook

```json
{
 "id":"3"
}
```

---

### GetAllBooks

```json
{}
```

---

### UpdateBook

```json
{
 "id":"3",
 "title":"Clean Code Updated",
 "author":"Robert Martin",
 "category":"Software Engineering"
}
```

---

### DeleteBook

```json
{
 "id":"3"
}
```

---

### SearchBooks

```json
{
 "query":"clean"
}
```

---

### GetAvailableBooks

```json
{}
```

---

### CountBooks

```json
{}
```

---

# 5. Borrow Service gRPC Testing

### CreateBorrow

```json
{
 "user_id":"1",
 "book_id":"3",
 "due_date":"2026-12-12"
}
```

---

### GetBorrow

```json
{
 "id":"1"
}
```

---

### GetAllBorrows

```json
{}
```

---

### ReturnBorrow

```json
{
 "id":"1"
}
```

---

### ExtendBorrowPeriod

```json
{
 "id":"1"
}
```

---

### CancelBorrow

```json
{
 "id":"1"
}
```

---

### GetBorrowsByUserID

```json
{
 "user_id":"1"
}
```

---

### GetBorrowsByBookID

```json
{
 "book_id":"3"
}
```

---

### GetActiveBorrows

```json
{}
```

---

### GetOverdueBorrows

```json
{}
```

---

### CountBorrows

```json
{}
```

---

### CheckBorrowExists

```json
{
 "id":"1"
}
```

---

# 6. Screenshot Section

For every test include:

```txt
Method name
Request JSON
Response JSON
Status OK
```

Example:

```txt
Figure 1. Testing CreateBook gRPC method in Postman.
```

---

# 7. Conclusion

All 36 gRPC methods were tested successfully. The tests validated communication between API Gateway and microservices. CRUD operations, borrow management, search functionality, and statistics endpoints operated correctly.
