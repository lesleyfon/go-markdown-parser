import z from "zod";
import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient, UseMutationResult } from "@tanstack/react-query";

export const Route = createFileRoute("/login")({
	component: Login,
});

const zodEmailValidator = z.string().email("Invalid Email").min(2, {
	message: "Email must be at least 2 characters.",
});
const emailPasswordValidator = z.string().min(6, {
	message: "Password must be of length 6 or greater",
});

async function loginFn({ email, password }: { email: string; password: string }) {
	const myHeaders = new Headers();
	myHeaders.append("Content-Type", "application/json");

	const requestOptions = {
		method: "POST",
		headers: myHeaders,
		body: JSON.stringify({
			email,
			password,
		}),
		redirect: "follow" as RequestRedirect,
	};

	const response = await fetch("http://0.0.0.0:8080/auth/v1/login", requestOptions);

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.message);
	}

	const data = await response.json();

	return data;
}

function AuthenticationStatus({
	mutation,
}: {
	mutation: UseMutationResult<unknown, Error, { email: string; password: string }, unknown>;
}) {
	return (
		<div>
			<p className="text-blue-300 text-left">{mutation.isPending && "Loading..."}</p>
			<p className="text-red-300 text-left">{mutation.isError && mutation.error.message}</p>
			<p className="text-green-300 text-left">{mutation.isSuccess && "Success"}</p>
		</div>
	);
}

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
			<h1>login</h1>
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
				<form.Field
					name="email"
					validators={{
						onChange: zodEmailValidator,
					}}
					children={(field) => (
						<>
							<input
								type="email"
								className=" outline"
								name={field.name}
								value={field.state.value}
								onBlur={field.handleBlur}
								onChange={(e) => field.handleChange(e.target.value)}
							/>
							<em role="alert" className=" text-red-300 text-left">
								{field.state.meta.errors.join(", ")}{" "}
							</em>
						</>
					)}
				/>
				<form.Field
					name="password"
					validators={{
						onChange: emailPasswordValidator,
					}}
					children={(field) => (
						<>
							<input
								type="password"
								className=" outline"
								name={field.name}
								value={field.state.value}
								onBlur={field.handleBlur}
								onChange={(e) => field.handleChange(e.target.value)}
							/>
							<em role="alert" className=" text-red-300 text-left">
								{field.state.meta.errors.join(", ")}{" "}
							</em>
						</>
					)}
				/>
				<button type="submit">Submit</button>
				<AuthenticationStatus mutation={mutation} />
			</form>
		</div>
	);
}
