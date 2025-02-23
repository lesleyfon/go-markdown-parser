import { ReactFormExtendedApi } from "@tanstack/react-form";
import { memo } from "react";

import { inputFieldData } from "@/lib/constants";

const AuthFormFieldGroup = memo(
	({ form }: { form: ReactFormExtendedApi<{ email: string; password: string }> }) => {
		return inputFieldData.map((fieldData) => (
			<form.Field
				name={fieldData.name as "email" | "password"}
				validators={{
					onChange: fieldData.validators.onChange,
				}}
				children={(field) => (
					<>
						<input
							type={fieldData.type}
							className="outline"
							name={field.name}
							value={field.state.value}
							onBlur={field.handleBlur}
							onChange={(e) => field.handleChange(e.target.value)}
						/>
						<em role="alert" className=" text-red-300 text-left">
							{field.state.meta.errors.join(", ")}
						</em>
					</>
				)}
			/>
		));
	}
);
export default AuthFormFieldGroup;
