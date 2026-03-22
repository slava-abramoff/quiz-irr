import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { QuizPage } from './pages/QuizPage'
import './App.css'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/:testId?" element={<QuizPage />} />
      </Routes>
    </BrowserRouter>
  )
}
