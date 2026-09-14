import { type FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError } from '../../types'

// Implements HU-04 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to register new books... so that the library
// catalog is populated."
export function BookFormPage() {
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [author, setAuthor] = useState('')
  const [isbn, setIsbn] = useState('')
  const [category, setCategory] = useState('')
  const [year, setYear] = useState('')
  const [totalCopies, setTotalCopies] = useState('1')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setError(null)
    setIsSubmitting(true)
    try {
      await api.post('/books', {
        title,
        author,
        isbn,
        category,
        year: Number(year),
        totalCopies: Number(totalCopies),
      })
      navigate('/books')
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to register the book. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="mx-auto max-w-lg space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Register book</h1>
        <p className="text-sm text-slate-500">HU-04 — populates the catalog with a new title.</p>
      </div>

      <Card>
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="title" className="mb-1 block text-sm font-medium text-slate-700">Title</label>
            <input id="title" required value={title} onChange={(e) => setTitle(e.target.value)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100" />
          </div>

          <div>
            <label htmlFor="author" className="mb-1 block text-sm font-medium text-slate-700">Author</label>
            <input id="author" required value={author} onChange={(e) => setAuthor(e.target.value)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100" />
          </div>

          <div>
            <label htmlFor="isbn" className="mb-1 block text-sm font-medium text-slate-700">ISBN</label>
            <input id="isbn" required value={isbn} onChange={(e) => setIsbn(e.target.value)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100" />
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="col-span-2">
              <label htmlFor="category" className="mb-1 block text-sm font-medium text-slate-700">Category</label>
              <input id="category" required value={category} onChange={(e) => setCategory(e.target.value)}
                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100" />
            </div>
            <div>
              <label htmlFor="year" className="mb-1 block text-sm font-medium text-slate-700">Year</label>
              <input id="year" type="number" required value={year} onChange={(e) => setYear(e.target.value)}
                className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100" />
            </div>
          </div>

          <div>
            <label htmlFor="totalCopies" className="mb-1 block text-sm font-medium text-slate-700">Total copies</label>
            <input id="totalCopies" type="number" min={1} required value={totalCopies} onChange={(e) => setTotalCopies(e.target.value)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100" />
          </div>

          {error && (
            <p role="alert" className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-error-600">
              {error}
            </p>
          )}

          <div className="flex justify-end gap-3">
            <Button type="button" variant="secondary" onClick={() => navigate('/books')}>Cancel</Button>
            <Button type="submit" isLoading={isSubmitting}>Register book</Button>
          </div>
        </form>
      </Card>
    </div>
  )
}
