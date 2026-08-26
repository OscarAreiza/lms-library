import { type FormEvent, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError, Book, Paginated } from '../../types'

// Implements HU-05 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to search the catalog... and see real-time
// availability."
export function BooksListPage() {
  const [books, setBooks] = useState<Book[]>([])
  const [search, setSearch] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  async function load(query: string) {
    setIsLoading(true)
    setError(null)
    try {
      const { data } = await api.get<Paginated<Book>>('/books', { params: { search: query } })
      setBooks(data.data ?? [])
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to load the catalog.')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    load('')
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  function handleSearch(event: FormEvent) {
    event.preventDefault()
    load(search)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">Catalog</h1>
          <p className="text-sm text-slate-500">Search the book catalog and check availability.</p>
        </div>
        <Link to="/books/new">
          <Button>Register book</Button>
        </Link>
      </div>

      <form onSubmit={handleSearch} className="flex gap-2">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by title, author, or ISBN"
          className="w-full max-w-sm rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
        />
        <Button type="submit" variant="secondary">Search</Button>
      </form>

      {error && (
        <p role="alert" className="rounded-lg bg-rose-50 px-3 py-2 text-sm text-error-600">
          {error}
        </p>
      )}

      <Card className="overflow-x-auto p-0">
        <table className="w-full text-left text-sm">
          <thead className="border-b border-slate-200 bg-slate-50 text-slate-500">
            <tr>
              <th className="px-4 py-3 font-medium">Title</th>
              <th className="px-4 py-3 font-medium">Author</th>
              <th className="px-4 py-3 font-medium">ISBN</th>
              <th className="px-4 py-3 font-medium">Category</th>
              <th className="px-4 py-3 font-medium">Availability</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {isLoading && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-slate-400">Loading...</td></tr>
            )}
            {!isLoading && books.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-slate-400">No books found.</td></tr>
            )}
            {books.map((b) => (
              <tr key={b.id}>
                <td className="px-4 py-3">{b.title}</td>
                <td className="px-4 py-3 text-slate-500">{b.author}</td>
                <td className="px-4 py-3 text-slate-500">{b.isbn}</td>
                <td className="px-4 py-3 text-slate-500">{b.category}</td>
                <td className="px-4 py-3">
                  <span
                    className={
                      b.availableCopies > 0
                        ? 'rounded-full bg-emerald-50 px-2 py-0.5 text-xs text-success-600'
                        : 'rounded-full bg-rose-50 px-2 py-0.5 text-xs text-error-600'
                    }
                  >
                    {b.availableCopies} / {b.totalCopies} available
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </Card>
    </div>
  )
}
