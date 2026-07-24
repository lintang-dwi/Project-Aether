import React, { useEffect } from 'react';
import Editor from '@monaco-editor/react';
import { useStore } from '../store/useStore';
import { FileCode, Folder, Code } from 'lucide-react';

const mockFiles = [
  'main.go',
  'assets.go',
  'internal/coordinator/coordinator.go',
  'internal/eventbus/eventbus.go',
  'internal/storage/sqlite.go',
  'internal/parser/treesitter_go.go',
  'internal/context/builder.go',
  'internal/action/operations.go',
  'internal/git/engine.go',
  'internal/security/boundary.go',
];

export const WorkspaceView: React.FC = () => {
  const { selectedFile, setSelectedFile, fileContent, setFileContent, discoveredFiles, fetchFileContent, scanWorkspace } = useStore();
  const activeFile = selectedFile || 'main.go';

  useEffect(() => {
    fetchFileContent(activeFile);
  }, [activeFile, fetchFileContent]);

  const fileList = discoveredFiles.length > 0 ? discoveredFiles.map((f) => f.path) : mockFiles;

  const handleSelect = (file: string) => {
    setSelectedFile(file);
  };

  const getLanguage = (filename: string) => {
    if (filename.endsWith('.go')) return 'go';
    if (filename.endsWith('.ts') || filename.endsWith('.tsx')) return 'typescript';
    if (filename.endsWith('.js') || filename.endsWith('.jsx')) return 'javascript';
    if (filename.endsWith('.json')) return 'json';
    if (filename.endsWith('.md')) return 'markdown';
    return 'go';
  };

  return (
    <div className="h-[calc(100vh-4rem)] flex">
      {/* File Tree Explorer */}
      <div className="w-64 glass-panel border-r border-white/10 p-4 space-y-3 z-10 overflow-y-auto shrink-0 bg-[#12141c]">
        <div className="flex items-center justify-between text-xs font-mono text-gray-400 uppercase tracking-wider">
          <div className="flex items-center gap-2">
            <Folder className="w-4 h-4 text-indigo-400" />
            <span>Workspace Files</span>
          </div>
          <span className="text-[10px] text-cyan-400 font-bold">{fileList.length}</span>
        </div>

        <div className="space-y-1">
          {fileList.map((file) => {
            const isSelected = activeFile === file;
            return (
              <button
                key={file}
                onClick={() => handleSelect(file)}
                className={`w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-xs font-mono transition-all ${
                  isSelected
                    ? 'bg-indigo-600/20 text-cyan-300 border border-indigo-500/30 font-semibold'
                    : 'text-gray-400 hover:text-gray-200 hover:bg-white/5'
                }`}
              >
                <FileCode className="w-3.5 h-3.5 text-cyan-400 shrink-0" />
                <span className="truncate">{file}</span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Monaco Editor Canvas */}
      <div className="flex-1 h-full flex flex-col bg-[#0b0d14]">
        {/* Editor Top Bar */}
        <div className="h-10 glass-panel border-b border-white/10 px-4 flex items-center justify-between text-xs font-mono text-gray-300 shrink-0 bg-[#12141c]">
          <div className="flex items-center gap-2">
            <Code className="w-4 h-4 text-cyan-400" />
            <span className="font-semibold text-gray-200">{activeFile}</span>
          </div>
          <span className="px-2 py-0.5 rounded bg-indigo-500/10 border border-indigo-500/20 text-[10px] text-indigo-300 font-mono">
            AST Parsed · Tree-sitter {getLanguage(activeFile).toUpperCase()}
          </span>
        </div>

        <div className="flex-1">
          <Editor
            height="100%"
            language={getLanguage(activeFile)}
            theme="vs-dark"
            value={fileContent}
            onChange={(v) => setFileContent(v || '')}
            options={{
              fontSize: 13,
              fontFamily: 'JetBrains Mono, monospace',
              minimap: { enabled: false },
              scrollBeyondLastLine: false,
              automaticLayout: true,
            }}
          />
        </div>
      </div>
    </div>
  );
};
