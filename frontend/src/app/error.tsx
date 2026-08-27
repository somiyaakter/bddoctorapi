'use client'

import { useEffect } from 'react'

export default function Error({
	error,
	reset,
}: {
	error: Error & { digest?: string }
	reset: () => void
}) {
	useEffect(() => {
		console.error('Application error:', error)
	}, [error])

	return (
		<main className="mx-auto flex min-h-[60vh] max-w-md items-center justify-center px-4 py-16 sm:px-6">
			<div className="w-full text-center">
				<div className="mx-auto mb-5 flex h-14 w-14 items-center justify-center rounded-full bg-red-50">
					<span className="text-2xl">!</span>
				</div>

				<h2 className="font-display text-xl font-semibold text-ink sm:text-2xl">
					Something went wrong
				</h2>

				<p className="mt-2 text-sm leading-6 text-slate">
					Sorry, we could not load the information. Please try again in a
					moment.
				</p>

				<button
					type="button"
					onClick={() => reset()}
					className="mt-6 rounded-md bg-brand px-5 py-2.5 text-sm font-medium text-white transition hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand focus:ring-offset-2"
				>
					Try Again
				</button>
			</div>
		</main>
	)
}