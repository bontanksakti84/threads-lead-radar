function getTier(score) {
  if (score >= 80) return "hot";
  if (score >= 60) return "warm";
  if (score >= 40) return "potential";
  return "low";
}

function formatIntent(intent) {
  if (!intent) return "-";

  return intent
    .split("_")
    .map(
      (word) =>
        word.charAt(0).toUpperCase() + word.slice(1)
    )
    .join(" ");
}

export default function LeadTable({
  posts,
  selectedPost,
  onSelect,
}) {
  return (
    <div className="table-container">
      <table>
        <thead>
          <tr>
            <th>Score</th>
            <th>Username</th>
            <th>Content</th>
            <th>Category</th>
            <th>Intent</th>
            <th>Signals</th>
            <th>Status</th>
          </tr>
        </thead>

        <tbody>
          {posts.length === 0 ? (
            <tr>
              <td colSpan="7" className="empty-state">
                No leads found.
              </td>
            </tr>
          ) : (
            posts.map((post) => {
              const tier = getTier(post.lead_score);

              return (
                <tr
                  key={post.id}
                  className={
                    selectedPost?.id === post.id
                      ? "selected-row"
                      : ""
                  }
                  onClick={() => onSelect(post)}
                >
                  <td>
                    <span
                      className={`score-badge ${tier}`}
                    >
                      {post.lead_score}
                    </span>
                  </td>

                  <td>
                    <strong>
                      @{post.username || "unknown"}
                    </strong>
                  </td>

                  <td className="content-cell">
                    {post.content}
                  </td>

                  <td>
                    <span className="category-badge">
                      {post.category || "other"}
                    </span>
                  </td>

                  <td>
                    {formatIntent(post.intent)}
                  </td>

                  <td>
                    <div className="signals">
                      {post.has_budget && (
                        <span className="signal">
                          💰 Budget
                        </span>
                      )}

                      {post.has_urgency && (
                        <span className="signal">
                          ⚡ Urgent
                        </span>
                      )}

                      {post.needs_developer && (
                        <span className="signal">
                          👨‍💻 Developer
                        </span>
                      )}
                    </div>
                  </td>

                  <td>
                    <span
                      className={`status-badge ${
                        post.lead?.status || "new"
                      }`}
                    >
                      {post.lead?.status || "new"}
                    </span>
                  </td>
                </tr>
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
}
