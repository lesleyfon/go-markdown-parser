import React from "react";
import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarHeader,
	SidebarProvider,
} from "@/components/ui/sidebar";

export function AppSideBar() {
	return (
		<Sidebar>
			<SidebarHeader>
				<p>App</p>
			</SidebarHeader>
			<SidebarContent>
				<div>
					<p>Sidebar content</p>
				</div>
			</SidebarContent>
			<SidebarFooter />
		</Sidebar>
	);
}

export default function AppSideBarWrapper({ children }: { children?: React.ReactNode }) {
	return <SidebarProvider>{children}</SidebarProvider>;
}
