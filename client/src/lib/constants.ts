import { emailValidator, passwordValidator } from "./formValidations";

export const inputFieldData = [
	{
		name: "email",
		placeholder: "Email",
		type: "email",
		validators: {
			onChange: emailValidator,
		},
	},
	{
		name: "password",
		placeholder: "Password",
		type: "password",
		validators: {
			onChange: passwordValidator,
		},
	},
];
