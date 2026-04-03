"use client";

import { useState, useEffect } from "react";
import SearchBar from "@/components/SearchBar";
import ServerCard from "@/components/ServerCard";
import InstallCommand from "@/components/InstallCommand";
import { searchServers, getAllServers, ServerEntry } from "@/lib/api";

export default function Home() {
  const [servers, setServers] = useState<ServerEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [query, setQuery] = useState("");

  useEffect(() => {
    getAllServers(24).then((data) => {
      setServers(data);
      setLoading(false);
    }).catch(() => setLoading(false));
  }, []);

  const handleSearch = async (q: string) => {
    setQuery(q);
    setLoading(true);
    const results = await searchServers(q, 30);
    setServers(results);
    setLoading(false);
  };

  return (
    <div className="min-h-screen bg-black text-white">
      {/* Hero */}
      <header className="pt-20 pb-16 px-6">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
            MCP Hub
          </h1>
          <p className="text-xl text-gray-400 mb-8">
            The package manager for Model Context Protocol servers
          </p>
          <SearchBar onSearch={handleSearch} />
          <div className="mt-6">
            <InstallCommand command="mcphub search filesystem" />
          </div>
        </div>
      </header>

      {/* Stats */}
      <div className="max-w-6xl mx-auto px-6 mb-8">
        <div className="flex items-center gap-6 text-sm text-gray-500">
          <span>{servers.length} servers {query ? `matching "${query}"` : "available"}</span>
        </div>
      </div>

      {/* Server Grid */}
      <main className="max-w-6xl mx-auto px-6 pb-20">
        {loading ? (
          <div className="text-center py-20 text-gray-500">Loading servers...</div>
        ) : servers.length === 0 ? (
          <div className="text-center py-20 text-gray-500">
            No servers found. Try a different search.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {servers.map((entry, i) => (
              <ServerCard key={entry.server.name + i} entry={entry} />
            ))}
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800 py-8 px-6">
        <div className="max-w-6xl mx-auto flex items-center justify-between text-sm text-gray-500">
          <span>MCP Hub</span>
          <div className="flex items-center gap-4">
            <a href="https://github.com/Ricardo-M-L/mcphub" className="hover:text-white transition-colors">
              GitHub
            </a>
            <a href="https://modelcontextprotocol.io" className="hover:text-white transition-colors">
              MCP Spec
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}
