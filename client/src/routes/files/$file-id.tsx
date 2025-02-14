import { createFileRoute, useParams } from "@tanstack/react-router";

export const Route = createFileRoute("/files/$file-id")({
	component: RouteComponent,
});

function RouteComponent() {
	const fileId = useParams({
		select: (params) => params["file-id"],
		from: "/files/$file-id",
	});
	console.log(fileId);

	return <div>FileId: {fileId}!</div>;
}
