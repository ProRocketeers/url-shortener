import { type FieldValues } from "react-hook-form"
import { type InputHTMLAttributes, useState } from "react"
import { useLocale } from "next-intl"
import { FormFieldController, type FormFieldProps } from "./FormFieldController"
import { Label } from "@/components/ui/reusable/Label"
import { Input } from "@/components/ui/reusable/Input"
import { Button } from "@/components/ui/reusable/Button"

type PresetOption = "1h" | "1d" | "7d" | "1m" | "1y" | "custom"

type Props<T extends FieldValues> = Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> & Omit<FormFieldProps<T>, "children"> & {
	checkboxLabel: string
	dateLabel?: string
	timeLabel?: string
}

export const ExpirationController = <T extends FieldValues>({ 
	label, 
	control, 
	name, 
	checkboxLabel,
	dateLabel = "Date",
	timeLabel = "Time",
	...props 
}: Props<T>) => {
	const locale = useLocale()
	const localeTag = locale === "cs" ? "cs-CZ" : locale === "en" ? "en-US" : locale
	const [selectedOption, setSelectedOption] = useState<PresetOption>("1d")

	return (
		<FormFieldController<T> control={control} name={name} label={label}>
			{(field) => {
				const isEnabled = Boolean(field.value)

				const toDate = (value?: string) => {
					if (!value) {
						return null
					}

					const parsed = new Date(value)
					if (Number.isNaN(parsed.getTime())) {
						return null
					}

					return parsed
				}

				const toDateInputValue = (value?: string) => {
					const date = toDate(value)
					if (!date) {
						return ""
					}

					const year = date.getFullYear()
					const month = String(date.getMonth() + 1).padStart(2, '0')
					const day = String(date.getDate()).padStart(2, '0')
					return `${year}-${month}-${day}`
				}

				const toCzDateInputValue = (value?: string) => {
					const date = toDate(value)
					if (!date) {
						return ""
					}

					const day = String(date.getDate()).padStart(2, '0')
					const month = String(date.getMonth() + 1).padStart(2, '0')
					const year = date.getFullYear()
					return `${day}.${month}.${year}`
				}

				const toTimeInputValue = (value?: string) => {
					const date = toDate(value)
					if (!date) {
						return ""
					}

					const hours = String(date.getHours()).padStart(2, '0')
					const minutes = String(date.getMinutes()).padStart(2, '0')
					return `${hours}:${minutes}`
				}

				const formatDisplayValue = (value?: string) => {
					const date = toDate(value)
					if (!date) {
						return ""
					}

					return new Intl.DateTimeFormat(localeTag, {
						dateStyle: 'medium',
						timeStyle: 'short',
					}).format(date)
				}

				const handleCheckboxChange = (checked: boolean) => {
					if (!checked) {
						field.onChange(undefined)
						setSelectedOption("1d")
						return
					}

					if (!field.value) {
						const tomorrow = new Date()
						tomorrow.setDate(tomorrow.getDate() + 1)
						tomorrow.setHours(9, 0, 0, 0)
						field.onChange(tomorrow.toISOString())
						setSelectedOption("1d")
					}
				}

				const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
					const dateValue = e.target.value
					if (!dateValue) {
						field.onChange(undefined)
						return
					}

					const source = toDate(field.value) ?? new Date()

					if (locale === "cs") {
						const matched = dateValue.trim().match(/^(\d{1,2})[.\/-](\d{1,2})[.\/-](\d{4})$/)
						if (!matched) {
							return
						}

						const day = Number(matched[1])
						const month = Number(matched[2])
						const year = Number(matched[3])
						source.setFullYear(year, month - 1, day)
					} else {
						const [year, month, day] = dateValue.split('-').map(Number)
						source.setFullYear(year, month - 1, day)
					}

					field.onChange(source.toISOString())
				}

				const handleTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
					const timeValue = e.target.value
					if (!timeValue) {
						return
					}
					
					const currentDate = toDate(field.value) ?? new Date()
					const [hours, minutes] = timeValue.split(':').map(Number)
					currentDate.setHours(hours, minutes, 0, 0)
					field.onChange(currentDate.toISOString())
				}

				const applyPreset = (hoursToAdd: number, option: PresetOption) => {
					const base = new Date()
					base.setMinutes(0, 0, 0)
					base.setHours(base.getHours() + hoursToAdd)
					field.onChange(base.toISOString())
					setSelectedOption(option)
				}

				const enableCustomInput = () => {
					setSelectedOption("custom")
					if (!field.value) {
						const tomorrow = new Date()
						tomorrow.setDate(tomorrow.getDate() + 1)
						tomorrow.setHours(9, 0, 0, 0)
						field.onChange(tomorrow.toISOString())
					}
				}

				const getOptionButtonClassName = (option: PresetOption) => {
					if (selectedOption === option) {
						return "border-[#051641] bg-[#051641] text-white hover:bg-[#0a3d7a] hover:text-white"
					}

					return "border-slate-200 bg-white text-slate-700"
				}

				return (
					<div className="space-y-3">
						<div className="flex items-center space-x-2">
							<input
								type="checkbox"
								id={`${name}-checkbox`}
								checked={isEnabled}
								onChange={(e) => handleCheckboxChange(e.target.checked)}
								className="h-4 w-4 rounded border-slate-300 text-[#051641] focus:ring-[#051641]"
							/>
							<Label htmlFor={`${name}-checkbox`} className="text-sm font-medium cursor-pointer">
								{checkboxLabel}
							</Label>
						</div>

						{isEnabled && (
							<div className="space-y-4 rounded-md border border-slate-200 bg-slate-50 p-3">
								<div className="flex flex-wrap gap-2">
									<Button type="button" variant="outline" size="sm" className={getOptionButtonClassName("1h")} onClick={() => applyPreset(1, "1h")}>
										1 hodina
									</Button>
									<Button type="button" variant="outline" size="sm" className={getOptionButtonClassName("1d")} onClick={() => applyPreset(24, "1d")}>
										1 den
									</Button>
									<Button type="button" variant="outline" size="sm" className={getOptionButtonClassName("7d")} onClick={() => applyPreset(24 * 7, "7d")}>
										7 dni
									</Button>
									<Button type="button" variant="outline" size="sm" className={getOptionButtonClassName("1m")} onClick={() => applyPreset(24 * 30, "1m")}>
										1 mesic
									</Button>
									<Button type="button" variant="outline" size="sm" className={getOptionButtonClassName("1y")} onClick={() => applyPreset(24 * 365, "1y")}>
										1 rok
									</Button>
									<Button type="button" variant="outline" size="sm" className={getOptionButtonClassName("custom")} onClick={enableCustomInput}>
										Vlastni
									</Button>
								</div>

								<div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
									<div>
										<Label className="text-xs text-slate-600 mb-1.5 block">
											{dateLabel}
										</Label>
										{locale === "cs" ? (
											<Input
												type="text"
												placeholder="dd.mm.rrrr"
												value={toCzDateInputValue(field.value)}
												onChange={handleDateChange}
												disabled={selectedOption !== "custom"}
												className="border-slate-200 bg-white"
											/>
										) : (
											<input
												type="date"
												lang={localeTag}
												key={`date-${localeTag}`}
												value={toDateInputValue(field.value)}
												onChange={handleDateChange}
												disabled={selectedOption !== "custom"}
												className="flex h-9 w-full rounded-md border border-slate-200 bg-white px-3 py-1 text-sm shadow-xs transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#051641] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
											/>
										)}
									</div>
									<div>
										<Label htmlFor={`${name}-time`} className="text-xs text-slate-600 mb-1.5 block">
											{timeLabel}
										</Label>
										<input
											id={`${name}-time`}
											type="time"
											lang={localeTag}
											value={toTimeInputValue(field.value)}
											onChange={handleTimeChange}
											disabled={selectedOption !== "custom"}
											className="flex h-9 w-full rounded-md border border-slate-200 bg-white px-3 py-1 text-sm shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#051641] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
										/>
									</div>
								</div>
							</div>
						)}
					</div>
				)
			}}
		</FormFieldController>
	)
}
