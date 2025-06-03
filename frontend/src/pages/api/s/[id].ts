import type { APIRoute } from "astro";
import api from "../../../utils/api";

export const GET: APIRoute = async ({ params }) => {
  const code = params.id;

  try {
    const res = await api(`s/${code}`)();
    const data = await res.json();

    return new Response(JSON.stringify(data), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    return new Response(JSON.stringify({ error: "Error interno" }), {
      status: 500,
    });
  }
};
