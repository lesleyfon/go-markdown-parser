import { QueryClient, QueryClientProvider, useQuery } from "@tanstack/react-query";
import { createRootRoute, Link, Outlet, useNavigate } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

import { authenticateUserFn } from "@/apis/authAPIFn";
import AppSideBarWrapper, { AppSideBar } from "@/components/app-sidebar";
import { AuthResponse, UnAuthenticateUserResponse } from "@/types";

const queryClient = new QueryClient();

function NavHeader() {
	const navigate = useNavigate();
	const { data: authQueryData } = useQuery<AuthResponse | UnAuthenticateUserResponse>({
		queryKey: ["authenticate"],
		queryFn: authenticateUserFn,
	});

	const isAuthenticated = authQueryData?.isAuthenticated;

	function handleLogOut() {
		localStorage.removeItem("auth-token");

		navigate({
			to: "/",
			reloadDocument: true,
		});
	}

	return (
		<div className="p-2 flex gap-2 border-b border-gray-200">
			<Link to="/" className="[&.active]:font-bold">
				Home
			</Link>{" "}
			{isAuthenticated ? (
				<button className=" cursor-pointer" onClick={handleLogOut}>
					Logout
				</button>
			) : (
				<>
					<Link to="/login" className="[&.active]:font-bold">
						To login
					</Link>
					<Link to="/signup" className="[&.active]:font-bold">
						To signup
					</Link>
				</>
			)}
		</div>
	);
}

export const Route = createRootRoute({
	component: () => (
		<AppSideBarWrapper>
			<QueryClientProvider client={queryClient}>
				<AppSideBar />
				<main className="flex-1">
					<NavHeader />
					<Outlet />
				</main>
			</QueryClientProvider>
			<TanStackRouterDevtools />
		</AppSideBarWrapper>
	),
});
