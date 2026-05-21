import { useEffect, useState } from "react";
import api from "../services/api";

export default function Books() {
  const [books, setBooks] = useState([]);
  const [stats, setStats] = useState(null);
  const [search, setSearch] = useState("");
  const [form, setForm] = useState({
    title: "",
    author: "",
    category: "",
    available: true,
  });

  useEffect(() => {
    loadBooks();
    loadStats();
  }, []);

  async function loadBooks() {
    const response = await api.get("/books");
    setBooks(response.data.books || []);
  }

  async function loadStats() {
    const response = await api.get("/books/stats");
    setStats(response.data);
  }

  async function createBook(e) {
    e.preventDefault();

    await api.post("/books", form);

    setForm({
      title: "",
      author: "",
      category: "",
      available: true,
    });

    await loadBooks();
    await loadStats();
  }

  async function deleteBook(id) {
    await api.delete(`/books/${id}`);
    await loadBooks();
    await loadStats();
  }

  async function markAvailable(id) {
    await api.put(`/books/${id}/available`);
    await loadBooks();
    await loadStats();
  }

  async function markUnavailable(id) {
    await api.put(`/books/${id}/unavailable`);
    await loadBooks();
    await loadStats();
  }

  async function searchBooks(e) {
    e.preventDefault();

    if (!search.trim()) {
      await loadBooks();
      return;
    }

    const response = await api.get(`/books/search?q=${search}`);
    setBooks(response.data.books || []);
  }

  async function updateBook(book) {
    const newTitle = prompt("New title:", book.title);
    const newAuthor = prompt("New author:", book.author);
    const newCategory = prompt("New category:", book.category);

    if (!newTitle || !newAuthor || !newCategory) return;

    await api.put(`/books/${book.id}`, {
      title: newTitle,
      author: newAuthor,
      category: newCategory,
      available: book.available,
    });

    await loadBooks();
    await loadStats();
  }

  return (
    <div style={{ padding: "20px" }}>
      <h1>E-Library</h1>

      {stats && (
        <div style={{ marginBottom: "20px" }}>
          <h3>Total books: {stats.total_books}</h3>
          <p>Available: {stats.available_books}</p>
          <p>Unavailable: {stats.unavailable_books || 0}</p>
        </div>
      )}

      <form onSubmit={searchBooks} style={{ marginBottom: "20px" }}>
        <input
          placeholder="Search books"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <button type="submit">Search</button>
        <button type="button" onClick={loadBooks}>
          Reset
        </button>
      </form>

      <form onSubmit={createBook} style={{ marginBottom: "30px" }}>
        <h2>Create Book</h2>

        <input
          placeholder="Title"
          value={form.title}
          onChange={(e) => setForm({ ...form, title: e.target.value })}
          required
        />

        <input
          placeholder="Author"
          value={form.author}
          onChange={(e) => setForm({ ...form, author: e.target.value })}
          required
        />

        <input
          placeholder="Category"
          value={form.category}
          onChange={(e) => setForm({ ...form, category: e.target.value })}
          required
        />

        <button type="submit">Create</button>
      </form>

      {books.map((book) => (
        <div
          key={book.id}
          style={{
            border: "1px solid gray",
            padding: "15px",
            margin: "10px",
            borderRadius: "10px",
          }}
        >
          <h2>{book.title}</h2>
          <p>Author: {book.author}</p>
          <p>Category: {book.category}</p>
          <p>Available: {book.available ? "Yes" : "No"}</p>

          <button onClick={() => updateBook(book)}>Update</button>
          <button onClick={() => deleteBook(book.id)}>Delete</button>
          <button onClick={() => markAvailable(book.id)}>Available</button>
          <button onClick={() => markUnavailable(book.id)}>Unavailable</button>
        </div>
      ))}
    </div>
  );
}