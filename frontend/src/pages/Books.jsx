import { useEffect, useState } from "react";
import api from "../services/api";

export default function Books() {
  const [books, setBooks] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    loadBooks();
  }, []);

  async function loadBooks() {
    try {
      const response = await api.get("/books");

      console.log(response.data);

      setBooks(response.data.books || []);
    } catch (err) {
      console.error(err);
      setError(err.message);
    }
  }

  return (
    <div>
      <h1>E-Library</h1>

      {error && <p>{error}</p>}

      {books.map((book) => (
        <div
          key={book.id}
          style={{
            border: "1px solid gray",
            padding: "15px",
            margin: "10px",
            borderRadius: "10px"
          }}
        >
          <h2>{book.title}</h2>
          <p>Author: {book.author}</p>
          <p>Category: {book.category}</p>
        </div>
      ))}
    </div>
  );
}