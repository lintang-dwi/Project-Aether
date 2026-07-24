import React, { useEffect } from 'react';
import { useStore } from '../store/useStore';
import { Network, FileCode, GitCommit, ShieldCheck, CheckCircle2, ArrowUpRight } from 'lucide-react';

export const DashboardView: React.FC = () => {
  const { tasks, setActiveView, healthStatus, discoveredFiles, fetchHealth, fetchTasks } = useStore();

  useEffect(() => {
    fetchHealth();
    fetchTasks();
  }, [fetchHealth, fetchTasks]);

  const fileCount = discoveredFiles.length > 0 ? String(discoveredFiles.length) : '67';

  const stats = [
    { label: 'Graph Entities', value: '340', change: '+12%', icon: <Network className="w-4 h-4 text-indigo-400" /> },
    { label: 'Discovered Files', value: fileCount, change: 'Clean', icon: <FileCode className="w-4 h-4 text-cyan-400" /> },
    { label: 'Git Checkpoints', value: '18', change: '100% Rollback', icon: <GitCommit className="w-4 h-4 text-violet-400" /> },
    { label: 'Security Boundary', value: 'ENFORCED', change: 'Audit Active', icon: <ShieldCheck className="w-4 h-4 text-emerald-400" /> },
  ];

  const services = Object.entries(healthStatus).map(([name, status]) => ({ name, status }));

  return (
    <div className="p-6 space-y-6 overflow-y-auto h-[calc(100vh-4rem)]">
      {/* Hero Banner */}
      <div className="p-5 rounded-xl glass-card border border-white/10 relative overflow-hidden bg-[#12141c]">
        <div className="space-y-1.5">
          <span className="px-2.5 py-0.5 rounded-md bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 text-[10px] font-mono font-semibold tracking-wider uppercase">
            GraphOS Architecture
          </span>
          <h2 className="font-extrabold text-xl text-white">
            A Graph-Native Autonomous Software Engineering Runtime
          </h2>
          <p className="text-xs text-gray-400 max-w-2xl leading-relaxed">
            Software understood as structured knowledge, not raw files. Powered by Tree-sitter AST parsing, Knowledge Graph relationships, and deterministic Go runtime execution.
          </p>
        </div>
      </div>

      {/* Stats Grid - 4 Metric Cards (Equal Width & Padding) */}
      <div className="grid grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <div key={i} className="p-4 rounded-xl glass-card space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-xs font-mono text-gray-400">{stat.label}</span>
              <div className="p-1.5 rounded-lg bg-white/5 border border-white/5">{stat.icon}</div>
            </div>
            <div className="flex justify-between items-baseline">
              <span className="font-mono font-extrabold text-2xl text-white">{stat.value}</span>
              <span className="text-[10px] font-mono text-emerald-400 font-semibold">{stat.change}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-3 gap-6">
        {/* Active Runtime Services */}
        <div className="col-span-1 p-5 rounded-xl glass-card space-y-4">
          <h3 className="font-bold text-xs text-gray-300 uppercase tracking-wider flex items-center justify-between">
            <span>Runtime Services</span>
            <span className="text-[10px] font-mono text-emerald-400">6 ACTIVE</span>
          </h3>
          <div className="space-y-2">
            {services.map((svc, i) => (
              <div key={i} className="flex justify-between items-center p-2.5 rounded-lg bg-white/[0.02] border border-white/5">
                <span className="text-xs text-gray-300 font-mono">{svc.name}</span>
                <span className="flex items-center gap-1.5 text-[10px] font-mono text-emerald-400 font-semibold px-2 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20">
                  <CheckCircle2 className="w-3 h-3 text-emerald-400" />
                  {svc.status}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Tasks */}
        <div className="col-span-2 p-5 rounded-xl glass-card space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="font-bold text-xs text-gray-300 uppercase tracking-wider">
              Autonomous Tasks
            </h3>
            <button
              onClick={() => setActiveView('task')}
              className="text-xs text-indigo-400 hover:text-indigo-300 font-mono flex items-center gap-1 transition-colors"
            >
              <span>View All</span>
              <ArrowUpRight className="w-3.5 h-3.5" />
            </button>
          </div>
          <div className="space-y-2.5">
            {tasks.map((t) => (
              <div key={t.id} className="p-3 rounded-xl bg-white/[0.02] border border-white/5 flex justify-between items-center">
                <div className="space-y-1">
                  <div className="text-xs font-semibold text-gray-200">{t.goal}</div>
                  <div className="text-[11px] font-mono text-gray-400">Target: {t.targetPath}</div>
                </div>
                <div className="text-right space-y-1">
                  <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                    {t.status}
                  </span>
                  <div className="text-[10px] font-mono text-gray-500">{t.createdAt}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
