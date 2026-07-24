import React, { useState, useEffect } from 'react';
import { ReactFlow, Background, Controls, Node, Edge } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { FileCode, FunctionSquare, Layers, Box, Info } from 'lucide-react';
import { useStore } from '../store/useStore';

const initialNodes: Node[] = [
  {
    id: 'file:main.go',
    type: 'default',
    data: { label: 'file: main.go (Go)', type: 'FILE', path: 'main.go' },
    position: { x: 250, y: 50 },
    style: { background: '#12141c', color: '#60a5fa', border: '1px solid #3b82f6', borderRadius: '12px', padding: '10px 14px', fontSize: '12px', fontFamily: 'JetBrains Mono, monospace' },
  },
  {
    id: 'pkg:coordinator',
    type: 'default',
    data: { label: 'package: coordinator', type: 'PACKAGE', path: 'internal/coordinator' },
    position: { x: 100, y: 180 },
    style: { background: '#12141c', color: '#22d3ee', border: '1px solid #22d3ee', borderRadius: '12px', padding: '10px 14px', fontSize: '12px', fontFamily: 'JetBrains Mono, monospace' },
  },
  {
    id: 'pkg:eventbus',
    type: 'default',
    data: { label: 'package: eventbus', type: 'PACKAGE', path: 'internal/eventbus' },
    position: { x: 400, y: 180 },
    style: { background: '#12141c', color: '#22d3ee', border: '1px solid #22d3ee', borderRadius: '12px', padding: '10px 14px', fontSize: '12px', fontFamily: 'JetBrains Mono, monospace' },
  },
  {
    id: 'fn:StartAll',
    type: 'default',
    data: { label: 'func: StartAll(ctx)', type: 'FUNCTION', path: 'internal/coordinator/coordinator.go' },
    position: { x: 50, y: 320 },
    style: { background: '#12141c', color: '#c084fc', border: '1px solid #8b5cf6', borderRadius: '12px', padding: '10px 14px', fontSize: '12px', fontFamily: 'JetBrains Mono, monospace' },
  },
  {
    id: 'fn:Publish',
    type: 'default',
    data: { label: 'func: Publish(evt)', type: 'FUNCTION', path: 'internal/eventbus/eventbus.go' },
    position: { x: 350, y: 320 },
    style: { background: '#12141c', color: '#c084fc', border: '1px solid #8b5cf6', borderRadius: '12px', padding: '10px 14px', fontSize: '12px', fontFamily: 'JetBrains Mono, monospace' },
  },
];

const initialEdges: Edge[] = [
  { id: 'e1', source: 'file:main.go', target: 'pkg:coordinator', animated: true, style: { stroke: '#6366f1', strokeWidth: 1.5 } },
  { id: 'e2', source: 'file:main.go', target: 'pkg:eventbus', animated: true, style: { stroke: '#6366f1', strokeWidth: 1.5 } },
  { id: 'e3', source: 'pkg:coordinator', target: 'fn:StartAll', style: { stroke: '#22d3ee', strokeWidth: 1.5 } },
  { id: 'e4', source: 'pkg:eventbus', target: 'fn:Publish', style: { stroke: '#22d3ee', strokeWidth: 1.5 } },
  { id: 'e5', source: 'fn:StartAll', target: 'fn:Publish', animated: true, style: { stroke: '#8b5cf6', strokeWidth: 1.5, strokeDasharray: '4 4' } },
];

