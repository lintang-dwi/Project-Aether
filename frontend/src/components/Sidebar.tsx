import React from 'react';
import { useStore, ViewType } from '../store/useStore';
import { 
  LayoutDashboard, 
  Network, 
  Code2, 
  Bot, 
  ShieldCheck, 
  Cpu
} from 'lucide-react';

export const Sidebar: React.FC = () => {
  const { activeView, setActiveView } = useStore();

  const navItems: { id: ViewType; label: string; icon: React.ReactNode }[] = [
    { id: 'dashboard', label: 'Dashboard', icon: <LayoutDashboard size={18} strokeWidth={1.75} /> },
    { id: 'graph', label: 'Knowledge Graph', icon: <Network size={18} strokeWidth={1.75} /> },
    { id: 'workspace', label: 'Code Workspace', icon: <Code2 size={18} strokeWidth={1.75} /> },
    { id: 'task', label: 'Autonomous Tasks', icon: <Bot size={18} strokeWidth={1.75} /> },
    { id: 'security', label: 'Security & Logs', icon: <ShieldCheck size={18} strokeWidth={1.75} /> },
  ];

  return (
    <aside className="w-64 glass-panel border-r border-white/10 flex flex-col justify-between h-screen p-4 z-20 shrink-0">
      <div>
        {/* Brand Header */}
        <div className="flex items-center gap-3 px-3 py-3.5 mb-5 border-b border-white/5">
          <div className="w-8 h-8 rounded-xl bg-gradient-to-br from-indigo-500 to-violet-600 p-0.5 flex items-center justify-center">
            <div className="w-full h-full bg-[#0b0d14] rounded-[10px] flex items-center justify-center">
              <Cpu className="w-4 h-4 text-cyan-400" />
            </div>
          </div>
          <div>
            <h1 className="font-extrabold text-sm tracking-wider bg-gradient-to-r from-white via-indigo-100 to-cyan-300 bg-clip-text text-transparent">
              AETHER
            </h1>
            <p className="text-[10px] font-mono text-gray-400 uppercase tracking-widest">
              GraphOS Desktop
            </p>
          </div>
        </div>

        {/* Navigation List */}
        <nav className="space-y-1">
          {navItems.map((item) => {
            const isActive = activeView === item.id;
            return (
              <button
                key={item.id}
                onClick={() => setActiveView(item.id)}
                className={`w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-xs transition-all duration-150 ${
                  isActive
                    ? 'bg-indigo-600/20 text-white font-semibold border border-indigo-500/30'
                    : 'text-gray-400 font-medium hover:text-gray-200 hover:bg-white/5'
                }`}
              >
                <span className={isActive ? 'text-indigo-400' : 'text-gray-400'}>
                  {item.icon}
                </span>
                <span>{item.label}</span>
              </button>
            );
          })}
        </nav>
      </div>

      {/* Footer Info */}
      <div className="p-3 rounded-xl bg-white/[0.03] border border-white/5 text-[11px] font-mono text-gray-400">
        <div className="flex justify-between items-center mb-1">
          <span>RUNTIME</span>
          <span className="flex items-center gap-1.5 text-emerald-400 font-semibold">
            <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
            ACTIVE
          </span>
        </div>
        <div className="text-[10px] text-gray-500">v0.1.0 · Local First</div>
      </div>
    </aside>
  );
};
