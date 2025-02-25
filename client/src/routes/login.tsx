import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";

import { loginFn } from "@/apis/authAPIFn";
import AuthFormFieldGroup from "@/components/auth-form-field-group";
import { AuthenticationStatus } from "@/components/authentication-status";
import { handleSignupLoginFormSubmit } from "@/lib/utils";

export const Route = createFileRoute("/login")({
	component: Login,
});

export default function Login() {
	const navigate = useNavigate();

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
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: ["login"] });

			const token = data.token;
			localStorage.setItem("auth-token", token);

			// Navigate to home
			navigate({ to: "/", reloadDocument: true });
		},
	});

	return (
		<div>
			<h1>Login</h1>
			<form
				onSubmit={(e) => handleSignupLoginFormSubmit(e, form, mutation)}
				className=" flex flex-col gap-2"
			>
				<AuthFormFieldGroup form={form} />
				<button type="submit">Submit</button>
				<AuthenticationStatus mutation={mutation} />
			</form>
		</div>
	);
}
