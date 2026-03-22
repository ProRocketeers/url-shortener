"use client"

import { useState } from "react"
import { Button } from "@/components/ui/reusable/Button"
import { useTranslations } from 'next-intl'
import { useLocale } from 'next-intl'
import { useRouter } from '@/i18n/routing'
import { ArrowLeft, Link as LinkIcon, ExternalLink, Copy, Check, CircleAlert } from 'lucide-react'
import { useShortLinkInfoBySlug } from "@/hooks/api/url-shortener"

type DetailPageProps = {
	slug: string
}

export function DetailPage({ slug }: DetailPageProps) {
	const router = useRouter()
	const locale = useLocale()
	const t = useTranslations('common')
	const tDetail = useTranslations('detail')
	const [copiedField, setCopiedField] = useState<"original" | "short" | null>(null)
	const { data, isLoading, isError } = useShortLinkInfoBySlug(slug)

	const handleCopy = async (value: string, field: "original" | "short") => {
		await navigator.clipboard.writeText(value)
		setCopiedField(field)
		setTimeout(() => setCopiedField(null), 1500)
	}

	if (isLoading) {
		return (
			<div className="w-full rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
				<p className="text-slate-600 text-center">Loading...</p>
			</div>
		)
	}

	if (isError || !data) {
		return (
			<div className="w-full max-w-2xl rounded-2xl border border-red-200 bg-white p-6 shadow-sm sm:p-8">
				<div className="flex flex-col items-center justify-center gap-5 text-center">
					<div className="flex h-14 w-14 items-center justify-center rounded-full bg-red-50 text-red-600">
						<CircleAlert className="h-7 w-7" />
					</div>
					<div className="space-y-2">
						<h1 className="text-2xl font-semibold text-slate-900">{tDetail('loadErrorTitle')}</h1>
						<p className="max-w-xl text-sm leading-6 text-slate-600 sm:text-base">
							{tDetail('loadErrorDescription')}
						</p>
					</div>
					<Button variant="outline" className="border-slate-200" onClick={() => router.push('/')}>
						<ArrowLeft className="h-4 w-4" />
						{tDetail('backToForm')}
					</Button>
				</div>
			</div>
		)
	}

	const expiresAtDisplay = (() => {
		if (!data.expiresAt) {
			return tDetail('noExpiration')
		}

		const parsed = new Date(data.expiresAt)
		if (Number.isNaN(parsed.getTime())) {
			return tDetail('noExpiration')
		}

		return new Intl.DateTimeFormat(locale, {
			dateStyle: 'medium',
			timeStyle: 'short',
		}).format(parsed)
	})()

	const createdAtDisplay = (() => {
		if (!data.createdAt) {
			return '-'
		}

		const parsed = new Date(data.createdAt)
		if (Number.isNaN(parsed.getTime())) {
			return '-'
		}

		return new Intl.DateTimeFormat(locale, {
			dateStyle: 'medium',
			timeStyle: 'short',
		}).format(parsed)
	})()

	return (
		<div className="w-full max-w-4xl flex flex-col gap-8">
			{/* Header */}
			<div className="flex flex-col items-center text-center gap-6">
				<div className="bg-gradient-to-br from-[#051641] to-[#0a3d7a] rounded-2xl w-16 h-16 flex items-center justify-center shadow-md">
					<LinkIcon className="w-8 h-8 text-white" />
				</div>
				<div>
					<h1 className="text-3xl font-bold text-slate-700">{tDetail('title')}</h1>
					<p className="mt-2 text-sm text-slate-600">
						Slug: {slug}
					</p>
				</div>
			</div>

			{/* Detail Card */}
			<div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
				<div className="space-y-6">
					{/* Placeholder Info */}
					<div className="rounded-lg border border-blue-200 bg-blue-50 p-4">
						<p className="text-sm text-blue-800">
							📊 {tDetail('statsPlaceholder') || 'Statistics and detailed information will be implemented here in the future.'}
						</p>
					</div>

					{/* Original URL Section - Placeholder */}
					<div>
						<h2 className="text-sm font-semibold text-slate-700 mb-3">{t('originalUrl')}</h2>
						<div className="flex items-center gap-2">
							<div className="flex-1 border border-slate-200 rounded-md px-3 py-2 bg-slate-50">
								<p className="text-sm text-slate-600 truncate">{data.originalUrl}</p>
							</div>
							<Button
								type="button"
								variant="outline"
								size="icon"
								className="border-slate-200"
								onClick={() => handleCopy(data.originalUrl, "original")}
							>
								{copiedField === "original" ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
							</Button>
							<a href={data.originalUrl} target="_blank" rel="noopener noreferrer">
								<Button variant="outline" size="icon" className="border-slate-200">
									<ExternalLink className="h-4 w-4" />
								</Button>
							</a>
						</div>
					</div>

					{/* Short URL Section - Placeholder */}
					<div>
						<h2 className="text-sm font-semibold text-slate-700 mb-3">{t('shortUrl')}</h2>
						<div className="flex items-center gap-2">
							<div className="flex-1 border border-slate-200 rounded-md px-3 py-2 bg-slate-50">
								<p className="text-sm font-mono text-slate-600">{data.shortUrl}</p>
							</div>
							<Button
								type="button"
								variant="outline"
								size="icon"
								className="border-slate-200"
								onClick={() => handleCopy(data.shortUrl, "short")}
							>
								{copiedField === "short" ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
							</Button>
							<a href={data.shortUrl} target="_blank" rel="noopener noreferrer">
								<Button variant="outline" size="icon" className="border-slate-200">
									<ExternalLink className="h-4 w-4" />
								</Button>
							</a>
						</div>
					</div>

					{/* Stats Grid - Placeholder */}
					<div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
						<div className="rounded-lg border border-slate-200 bg-slate-50 p-4">
							<p className="text-sm text-slate-600 mb-1">{tDetail('clicks')}</p>
							<p className="text-2xl font-bold text-slate-900">{data.clickCount}</p>
						</div>
						<div className="rounded-lg border border-slate-200 bg-slate-50 p-4">
							<p className="text-sm text-slate-600 mb-1">{tDetail('createdAt')}</p>
							<p className="text-sm font-medium text-slate-900">{createdAtDisplay}</p>
						</div>
						<div className="rounded-lg border border-slate-200 bg-slate-50 p-4">
							<p className="text-sm text-slate-600 mb-1">{tDetail('expiresAt')}</p>
							<p className="text-sm font-medium text-slate-900">{expiresAtDisplay}</p>
						</div>
					</div>
				</div>
			</div>

			{/* Back Button */}
			<Button
				variant="outline"
				onClick={() => router.push('/')}
				className="self-center w-full border-slate-200"
			>
				<ArrowLeft className="h-4 w-4" />
				{t('back') || 'Back'}
			</Button>
		</div>
	)
}
