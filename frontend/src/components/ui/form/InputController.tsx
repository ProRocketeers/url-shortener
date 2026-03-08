import { type FieldValues } from "react-hook-form"
import { Input } from "@/components/ui/reusable/Input"
import { type InputHTMLAttributes } from "react"
import { FormFieldController, type FormFieldProps } from "./FormFieldController"

type Props<T extends FieldValues> = InputHTMLAttributes<HTMLInputElement> & Omit<FormFieldProps<T>, "children">;

export const InputController = <T extends FieldValues>({ label, control, name, ...props }: Props<T>) => {
	return (
		<FormFieldController<T> control={control} name={name} label={label}>
			{(field) => {
				const { value, ...fieldWithoutValue } = field
				// Convert undefined/null to empty string for controlled inputs
				const displayValue = value === undefined || value === null ? "" : value
				
				return (
					<Input
						placeholder={props.placeholder}
						{...fieldWithoutValue}
						{...props}
						value={displayValue}
						onChange={(e) => {
							const newValue = props.type === "number" ? parseFloat(e.target.value) : e.target.value
							field.onChange(newValue || undefined)
						}}
					/>
				)
			}}
		</FormFieldController>
	)
}
