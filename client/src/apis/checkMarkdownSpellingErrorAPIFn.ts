import { API_BASE_URL } from "@/lib/constants";

export const FILE_KEY = "markdownfile";
export const checkMarkdownSpellingErrorAPIFn = async ({ file }: { file: File }) => {
	const formdata = new FormData();
	const filename = file.name;

  formdata.append(FILE_KEY, file, filename);

  const requestOptions = {
		method: "POST",
		body: formdata,
		redirect: "follow" as RequestRedirect,
  } as RequestInit;
  
	const token = localStorage.getItem("auth-token");
	if (token) {
		requestOptions.headers = {
			Authorization: `Bearer ${token}`,
		};
	}

	const response = await fetch(`${API_BASE_URL}/api/v1/markdown`, requestOptions);

  if (!response.ok) {
    const errorResponse = await response.json();
    throw Error(errorResponse.message);
  }

	const parseMarkdownHtml: string = await response.text();

	return parseMarkdownHtml;
};
