import { useEffect, useState } from 'react'

import { Button } from '../../components/ui/Button'
import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError, Loan, Paginated } from '../../types'

// Implements HU-07 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want to register the return of a loaned book...
// so that the loan cycle is closed" — plus the history half of the same HU.
export function LoansListPage() {
  const [loans, setLoans] = useState<Loan[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  async function load() {
    setIsLoading(true)
    setError(null)
    try {
      const { data } = await api.get<Paginated<Loan>>('/loans')
      setLoans(data.data ?? [])
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to load loans.')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  async function registerReturn(id: string) {
    try {
      await api.post(`/loans/${id}/return`)
      load()
    } catch (err) {
      const apiError = (err as { response?: { data?: ApiError } })?.response?.data
      setError(apiError?.message ?? 'Unable to register the return.')
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Loans</h1>
        <p className="text-sm text-slate-500">Active loans and full history. Registration lands with HU-06.</p>
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
              <th className="px-4 py-3 font-medium">Student ID</th>
              <th className="px-4 py-3 font-medium">Book ID</th>
              <th className="px-4 py-3 font-medium">Due date</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {isLoading && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-slate-400">Loading...</td></tr>
            )}
            {!isLoading && loans.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-slate-400">No loans yet.</td></tr>
            )}
            {loans.map((l) => (
              <tr key={l.id}>
                <td className="px-4 py-3 text-slate-500">{l.studentId}</td>
                <td className="px-4 py-3 text-slate-500">{l.bookId}</td>
                <td className="px-4 py-3 text-slate-500">{new Date(l.dueDate).toLocaleDateString()}</td>
                <td className="px-4 py-3">
                  {l.status === 'ACTIVE' ? (
                    <span className="rounded-full bg-emerald-50 px-2 py-0.5 text-xs text-success-600">Active</span>
                  ) : l.wasLate ? (
                    <span className="rounded-full bg-rose-50 px-2 py-0.5 text-xs text-error-600">Returned (late)</span>
                  ) : (
                    <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-500">Returned</span>
                  )}
                </td>
                <td className="px-4 py-3">
                  {l.status === 'ACTIVE' && (
                    <Button variant="secondary" onClick={() => registerReturn(l.id)}>Register return</Button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </Card>
    </div>
  )
}
