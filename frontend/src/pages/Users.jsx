import { useEffect, useState } from "react";
import api from "../services/api";

const card = {
  border: "1px solid #3b3b46",
  padding: "20px",
  margin: "14px auto",
  borderRadius: "16px",
  maxWidth: "900px",
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

export default function Users() {
  const [users, setUsers] = useState([]);
  const [count, setCount] = useState(0);
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    role: "student",
  });

  useEffect(() => {
    loadUsers();
  }, []);

  async function loadUsers() {
    const usersRes = await api.get("/users");
    const countRes = await api.get("/users/count");

    setUsers(usersRes.data.users || []);
    setCount(countRes.data.count || 0);
  }

  async function registerUser(e) {
    e.preventDefault();

    await api.post("/users/register", form);

    setForm({
      name: "",
      email: "",
      password: "",
      role: "student",
    });

    await loadUsers();
  }

  async function deleteUser(id) {
    await api.delete(`/users/${id}`);
    await loadUsers();
  }

  async function updateName(user) {
    const name = prompt("New name:", user.name);
    if (!name) return;

    await api.put(`/users/${user.id}/name`, { name });
    await loadUsers();
  }

  async function updateRole(user) {
    const role = prompt("New role:", user.role);
    if (!role) return;

    await api.put(`/users/${user.id}/role`, { role });
    await loadUsers();
  }

  async function changePassword(user) {
    const password = prompt("New password:");
    if (!password) return;

    await api.put(`/users/${user.id}/password`, { password });
    alert("Password changed");
  }

  return (
    <div style={{ padding: "20px", color: "white" }}>
      <h1 style={{ textAlign: "center", fontSize: "56px" }}>Users</h1>
      <h2 style={{ textAlign: "center", color: "#aab" }}>Total users: {count}</h2>

      <form onSubmit={registerUser} style={card}>
        <h2>Create User</h2>

        <input
          style={input}
          placeholder="Name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          required
        />

        <input
          style={input}
          placeholder="Email"
          value={form.email}
          onChange={(e) => setForm({ ...form, email: e.target.value })}
          required
        />

        <input
          style={input}
          placeholder="Password"
          type="password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          required
        />

        <input
          style={input}
          placeholder="Role"
          value={form.role}
          onChange={(e) => setForm({ ...form, role: e.target.value })}
          required
        />

        <button style={{ ...button, background: "#4f46e5", color: "white" }}>
          Register
        </button>
      </form>

      {users.map((user) => (
        <div key={user.id} style={card}>
          <h2>{user.name}</h2>
          <p>Email: {user.email}</p>
          <p>Role: {user.role}</p>
          <p>ID: {user.id}</p>

          <button style={{ ...button, background: "#2563eb", color: "white" }} onClick={() => updateName(user)}>
            Update Name
          </button>

          <button style={{ ...button, background: "#7c3aed", color: "white" }} onClick={() => updateRole(user)}>
            Update Role
          </button>

          <button style={{ ...button, background: "#059669", color: "white" }} onClick={() => changePassword(user)}>
            Change Password
          </button>

          <button style={{ ...button, background: "#dc2626", color: "white" }} onClick={() => deleteUser(user.id)}>
            Delete
          </button>
        </div>
      ))}
    </div>
  );
}