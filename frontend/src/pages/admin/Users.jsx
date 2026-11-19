import { useCallback, useEffect, useState } from "react";
import { apiRequest } from "../../services/api";

function RoleBadge({ role }) {
  const styles = {
    admin: "bg-purple-100 text-purple-700",
    staff: "bg-blue-100 text-blue-700",
    customer: "bg-green-100 text-green-700",
  };

  return (
    <span
      className={`rounded-full px-3 py-1 text-xs font-bold ${
        styles[role?.toLowerCase()] || "bg-slate-100 text-slate-600"
      }`}
    >
      {role || "-"}
    </span>
  );
}

export default function Users() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const loadUsers = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const result = await apiRequest("/users");
      setUsers(Array.isArray(result) ? result : []);
    } catch (err) {
      setError(err.message || "Gagal mengambil data pengguna");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadUsers();
  }, [loadUsers]);

  const adminCount = users.filter(
    (user) => user.role?.toLowerCase() === "admin",
  ).length;

  const staffCount = users.filter(
    (user) => user.role?.toLowerCase() === "staff",
  ).length;

  const customerCount = users.filter(
    (user) => user.role?.toLowerCase() === "customer",
  ).length;

  return (
    <div>
      {/* Header */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <p className="font-bold uppercase tracking-wide text-blue-600">
            Manajemen Pengguna
          </p>

          <h1 className="mt-2 text-3xl font-black text-slate-900">Pengguna</h1>

          <p className="mt-2 text-slate-500">Kelola data pengguna DriveRent.</p>
        </div>

        <button
          onClick={loadUsers}
          disabled={loading}
          className="rounded-xl border border-slate-200 bg-white px-5 py-3 font-semibold text-slate-700 hover:bg-slate-50 disabled:opacity-50"
        >
          {loading ? "Memuat..." : "Refresh"}
        </button>
      </div>

      {/* Error */}
      {error && (
        <div className="mt-6 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {/* Statistics */}
      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Total Pengguna" value={users.length} />
        <StatCard title="Admin" value={adminCount} />
        <StatCard title="Staff" value={staffCount} />
        <StatCard title="Customer" value={customerCount} />
      </div>

      {/* User Table */}
      <div className="mt-6 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 p-5">
          <h2 className="font-bold text-slate-900">
            Daftar Pengguna ({users.length})
          </h2>
        </div>

        <div className="overflow-x-auto">
          {loading ? (
            <p className="p-6 text-slate-500">Memuat data pengguna...</p>
          ) : error && users.length === 0 ? (
            <p className="p-6 text-slate-500">
              Data pengguna tidak dapat ditampilkan.
            </p>
          ) : users.length === 0 ? (
            <p className="p-6 text-slate-500">Belum ada pengguna.</p>
          ) : (
            <table className="w-full min-w-[750px] text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
                <tr>
                  <th className="px-5 py-4">ID</th>
                  <th className="px-5 py-4">Nama</th>
                  <th className="px-5 py-4">Email</th>
                  <th className="px-5 py-4">Role</th>
                  <th className="px-5 py-4">Dibuat</th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {users.map((user) => (
                  <tr key={user.id} className="transition hover:bg-slate-50">
                    <td className="px-5 py-4 font-bold text-slate-900">
                      #{user.id}
                    </td>

                    <td className="px-5 py-4 font-semibold text-slate-800">
                      {user.name || "-"}
                    </td>

                    <td className="px-5 py-4 text-slate-600">
                      {user.email || "-"}
                    </td>

                    <td className="px-5 py-4">
                      <RoleBadge role={user.role} />
                    </td>

                    <td className="px-5 py-4 text-slate-500">
                      {user.created_at
                        ? new Date(user.created_at).toLocaleDateString("id-ID")
                        : "-"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}

function StatCard({ title, value }) {
  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-5">
      <p className="text-sm text-slate-500">{title}</p>
      <h2 className="mt-2 text-2xl font-black text-slate-900">{value}</h2>
    </div>
  );
}
