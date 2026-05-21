import { BrowserRouter, Link, Route, Routes } from "react-router-dom";
import Books from "./pages/Books";
import Users from "./pages/Users";
import Borrows from "./pages/Borrows";

function App() {
  return (
    <BrowserRouter>
      <nav>
        <Link to="/">📚 Books</Link>
        <Link to="/users">👤 Users</Link>
        <Link to="/borrows">📖 Borrows</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Books />} />
        <Route path="/users" element={<Users />} />
        <Route path="/borrows" element={<Borrows />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;