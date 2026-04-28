import { siteConfig } from "@/config/site";

export const applyMetadata = () => {
	document.title = siteConfig.metadata.title;

	let description = document.querySelector<HTMLMetaElement>(
		'meta[name="description"]',
	);

	if (!description) {
		description = document.createElement("meta");
		description.name = "description";
		document.head.append(description);
	}

	description.content = siteConfig.metadata.description;
};
