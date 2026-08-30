import { Card } from '../components/ui/Card'

// Placeholder overview — real counts are wired once HU-07/HU-08 (loans/returns/
// overdue) endpoints exist. Layout matches library-docs/12-ux-ui/navigation-map.md.
const stats = [
  { label: 'Active loans', value: '—' },
  { label: 'Overdue loans', value: '—' },
  { label: 'Registered students', value: '—' },
  { label: 'Books in catalog', value: '—' },
]

export function DashboardPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">Dashboard</h1>
        <p className="text-sm text-slate-500">Overview of today's library activity.</p>
      </div>

      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.label}>
            <p className="text-sm font-medium text-slate-500">{stat.label}</p>
            <p className="mt-2 text-3xl font-semibold text-slate-900">{stat.value}</p>
          </Card>
        ))}
      </div>
    </div>
  )
}
