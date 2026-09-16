import { useEffect, useState } from "react";
import {
  ExternalLink,
  Save,
  CheckCircle2,
} from "lucide-react";

import {
  updateLeadStatus,
  updateLeadNotes,
} from "../lib/api";

const STATUSES = [
  "new",
  "contacted",
  "replied",
  "qualified",
  "won",
  "lost",
];

function formatStatus(status) {
  if (!status) return "New";

  return status
    .split("_")
    .map(
      (word) =>
        word.charAt(0).toUpperCase() +
        word.slice(1)
    )
    .join(" ");
}

export default function LeadDetail({
  post,
  onUpdated,
}) {
  const [notes, setNotes] = useState("");
  const [savingStatus, setSavingStatus] =
    useState(false);
  const [savingNotes, setSavingNotes] =
    useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    setNotes(post?.lead?.notes || "");
    setSaved(false);
  }, [post]);

  if (!post) {
    return (
      <aside className="detail-panel empty-detail">
        <p>Select a lead to see details.</p>
      </aside>
    );
  }

  async function handleStatusChange(event) {
    const status = event.target.value;

    try {
      setSavingStatus(true);

      await updateLeadStatus(
        post.id,
        status
      );

      await onUpdated();
    } catch (error) {
      console.error(error);
      alert(error.message);
    } finally {
      setSavingStatus(false);
    }
  }

  async function handleSaveNotes() {
    try {
      setSavingNotes(true);
      setSaved(false);

      await updateLeadNotes(
        post.id,
        notes
      );

      setSaved(true);

      await onUpdated();

      setTimeout(() => {
        setSaved(false);
      }, 2000);
    } catch (error) {
      console.error(error);
      alert(error.message);
    } finally {
      setSavingNotes(false);
    }
  }

  return (
    <aside className="detail-panel">
      <div className="detail-header">
        <div>
          <div className="detail-username">
            @{post.username || "unknown"}
          </div>

          <div className="detail-category">
            {post.category || "other"}
          </div>
        </div>

        <div className="detail-score">
          {post.lead_score}
        </div>
      </div>

      <section>
        <h3>Original Post</h3>

        <p className="original-content">
          {post.content}
        </p>

        {post.post_url && (
          <a
            href={post.post_url}
            target="_blank"
            rel="noreferrer"
            className="threads-link"
          >
            <ExternalLink size={14} />
            Open Threads
          </a>
        )}
      </section>

      <section>
        <h3>AI Summary</h3>

        <p>
          {post.ai_summary ||
            "No AI summary available."}
        </p>
      </section>

      <section>
        <h3>Intent</h3>

        <div className="detail-intent">
          {formatStatus(post.intent)}
        </div>
      </section>

      <section>
        <h3>Signals</h3>

        <div className="detail-signals">
          {post.has_budget && (
            <span>💰 Budget</span>
          )}

          {post.has_urgency && (
            <span>⚡ Urgent</span>
          )}

          {post.needs_developer && (
            <span>👨‍💻 Developer</span>
          )}

          {!post.has_budget &&
            !post.has_urgency &&
            !post.needs_developer && (
              <span className="muted">
                No strong signals
              </span>
            )}
        </div>
      </section>

      <section>
        <h3>Matched Keywords</h3>

        <div className="keyword-list">
          {(post.matched_keywords || []).length >
          0 ? (
            post.matched_keywords.map(
              (keyword) => (
                <span
                  className="keyword"
                  key={keyword}
                >
                  {keyword}
                </span>
              )
            )
          ) : (
            <span className="muted">
              No matched keywords
            </span>
          )}
        </div>
      </section>

      <section>
        <h3>Lead Status</h3>

        <select
          value={post.lead?.status || "new"}
          onChange={handleStatusChange}
          disabled={savingStatus}
        >
          {STATUSES.map((status) => (
            <option
              key={status}
              value={status}
            >
              {formatStatus(status)}
            </option>
          ))}
        </select>
      </section>

      <section>
        <h3>Notes</h3>

        <textarea
          value={notes}
          onChange={(event) =>
            setNotes(event.target.value)
          }
          placeholder="Catat hasil follow-up, nomor WhatsApp, kebutuhan client, budget, next action..."
          rows={6}
        />

        <button
          className="primary-button"
          onClick={handleSaveNotes}
          disabled={savingNotes}
        >
          {saved ? (
            <>
              <CheckCircle2 size={15} />
              Saved
            </>
          ) : savingNotes ? (
            "Saving..."
          ) : (
            <>
              <Save size={15} />
              Save Notes
            </>
          )}
        </button>
      </section>
    </aside>
  );
}