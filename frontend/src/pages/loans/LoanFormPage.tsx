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
// The Administrator identifies the student and book by the natural keys
// they actually have on hand — document ID, ISBN — never the internal UUID,
// which is never shown anywhere in the UI. circulation-service resolves
// these server-side against membership-service/backend before registering
// the loan (library-docs/09-microservices/data-ownership-matrix.md).
export function LoanFormPage() {
  const navigate = useNavigate()
  const [studentDocumentId, setStudentDocumentId] = useState('')
  const [bookIsbn, setBookIsbn] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setError(null)
    setIsSubmitting(true)
    try {
      await api.post('/loans', { studentDocumentId, bookIsbn })
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
            <label htmlFor="studentDocumentId" className="mb-1 block text-sm font-medium text-slate-700">
              Student document ID
            </label>
            <input
              id="studentDocumentId"
              required
              value={studentDocumentId}
              onChange={(e) => setStudentDocumentId(e.target.value)}
              placeholder="e.g. 1075300000"
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
            />
          </div>

          <div>
            <label htmlFor="bookIsbn" className="mb-1 block text-sm font-medium text-slate-700">
              Book ISBN
            </label>
            <input
              id="bookIsbn"
              required
              value={bookIsbn}
              onChange={(e) => setBookIsbn(e.target.value)}
              placeholder="e.g. 978-0132350884"
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
