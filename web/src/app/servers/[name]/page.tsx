"use client";

import { useState, useEffect, use } from "react";
import { ServerEntry, shortName, getInstallCommand } from "@/lib/api";
import InstallCommand from "@/components/InstallCommand";
import Link from "next/link";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "https://registry.modelcontextprotocol.io/v0.1";

export default function ServerPage({ params }: { params: Promise<{ name: string }> }) {
  const { name } = use(params);
  const [entry, setEntry] = useState<ServerEntry | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const decodedName = decodeURIComponent(name);
    fetch(`${API_BASE}/servers?name=${encodeURIComponent(decodedName)}&version=latest`)
      .then((r) => r.json())
      .then((data) => {
        if (data.servers?.length > 0) {
          setEntry(data.servers[0]);
        }
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [name]);

  if (loading) {
    return (
      <div className="min-h-screen bg-black text-white flex items-center justify-center">
        Loading...
      </div>
    );
  }

  if (!entry) {
    return (
      <div className="min-h-screen bg-black text-white flex items-center justify-center">
        Server not found
      </div>
    );
  }

  const server = entry.server;
  const sName = shortName(server.name);

  return (
    <div className="min-h-screen bg-black text-white">
      <header className="border-b border-gray-800 py-4 px-6">
        <div className="max-w-4xl mx-auto">
          <Link href="/" className="text-blue-400 hover:text-blue-300 text-sm">
            &larr; Back to search
          </Link>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-12">
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <h1 className="text-3xl font-bold">{sName}</h1>
            {server.version && (
              <span className="text-sm text-gray-500 bg-gray-800 px-3 py-1 rounded-full">
                v{server.version}
              </span>
            )}
          </div>
          <p className="text-gray-400 text-lg">{server.description}</p>
        </div>

        <div className="mb-8">
          <h2 className="text-lg font-semibold mb-3">Install</h2>
          <InstallCommand command={getInstallCommand(server)} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
            <h3 className="text-sm font-semibold text-gray-400 mb-3 uppercase">Details</h3>
            <dl className="space-y-2 text-sm">
              <div className="flex justify-between">
                <dt className="text-gray-500">Registry Name</dt>
                <dd className="text-white font-mono text-xs">{server.name}</dd>
              </div>
              {server.repository && (
                <div className="flex justify-between">
                  <dt className="text-gray-500">Repository</dt>
                  <dd>
                    <a href={server.repository.url} className="text-blue-400 hover:underline text-xs" target="_blank">
                      {server.repository.url.replace("https://github.com/", "")}
                    </a>
                  </dd>
                </div>
              )}
              {server.websiteUrl && (
                <div className="flex justify-between">
                  <dt className="text-gray-500">Website</dt>
                  <dd>
                    <a href={server.websiteUrl} className="text-blue-400 hover:underline text-xs" target="_blank">
                      {server.websiteUrl}
                    </a>
                  </dd>
                </div>
              )}
            </dl>
          </div>

          {server.packages && server.packages.length > 0 && (
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
              <h3 className="text-sm font-semibold text-gray-400 mb-3 uppercase">Packages</h3>
              <div className="space-y-3">
                {server.packages.map((pkg, i) => (
                  <div key={i} className="text-sm">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="px-2 py-0.5 bg-blue-900/50 text-blue-400 rounded text-xs">
                        {pkg.registryType}
                      </span>
                      <span className="px-2 py-0.5 bg-gray-800 text-gray-400 rounded text-xs">
                        {pkg.transport.type}
                      </span>
                    </div>
                    <code className="text-gray-300 text-xs">{pkg.identifier}</code>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {server.packages?.[0]?.environmentVariables && server.packages[0].environmentVariables.length > 0 && (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
            <h3 className="text-sm font-semibold text-gray-400 mb-3 uppercase">Environment Variables</h3>
            <div className="space-y-3">
              {server.packages[0].environmentVariables.map((env, i) => (
                <div key={i} className="flex items-start gap-3 text-sm">
                  <code className="text-yellow-400 font-mono shrink-0">{env.name}</code>
                  {env.isRequired && (
                    <span className="text-red-400 text-xs shrink-0">required</span>
                  )}
                  <span className="text-gray-500">{env.description}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
