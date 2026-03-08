import { apiClient } from "@/api/client"

export type CreateShortLinkRequest = {
	originalUrl: string
	slug?: string
	expiresAt?: string // RFC3339 format, e.g., "2036-03-09T12:00:00Z"
}

export type ShortenUrlResponse = {
	shortUrl: string
	slug: string
}

export type ShortLinkInfoResponse = {
	originalUrl: string
	shortUrl: string
	clickCount: number
	expiresAt?: string
	createdAt?: string
	updatedAt?: string
}

export type ShortLinkResponse = {
	id: number
	originalUrl: string
	shortUrl: string
	slug: string
	expiresAt?: string
}

export type ShortLinkDetailResponse = {
	id: number
	originalUrl: string
	shortUrl: string
	slug: string
	expiresAt?: string
	// Additional stats will be added here in future
}

export const createShortLink = async (payload: CreateShortLinkRequest): Promise<ShortenUrlResponse> => {
	const response = await apiClient.post<ShortenUrlResponse>("/v1/shorten", payload)
	return response.data
}

export const getShortLinkInfoBySlug = async (slug: string): Promise<ShortLinkInfoResponse> => {
	const response = await apiClient.get<ShortLinkInfoResponse>(`/v1/info/${slug}`)
	return response.data
}
