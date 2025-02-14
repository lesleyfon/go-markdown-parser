import { API_BASE_URL } from "@/lib/constants";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useParams } from "@tanstack/react-router";
import { Loader2 } from "lucide-react";

export interface FileDetail {
	created_at: string;
	file_content: string;
	file_name: string;
	html_content: string;
	updated_at: string;
	user_id: string;
}

export interface FileResponse {
	file?: FileDetail;
	message: string;
	status: number;
}

export const Route = createFileRoute("/files/$file-id")({
	component: FileDetails,
});

function FileDetails() {
	const fileId = useParams({
		select: (params) => params["file-id"],
		from: "/files/$file-id",
	});

	const { data, isLoading, isError, error } = useQuery({
		queryKey: ["file-by-id", fileId],
		queryFn: async (): Promise<FileResponse> => {
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
		},
		enabled: !!fileId,
	});

	if (isLoading) {
		return (
			<div className="flex items-center justify-center h-full">
				<Loader2 className="animate-spin" />
			</div>
		);
	}

	if (isError) {
		return (
			<div className="flex items-center justify-center h-full">
				<div className="text-red-500">Error: {error?.message}</div>
			</div>
		);
	}

	return (
		<div>
			<h1 className="text-2xl font-bold mb-4">
				{data?.file?.file_name ? `File Details: ${data?.file?.file_name}` : data?.message}
			</h1>
			{data?.file?.html_content ? (
				<HtmlFileViewer htmlContent={data?.file.html_content} />
			) : null}
		</div>
	);
}

function HtmlFileViewer({ htmlContent }: { htmlContent: string }) {
	const fromBase64 = atob(htmlContent);
	return <iframe srcDoc={fromBase64} className="w-full h-[80vh]" />;
}
