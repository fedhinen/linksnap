export type ShortUrl = {
  id: number;
  userId: string;
  originalUrl: string;
  shortenedCode: string;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
};
