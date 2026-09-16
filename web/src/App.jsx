import { useEffect, useState } from "react";
import {
  RefreshCw,
  Radar,
} from "lucide-react";

import DashboardStats from "./components/DashboardStats";
import FilterBar from "./components/FilterBar";
import LeadTable from "./components/LeadTable";
import LeadDetail from "./components/LeadDetail";

import {
  getDashboardStats,
  getPosts,
} from "./lib/api";

import "./App.css";

const DEFAULT_FILTERS = {
  search: "",
  tier: "",
  category: "",
  intent: "",
};

function App() {
  const [stats, setStats] = useState(null);
  const [posts, setPosts] = useState([]);
  const [selectedPost, setSelectedPost] =
    useState(null);

  const [filters, setFilters] =
    useState(DEFAULT_FILTERS);

  const [loading, setLoading] =
    useState(true);

  const [error, setError] =
    useState("");


  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  const LIMIT = 5;


  async function loadDashboard() {
    try {
      setLoading(true);
      setError("");

      const [statsData, postsData] =
        await Promise.all([
          getDashboardStats(),
          getPosts({
            ...filters,
            min_score: 40,
            page,
            limit: LIMIT,
          })
        ]);

      setStats(statsData);
      setPosts(postsData.data || []);
      setTotal(postsData.total || 0);
      setTotalPages(postsData.total_pages || 1);

      setSelectedPost((current) => {
        if (!current) return null;

        return (
          postsData.data?.find(
            (post) => post.id === current.id
          ) || null
        );
      });
    } catch (err) {
      console.error(err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
      loadDashboard();
    }, [
      page,
      filters.search,
      filters.tier,
      filters.category,
      filters.intent,
    ]);

  function handleFilterChange(name, value) {
    setFilters((current) => ({
      ...current,
      [name]: value,
    }));

    setPage(1);
  }

  function resetFilters() {
    setFilters(DEFAULT_FILTERS);
    setPage(1);
  }

  return (
    <div className="app">
      <header className="topbar">
        <div className="brand">
          <div className="brand-icon">
            <Radar size={22} />
          </div>

          <div>
            <h1>Threads Lead Radar</h1>
            <span>
              Ichi Solution Lead Intelligence
            </span>
          </div>
        </div>

        <button
          className="refresh-button"
          onClick={loadDashboard}
          disabled={loading}
        >
          <RefreshCw size={17} />

          {loading
            ? "Loading..."
            : "Refresh"}
        </button>
      </header>

      <main>
        {error && (
          <div className="error-banner">
            {error}
          </div>
        )}

        <DashboardStats stats={stats} />

        <FilterBar
          filters={filters}
          onChange={handleFilterChange}
          onReset={resetFilters}
        />

        <div className="content-layout">
          <div className="table-panel">
            <div className="panel-header">
              <div>
                <h2>Lead Radar</h2>

                <span>
                  {posts.length} leads shown
                </span>
              </div>
            </div>

            <LeadTable
              posts={posts}
              selectedPost={selectedPost}
              onSelect={setSelectedPost}
            />

            <div className="pagination">
              <span>
                {total === 0
                  ? "No leads"
                  : `Showing ${
                      (page - 1) * LIMIT + 1
                    }–${Math.min(page * LIMIT, total)} of ${total}`}
              </span>

              <div className="pagination-controls">
                <button
                  disabled={page <= 1}
                  onClick={() =>
                    setPage((current) => current - 1)
                  }
                >
                  ← Previous
                </button>

                <span>
                  Page {page} of {totalPages}
                </span>

                <button
                  disabled={page >= totalPages}
                  onClick={() =>
                    setPage((current) => current + 1)
                  }
                >
                  Next →
                </button>
              </div>
            </div>
          </div>

          <LeadDetail
            post={selectedPost}
            onUpdated={loadDashboard}
          />
        </div>
      </main>
    </div>
  );
}

export default App;