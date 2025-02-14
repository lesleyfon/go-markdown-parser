import React, { useMemo } from "react";
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
import { useQuery } from "@tanstack/react-query";
import { File, Loader2 } from "lucide-react";
import { Link, useLocation } from "@tanstack/react-router";
import { API_BASE_URL } from "@/lib/constants";

export interface File {
	_id: string;
	created_at: string;
	file_content: string;
	file_name: string;
	updated_at: string;
	user_id: string;
}

export interface FilesResponse {
	files: File[];
	message: string;
	status: number;
}

export function AppSideBar() {
	const token = localStorage?.getItem("auth-token") ?? "";
	const location = useLocation();

	const authQuery = useQuery({
		queryKey: ["authenticate"],
		queryFn: async () => {
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
		},
	});

	const filesQuery = useQuery<FilesResponse>({
		queryKey: ["all-files", authQuery.data?.isAuthenticated],
		queryFn: async () => {
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
		},
		enabled: authQuery.data?.isAuthenticated,
	});

	const sideBarMenuItems = useMemo(
		() => (
			<SidebarGroupContent>
				<SidebarMenu>
					{filesQuery?.data?.files.map((item) => (
						<SidebarMenuItem key={item.file_name}>
							<SidebarMenuButton asChild>
								<Link to="/files/$file-id" params={{ "file-id": item._id }}>
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

	// Add loading states
	if (authQuery.isLoading) {
		return (
			<div className="flex items-center justify-center h-full">
				<Loader2 className="animate-spin" />
			</div>
		);
	}

	if (authQuery.error || filesQuery.error) {
		return (
			<div className="p-4 text-red-500">
				Error: {authQuery.error?.message || filesQuery.error?.message}
			</div>
		);
	}

	// DO not render sidebar if a user is not authenticated OR if they are on the login or signup page
	const pathname = location.pathname;
	if (!authQuery.data?.isAuthenticated || pathname === "/login" || pathname === "/signup") {
		return null;
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
