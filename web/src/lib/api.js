const API_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080";

export async function getPosts(params = {}) {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== "") {
      searchParams.set(key, value);
    }
  });

  const response = await fetch(
    `${API_URL}/api/posts?${searchParams.toString()}`
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch posts: ${response.status}`);
  }

  return response.json();
}

export async function getDashboardStats() {
  const response = await fetch(`${API_URL}/api/stats`);

  if (!response.ok) {
    throw new Error(`Failed to fetch stats: ${response.status}`);
  }

  return response.json();
}

export async function updateLeadStatus(postId, status) {
  const response = await fetch(
    `${API_URL}/api/posts/${postId}/lead/status`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        status,
      }),
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to update lead status: ${response.status}`);
  }

  return response.json();
}

export async function updateLeadNotes(postId, notes) {
  const response = await fetch(
    `${API_URL}/api/posts/${postId}/lead/notes`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        notes,
      }),
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to update lead notes: ${response.status}`);
  }

  return response.json();
}
