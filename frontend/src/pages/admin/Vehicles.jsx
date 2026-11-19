import { useEffect, useState } from "react";
import { vehicleApi } from "../../services/api";

const API_ORIGIN = (
  import.meta.env.VITE_API_URL || "http://localhost:8080/api"
).replace(/\/api\/?$/, "");

const emptyForm = {
  name: "",
  brand: "",
  model: "",
  year: new Date().getFullYear(),
  license_plate: "",
  category: "MPV",
  transmission: "Automatic",
  fuel_type: "Bensin",
  capacity: 7,
  price_per_day: "",
  description: "",
  status: "available",
};

function getImageUrl(image) {
  if (!image) return "";

  if (image.startsWith("http://") || image.startsWith("https://")) {
    return image;
  }

  return `${API_ORIGIN}${image.startsWith("/") ? "" : "/"}${image}`;
}

export default function Vehicles() {
  const [vehicles, setVehicles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [form, setForm] = useState(emptyForm);

  const [selectedImage, setSelectedImage] = useState(null);
  const [currentImage, setCurrentImage] = useState("");
  const [previewUrl, setPreviewUrl] = useState("");

  // Bersihkan URL preview sementara ketika berubah atau komponen ditutup.
  useEffect(() => {
    if (!selectedImage) {
      setPreviewUrl("");
      return;
    }

    const objectUrl = URL.createObjectURL(selectedImage);
    setPreviewUrl(objectUrl);

    return () => URL.revokeObjectURL(objectUrl);
  }, [selectedImage]);

  async function loadVehicles() {
    setLoading(true);
    setError("");

    try {
      const result = await vehicleApi.all();
      setVehicles(Array.isArray(result) ? result : []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadVehicles();
  }, []);

  function handleChange(e) {
    const { name, value } = e.target;

    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  }

  function resetForm() {
    setForm(emptyForm);
    setSelectedImage(null);
    setCurrentImage("");
    setEditingId(null);
  }

  function closeForm() {
    if (saving) return;

    setShowForm(false);
    resetForm();
    setError("");
  }

  function openCreateForm() {
    resetForm();
    setError("");
    setSuccess("");
    setShowForm(true);
  }

  function openEditForm(vehicle) {
    setEditingId(vehicle.id);

    setForm({
      name: vehicle.name || "",
      brand: vehicle.brand || "",
      model: vehicle.model || "",
      year: vehicle.year || "",
      license_plate: vehicle.license_plate || "",
      category: vehicle.category || "MPV",
      transmission: vehicle.transmission || "Automatic",
      fuel_type: vehicle.fuel_type || "Bensin",
      capacity: vehicle.capacity || "",
      price_per_day: vehicle.price_per_day || "",
      description: vehicle.description || "",
      status: vehicle.status || "available",
    });

    setSelectedImage(null);
    setCurrentImage(vehicle.image || "");
    setError("");
    setSuccess("");
    setShowForm(true);
  }

  function handleImageChange(e) {
    const file = e.target.files?.[0];

    if (!file) return;

    const allowedTypes = ["image/jpeg", "image/png", "image/webp"];

    if (!allowedTypes.includes(file.type)) {
      setError("Format foto harus JPG, PNG, atau WebP.");
      e.target.value = "";
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      setError("Ukuran foto maksimal 5 MB.");
      e.target.value = "";
      return;
    }

    setError("");
    setSelectedImage(file);
  }

  async function handleSubmit(e) {
    e.preventDefault();

    setSaving(true);
    setError("");
    setSuccess("");

    const payload = new FormData();

    payload.append("name", form.name.trim());
    payload.append("brand", form.brand.trim());
    payload.append("model", form.model.trim());
    payload.append("year", form.year);
    payload.append("license_plate", form.license_plate.trim());
    payload.append("category", form.category);
    payload.append("transmission", form.transmission);
    payload.append("fuel_type", form.fuel_type);
    payload.append("capacity", form.capacity);
    payload.append("price_per_day", form.price_per_day);
    payload.append("description", form.description.trim());
    payload.append("status", form.status);

    if (selectedImage) {
      payload.append("image", selectedImage);
    }

    try {
      if (editingId) {
        await vehicleApi.update(editingId, payload);
        setSuccess("Kendaraan berhasil diperbarui.");
      } else {
        await vehicleApi.create(payload);
        setSuccess("Kendaraan berhasil ditambahkan.");
      }

      setShowForm(false);
      resetForm();

      await loadVehicles();
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(vehicle) {
    const confirmed = window.confirm(`Yakin ingin menghapus ${vehicle.name}?`);

    if (!confirmed) return;

    setError("");
    setSuccess("");

    try {
      await vehicleApi.remove(vehicle.id);
      setSuccess("Kendaraan berhasil dihapus.");
      await loadVehicles();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleStatusChange(vehicle, status) {
    setError("");
    setSuccess("");

    try {
      await vehicleApi.updateStatus(vehicle.id, status);
      setSuccess("Status kendaraan berhasil diperbarui.");
      await loadVehicles();
    } catch (err) {
      setError(err.message);
      await loadVehicles();
    }
  }

  return (
    <div className="space-y-6">
      {/* HEADER */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <p className="text-sm font-bold uppercase tracking-wider text-blue-700">
            Manajemen Armada
          </p>

          <h1 className="mt-2 text-3xl font-black text-slate-900">Kendaraan</h1>

          <p className="mt-2 text-sm text-slate-500">
            Tambahkan, edit, dan kelola kendaraan DriveRent.
          </p>
        </div>

        <button
          onClick={openCreateForm}
          className="rounded-xl bg-blue-700 px-5 py-3 font-bold text-white transition hover:bg-blue-800"
        >
          + Tambah Kendaraan
        </button>
      </div>

      {/* NOTIFICATION */}
      {error && (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      )}

      {success && (
        <div className="rounded-xl border border-green-200 bg-green-50 p-4 text-sm text-green-700">
          {success}
        </div>
      )}

      {/* FORM */}
      {showForm && (
        <form
          onSubmit={handleSubmit}
          className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm sm:p-7"
        >
          <div className="mb-6 flex items-center justify-between">
            <div>
              <h2 className="text-xl font-bold text-slate-900">
                {editingId ? "Edit Kendaraan" : "Tambah Kendaraan"}
              </h2>

              <p className="mt-1 text-sm text-slate-500">
                Lengkapi informasi dan foto kendaraan.
              </p>
            </div>

            <button
              type="button"
              onClick={closeForm}
              disabled={saving}
              className="rounded-lg px-3 py-2 text-sm font-semibold text-slate-500 hover:bg-slate-100"
            >
              Tutup
            </button>
          </div>

          {/* IMAGE UPLOAD */}
          <div className="mb-7">
            <label className="mb-2 block text-sm font-semibold text-slate-700">
              Foto Kendaraan
            </label>

            <div className="grid gap-5 sm:grid-cols-[220px_1fr]">
              <div className="flex aspect-[4/3] items-center justify-center overflow-hidden rounded-xl border border-dashed border-slate-300 bg-slate-50">
                {previewUrl || currentImage ? (
                  <img
                    src={previewUrl || getImageUrl(currentImage)}
                    alt="Preview kendaraan"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="text-center text-slate-400">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="mx-auto h-10 w-10"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={1.5}
                        d="M3 16l5-5 4 4 4-6 5 7M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                      />
                    </svg>

                    <p className="mt-2 text-xs">Preview foto kendaraan</p>
                  </div>
                )}
              </div>

              <div className="flex flex-col justify-center">
                <input
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  onChange={handleImageChange}
                  className="block w-full cursor-pointer rounded-xl border border-slate-200 text-sm text-slate-600 file:mr-4 file:cursor-pointer file:border-0 file:bg-blue-50 file:px-4 file:py-3 file:font-semibold file:text-blue-700 hover:file:bg-blue-100"
                />

                <p className="mt-3 text-xs leading-5 text-slate-500">
                  Format yang didukung: JPG, PNG, dan WebP.
                  <br />
                  Ukuran maksimal: 5 MB.
                </p>

                {selectedImage && (
                  <div className="mt-3 flex items-center gap-2">
                    <p className="text-xs font-medium text-green-700">
                      Foto baru: {selectedImage.name}
                    </p>

                    <button
                      type="button"
                      onClick={() => setSelectedImage(null)}
                      className="text-xs font-semibold text-red-600 hover:underline"
                    >
                      Batalkan
                    </button>
                  </div>
                )}

                {!selectedImage && currentImage && (
                  <p className="mt-3 text-xs text-slate-500">
                    Foto lama akan tetap digunakan jika kamu tidak memilih foto
                    baru.
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* VEHICLE FIELDS */}
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <Field label="Nama Kendaraan">
              <input
                required
                name="name"
                value={form.name}
                onChange={handleChange}
                placeholder="Toyota Avanza"
                className="input-admin"
              />
            </Field>

            <Field label="Brand">
              <input
                required
                name="brand"
                value={form.brand}
                onChange={handleChange}
                placeholder="Toyota"
                className="input-admin"
              />
            </Field>

            <Field label="Model">
              <input
                required
                name="model"
                value={form.model}
                onChange={handleChange}
                placeholder="Avanza Veloz"
                className="input-admin"
              />
            </Field>

            <Field label="Tahun">
              <input
                required
                type="number"
                name="year"
                value={form.year}
                onChange={handleChange}
                className="input-admin"
              />
            </Field>

            <Field label="Nomor Polisi">
              <input
                required
                name="license_plate"
                value={form.license_plate}
                onChange={handleChange}
                placeholder="H 1234 AB"
                className="input-admin"
              />
            </Field>

            <Field label="Kategori">
              <select
                name="category"
                value={form.category}
                onChange={handleChange}
                className="input-admin"
              >
                <option>MPV</option>
                <option>SUV</option>
                <option>Hatchback</option>
                <option>Sedan</option>
                <option>City Car</option>
                <option>Van</option>
              </select>
            </Field>

            <Field label="Transmisi">
              <select
                name="transmission"
                value={form.transmission}
                onChange={handleChange}
                className="input-admin"
              >
                <option>Automatic</option>
                <option>Manual</option>
              </select>
            </Field>

            <Field label="Bahan Bakar">
              <select
                name="fuel_type"
                value={form.fuel_type}
                onChange={handleChange}
                className="input-admin"
              >
                <option>Bensin</option>
                <option>Diesel</option>
                <option>Hybrid</option>
                <option>Electric</option>
              </select>
            </Field>

            <Field label="Kapasitas Penumpang">
              <input
                required
                type="number"
                min="1"
                name="capacity"
                value={form.capacity}
                onChange={handleChange}
                className="input-admin"
              />
            </Field>

            <Field label="Harga Sewa per Hari (Rp)">
              <input
                required
                type="number"
                min="1"
                name="price_per_day"
                value={form.price_per_day}
                onChange={handleChange}
                placeholder="350000"
                className="input-admin"
              />
            </Field>

            <Field label="Status">
              <select
                name="status"
                value={form.status}
                onChange={handleChange}
                className="input-admin"
              >
                <option value="available">Available</option>
                <option value="rented">Rented</option>
                <option value="maintenance">Maintenance</option>
              </select>
            </Field>
          </div>

          <div className="mt-5">
            <Field label="Deskripsi">
              <textarea
                name="description"
                value={form.description}
                onChange={handleChange}
                rows="3"
                placeholder="Deskripsi kendaraan..."
                className="input-admin"
              />
            </Field>
          </div>

          <div className="mt-6 flex flex-wrap gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-xl bg-blue-700 px-6 py-3 font-bold text-white hover:bg-blue-800 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {saving
                ? "Menyimpan..."
                : editingId
                  ? "Simpan Perubahan"
                  : "Simpan Kendaraan"}
            </button>

            <button
              type="button"
              onClick={closeForm}
              disabled={saving}
              className="rounded-xl border border-slate-300 px-6 py-3 font-bold text-slate-600 hover:bg-slate-50"
            >
              Batal
            </button>
          </div>
        </form>
      )}

      {/* VEHICLE TABLE */}
      <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="border-b border-slate-100 p-5">
          <h2 className="font-bold text-slate-900">
            Daftar Armada ({vehicles.length})
          </h2>
        </div>

        {loading ? (
          <p className="p-8 text-center text-slate-500">Memuat kendaraan...</p>
        ) : vehicles.length === 0 ? (
          <p className="p-8 text-center text-slate-500">Belum ada kendaraan.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[950px] text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                <tr>
                  <th className="px-5 py-4">Foto</th>
                  <th className="px-5 py-4">Kendaraan</th>
                  <th className="px-5 py-4">Nomor Polisi</th>
                  <th className="px-5 py-4">Harga / Hari</th>
                  <th className="px-5 py-4">Status</th>
                  <th className="px-5 py-4">Aksi</th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {vehicles.map((vehicle) => (
                  <tr key={vehicle.id}>
                    <td className="px-5 py-3">
                      {vehicle.image ? (
                        <img
                          src={getImageUrl(vehicle.image)}
                          alt={vehicle.name}
                          className="h-14 w-20 rounded-lg object-cover"
                        />
                      ) : (
                        <div className="flex h-14 w-20 items-center justify-center rounded-lg bg-slate-100 text-xs text-slate-400">
                          No Image
                        </div>
                      )}
                    </td>

                    <td className="px-5 py-4">
                      <p className="font-bold text-slate-900">{vehicle.name}</p>
                      <p className="mt-1 text-xs text-slate-500">
                        {vehicle.brand} · {vehicle.year}
                      </p>
                    </td>

                    <td className="px-5 py-4">{vehicle.license_plate}</td>

                    <td className="px-5 py-4 font-semibold">
                      Rp
                      {Number(vehicle.price_per_day).toLocaleString("id-ID")}
                    </td>

                    <td className="px-5 py-4">
                      <select
                        value={vehicle.status}
                        onChange={(e) =>
                          handleStatusChange(vehicle, e.target.value)
                        }
                        className="rounded-lg border border-slate-200 px-2 py-2 text-xs font-semibold"
                      >
                        <option value="available">Available</option>
                        <option value="rented">Rented</option>
                        <option value="maintenance">Maintenance</option>
                      </select>
                    </td>

                    <td className="px-5 py-4">
                      <div className="flex gap-2">
                        <button
                          onClick={() => openEditForm(vehicle)}
                          className="rounded-lg bg-blue-50 px-3 py-2 text-xs font-bold text-blue-700 hover:bg-blue-100"
                        >
                          Edit
                        </button>

                        <button
                          onClick={() => handleDelete(vehicle)}
                          className="rounded-lg bg-red-50 px-3 py-2 text-xs font-bold text-red-700 hover:bg-red-100"
                        >
                          Hapus
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

function Field({ label, children }) {
  return (
    <label className="mb-5 block">
      <span className="mb-2 block text-sm font-semibold text-slate-700">
        {label}
      </span>

      {children}
    </label>
  );
}
