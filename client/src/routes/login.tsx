import z from "zod";
import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "@tanstack/react-form";

export const Route = createFileRoute("/login")({
	component: Login,
});

const zodEmailValidator = z.string().email("Invalid Email").min(2, {
	message: "Email must be at least 2 characters.",
});
const emailPasswordValidator = z.string().min(6, {
	message: "Password must be of length 6 or greater",
});

export default function Login() {
	const form = useForm({
		defaultValues: {
			email: "",
			password: "",
		},
	});
	return (
		<div>
			<h1>login</h1>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					e.stopPropagation();
					console.log(form.getFieldValue("email"));
					console.log(form.getFieldValue("password"));
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
							<em role="alert">{field.state.meta.errors.join(", ")} </em>
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
							<em role="alert">{field.state.meta.errors.join(", ")} </em>
						</>
					)}
				/>
				<button type="submit">Submit</button>
			</form>
		</div>
	);
}
