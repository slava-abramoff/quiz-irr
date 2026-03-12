import { Routes, Route } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import TestsPage from './pages/TestsPage'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<TestsPage />} />
      <Route path="/login" element={<LoginPage />} />
    </Routes>
  )
}
