import { DetailPage } from "@/components/client/url-shortener/DetailPage"

type DetailPageProps = {
	params: Promise<{ slug: string; locale: string }>
}

export default async function Detail({ params }: DetailPageProps) {
	const { slug } = await params

	return <DetailPage slug={slug} />
}
