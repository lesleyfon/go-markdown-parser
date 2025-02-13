import z from "zod";
export const emailValidator = z.string().email("Invalid Email").min(2, {
	message: "Email must be at least 2 characters.",
});
export const passwordValidator = z.string().min(6, {
	message: "Password must be of length 6 or greater",
});

export const documentValidator = z
	.instanceof(File)
	.refine((file) => ["text/markdown"].includes(file.type), { message: "Invalid file type" });
