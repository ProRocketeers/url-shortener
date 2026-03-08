import { zodResolver } from "@hookform/resolvers/zod"
import { ReactNode } from "react"
import {
	FormProvider,
	useForm,
	DefaultValues,
	type SubmitHandler,
	type UseFormReturn,
	type UseFormProps,
} from "react-hook-form"
import { z } from "zod"

type Props<T extends z.ZodObject<z.ZodRawShape>> = {
	children: (form: UseFormReturn<z.input<T>, unknown, z.output<T>>) => ReactNode;
	schema: T;
	defaultValues?: DefaultValues<z.input<T>>;
	onSubmit: SubmitHandler<z.output<T>>;
} & UseFormProps<z.input<T>>;

export const FormContext = <T extends z.ZodObject<z.ZodRawShape>>({
	children,
	schema,
	defaultValues,
	onSubmit,
	...formProps
}: Props<T>) => {
	const form = useForm<z.input<T>, unknown, z.output<T>>({
		resolver: zodResolver(schema),
		defaultValues,
		...formProps,
	})

	return (
		<FormProvider {...form}>
			<form
				onSubmit={(e) => {
					e.preventDefault()
					form.handleSubmit(onSubmit)(e)
				}}
			>
				{children(form)}
			</form>
		</FormProvider>
	)
}
