const rawBase = import.meta.env.BASE_URL || "/";

export const appBase = rawBase.endsWith("/") ? rawBase : `${rawBase}/`;
const appBasePrefix = appBase.slice(0, -1);

export const appPath = (path = "/") => {
	const normalizedPath = path.startsWith("/") ? path.slice(1) : path;
	return `${appBase}${normalizedPath}`;
};

export const apiPath = (path: string) => appPath(path);

export const routePath = (pathname = window.location.pathname) => {
	if (appBase !== "/" && pathname === appBasePrefix) {
		return "/";
	}
	if (appBase !== "/" && pathname.startsWith(appBase)) {
		const stripped = pathname.slice(appBase.length - 1);
		return stripped || "/";
	}
	return pathname || "/";
};
