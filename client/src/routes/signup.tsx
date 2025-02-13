import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { signupFn } from "../lib/authAPIFn";
import { AuthenticationStatus } from "../components/authentication-status";
import AuthFormFieldGroup from "@/components/auth-form-field-group";

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

	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: signupFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["signup"] });
		},
	});

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
				<AuthFormFieldGroup form={form} />
				<button type="submit">Submit</button>
				<AuthenticationStatus mutation={mutation} />
			</form>
		</div>
	);
}
