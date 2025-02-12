export async function loginFn({ email, password }: { email: string; password: string }) {
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