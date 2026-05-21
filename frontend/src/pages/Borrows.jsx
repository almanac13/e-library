import { useEffect, useState } from "react";
import api from "../services/api";

const card = {
  border: "1px solid #3b3b46",
  padding: "20px",
  margin: "14px auto",
  borderRadius: "16px",
  maxWidth: "950px",
  background: "#181820",
};

const input = {
  padding: "10px",
  margin: "6px",
  borderRadius: "8px",
  border: "1px solid #555",
  background: "#111",
  color: "white",
};

const button = {
  padding: "10px 14px",
  margin: "6px",
  borderRadius: "8px",
  border: "none",
  cursor: "pointer",
};

export default function Borrows() {
  const [borrows, setBorrows] = useState([]);
  const [count, setCount] = useState(0);
  const [form, setForm] = useState({
    user_id: "",
    book_id: "",
    due_date: "",
  });

  useEffect(() => {
    loadBorrows();
  }, []);

  async function loadBorrows() {
    const borrowsRes = await api.get("/borrows");
    const countRes = await api.get("/borrows/count");

    setBorrows(borrowsRes.data.borrows || []);
    setCount(countRes.data.count || 0);
  }

  async function createBorrow(e) {
    e.preventDefault();

    await api.post("/borrows", form);

    setForm({
      user_id: "",
      book_id: "",
      due_date: "",
    });

    await loadBorrows();
  }

  async function returnBorrow(id) {
    await api.put(`/borrows/${id}/return`);
    await loadBorrows();
  }

  async function extendBorrow(id) {
    const newDueDate = prompt("New due date, example: 2026-06-01");
    if (!newDueDate) return;

    await api.put(`/borrows/${id}/extend`, {
      new_due_date: newDueDate,
    });

    await loadBorrows();
  }

  async function cancelBorrow(id) {
    await api.put(`/borrows/${id}/cancel`);
    await loadBorrows();
  }

  async function loadActive() {
    const response = await api.get("/borrows/active");
    setBorrows(response.data.borrows || []);
  }

  async function loadOverdue() {
    const response = await api.get("/borrows/overdue");
    setBorrows(response.data.borrows || []);
  }

  return (
    <div style={{ padding: "20px", color: "white" }}>
      <h1 style={{ textAlign: "center", fontSize: "56px" }}>Borrows</h1>
      <h2 style={{ textAlign: "center", color: "#aab" }}>
        Total borrows: {count}
      </h2>

      <form onSubmit={createBorrow} style={card}>
        <h2>Create Borrow</h2>

        <input
          style={input}
          placeholder="User ID"
          value={form.user_id}
          onChange={(e) => setForm({ ...form, user_id: e.target.value })}
          required
        />

        <input
          style={input}
          placeholder="Book ID"
          value={form.book_id}
          onChange={(e) => setForm({ ...form, book_id: e.target.value })}
          required
        />

        <input
          style={input}
          placeholder="Due date: 2026-06-01"
          value={form.due_date}
          onChange={(e) => setForm({ ...form, due_date: e.target.value })}
          required
        />

        <button style={{ ...button, background: "#4f46e5", color: "white" }}>
          Create Borrow
        </button>

        <button
          type="button"
          style={{ ...button, background: "#2563eb", color: "white" }}
          onClick={loadBorrows}
        >
          All
        </button>

        <button
          type="button"
          style={{ ...button, background: "#059669", color: "white" }}
          onClick={loadActive}
        >
          Active
        </button>

        <button
          type="button"
          style={{ ...button, background: "#f59e0b", color: "black" }}
          onClick={loadOverdue}
        >
          Overdue
        </button>
      </form>

      {borrows.map((borrow) => (
        <div key={borrow.id} style={card}>
          <h2>{borrow.id}</h2>
          <p>User ID: {borrow.user_id}</p>
          <p>Book ID: {borrow.book_id}</p>
          <p>Status: {borrow.status}</p>
          <p>Borrow date: {borrow.borrow_date}</p>
          <p>Due date: {borrow.due_date}</p>
          <p>Return date: {borrow.return_date || "Not returned"}</p>

          <button
            style={{ ...button, background: "#059669", color: "white" }}
            onClick={() => returnBorrow(borrow.id)}
          >
            Return
          </button>

          <button
            style={{ ...button, background: "#2563eb", color: "white" }}
            onClick={() => extendBorrow(borrow.id)}
          >
            Extend
          </button>

          <button
            style={{ ...button, background: "#dc2626", color: "white" }}
            onClick={() => cancelBorrow(borrow.id)}
          >
            Cancel
          </button>
        </div>
      ))}
    </div>
  );
}