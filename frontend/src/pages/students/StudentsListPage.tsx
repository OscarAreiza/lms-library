import { type FormEvent, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError, Paginated, Student } from '../../types'

// Implements HU-03 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to search, edit, or deactivate student records."
export function StudentsListPage() {
  const [students, setStudents] = useState<Student[]>([])
  const [search, setSearch] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editForm, setEditForm] = useState({ fullName: '', email: '', phone: '' })

  async function load(query: string) {
    setIsLoading(true)
    setError(null)
    try {
      const { data } = await api.get<Paginated<Student>>('/students', { params: { search: query } })
      setStudents(data.data ?? [])
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to load students.')
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

  function startEdit(student: Student) {
    setEditingId(student.id)
    setEditForm({ fullName: student.fullName, email: student.email, phone: student.phone ?? '' })
  }

  async function saveEdit(id: string) {
    try {
      await api.patch(`/students/${id}`, editForm)
      setEditingId(null)
      load(search)
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to update the student.')
    }
  }

  async function deactivate(id: string) {
    if (!window.confirm('Deactivate this student? This is blocked if they have active loans or a suspension.')) {
      return
    }
    try {
      await api.post(`/students/${id}/deactivate`)
      load(search)
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to deactivate the student.')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">Students</h1>
          <p className="text-sm text-slate-500">Search, register, edit, and deactivate student records.</p>
        </div>
        <Link to="/students/new">
          <Button>Register student</Button>
        </Link>
      </div>

      <form onSubmit={handleSearch} className="flex gap-2">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by name or document ID"
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
              <th className="px-4 py-3 font-medium">Name</th>
              <th className="px-4 py-3 font-medium">Document ID</th>
              <th className="px-4 py-3 font-medium">Email</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {isLoading && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-slate-400">Loading...</td></tr>
            )}
            {!isLoading && students.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-slate-400">No students found.</td></tr>
            )}
            {students.map((s) => (
              <tr key={s.id}>
                {editingId === s.id ? (
                  <>
                    <td className="px-4 py-2">
                      <input
                        value={editForm.fullName}
                        onChange={(e) => setEditForm((f) => ({ ...f, fullName: e.target.value }))}
                        className="w-full rounded border border-slate-300 px-2 py-1"
                      />
                    </td>
                    <td className="px-4 py-2 text-slate-500">{s.documentId}</td>
                    <td className="px-4 py-2">
                      <input
                        value={editForm.email}
                        onChange={(e) => setEditForm((f) => ({ ...f, email: e.target.value }))}
                        className="w-full rounded border border-slate-300 px-2 py-1"
                      />
                    </td>
                    <td className="px-4 py-2 text-slate-500">—</td>
                    <td className="px-4 py-2">
                      <div className="flex gap-2">
                        <Button onClick={() => saveEdit(s.id)}>Save</Button>
                        <Button variant="ghost" onClick={() => setEditingId(null)}>Cancel</Button>
                      </div>
                    </td>
                  </>
                ) : (
                  <>
                    <td className="px-4 py-3">{s.fullName}</td>
                    <td className="px-4 py-3 text-slate-500">{s.documentId}</td>
                    <td className="px-4 py-3 text-slate-500">{s.email}</td>
                    <td className="px-4 py-3">
                      {s.deactivatedAt ? (
                        <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-500">Deactivated</span>
                      ) : s.suspendedUntil ? (
                        <span className="rounded-full bg-amber-50 px-2 py-0.5 text-xs text-warning-500">Suspended</span>
                      ) : (
                        <span className="rounded-full bg-emerald-50 px-2 py-0.5 text-xs text-success-600">Active</span>
                      )}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex gap-2">
                        <Button variant="ghost" onClick={() => startEdit(s)}>Edit</Button>
                        {!s.deactivatedAt && (
                          <Button variant="danger" onClick={() => deactivate(s.id)}>Deactivate</Button>
                        )}
                      </div>
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
