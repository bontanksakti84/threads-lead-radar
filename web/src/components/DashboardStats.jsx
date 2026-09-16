import {
  Flame,
  CircleDot,
  Target,
  FileText,
} from "lucide-react";

function StatCard({ icon: Icon, label, value, className = "" }) {
  return (
    <div className={`stat-card ${className}`}>
      <div className="stat-icon">
        <Icon size={20} />
      </div>

      <div>
        <div className="stat-label">{label}</div>
        <div className="stat-value">{value ?? 0}</div>
      </div>
    </div>
  );
}

export default function DashboardStats({ stats }) {
  return (
    <div className="stats-grid">
      <StatCard
        icon={FileText}
        label="Total Posts"
        value={stats?.total_posts}
      />

      <StatCard
        icon={Flame}
        label="Hot Leads"
        value={stats?.hot}
      />

      <StatCard
        icon={CircleDot}
        label="Warm Leads"
        value={stats?.warm}
      />

      <StatCard
        icon={Target}
        label="Potential"
        value={stats?.potential}
      />
    </div>
  );
}
