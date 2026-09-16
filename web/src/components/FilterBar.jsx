export default function FilterBar({
  filters,
  onChange,
  onReset,
}) {
  return (
    <div className="filter-bar">
      <input
        type="text"
        placeholder="Search developer, website, automation..."
        value={filters.search}
        onChange={(e) =>
          onChange("search", e.target.value)
        }
      />

      <select
        value={filters.tier}
        onChange={(e) =>
          onChange("tier", e.target.value)
        }
      >
        <option value="">All tiers</option>
        <option value="hot">Hot</option>
        <option value="warm">Warm</option>
        <option value="potential">Potential</option>
        <option value="low">Low</option>
      </select>

      <select
        value={filters.category}
        onChange={(e) =>
          onChange("category", e.target.value)
        }
      >
        <option value="">All categories</option>
        <option value="website">Website</option>
        <option value="mobile_app">Mobile App</option>
        <option value="custom_software">
          Custom Software
        </option>
        <option value="pos">POS</option>
        <option value="hrm">HRM</option>
        <option value="crm">CRM</option>
        <option value="automation">Automation</option>
        <option value="ai">AI</option>
        <option value="ecommerce">E-commerce</option>
      </select>

      <select
        value={filters.intent}
        onChange={(e) =>
          onChange("intent", e.target.value)
        }
      >
        <option value="">All intents</option>
        <option value="hire_developer">
          Hire Developer
        </option>
        <option value="looking_for_jasa">
          Looking for Service
        </option>
        <option value="project_inquiry">
          Project Inquiry
        </option>
        <option value="recommendation">
          Recommendation
        </option>
      </select>

      <button onClick={onReset}>
        Reset
      </button>
    </div>
  );
}
