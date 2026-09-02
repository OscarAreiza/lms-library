import { clsx } from 'clsx'
import type { PropsWithChildren } from 'react'

export function Card({ children, className }: PropsWithChildren<{ className?: string }>) {
  return (
    <div className={clsx('rounded-xl border border-slate-200 bg-white p-6 shadow-sm', className)}>
      {children}
    </div>
  )
}
