import { API_BASE_URL } from "@/lib/constants";

export interface AuthResponse {
	email: string;
	refreshedToken: string;
	token: string;
	userId: string;
}

const myHeaders = new Headers();
	myHeaders.append("Content-Type", "application/json");

const requestOptions = {
	method: "POST",
	headers: myHeaders,
	
	redirect: "follow" as RequestRedirect,
};

export async function loginFn({ email, password }: { email: string; password: string }): Promise<AuthResponse> {
	const body = JSON.stringify({
		email,
		password,
	});


	const response = await fetch(`${API_BASE_URL}/auth/v1/login`	, {...requestOptions, body});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.message);
	}

	const data = await response.json();

	return data;
}

export async function signupFn({ email, password }: { email: string; password: string }): Promise<AuthResponse> {

	const body = JSON.stringify({
		email,
		password,
	});

	const response = await fetch(`${API_BASE_URL}/auth/v1/signup`, {
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


 interface UnAuthenticateUserResponse {
	message: string;
	status: number;
	isAuthenticated: boolean;
}

export async function authenticateUserFn(): Promise<AuthResponse | UnAuthenticateUserResponse> {
	const token = localStorage.getItem("auth-token");
	if (!token) {
		return {
			message: "No token found, please login again",
			status: 401,
			isAuthenticated: false,
		};
	}
	const res = await fetch(`${API_BASE_URL}/auth/v1/authenticate`, {
		headers: {
			Authorization: `Bearer ${token}`,
		},
	});

	if (!res.ok) {
		const error = await res.json();
		throw new Error(error.message);
	}
	const data = await res.json();

	return data;
}