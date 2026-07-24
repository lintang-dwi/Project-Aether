import React, { useState, useEffect } from 'react';
import { useStore } from '../store/useStore';
import { ShieldCheck, Lock, Terminal, CheckCircle2, XCircle } from 'lucide-react';

const formatLogLine = (log: string) => {
  if (log.startsWith('[WORKSPACE]')) {
    return <><span className="text-cyan-400 font-bold">[WORKSPACE]</span> {log.slice(11)}</>;
  }
  if (log.startsWith('[SYSTEM]')) {
    return <><span className="text-indigo-400 font-bold">[SYSTEM]</span> {log.slice(8)}</>;
  }
  if (log.startsWith('[EVENTBUS]')) {
    return <><span className="text-violet-400 font-bold">[EVENTBUS]</span> {log.slice(10)}</>;
  }
  if (log.startsWith('[STORAGE]')) {
    return <><span className="text-amber-300 font-bold">[STORAGE]</span> {log.slice(9)}</>;
  }
  if (log.startsWith('[SECURITY]')) {
    return <><span className="text-emerald-400 font-bold">[SECURITY]</span> {log.slice(10)}</>;
  }
  if (log.startsWith('[TASK]')) {
    return <><span className="text-cyan-300 font-bold">[TASK]</span> {log.slice(6)}</>;
  }
  if (log.startsWith('[CONTEXT]')) {
    return <><span className="text-indigo-300 font-bold">[CONTEXT]</span> {log.slice(9)}</>;
  }
  if (log.startsWith('[PLANNER]')) {
    return <><span className="text-violet-300 font-bold">[PLANNER]</span> {log.slice(9)}</>;
  }
  if (log.startsWith('[ACTION]')) {
    return <><span className="text-emerald-300 font-bold">[ACTION]</span> {log.slice(8)}</>;
  }
  return log;
};

export const SecurityView: React.FC = () => {
  const { logs, auditLogs, allowWrite, allowExecute, togglePermissionGuard, fetchAuditTrail, fetchLogs } = useStore();

  useEffect(() => {
    fetchAuditTrail();
    fetchLogs();
  }, [fetchAuditTrail, fetchLogs]);

  return (
    <div className="p-6 space-y-6 overflow-y-auto h-[calc(100vh-4rem)]">
      {/* Security Boundary Controls */}
      <div className="grid grid-cols-2 gap-6">
        <div className="p-5 rounded-xl glass-card space-y-4 bg-[#12141c]">
          <div className="flex items-center gap-2 text-xs font-mono text-emerald-400 uppercase tracking-wider">
            <ShieldCheck className="w-4 h-4 text-emerald-400" />
            <span>Permission Guards</span>
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 rounded-lg bg-white/[0.02] border border-white/5">
              <div>
                <div className="text-xs font-semibold text-gray-200">File Write Permission</div>
                <div className="text-[10px] font-mono text-gray-400">Allow Action Processor to modify files</div>
              </div>
              <button
                onClick={() => togglePermissionGuard('WRITE', !allowWrite)}
                className={`px-3 py-1 rounded-lg text-xs font-mono font-bold transition-all ${
                  allowWrite
                    ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                    : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'
                }`}
              >
                {allowWrite ? 'ALLOWED' : 'BLOCKED'}
              </button>
            </div>

            <div className="flex items-center justify-between p-3 rounded-lg bg-white/[0.02] border border-white/5">
              <div>
                <div className="text-xs font-semibold text-gray-200">Command Execution Permission</div>
                <div className="text-[10px] font-mono text-gray-400">Allow Tool Sandbox to run shell commands</div>
              </div>
              <button
                onClick={() => togglePermissionGuard('EXECUTE', !allowExecute)}
                className={`px-3 py-1 rounded-lg text-xs font-mono font-bold transition-all ${
                  allowExecute
                    ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                    : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'
                }`}
              >
                {allowExecute ? 'ALLOWED' : 'BLOCKED'}
              </button>
            </div>
          </div>
        </div>

        {/* Audit Log */}
        <div className="p-5 rounded-xl glass-card space-y-4 bg-[#12141c]">
          <h3 className="font-bold text-xs text-gray-300 uppercase tracking-wider">
            Audit Trail
          </h3>
          <div className="space-y-2">
            {auditLogs.map((log, i) => (
              <div key={i} className="flex items-center justify-between p-2.5 rounded-lg bg-white/[0.02] border border-white/5 text-xs font-mono">
                <div className="flex items-center gap-2">
                  {log.allowed ? (
                    <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />
                  ) : (
                    <XCircle className="w-3.5 h-3.5 text-rose-400" />
                  )}
                  <span className="text-gray-500">{log.time}</span>
                  <span className="text-indigo-300 font-bold">[{log.perm}]</span>
                  <span className="text-gray-200 truncate max-w-[180px]">{log.target}</span>
                </div>
                <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${log.allowed ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}`}>
                  {log.allowed ? 'PASS' : 'DENIED'}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Live Log Stream */}
      <div className="p-5 rounded-xl glass-card space-y-3 bg-[#12141c]">
        <div className="flex items-center gap-2 text-xs font-mono text-cyan-400 uppercase tracking-wider">
          <Terminal className="w-4 h-4 text-cyan-400" />
          <span>Runtime Event & System Log Stream</span>
        </div>
        <div className="p-4 rounded-xl bg-[#0b0d14] border border-white/10 font-mono text-xs text-gray-300 space-y-2 max-h-64 overflow-y-auto">
          {logs.map((log, i) => (
            <div key={i} className="leading-relaxed border-b border-white/[0.03] pb-1.5 last:border-b-0 last:pb-0">
              {formatLogLine(log)}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
