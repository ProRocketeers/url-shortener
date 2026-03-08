import { type ReactNode } from "react"
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "./Form"
import type { Control, ControllerRenderProps, FieldValues, Path } from "react-hook-form"

export type FormFieldProps<T extends FieldValues> = {
	children: (field: ControllerRenderProps<T>) => ReactNode;
	control: Control<T>;
	name: Path<T>;
	description?: string;
	label: string;
};

export const FormFieldController = <T extends FieldValues>({
	children,
	control,
	name,
	description,
	label,
}: FormFieldProps<T>) => {
	return (
		<FormField
			control={control}
			name={name}
			render={(field) => (
				<FormItem>
					<FormLabel>{label}</FormLabel>
					<FormControl>{children(field.field)}</FormControl>
					{description && <FormDescription>{description}</FormDescription>}
					<FormMessage />
				</FormItem>
			)}
		/>
	)
}
