import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FieldApi, useForm } from "@tanstack/react-form";
import { useMutation, UseMutationResult } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { documentValidator } from "@/lib/formValidations";
import { checkMarkdownSpellingErrorAPIFn, FILE_KEY } from "@/apis/checkMarkdownSpellingErrorAPIFn";

export const Route = createFileRoute("/")({
	component: Index,
});

type ChangeInputParamType = FieldApi<
	{
		markdownfile: File | undefined;
	},
	"markdownfile",
	undefined,
	undefined,
	File
>;

function Index() {
	const [srcDoc, setSrcDoc] = useState<string>("");

	const form = useForm({
		defaultValues: { [FILE_KEY]: undefined } as { [FILE_KEY]: File | undefined },
	});

	const mutation = useMutation({
		mutationFn: checkMarkdownSpellingErrorAPIFn,
		onSuccess: (data) => setSrcDoc(data),
	});

	const validators = {
		onChange: documentValidator,
		onSubmit: documentValidator,
	};

	function handleFileInputChange(
		field: ChangeInputParamType,
		e: React.ChangeEvent<HTMLInputElement>
	) {
		const file = e?.target?.files?.[0];
		if (file) {
			field.handleChange(file);
		}
	}

	function handleFormSubmit(e: React.FormEvent<HTMLFormElement>) {
		e.preventDefault();
		e.stopPropagation();
		const file = form.getFieldValue(FILE_KEY);

		if (!form.state.canSubmit) return;

		if (file) {
			mutation.mutate({ file });
		}
	}
	return (
		<div className="container mx-auto p-6">
			<h1 className="text-2xl font-bold mb-4">Markdown Spell Checker</h1>
			<p className="text-gray-600 mb-6">
				Upload a markdown file (.md) to check for spelling errors. The results will be
				displayed below.
			</p>

			<div className="flex flex-col gap-6">
				<section className="w-full">
					<form
						className="flex flex-col sm:flex-row gap-4 items-start"
						onSubmit={handleFormSubmit}
					>
						<form.Field
							name={FILE_KEY}
							validators={validators}
							children={(field) => (
								<div>
									<Input
										className="w-3/4"
										type="file"
										name={FILE_KEY}
										accept=".md"
										placeholder=""
										onBlur={field.handleBlur}
										onChange={(e) => handleFileInputChange(field, e)}
									/>
									<em role="alert" className="text-red-300 text-sm mt-1 block">
										{field.state.meta.errors.join(", ")}
									</em>
									<FileAPIMutationStatus mutation={mutation} />
								</div>
							)}
						/>
						<Button type="submit" disabled={mutation.isPending} className="min-w-32">
							{mutation.isPending ? "Analyzing..." : "Check spelling"}
						</Button>
					</form>
				</section>

				{srcDoc.length > 0 && !mutation.isPending ? (
					<section className="w-full border rounded-xs overflow-hidden ml-4">
						<iframe
							title="Spell Check Results"
							srcDoc={srcDoc}
							className="w-full h-[70vh]"
						/>
					</section>
				) : null}
			</div>
		</div>
	);
}

export function FileAPIMutationStatus({
	mutation,
}: {
	mutation: UseMutationResult<unknown, Error, { file: File }, unknown>;
}) {
	return (
		<div>
			<p className="text-blue-300 text-left">
				{mutation.isPending ? "Analyzing your markdown file..." : null}
			</p>
			<p className="text-red-300 text-left capitalize">
				{mutation.isError ? mutation.error.message : null}
			</p>
			<p className="text-green-300 text-left capitalize">
				{mutation.isSuccess ? "Success" : null}
			</p>
		</div>
	);
}
