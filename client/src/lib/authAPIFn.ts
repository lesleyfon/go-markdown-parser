const myHeaders = new Headers();
	myHeaders.append("Content-Type", "application/json");

const requestOptions = {
	method: "POST",
	headers: myHeaders,
	
	redirect: "follow" as RequestRedirect,
};

export async function loginFn({ email, password }: { email: string; password: string }) {
	const body = JSON.stringify({
		email,
		password,
	});


	const response = await fetch("http://0.0.0.0:8080/auth/v1/login"	, {...requestOptions, body});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.message);
	}

	const data = await response.json();

	return data;
}

export async function signupFn({ email, password }: { email: string; password: string }) {

	const body = JSON.stringify({
		email,
		password,
	});

	const response = await fetch("http://0.0.0.0:8080/auth/v1/signup", {
		...requestOptions,
		body,
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.message);
	}

	const data = await response.json();

	return data;
}