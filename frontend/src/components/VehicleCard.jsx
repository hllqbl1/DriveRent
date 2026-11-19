import { Link } from "react-router";

const fallbackImages = {
  avanza: "/images/avanza.jpg",
  brio: "/images/brio.jpg",
  xpander: "/images/xpander.jpg",
  innova: "/images/innova.jpg",
};

export default function VehicleCard({ vehicle }) {
  const name = vehicle.name || "Kendaraan";
  const normalizedName = name.toLowerCase();

  const fallbackImage =
    Object.entries(fallbackImages).find(([key]) =>
      normalizedName.includes(key),
    )?.[1] || "/images/default-car.jpg";

  const imageUrl = vehicle.image || fallbackImage;

  const price = Number(vehicle.price_per_day || 0).toLocaleString("id-ID");

  return (
    <article className="group overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm transition duration-300 hover:-translate-y-1 hover:shadow-xl">
      <div className="relative h-56 overflow-hidden bg-slate-100">
        <img
          src={imageUrl}
          alt={name}
          className="h-full w-full object-cover transition duration-500 group-hover:scale-105"
          onError={(event) => {
            event.currentTarget.onerror = null;
            event.currentTarget.src = "/images/default-car.jpg";
          }}
        />

        <span className="absolute left-3 top-3 rounded-full bg-white px-3 py-1 text-xs font-bold text-emerald-700 shadow">
          {vehicle.status}
        </span>
      </div>

      <div className="p-5">
        <p className="text-xs font-bold uppercase tracking-wider text-blue-700">
          {vehicle.brand} · {vehicle.category}
        </p>

        <h3 className="mt-2 text-xl font-bold text-slate-900">{name}</h3>

        <div className="mt-3 flex flex-wrap gap-3 text-sm text-slate-500">
          <span>{vehicle.transmission}</span>
          <span>{vehicle.capacity} penumpang</span>
          <span>{vehicle.year}</span>
        </div>

        <div className="mt-5 flex items-center justify-between border-t border-slate-100 pt-4">
          <div>
            <p className="text-xs text-slate-500">Harga sewa mulai</p>
            <p className="text-lg font-extrabold text-slate-900">
              Rp{price}
              <span className="text-xs font-medium text-slate-500">
                {" "}
                / hari
              </span>
            </p>
          </div>

          <Link
            to={`/vehicles/${vehicle.id}`}
            className="rounded-xl bg-blue-700 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-blue-800"
          >
            Lihat detail
          </Link>
        </div>
      </div>
    </article>
  );
}
