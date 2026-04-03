const API_BASE = process.env.NEXT_PUBLIC_API_URL || "https://registry.modelcontextprotocol.io/v0.1";

export interface ServerDetail {
  name: string;
  description: string;
  title?: string;
  version?: string;
  repository?: { url: string; source?: string };
  websiteUrl?: string;
  packages?: Package[];
  remotes?: Remote[];
}

export interface Package {
  registryType: string;
  identifier: string;
  version?: string;
  runtimeHint?: string;
  transport: { type: string; url?: string };
  environmentVariables?: EnvVar[];
}

export interface Remote {
  type: string;
  url: string;
}

export interface EnvVar {
  name: string;
  description?: string;
  isRequired?: boolean;
  isSecret?: boolean;
}

export interface ServerEntry {
  server: ServerDetail;
}

export interface SearchResponse {
  servers: ServerEntry[];
  metadata: { nextCursor?: string; count?: number };
}

export async function searchServers(query: string, limit = 20): Promise<ServerEntry[]> {
  const params = new URLSearchParams({
    search: query,
    version: "latest",
    limit: String(limit),
  });
  const res = await fetch(`${API_BASE}/servers?${params}`);
  if (!res.ok) throw new Error(`Search failed: ${res.status}`);
  const data: SearchResponse = await res.json();
  return data.servers || [];
}

export async function getAllServers(limit = 96): Promise<ServerEntry[]> {
  const params = new URLSearchParams({
    version: "latest",
    limit: String(limit),
  });
  const res = await fetch(`${API_BASE}/servers?${params}`);
  if (!res.ok) throw new Error(`Failed to fetch servers: ${res.status}`);
  const data: SearchResponse = await res.json();
  return data.servers || [];
}

export function shortName(name: string): string {
  const parts = name.split("/");
  let short = parts[parts.length - 1];
  for (const prefix of ["server-", "mcp-", "mcp_"]) {
    if (short.startsWith(prefix)) {
      return short.slice(prefix.length);
    }
  }
  return short;
}

export function getInstallCommand(server: ServerDetail): string {
  return `mcphub install ${server.name}`;
}
