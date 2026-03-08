import { z } from "zod"

export const createUrlShortenerSchema = (t: (key: string) => string) =>
	z.object({
		originalUrl: z
			.string()
			.min(1, t("urlRequired"))
			.url(t("urlInvalid")),
		slug: z
			.string()
			.optional()
			.refine((val) => !val || /^[a-zA-Z0-9-_]+$/.test(val), {
				message: t("slugInvalid"),
			}),
		expiresAt: z.string().optional(),
	})

export type UrlShortenerFormValues = z.infer<ReturnType<typeof createUrlShortenerSchema>>

export const defaultValues: UrlShortenerFormValues = {
	originalUrl: "",
	slug: undefined,
	expiresAt: undefined,
}
