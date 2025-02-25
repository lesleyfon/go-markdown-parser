import { AuthenticationStatusProps } from "@/types";

export function AuthenticationStatus({ mutation }: AuthenticationStatusProps) {
	return (
		<div>
			<p className="text-blue-300 text-left">{mutation.isPending && "Loading..."}</p>
			<p className="text-red-300 text-left">{mutation.isError && mutation.error.message}</p>
			<p className="text-green-300 text-left">{mutation.isSuccess && "Success"}</p>
		</div>
	);
}
