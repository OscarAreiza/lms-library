import { useEffect, useState } from 'react'

import { Button } from '../../components/ui/Button'
import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError, Book, Paginated } from '../../types'

// Implements HU-09 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to edit a registered book's information...
// so that the catalog stays accurate." ISBN is immutable — not part of the
// edit form (02-domain/entities-and-rules.md modeling note on Book).
export function BooksListPage() {
  const [books, setBooks] = useState<Book[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editForm, setEditForm] = useState({ title: '', author: '', category: '', year: '' })

  async function load() {
    setIsLoading(true)
    setError(null)
    try {
      const { data } = await api.get<Paginated<Book>>('/books')
      setBooks(data.data ?? [])
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to load the catalog.')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  function startEdit(book: Book) {
    setEditingId(book.id)
    setEditForm({ title: book.title, author: book.author, category: book.category, year: String(book.year) })
  }

  async function saveEdit(id: string) {
    try {
      await api.patch(`/books/${id}`, { ...editForm, year: Number(editForm.year) })
      setEditingId(null)
      load()
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to update the book.')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Catalog</h1>
        <p className="text-sm text-slate-500">Edit book information. ISBN is fixed at registration.</p>
      </div>

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
              <th className="px-4 py-3 font-medium">Year</th>
              <th className="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {isLoading && (
              <tr><td colSpan={6} className="px-4 py-6 text-center text-slate-400">Loading...</td></tr>
            )}
            {!isLoading && books.length === 0 && (
              <tr><td colSpan={6} className="px-4 py-6 text-center text-slate-400">No books found.</td></tr>
            )}
            {books.map((b) => (
              <tr key={b.id}>
                {editingId === b.id ? (
                  <>
                    <td className="px-4 py-2">
                      <input value={editForm.title} onChange={(e) => setEditForm((f) => ({ ...f, title: e.target.value }))}
                        className="w-full rounded border border-slate-300 px-2 py-1" />
                    </td>
                    <td className="px-4 py-2">
                      <input value={editForm.author} onChange={(e) => setEditForm((f) => ({ ...f, author: e.target.value }))}
                        className="w-full rounded border border-slate-300 px-2 py-1" />
                    </td>
                    <td className="px-4 py-2 text-slate-400">{b.isbn} (fixed)</td>
                    <td className="px-4 py-2">
                      <input value={editForm.category} onChange={(e) => setEditForm((f) => ({ ...f, category: e.target.value }))}
                        className="w-full rounded border border-slate-300 px-2 py-1" />
                    </td>
                    <td className="px-4 py-2">
                      <input type="number" value={editForm.year} onChange={(e) => setEditForm((f) => ({ ...f, year: e.target.value }))}
                        className="w-20 rounded border border-slate-300 px-2 py-1" />
                    </td>
                    <td className="px-4 py-2">
                      <div className="flex gap-2">
                        <Button onClick={() => saveEdit(b.id)}>Save</Button>
                        <Button variant="ghost" onClick={() => setEditingId(null)}>Cancel</Button>
                      </div>
                    </td>
                  </>
                ) : (
                  <>
                    <td className="px-4 py-3">{b.title}</td>
                    <td className="px-4 py-3 text-slate-500">{b.author}</td>
                    <td className="px-4 py-3 text-slate-500">{b.isbn}</td>
                    <td className="px-4 py-3 text-slate-500">{b.category}</td>
                    <td className="px-4 py-3 text-slate-500">{b.year}</td>
                    <td className="px-4 py-3">
                      <Button variant="ghost" onClick={() => startEdit(b)}>Edit</Button>
                    </td>
                  </>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </Card>
    </div>
  )
}
