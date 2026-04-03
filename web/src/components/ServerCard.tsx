import { ServerEntry, shortName } from "@/lib/api";
import Link from "next/link";

interface ServerCardProps {
  entry: ServerEntry;
}

export default function ServerCard({ entry }: ServerCardProps) {
  const server = entry.server;
  const name = shortName(server.name);
  const transport = server.remotes?.[0]?.type || server.packages?.[0]?.transport?.type || "stdio";
  const registryType = server.packages?.[0]?.registryType || "remote";

  return (
    <Link
      href={`/servers/${encodeURIComponent(server.name)}`}
      className="block p-5 bg-gray-900 border border-gray-800 rounded-xl hover:border-gray-600 transition-colors"
    >
      <div className="flex items-start justify-between mb-2">
        <h3 className="text-lg font-semibold text-white">{name}</h3>
        {server.version && (
          <span className="text-xs text-gray-500 bg-gray-800 px-2 py-1 rounded">
            v{server.version}
          </span>
        )}
      </div>
      <p className="text-gray-400 text-sm mb-3 line-clamp-2">{server.description}</p>
      <div className="flex items-center gap-2">
        <span className="text-xs px-2 py-0.5 rounded bg-blue-900/50 text-blue-400">
          {transport}
        </span>
        <span className="text-xs px-2 py-0.5 rounded bg-gray-800 text-gray-400">
          {registryType}
        </span>
      </div>
    </Link>
  );
}
