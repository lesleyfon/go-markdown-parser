import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import AppSideBarWrapper, { AppSideBar } from "@/components/app-sidebar";

const queryClient = new QueryClient();

export const Route = createRootRoute({
	component: () => (
		<AppSideBarWrapper>
			<QueryClientProvider client={queryClient}>
				<AppSideBar />
				<main className="flex-1">
					<div className="p-2 flex gap-2">
						<Link to="/" className="[&.active]:font-bold">
							Home
						</Link>{" "}
						<Link to="/login" className="[&.active]:font-bold">
							To login
						</Link>
						<Link to="/signup" className="[&.active]:font-bold">
							To signup
						</Link>
					</div>
					<Outlet />
				</main>
			</QueryClientProvider>
			<TanStackRouterDevtools />
		</AppSideBarWrapper>
	),
});
