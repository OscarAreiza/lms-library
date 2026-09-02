import { Link } from 'react-router-dom'

import { Button } from '../../components/ui/Button'
import { EmptyState } from '../../components/ui/EmptyState'

export function StudentsListPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900">Students</h1>
          <p className="text-sm text-slate-500">Search, register, and manage student records.</p>
        </div>
        <Link to="/students/new">
          <Button>Register student</Button>
        </Link>
      </div>
      <EmptyState
        title="Student search coming soon"
        description="Listing, editing, and deactivating existing students lands with HU-03. Registration (HU-02) is already available above."
        hu="HU-03"
      />
    </div>
  )
}
