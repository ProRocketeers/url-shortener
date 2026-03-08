"use client"

import { Button } from "@/components/ui/reusable/Button"
import Image from "next/image"
import { useState, useMemo } from "react"
import { FormContext } from "@/components/ui/form"
import { createUrlShortenerSchema, defaultValues, type UrlShortenerFormValues } from "@/utils/schemas/urlShortenerSchema"
import { useCreateShortLink } from "@/hooks/api/url-shortener"
import { useTranslations } from 'next-intl'
import { Link as LinkIcon, Copy, Check, ExternalLink } from 'lucide-react'
import { InputController } from "@/components/ui/form/InputController"
import { ExpirationController } from "@/components/ui/form/ExpirationController"
import type { ShortenUrlResponse } from "@/api/url-shortener"
import { useRouter } from '@/i18n/routing'

export function UrlShortenerForm() {
	const t = useTranslations('common')
	const tForm = useTranslations('form')
	const tValidation = useTranslations('validation')
	const router = useRouter()
	const [error, setError] = useState("")
	const [shortLink, setShortLink] = useState<ShortenUrlResponse & { originalUrl: string } | null>(null)
	const [copied, setCopied] = useState(false)

	// Create schema with translations
	const schema = useMemo(() => createUrlShortenerSchema(tValidation), [tValidation])

	const { mutateAsync: createShortLink, isPending } = useCreateShortLink()

	const handleSubmit = async (formValues: UrlShortenerFormValues) => {
		setError("")
		setShortLink(null)
		setCopied(false)

		try {
			const payload = {
				originalUrl: formValues.originalUrl,
				...(formValues.slug && { slug: formValues.slug }),
				...(formValues.expiresAt && { expiresAt: formValues.expiresAt }),
			}

			const result = await createShortLink(payload)
			// Redirect to detail page
			router.push(`/detail/${result.slug}`)
		} catch {
			setError(t('error'))
		}
	}

	const handleCopy = async () => {
		if (shortLink?.shortUrl) {
			await navigator.clipboard.writeText(shortLink.shortUrl)
			setCopied(true)
			setTimeout(() => setCopied(false), 2000)
		}
	}

	return (
		<div className="w-full max-w-4xl">
			{/* Header with icon and title */}
			<div className="mb-6 flex flex-col items-center text-center gap-6">
				<div className="bg-gradient-to-br from-[#051641] to-[#0a3d7a] rounded-2xl w-16 h-16 flex items-center justify-center shadow-md">
					<LinkIcon className="w-8 h-8 text-white" />
				</div>
				<div>
					<h1 className="text-3xl font-bold text-slate-700">{t('appTitle')}</h1>
					<p className="mt-2 text-sm text-slate-600 max-w-md">
						{t('appDescription')}
					</p>
				</div>
			</div>

			{/* Form Card */}
			<div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
				<FormContext
					schema={schema}
					defaultValues={defaultValues}
					onSubmit={handleSubmit}
					mode="onChange"
					reValidateMode="onChange"
				>
					{(form) => {
						return (
							<>
								<div className="mt-6 grid gap-4">
									<InputController
										control={form.control}
										name="originalUrl"
										label={tForm('urlLabel')}
										placeholder={tForm('urlPlaceholder')}
										type="url"
									/>

									<InputController
										control={form.control}
										name="slug"
										label={tForm('customSlugLabel')}
										placeholder={tForm('customSlugPlaceholder')}
										type="text"
									/>

									<ExpirationController
										control={form.control}
										name="expiresAt"
										label={tForm('expiresAtLabel')}
										checkboxLabel={tForm('setExpiration')}
										placeholder={tForm('expiresAtPlaceholder')}
										dateLabel={tForm('dateLabel')}
										timeLabel={tForm('timeLabel')}
									/>

									<div className="flex items-center gap-3">
										<Button type="submit" variant="secondary" disabled={isPending}>
											{isPending ? t('shortening') : t('shorten')}
										</Button>
									</div>

									{error && <p className="text-sm text-red-600">{error}</p>}
								</div>

								{/* Success Message */}
								{shortLink && (
									<div className="mt-6 rounded-xl border border-green-200 bg-green-50 p-6">
										<h2 className="text-lg font-semibold text-green-800 mb-4 flex items-center gap-2">
											<Check className="h-5 w-5" />
											{t('urlShortened')}
										</h2>
										
										<div className="space-y-4">
											<div>
												<p className="text-sm font-medium text-slate-700 mb-2">{t('originalUrl')}:</p>
												<div className="flex items-center gap-2">
													<p className="text-sm text-slate-600 flex-1 truncate">{shortLink.originalUrl}</p>
													<a 
														href={shortLink.originalUrl} 
														target="_blank" 
														rel="noopener noreferrer"
														className="text-slate-500 hover:text-slate-700"
													>
														<ExternalLink className="h-4 w-4" />
													</a>
												</div>
											</div>

											<div>
												<p className="text-sm font-medium text-slate-700 mb-2">{t('shortUrl')}:</p>
												<div className="flex items-center gap-2">
													<div className="flex-1 bg-white border border-slate-200 rounded-md px-3 py-2">
														<p className="text-sm font-mono text-slate-900">{shortLink.shortUrl}</p>
													</div>
													<Button
														type="button"
														variant="outline"
														size="icon"
														onClick={handleCopy}
														className="shrink-0"
													>
														{copied ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
													</Button>
													<a 
														href={shortLink.shortUrl} 
														target="_blank" 
														rel="noopener noreferrer"
													>
														<Button
															type="button"
															variant="outline"
															size="icon"
															className="shrink-0"
														>
															<ExternalLink className="h-4 w-4" />
														</Button>
													</a>
												</div>
											</div>
										</div>
									</div>
								)}
							</>
						)
					}}
				</FormContext>
			</div>
		</div>
	)
}
