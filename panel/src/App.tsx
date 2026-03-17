import { Routes, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import TestsPage from './pages/TestsPage';
import EditorPage from './pages/EditorPage';
import TestDetailsPage from './pages/TestDetailsPage';

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<TestsPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/editor/:id" element={<EditorPage />} />
      <Route path="/test/:id" element={<TestDetailsPage />} />
    </Routes>
  );
}
