import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

export const Route = createRootRoute({
	component: () => (
		<>
			<div className="p-2 flex gap-2">
				<Link to="/" className="[&.active]:font-bold">
					Home
				</Link>{" "}
				<Link to="/login" className="[&.active]:font-bold">
					To login
				</Link>
			</div>
			<QueryClientProvider client={queryClient}>
				<Outlet />
			</QueryClientProvider>
			<TanStackRouterDevtools />
		</>
	),
});
