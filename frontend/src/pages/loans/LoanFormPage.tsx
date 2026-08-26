import { type FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError } from '../../types'

// Implements HU-06 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to register a book loan to a registered
// student with a due date."
//
// Student/Book are entered by ID for now — a picker backed by HU-03/HU-05's
// search endpoints replaces these plain fields once this branch is merged
// alongside them.
export function LoanFormPage() {
  const navigate = useNavigate()
  const [studentId, setStudentId] = useState('')
  const [bookId, setBookId] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setError(null)
    setIsSubmitting(true)
    try {
      await api.post('/loans', { studentId, bookId })
      navigate('/loans')
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to register the loan. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Register loan</h1>
        <p className="text-sm text-slate-500">HU-06 — due date is set automatically to 7 days from today.</p>
      </div>

      <Card>
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="studentId" className="mb-1 block text-sm font-medium text-slate-700">
              Student ID
            </label>
            <input
              id="studentId"
              required
              value={studentId}
              onChange={(e) => setStudentId(e.target.value)}
              placeholder="Student UUID"
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
            />
          </div>

          <div>
            <label htmlFor="bookId" className="mb-1 block text-sm font-medium text-slate-700">
              Book ID
            </label>
            <input
              id="bookId"
              required
              value={bookId}
              onChange={(e) => setBookId(e.target.value)}
              placeholder="Book UUID"
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
            />
          </div>

          {error && (
            <p role="alert" className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-error-600">
              {error}
            </p>
          )}

          <div className="flex justify-end gap-3">
            <Button type="button" variant="secondary" onClick={() => navigate('/loans')}>
              Cancel
            </Button>
            <Button type="submit" isLoading={isSubmitting}>
              Register loan
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
