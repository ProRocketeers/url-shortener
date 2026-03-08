import { defineRouting } from 'next-intl/routing'
import { createNavigation } from 'next-intl/navigation'

export const routing = defineRouting({
	// Supported languages
	locales: ['cs', 'en'],

	// Default language
	defaultLocale: 'cs',

	// Prefix locale in URL (e.g., /en/page)
	localePrefix: 'as-needed', // 'cs' won't have prefix, 'en' will have /en
})

// Export navigation functions with localization
export const { Link, redirect, usePathname, useRouter, getPathname } =
	createNavigation(routing)
