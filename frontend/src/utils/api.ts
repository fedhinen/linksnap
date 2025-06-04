export const BASE_URL =
  import.meta.env.PUBLIC_API_URL ??
  import.meta.env.API_URL ??
  "http://localhost:5271/api";

const api = (path: string): ((options?: RequestInit) => Promise<Response>) => {
  const sanitizedPath = path.replace(/^\/+/, "");
  return (options?: RequestInit) => {
    console.log({ url: `${BASE_URL}/${sanitizedPath}`, options });
    return fetch(`${BASE_URL}/${sanitizedPath}`, options);
  };
};

export default api;
