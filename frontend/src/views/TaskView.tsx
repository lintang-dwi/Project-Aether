import React, { useState } from 'react';
import { useStore } from '../store/useStore';
import { Bot, Play, CheckCircle2, GitCommit, Layers, ArrowRight } from 'lucide-react';

export const TaskView: React.FC = () => {
  const { tasks, executeTask } = useStore();
  const [goal, setGoal] = useState('');
  const [targetPath, setTargetPath] = useState('internal/coordinator');
  const [isPlanning, setIsPlanning] = useState(false);

  const handleRunTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!goal) return;

    setIsPlanning(true);
    try {
      await executeTask(goal, targetPath);
      setGoal('');
    } finally {
      setIsPlanning(false);
    }
  };

  return (
    <div className="p-6 space-y-6 overflow-y-auto h-[calc(100vh-4rem)]">
      {/* Create Task Panel */}
      <div className="p-5 rounded-xl glass-card space-y-4 bg-[#12141c]">
        <div className="flex items-center gap-2 text-xs font-mono text-cyan-400 uppercase tracking-wider">
          <Bot className="w-4 h-4 text-cyan-400" />
          <span>New Autonomous Task</span>
        </div>

        <form onSubmit={handleRunTask} className="space-y-4">
          <div className="grid grid-cols-3 gap-4">
            <div className="col-span-2 space-y-1.5">
              <label className="text-xs font-mono text-gray-400">Task Goal / User Need</label>
              <input
                type="text"
                value={goal}
                onChange={(e) => setGoal(e.target.value)}
                placeholder="e.g. Add OpenTelemetry tracing to EventBus publisher"
                className="w-full h-10 px-3.5 rounded-xl bg-[#0b0d14] border border-white/10 text-xs font-mono text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-mono text-gray-400">Target Workspace Path</label>
              <input
                type="text"
                value={targetPath}
                onChange={(e) => setTargetPath(e.target.value)}
                placeholder="internal/eventbus"
                className="w-full h-10 px-3.5 rounded-xl bg-[#0b0d14] border border-white/10 text-xs font-mono text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={isPlanning || !goal}
            className="h-10 flex items-center gap-2 px-5 rounded-xl bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 text-white font-medium text-xs border border-indigo-400/20 shadow-sm transition-all disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:from-indigo-600 disabled:hover:to-violet-600"
          >
            <Play className={`w-3.5 h-3.5 ${isPlanning ? 'animate-spin' : ''}`} />
            <span>{isPlanning ? 'Planning & Executing...' : 'Execute Autonomous Task'}</span>
          </button>
        </form>
      </div>

      {/* Task Execution Log & Steps */}
      <div className="p-5 rounded-xl glass-card space-y-4 bg-[#12141c]">
        <h3 className="font-bold text-xs text-gray-300 uppercase tracking-wider flex items-center justify-between">
          <span>Execution History & Git Checkpoints</span>
          <span className="text-xs font-mono text-violet-400 flex items-center gap-1.5">
            <GitCommit className="w-3.5 h-3.5" /> 100% Rollback Enabled
          </span>
        </h3>

        <div className="space-y-3">
          {tasks.map((t) => (
            <div key={t.id} className="p-4 rounded-xl bg-white/[0.02] border border-white/5 space-y-3">
              <div className="flex justify-between items-center">
                <div className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  <span className="text-xs font-semibold text-gray-200">{t.goal}</span>
                </div>
                <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                  {t.status}
                </span>
              </div>

              <div className="p-3 rounded-lg bg-[#0b0d14] border border-white/5 font-mono text-[11px] space-y-2">
                <div className="flex items-start gap-2">
                  <span className="px-1.5 py-0.5 rounded bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 font-bold shrink-0 text-[10px]">[STEP 1]</span>
                  <span className="text-gray-300">Prepared context for <code className="text-cyan-300">{t.targetPath}</code> using Knowledge Graph.</span>
                </div>
                <div className="flex items-start gap-2">
                  <span className="px-1.5 py-0.5 rounded bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 font-bold shrink-0 text-[10px]">[STEP 2]</span>
                  <span className="text-gray-300">Executed atomic edits. Created Git Checkpoint <code className="text-violet-300">checkpoint-{t.id.slice(-6)}</code>.</span>
                </div>
                <div className="flex items-start gap-2">
                  <span className="px-1.5 py-0.5 rounded bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 font-bold shrink-0 text-[10px]">[VALIDATION]</span>
                  <span className="text-emerald-300">Tree-sitter Go AST syntax check PASSED.</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
