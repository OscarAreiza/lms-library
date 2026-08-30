// Mirrors library-docs/07-api/contracts/openapi/library-api.yaml component schemas.

export interface Student {
  id: string
  fullName: string
  documentId: string
  email: string
  phone: string | null
  suspendedUntil: string | null
  deactivatedAt: string | null
  createdAt: string
  updatedAt: string
}

export interface Book {
  id: string
  title: string
  author: string
  isbn: string
  category: string
  year: number
  totalCopies: number
  availableCopies: number
  createdAt: string
  updatedAt: string
}

export interface Loan {
  id: string
  studentId: string
  bookId: string
  loanDate: string
  dueDate: string
  returnDate: string | null
  status: 'ACTIVE' | 'RETURNED'
  wasLate: boolean | null
  createdAt: string
  updatedAt: string
}

export interface ApiError {
  error: string
  message: string
  correlationId?: string
}

export interface PaginatedMeta {
  page: number
  limit: number
  total: number
  totalPages: number
}

export interface Paginated<T> {
  data: T[]
  meta: PaginatedMeta
}
