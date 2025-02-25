import { emailValidator, passwordValidator } from "@/lib/formValidations";
import { AppEnv } from "@/types";



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
const APP_ENV = import.meta.env.VITE_APP_ENV
const parsedAppEnv = JSON.parse(APP_ENV) as AppEnv
export const API_BASE_URL = parsedAppEnv.API_BASE_URL
