import type { APIRoute } from "astro";
import api from "../../utils/api";

export const POST: APIRoute = async ({ request, locals }) => {
  const token = locals.authToken;

  if (!token) {
    return new Response(null, { status: 302, headers: { Location: "/" } });
  }

  const formData = await request.formData();
  const url = formData.get("url");

  if (typeof url !== "string" || !url.startsWith("http")) {
    return new Response("URL inválida", { status: 400 });
  }

  await api("shorturl/")({
    method: "POST",
    body: JSON.stringify({ url }),
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return new Response(null, {
    status: 302,
    headers: { Location: "/dashboard" }, // Redirige de vuelta al dashboard
  });
};
