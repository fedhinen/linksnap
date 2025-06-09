import type { APIRoute } from "astro";
import api from "../../../utils/api";

export const DELETE: APIRoute = async ({ params, locals }) => {
  const token = locals.authToken;

  if (!token) {
    return new Response(null, {
      status: 302,
      headers: { Location: "/dashboard" },
    });
  }

  const id = params.id;

  if (!id) {
    return new Response(null, {
      status: 302,
      headers: { Location: "/dashboard" },
    });
  }

  await api(`shorturl/${id}`)({
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  return new Response(JSON.stringify({ redirect: "/dashboard" }), {
    status: 200,
    headers: {
      "Content-Type": "application/json",
    },
  });
};
