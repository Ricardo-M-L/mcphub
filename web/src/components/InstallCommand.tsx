"use client";

import { useState } from "react";

interface InstallCommandProps {
  command: string;
}

export default function InstallCommand({ command }: InstallCommandProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="flex items-center bg-gray-900 border border-gray-700 rounded-lg overflow-hidden">
      <code className="flex-1 px-4 py-3 text-sm text-green-400 font-mono">
        $ {command}
      </code>
      <button
        onClick={handleCopy}
        className="px-4 py-3 text-gray-400 hover:text-white hover:bg-gray-800 transition-colors border-l border-gray-700"
      >
        {copied ? "Copied!" : "Copy"}
      </button>
    </div>
  );
}
