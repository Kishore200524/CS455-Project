export interface FeedbackInput {
  courseId: string;
  content: string;
}

export interface FeedbackResponse {
  id: string;
  status: string;
  createdAt: string;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

export async function submitFeedback(
  input: FeedbackInput,
  accessToken: string,
): Promise<FeedbackResponse> {
  const response = await fetch("/api/v1/feedback", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${accessToken}`,
    },
    body: JSON.stringify(input),
  });

  const result: unknown = await response.json();
  if (!response.ok) {
    const message =
      isRecord(result) && typeof result.error === "string"
        ? result.error
        : "We could not send your feedback. Please try again.";
    throw new Error(message);
  }

  if (
    !isRecord(result) ||
    typeof result.id !== "string" ||
    typeof result.status !== "string" ||
    typeof result.createdAt !== "string"
  ) {
    throw new Error("The server returned an unexpected response.");
  }

  return {
    id: result.id,
    status: result.status,
    createdAt: result.createdAt,
  };
}
