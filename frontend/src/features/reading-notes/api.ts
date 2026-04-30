import type { ReadingNoteDetail, ReadingNoteListResponse } from "./types";
import { apiPath } from "@/utils/app-base";

const request = async <T>(path: string): Promise<T> => {
	const response = await fetch(apiPath(path));
	if (!response.ok) {
		throw new Error(`Request failed: ${response.status}`);
	}
	return response.json() as Promise<T>;
};

export const listReadingNotes = () =>
	request<ReadingNoteListResponse>("/api/reading-notes");

export const getReadingNote = (slug: string) =>
	request<ReadingNoteDetail>(
		`/api/reading-notes/${encodeURIComponent(slug)}`,
	);
