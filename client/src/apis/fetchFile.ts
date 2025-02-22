import { API_BASE_URL } from "@/lib/constants";
import { FileResponse } from "@/routes/files/$file-id";

export async function fetchFile(fileId: string): Promise<FileResponse> {
	const token = localStorage.getItem("auth-token");
	if (!token) {
		return {
			message: "No token found, please login again",
			status: 401,
		};
	}

	const res = await fetch(`${API_BASE_URL}/api/v1/markdown/files/${fileId}`, {
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


export async function fetchAllFiles(){
	const token = localStorage.getItem("auth-token");
			if (!token) {
				throw new Error("No token found, please login again");
			}
			const res = await fetch(`${API_BASE_URL}/api/v1/markdown/files`, {
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