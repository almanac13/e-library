import { useEffect, useState } from "react";
import api from "../services/api";

export default function Borrows() {
  const [borrows, setBorrows] = useState([]);
  const [users, setUsers] = useState([]);
  const [books, setBooks] = useState([]);
  const [userId, setUserId] = useState("");
  const [bookId, setBookId] = useState("");
  const [dueDate, setDueDate] = useState("");

  const getArray = (data) => {
    if (Array.isArray(data)) return data;
    if (Array.isArray(data.data)) return data.data;
    if (Array.isArray(data.users)) return data.users;
    if (Array.isArray(data.books)) return data.books;
    if (Array.isArray(data.borrows)) return data.borrows;
    return [];
  };

  const loadBorrows = async () => {
    const res = await api.get("/borrows");
    setBorrows(getArray(res.data));
  };

  const loadUsers = async () => {
    const res = await api.get("/users");
    setUsers(getArray(res.data));
  };

  const loadBooks = async () => {
    const res = await api.get("/books");
    const allBooks = getArray(res.data);

    const availableBooks = allBooks.filter((book) => {
      const value = String(book.available).toLowerCase();
      return value === "true" || value === "yes" || value === "available";
    });

    setBooks(availableBooks);
  };

  const createBorrow = async () => {
  if (!userId || !bookId || !dueDate) {
    alert("Select user, book and due date");
    return;
  }

  try {
    const res = await api.post("/borrows", {
      user_id: userId,
      book_id: bookId,
      due_date: dueDate,
    });

    console.log("Created:", res.data);
    alert("Borrow created");

    setUserId("");
    setBookId("");
    setDueDate("");

    await loadBorrows();
    await loadBooks();
  } catch (err) {
    console.log("CREATE ERROR:", err.response?.data || err.message);
    alert(err.response?.data?.error || "Create borrow failed");
  }
};

  useEffect(() => {
    loadBorrows();
    loadUsers();
    loadBooks();
  }, []);

  return (
    <div>
      <h1>Borrows</h1>
      <h2>Total borrows: {borrows.length}</h2>

      <div className="card">
        <h2>Create Borrow</h2>

        <select value={userId} onChange={(e) => setUserId(e.target.value)}>
          <option value="">Select user</option>
          {users.map((user) => (
            <option key={user.id} value={user.id}>
              {user.name} — {user.email}
            </option>
          ))}
        </select>

        <select value={bookId} onChange={(e) => setBookId(e.target.value)}>
          <option value="">Select available book</option>
          {books.map((book) => (
            <option key={book.id} value={book.id}>
              {book.title} — {book.author}
            </option>
          ))}
        </select>

        <input
          type="date"
          value={dueDate}
          onChange={(e) => setDueDate(e.target.value)}
        />

        <button onClick={createBorrow}>Create Borrow</button>
      </div>

      {borrows.map((borrow) => (
        <div className="card" key={borrow.id}>
          <h2>Borrow ID: {borrow.id}</h2>
          <p>User ID: {borrow.user_id}</p>
          <p>Book ID: {borrow.book_id}</p>
          <p>Status: {borrow.status}</p>
          <p>Due Date: {borrow.due_date}</p>
        </div>
      ))}
    </div>
  );
}