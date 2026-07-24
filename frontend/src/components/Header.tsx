import React from 'react';
import { useStore } from '../store/useStore';
import { RefreshCw, Terminal, Sparkles, Activity } from 'lucide-react';

export const Header: React.FC = () => {
  const { activeView, isScanning, scanWorkspace } = useStore();

  const handleScan = () => {
    scanWorkspace();
  };

  const viewTitles: Record<string, string> = {
    dashboard: 'System Overview & Metrics',
    graph: 'Knowledge Graph Visualization',
    workspace: 'Codebase Explorer & AST Analysis',
    task: 'Autonomous AI Agent Task Execution',
    security: 'Security Boundary & Audit Trail',
  };

  return (
    <header className="h-16 glass-panel border-b border-white/10 px-6 flex items-center justify-between z-10 shrink-0">
      <div>
        <h2 className="font-bold text-sm text-gray-100 uppercase tracking-wider">
          {activeView}
        </h2>
        <p className="text-xs text-gray-400 font-mono">
          {viewTitles[activeView] || 'Project Aether Runtime'}
        </p>
      </div>

      <div className="flex items-center gap-2.5">
        {/* Model Badge */}
        <div className="h-8 flex items-center gap-2 px-3.5 rounded-xl bg-cyan-950/30 border border-cyan-500/20 text-cyan-300 text-xs font-mono">
          <Sparkles className="w-3.5 h-3.5 text-cyan-400" />
          <span>GPT-4o / Llama 3.2</span>
        </div>

        {/* Scan Button (Primary Action) */}
        <button
          onClick={handleScan}
          disabled={isScanning}
          className="h-8 flex items-center gap-2 px-4 rounded-xl bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 text-white font-medium text-xs border border-indigo-400/20 transition-all disabled:opacity-40 disabled:cursor-not-allowed"
        >
          <RefreshCw className={`w-3.5 h-3.5 ${isScanning ? 'animate-spin' : ''}`} />
          <span>{isScanning ? 'Scanning...' : 'Rescan Graph'}</span>
        </button>

        {/* Runtime Health Status Pill */}
        <div className="h-8 flex items-center gap-2 px-3.5 rounded-xl bg-emerald-950/30 border border-emerald-500/20 text-emerald-400 text-xs font-mono">
          <Activity className="w-3.5 h-3.5 text-emerald-400" />
          <span>Runtime: OK</span>
        </div>
      </div>
    </header>
  );
};
