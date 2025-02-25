import { ReactFormExtendedApi } from "@tanstack/react-form";
import { UseMutationResult } from "@tanstack/react-query";
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

import { SignupLoginFormData, AuthUser } from "@/types";

/**
 * Merges multiple class names using clsx and tailwind-merge.
 * 
 * @param inputs - The class values to merge
 * @returns The merged class names
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Handles the submission of a signup or login form.
 * 
 * @param event - The form submission event
 * @param form - The form instance
 * @param mutation - The mutation instance
 */
export const handleSignupLoginFormSubmit = (
	event: React.FormEvent<HTMLFormElement>,
	form: ReactFormExtendedApi<SignupLoginFormData, undefined>,
	mutation: UseMutationResult<AuthUser, Error, SignupLoginFormData, unknown>
): void => {
	event.preventDefault();
	event.stopPropagation();

	mutation.mutate({
		email: form.getFieldValue("email"),
		password: form.getFieldValue("password"),
	});
	form.handleSubmit();
};
