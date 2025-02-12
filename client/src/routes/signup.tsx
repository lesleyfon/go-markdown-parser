import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { signupFn } from "../lib/authAPIFn";
import { inputFieldData } from "../lib/constants";
import { useMemo } from "react";
import { AuthenticationStatus } from "../components/authentication-status";

export const Route = createFileRoute("/signup")({
	component: Signup,
});

export default function Signup() {
	const form = useForm({
		defaultValues: {
			email: "",
			password: "",
		},
	});

	// Access the client
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: signupFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["signup"] });
		},
	});

	const inputFields = useMemo(
		() =>
			inputFieldData.map((fieldData) => (
				<form.Field
					name={fieldData.name as "email" | "password"}
					validators={{
						onChange: fieldData.validators.onChange,
					}}
					children={(field) => (
						<>
							<input
								type={fieldData.type}
								className=" outline"
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
			)),
		[form]
	);

	return (
		<div>
			<h1>Signup</h1>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					e.stopPropagation();

					mutation.mutate({
						email: form.getFieldValue("email"),
						password: form.getFieldValue("password"),
					});
					form.handleSubmit();
				}}
				className=" flex flex-col gap-2"
			>
				{inputFields}
				<button type="submit">Submit</button>
				<AuthenticationStatus mutation={mutation} />
			</form>
		</div>
	);
}
