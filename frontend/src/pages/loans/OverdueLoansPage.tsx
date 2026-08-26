import { useEffect, useState } from 'react'

import { Card } from '../../components/ui/Card'
import { api } from '../../lib/api'
import type { ApiError, Loan, Paginated } from '../../types'

// Implements HU-08 (library-docs/04-requirements/user-stories.md):
// "As the administrator, I want the system to detect overdue loans... so that
// late returns are discouraged." Scenario 2 — the Overdue Loans report.
export function OverdueLoansPage() {
  const [loans, setLoans] = useState<Loan[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function load() {
      setIsLoading(true)
      setError(null)
      try {
        const { data } = await api.get<Paginated<Loan>>('/loans/overdue')
        setLoans(data.data ?? [])
      } catch (err) {
        const apiError = (err as { response?: { data?: ApiError } })?.response?.data
        setError(apiError?.message ?? 'Unable to load overdue loans.')
      } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Overdue loans</h1>
        <p className="text-sm text-slate-500">Active loans past their due date. A late return applies a 7-day suspension automatically.</p>
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
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {isLoading && (
              <tr><td colSpan={4} className="px-4 py-6 text-center text-slate-400">Loading...</td></tr>
            )}
            {!isLoading && loans.length === 0 && (
              <tr><td colSpan={4} className="px-4 py-6 text-center text-slate-400">No overdue loans right now.</td></tr>
            )}
            {loans.map((l) => (
              <tr key={l.id}>
                <td className="px-4 py-3 text-slate-500">{l.studentId}</td>
                <td className="px-4 py-3 text-slate-500">{l.bookId}</td>
                <td className="px-4 py-3 text-slate-500">{new Date(l.dueDate).toLocaleDateString()}</td>
                <td className="px-4 py-3">
                  <span className="rounded-full bg-amber-50 px-2 py-0.5 text-xs text-warning-500">Overdue</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </Card>
    </div>
  )
}
