import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { loginFn } from "../lib/authAPIFn";
import { AuthenticationStatus } from "../components/authentication-status";
import AuthFormFieldGroup from "@/components/auth-form-field-group";

export const Route = createFileRoute("/login")({
	component: Login,
});

export default function Login() {
	const form = useForm({
		defaultValues: {
			email: "",
			password: "",
		},
	});

	// Access the client
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: loginFn,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["login"] });
		},
	});

	return (
		<div>
			<h1>Login</h1>
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
