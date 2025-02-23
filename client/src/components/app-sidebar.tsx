import { useQuery } from "@tanstack/react-query";
import { Link, useLocation } from "@tanstack/react-router";
import { File, Loader2 } from "lucide-react";
import React, { useMemo } from "react";

import { authenticateUserFn } from "@/apis/authAPIFn";
import { fetchAllFiles } from "@/apis/fetchFile";
import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarGroupContent,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuButton,
	SidebarProvider,
	SidebarMenuItem,
} from "@/components/ui/sidebar";
import { AuthResponse, FilesResponse, UnAuthenticateUserResponse } from "@/types";

export function AppSideBar() {
	const location = useLocation();

	const {
		data: authQueryData,
		isLoading: authQueryIsLoading,
		error: authQueryError,
	} = useQuery<AuthResponse | UnAuthenticateUserResponse>({
		queryKey: ["authenticate"],
		queryFn: authenticateUserFn,
	});

	const filesQuery = useQuery<FilesResponse>({
		queryKey: ["all-files", authQueryData?.isAuthenticated],
		queryFn: fetchAllFiles,
		enabled: authQueryData?.isAuthenticated,
		retry: false,
	});

	const sideBarMenuItems = useMemo(
		() => (
			<SidebarGroupContent>
				<SidebarMenu>
					{filesQuery?.data?.files.map((item) => (
						<SidebarMenuItem key={item.file_name}>
							<SidebarMenuButton asChild>
								<Link
									to="/files/$file-id"
									params={{ "file-id": item._id }}
									className="[&.active]:font-bold"
								>
									<File />
									<span>{item.file_name}</span>
								</Link>
							</SidebarMenuButton>
						</SidebarMenuItem>
					))}
				</SidebarMenu>
			</SidebarGroupContent>
		),
		[filesQuery.data?.files]
	);
	// DO not render sidebar if a user is not authenticated OR if they are on the login or signup page
	const pathname = location.pathname;
	const sidebarPathBlacklist = ["/", "/login", "/signup"];
	if (
		sidebarPathBlacklist.includes(pathname) &&
		(authQueryIsLoading || !authQueryData?.isAuthenticated)
	) {
		return;
	}

	// Add loading states
	if (authQueryIsLoading) {
		return (
			<div className="flex items-center justify-center h-full">
				<Loader2 className="animate-spin" />
			</div>
		);
	}

	if (authQueryError || filesQuery.error) {
		return (
			<div className="p-4 text-red-500">
				Error: {authQueryError?.message || filesQuery.error?.message}
			</div>
		);
	}

	return (
		<Sidebar>
			<SidebarHeader>
				<h1 className="text-2xl font-bold">All Files</h1>
			</SidebarHeader>
			<SidebarContent>{sideBarMenuItems}</SidebarContent>
			<SidebarFooter />
		</Sidebar>
	);
}

export default function AppSideBarWrapper({ children }: { children?: React.ReactNode }) {
	return <SidebarProvider>{children}</SidebarProvider>;
}
