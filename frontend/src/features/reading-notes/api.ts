import type { ReadingNoteDetail, ReadingNoteListResponse } from "./types";

const request = async <T>(path: string): Promise<T> => {
	const response = await fetch(path);
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
