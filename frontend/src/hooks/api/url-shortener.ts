import { useMutation, useQuery } from "@tanstack/react-query"
import { createShortLink, getShortLinkInfoBySlug, type CreateShortLinkRequest } from "@/api/url-shortener"

export const useCreateShortLink = () => {
	return useMutation({
		mutationFn: (payload: CreateShortLinkRequest) => createShortLink(payload),
	})
}

export const useShortLinkInfoBySlug = (slug: string) => {
	return useQuery({
		queryKey: ["short-link-info", slug],
		queryFn: () => getShortLinkInfoBySlug(slug),
		enabled: !!slug,
		refetchInterval: (query) => query.state.data ? 3000 : false,
	})
}
