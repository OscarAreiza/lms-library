import { Navigate, Route, BrowserRouter as Router, Routes } from 'react-router-dom'

import { ProtectedRoute } from './components/layout/ProtectedRoute'
import { LoginPage } from './pages/LoginPage'
import { DashboardPage } from './pages/DashboardPage'
import { StudentsListPage } from './pages/students/StudentsListPage'
import { StudentFormPage } from './pages/students/StudentFormPage'
import { BooksListPage } from './pages/books/BooksListPage'
import { BookFormPage } from './pages/books/BookFormPage'
import { LoansListPage } from './pages/loans/LoansListPage'
import { OverdueLoansPage } from './pages/loans/OverdueLoansPage'
import { NotFoundPage } from './pages/NotFoundPage'

// Route tree mirrors library-docs/12-ux-ui/navigation-map.md. Detail/form routes
// (/students/new, /books/:id, /loans/new, ...) are added alongside the feat/HU-XX
// branch that implements each one.
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/login" element={<LoginPage />} />

        <Route path="/dashboard" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
        <Route path="/students" element={<ProtectedRoute><StudentsListPage /></ProtectedRoute>} />
        <Route path="/students/new" element={<ProtectedRoute><StudentFormPage /></ProtectedRoute>} />
        <Route path="/books" element={<ProtectedRoute><BooksListPage /></ProtectedRoute>} />
        <Route path="/books/new" element={<ProtectedRoute><BookFormPage /></ProtectedRoute>} />
        <Route path="/loans" element={<ProtectedRoute><LoansListPage /></ProtectedRoute>} />
        <Route path="/loans/overdue" element={<ProtectedRoute><OverdueLoansPage /></ProtectedRoute>} />

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Router>
  )
}

export default App
