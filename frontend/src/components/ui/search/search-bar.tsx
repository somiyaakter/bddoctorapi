'use client'

import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { useEffect, useRef, useState } from 'react'

const MIN_CHARS = 2
const DEBOUNCE_MS = 400

interface SearchBarProps {
	defaultValue?: string
	/** Where search results should navigate to. Defaults to the current
	 * page's path — pass this explicitly when rendering SearchBar
	 * somewhere that ISN'T the results page itself (e.g. the landing
	 * page), since usePathname() there would be "/" not "/doctors". */
	targetPath?: string
}

export function SearchBar({ defaultValue, targetPath }: SearchBarProps) {
	const [value, setValue] = useState(defaultValue ?? '')
	const router = useRouter()
	const pathname = usePathname()
	const searchParams = useSearchParams()
	const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

	const resolvedTarget = targetPath ?? pathname

	function pushQuery(q: string) {
		const params = new URLSearchParams(
			resolvedTarget === pathname ? searchParams.toString() : '',
		)
		if (q) {
			params.set('q', q)
		} else {
			params.delete('q')
		}
		params.delete('page')
		router.push(`${resolvedTarget}?${params.toString()}`)
	}

	useEffect(() => {
		if (debounceRef.current) clearTimeout(debounceRef.current)

		const trimmed = value.trim()
		if (trimmed.length > 0 && trimmed.length < MIN_CHARS) return

		debounceRef.current = setTimeout(() => {
			pushQuery(trimmed)
		}, DEBOUNCE_MS)

		return () => {
			if (debounceRef.current) clearTimeout(debounceRef.current)
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [value])

	function handleSubmit(e: React.FormEvent) {
		e.preventDefault()
		if (debounceRef.current) clearTimeout(debounceRef.current)
		pushQuery(value.trim())
	}

	return (
		<form onSubmit={handleSubmit} role="search" className="flex gap-2">
			<input
				type="search"
				value={value}
				onChange={e => setValue(e.target.value)}
				placeholder="Search by name, specialty, hospital, or BMDC number"
				aria-label="Search doctors"
				className="flex-1 rounded-md border border-hairline bg-white px-4 py-2 text-sm text-ink placeholder:text-slate focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
			/>
			<button
				type="submit"
				className="rounded-md bg-brand px-4 py-2 text-sm font-medium text-white hover:bg-brand-hover"
			>
				Search
			</button>
		</form>
	)
}
