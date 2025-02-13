export const FILE_KEY = "markdownfile";
export const checkMarkdownSpellingErrorAPIFn = async ({ file }: { file: File }) => {
	console.log(file);
	const formdata = new FormData();
	const filename = file.name;

  formdata.append(FILE_KEY, file, filename);

  const requestOptions = {
    method: "POST",
    body: formdata,
    redirect: "follow" as RequestRedirect,
  };

  const response = await fetch("http://0.0.0.0:8080/api/v1/markdown", requestOptions);

  if (!response.ok) {
    const errorResponse = await response.json();
    throw Error(errorResponse.message);
  }

	const parseMarkdownHtml: string = await response.text();

	return parseMarkdownHtml;
};
