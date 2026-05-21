# E-Library Microservices Project

## Postman API Testing Report

**Team Members:**
Bekbolat — User Service
Nurassyl — Book Service
Nurkhan — Borrow Service

**Project:** E-Library Microservices System
**Testing Tool:** Postman
**Backend:** Go + PostgreSQL + Docker + gRPC + API Gateway
**Base URL:**

```txt
http://localhost:8080/api
```

**Total tested endpoints:** 36

---

# 1. Purpose

The purpose of this report is to test and validate all API endpoints of the E-Library Microservices System using Postman.

The system contains three microservices:

* User Service
* Book Service
* Borrow Service

Each service contains 12 endpoints.

Total:

```txt
36 endpoints tested
```

---

# 2. Testing Environment

Operating System:

```txt
Windows 11
```

Backend:

```txt
Go
```

Database:

```txt
PostgreSQL
```

Containerization:

```txt
Docker + Docker Compose
```

Testing Tool:

```txt
Postman
```

API Gateway:

```txt
localhost:8080
```

---

# 3. User Service Testing

## Endpoint 1

Method:

```txt
POST
```

URL:

```txt
/api/users/register
```

Raw JSON:

```json
{
  "name": "Test User",
  "email": "test@test.com",
  "password": "123456",
  "role": "student"
}
```

Expected:

```txt
201 Created
```

---

## Endpoint 2

Method:

```txt
POST
```

URL:

```txt
/api/users/login
```

Raw JSON:

```json
{
  "email":"test@test.com",
  "password":"123456"
}
```

Expected:

```txt
200 OK
```

---

## Endpoint 3

Method:

```txt
GET
```

URL:

```txt
/api/users
```

Body:

```txt
No body
```

Expected:

```txt
200 OK
```

---

## Endpoint 4

Method:

```txt
GET
```

URL:

```txt
/api/users/{id}
```

Body:

```txt
No body
```

---

## Endpoint 5

Method:

```txt
PUT
```

URL:

```txt
/api/users/{id}
```

Raw JSON:

```json
{
 "name":"Updated User"
}
```

---

## Endpoint 6

Method:

```txt
PUT
```

URL:

```txt
/api/users/{id}/role
```

Raw JSON:

```json
{
 "role":"admin"
}
```

---

## Endpoint 7

Method:

```txt
PUT
```

URL:

```txt
/api/users/{id}/password
```

Raw JSON:

```json
{
 "password":"newpassword123"
}
```

---

## Endpoint 8

Method:

```txt
DELETE
```

URL:

```txt
/api/users/{id}
```

Body:

```txt
No body
```

---

## Endpoint 9

Method:

```txt
GET
```

URL:

```txt
/api/users/email/test@test.com
```

---

## Endpoint 10

Method:

```txt
GET
```

URL:

```txt
/api/users/role/student
```

---

## Endpoint 11

Method:

```txt
GET
```

URL:

```txt
/api/users/count
```

---

## Endpoint 12

Method:

```txt
GET
```

URL:

```txt
/api/users/exists/{id}
```

---

# 4. Book Service Testing

## Endpoint 13

Method:

```txt
POST
```

URL:

```txt
/api/books
```

Raw JSON:

```json
{
  "title":"Clean Code",
  "author":"Robert C Martin",
  "category":"Programming",
  "available":"Yes"
}
```

Expected:

```txt
201 Created
```

---

## Endpoint 14

Method:

```txt
GET
```

URL:

```txt
/api/books
```

---

## Endpoint 15

Method:

```txt
GET
```

URL:

```txt
/api/books/{id}
```

---

## Endpoint 16

Method:

```txt
PUT
```

URL:

```txt
/api/books/{id}
```

Raw JSON:

```json
{
 "title":"Clean Code Updated",
 "author":"Robert Martin",
 "category":"Software Engineering"
}
```

---

## Endpoint 17

Method:

```txt
DELETE
```

URL:

```txt
/api/books/{id}
```

---

## Endpoint 18

Method:

```txt
GET
```

URL:

```txt
/api/books/search?q=clean
```

---

## Endpoint 19

Method:

```txt
GET
```

URL:

```txt
/api/books/category/Programming
```

---

## Endpoint 20

Method:

```txt
PUT
```

URL:

```txt
/api/books/{id}/available
```

Raw JSON:

```json
{
 "available":"Yes"
}
```

---

## Endpoint 21

Method:

```txt
PUT
```

URL:

```txt
/api/books/{id}/unavailable
```

Raw JSON:

```json
{
 "available":"No"
}
```

---

## Endpoint 22

Method:

```txt
GET
```

URL:

```txt
/api/books/available
```

---

## Endpoint 23

Method:

```txt
GET
```

URL:

```txt
/api/books/stats
```

---

## Endpoint 24

Method:

```txt
GET
```

URL:

```txt
/api/books/count
```

---

# 5. Borrow Service Testing

## Endpoint 25

Method:

```txt
POST
```

URL:

```txt
/api/borrows
```

Raw JSON:

```json
{
 "user_id":"1",
 "book_id":"3",
 "due_date":"2026-12-12"
}
```

---

## Endpoint 26

Method:

```txt
GET
```

URL:

```txt
/api/borrows
```

---

## Endpoint 27

Method:

```txt
GET
```

URL:

```txt
/api/borrows/{id}
```

---

## Endpoint 28

Method:

```txt
PUT
```

URL:

```txt
/api/borrows/{id}/return
```

Raw JSON:

```json
{}
```

---

## Endpoint 29

Method:

```txt
PUT
```

URL:

```txt
/api/borrows/{id}/extend
```

Raw JSON:

```json
{
 "days":7
}
```

---

## Endpoint 30

Method:

```txt
DELETE
```

URL:

```txt
/api/borrows/{id}/cancel
```

---

## Endpoint 31

Method:

```txt
GET
```

URL:

```txt
/api/borrows/user/{userId}
```

---

## Endpoint 32

Method:

```txt
GET
```

URL:

```txt
/api/borrows/book/{bookId}
```

---

## Endpoint 33

Method:

```txt
GET
```

URL:

```txt
/api/borrows/active
```

---

## Endpoint 34

Method:

```txt
GET
```

URL:

```txt
/api/borrows/overdue
```

---

## Endpoint 35

Method:

```txt
GET
```

URL:

```txt
/api/borrows/count
```

---

## Endpoint 36

Method:

```txt
GET
```

URL:

```txt
/api/borrows/exists/{id}
```

---

# 6. Screenshot Section

For every endpoint include:

* Postman request URL
* Raw JSON body
* Response JSON
* Status code

Example:

```txt
Figure 1. Testing POST /api/users/register endpoint using Postman
```

---

# 7. Conclusion

All 36 endpoints of the E-Library Microservices System were tested using Postman. CRUD operations, search endpoints, count endpoints, update methods, and borrow management functionality were validated successfully.

Testing confirmed communication between User Service, Book Service, Borrow Service, PostgreSQL databases, and API Gateway.
