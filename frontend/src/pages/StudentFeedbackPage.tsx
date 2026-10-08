import { useState, type FormEvent } from "react";
import ApiStatus from "../components/ApiStatus";
import { useApiHealth } from "../hooks/useApiHealth";
import type { AuthSession } from "../services/auth";
import { submitFeedback, type FeedbackResponse } from "../services/feedback";

interface StudentFeedbackPageProps {
  session: AuthSession;
  onLogout: () => void;
}

function StudentFeedbackPage({ session, onLogout }: StudentFeedbackPageProps) {
  const apiHealth = useApiHealth();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [submitted, setSubmitted] = useState<FeedbackResponse | null>(null);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSubmitted(null);
    setIsSubmitting(true);

    const formElement = event.currentTarget;
    const form = new FormData(formElement);

    try {
      const result = await submitFeedback(
        {
          courseId: String(form.get("courseId") ?? ""),
          content: String(form.get("content") ?? ""),
        },
        session.accessToken,
      );
      setSubmitted(result);
      formElement.reset();
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "We could not send your feedback. Please try again.",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="min-h-screen bg-mist text-ink">
      <header className="border-b border-ink/10 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
          <a className="text-xl font-bold tracking-tight text-ocean" href="/">
            CampusEcho
          </a>
          <div className="flex items-center gap-3">
            <span className="hidden text-sm text-ink/60 sm:inline">
              {session.user.email}
            </span>
            <button
              className="rounded-lg border border-ink/15 px-4 py-2 text-sm font-semibold hover:bg-mist"
              onClick={onLogout}
              type="button"
            >
              Sign out
            </button>
          </div>
        </div>
      </header>

      <section className="mx-auto grid max-w-6xl gap-12 px-6 py-16 md:grid-cols-[1fr_1.1fr] md:py-24">
        <div className="max-w-lg">
          <p className="mb-4 text-sm font-bold uppercase tracking-[0.2em] text-ocean">
            Your voice matters
          </p>
          <h1 className="text-4xl font-bold leading-tight tracking-tight md:text-5xl">
            Help make your courses better.
          </h1>
          <p className="mt-6 text-lg leading-8 text-ink/70">
            Share thoughtful feedback with your campus. This starter form does
            not ask for your name or student ID.
          </p>
          <ApiStatus health={apiHealth} />
        </div>

        <div className="rounded-2xl border border-ink/10 bg-white p-6 shadow-sm md:p-9">
          <h2 className="text-2xl font-bold">Submit feedback</h2>
          <p className="mt-2 text-sm leading-6 text-ink/60">
            Keep feedback constructive. Your submission is stored without a
            student identity.
          </p>

          <form className="mt-8 space-y-5" onSubmit={handleSubmit}>
            <label className="block">
              <span className="mb-2 block text-sm font-semibold">Course ID</span>
              <input
                autoComplete="off"
                className="w-full rounded-lg border border-ink/15 px-4 py-3 outline-none transition placeholder:text-ink/35 focus:border-ocean focus:ring-2 focus:ring-ocean/15"
                maxLength={120}
                name="courseId"
                placeholder="e.g. CS455"
                required
              />
            </label>

            <label className="block">
              <span className="mb-2 block text-sm font-semibold">
                Your feedback
              </span>
              <textarea
                className="min-h-40 w-full resize-y rounded-lg border border-ink/15 px-4 py-3 outline-none transition placeholder:text-ink/35 focus:border-ocean focus:ring-2 focus:ring-ocean/15"
                maxLength={5000}
                name="content"
                placeholder="What is working well? What could be improved?"
                required
                rows={6}
              />
              <span className="mt-2 block text-right text-xs text-ink/50">
                Please keep your message under 5,000 characters.
              </span>
            </label>

            {error && (
              <p
                className="rounded-lg bg-rose-50 px-4 py-3 text-sm text-rose-800"
                role="alert"
              >
                {error}
              </p>
            )}
            {submitted && (
              <p
                className="rounded-lg bg-emerald-50 px-4 py-3 text-sm text-emerald-800"
                role="status"
              >
                Feedback received. Reference: {submitted.id}
              </p>
            )}

            <button
              className="w-full rounded-lg bg-ocean px-5 py-3 font-semibold text-white transition hover:bg-ocean/90 focus:outline-none focus:ring-2 focus:ring-ocean focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
              disabled={isSubmitting}
              type="submit"
            >
              {isSubmitting ? "Sending…" : "Send feedback"}
            </button>
          </form>
        </div>
      </section>

      <footer className="border-t border-ink/10 px-6 py-6 text-center text-sm text-ink/50">
        CampusEcho · A starter for anonymous and accountable campus feedback
      </footer>
    </main>
  );
}

export default StudentFeedbackPage;
