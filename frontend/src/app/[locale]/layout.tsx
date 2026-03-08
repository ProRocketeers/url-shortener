import type { Metadata } from "next"
import { Geist, Geist_Mono } from "next/font/google"
import "../globals.css"
import { QueryProvider } from "../../utils/providers/QueryProvider"
import { NextIntlClientProvider } from 'next-intl'
import { getMessages } from 'next-intl/server'
import { setRequestLocale } from 'next-intl/server'
import { routing } from '@/i18n/routing'
import { notFound } from 'next/navigation'
import { LanguageSwitcher } from "@/components/LanguageSwitcher"

const geistSans = Geist({
	variable: "--font-geist-sans",
	subsets: ["latin"],
})

const geistMono = Geist_Mono({
	variable: "--font-geist-mono",
	subsets: ["latin"],
})

export const metadata: Metadata = {
	title: "url-shortener",
	description: "By ProRocketeers",
}

export function generateStaticParams() {
	return routing.locales.map((locale) => ({ locale }))
}

export default async function RootLayout({
	children,
	params,
}: Readonly<{
	children: React.ReactNode
	params: Promise<{ locale: string }>
}>) {
	const { locale } = await params
	type SupportedLocale = (typeof routing.locales)[number]

	// Ensure that the incoming `locale` is valid
	if (!routing.locales.includes(locale as SupportedLocale)) {
		notFound()
	}

	// Enable static rendering
	setRequestLocale(locale)

	// Providing all messages to the client side
	const messages = await getMessages()

	return (
		<html lang={locale}>
			<body className={`${geistSans.variable} ${geistMono.variable} antialiased bg-slate-100`}>
				<NextIntlClientProvider messages={messages}>
					<QueryProvider>
						<div className="min-h-screen w-full px-4 py-10">
							<div className="mx-auto flex w-full max-w-5xl flex-col items-center gap-8">
								{children}
								<LanguageSwitcher />
							</div>
						</div>
					</QueryProvider>
				</NextIntlClientProvider>
			</body>
		</html>
	)
}