export const GraphView: React.FC = () => {
  const { graphNodes, graphEdges, fetchGraph } = useStore();
  const [selectedNode, setSelectedNode] = useState<Node | null>(initialNodes[0]);

  useEffect(() => {
    fetchGraph();
  }, [fetchGraph]);

  const activeNodes = graphNodes.length > 0 ? graphNodes.map((n, i) => ({
    id: n.id,
    type: 'default',
    data: { label: `${n.type.toLowerCase()}: ${n.name}`, type: n.type, path: n.path },
    position: { x: (i % 3) * 200 + 50, y: Math.floor(i / 3) * 150 + 50 },
    style: {
      background: '#12141c',
      color: n.type === 'FILE' ? '#60a5fa' : n.type === 'PACKAGE' ? '#22d3ee' : '#c084fc',
      border: `1px solid ${n.type === 'FILE' ? '#3b82f6' : n.type === 'PACKAGE' ? '#22d3ee' : '#8b5cf6'}`,
      borderRadius: '12px',
      padding: '10px 14px',
      fontSize: '12px',
      fontFamily: 'JetBrains Mono, monospace',
    },
  })) : initialNodes;

  const activeEdges = graphEdges.length > 0 ? graphEdges.map((e) => ({
    id: e.id,
    source: e.from,
    target: e.to,
    animated: e.type === 'CALLS' || e.type === 'DEPENDS_ON',
    style: { stroke: e.type === 'CALLS' ? '#8b5cf6' : '#6366f1', strokeWidth: 1.5 },
  })) : initialEdges;

  return (
    <div className="h-[calc(100vh-4rem)] flex relative">
      {/* React Flow Canvas */}
      <div className="flex-1 h-full bg-[#0b0d14] relative">
        <ReactFlow
          nodes={activeNodes}
          edges={activeEdges}
          onNodeClick={(_, node) => setSelectedNode(node)}
          fitView
        >
          <Background color="#1e2235" gap={24} size={1} />
          <Controls className="bg-[#12141c] border border-white/10 text-white rounded-xl" />
        </ReactFlow>

        {/* Node & Edge Legend Overlay */}
        <div className="absolute top-4 left-4 p-3 rounded-xl glass-card bg-[#12141c]/90 border border-white/10 text-xs font-mono space-y-2 z-10">
          <div className="text-[10px] text-gray-400 uppercase tracking-wider font-semibold border-b border-white/5 pb-1">
            Graph Legend
          </div>
          <div className="flex items-center gap-4 text-[11px]">
            <div className="flex items-center gap-1.5">
              <span className="w-2.5 h-2.5 rounded-full bg-blue-500"></span>
              <span className="text-gray-300">File</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="w-2.5 h-2.5 rounded-full bg-cyan-400"></span>
              <span className="text-gray-300">Package</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="w-2.5 h-2.5 rounded-full bg-violet-500"></span>
              <span className="text-gray-300">Function</span>
            </div>
          </div>
          <div className="flex items-center gap-4 text-[10px] text-gray-400 pt-1 border-t border-white/5">
            <div className="flex items-center gap-1.5">
              <span className="w-4 h-0.5 bg-indigo-500"></span>
              <span>Direct Dependency</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="w-4 h-0.5 bg-violet-400 border-b border-dashed border-violet-400"></span>
              <span>Indirect Ref</span>
            </div>
          </div>
        </div>
      </div>

      {/* Node Inspector Drawer */}
      {selectedNode && (
        <div className="w-80 glass-panel border-l border-white/10 p-5 space-y-4 z-10 overflow-y-auto shrink-0 bg-[#12141c]">
          <div className="flex items-center gap-2 text-xs font-mono text-cyan-400 uppercase tracking-wider border-b border-white/10 pb-3">
            <Info className="w-4 h-4 text-cyan-400" />
            <span>Entity Inspector</span>
          </div>

          <div className="space-y-3 divide-y divide-white/5">
            <div className="pb-2">
              <div className="text-[10px] font-mono text-gray-400 uppercase">Entity ID</div>
              <div className="text-xs font-mono font-bold text-gray-100 mt-0.5">{selectedNode.id}</div>
            </div>
            <div className="pt-3 pb-2">
              <div className="text-[10px] font-mono text-gray-400 uppercase mb-1">Type</div>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-indigo-500/10 text-indigo-300 border border-indigo-500/20">
                {selectedNode.data.type as string}
              </span>
            </div>
            <div className="pt-3 pb-2">
              <div className="text-[10px] font-mono text-gray-400 uppercase mb-0.5">Path</div>
              <div className="text-xs font-mono text-gray-300 break-all">{selectedNode.data.path as string}</div>
            </div>
            <div className="pt-3">
              <div className="text-[10px] font-mono text-gray-400 uppercase mb-1">Direct Dependencies</div>
              <div className="p-2.5 rounded-lg bg-black/40 border border-white/5 text-xs font-mono text-gray-300">
                - {selectedNode.id === 'file:main.go' ? 'internal/coordinator, internal/eventbus' : 'internal/observability'}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
