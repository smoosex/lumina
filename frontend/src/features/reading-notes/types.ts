export interface ReadingNoteListItem {
	slug: string;
	title: string;
	isbn: string;
	doubanId?: string;
	bookTitle: string;
	author?: string;
	translator?: string;
	publisher?: string;
	producer?: string;
	bookPubDate?: string;
	bookPages?: number;
	coverUrl?: string;
	doubanUrl?: string;
	excerpt?: string;
	tags: string[];
	rating?: number;
	readAt?: string;
	publishedAt?: string;
}

export interface ReadingNoteDetail extends ReadingNoteListItem {
	contentHtml: string;
	contentMd?: string;
}

export interface ReadingNoteListResponse {
	items: ReadingNoteListItem[];
}
