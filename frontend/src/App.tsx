import React from 'react';
import { Sidebar } from './components/Sidebar';
import { Header } from './components/Header';
import { useStore } from './store/useStore';
import { DashboardView } from './views/DashboardView';
import { GraphView } from './views/GraphView';
import { WorkspaceView } from './views/WorkspaceView';
import { TaskView } from './views/TaskView';
import { SecurityView } from './views/SecurityView';

export const App: React.FC = () => {
  const { activeView } = useStore();

  const renderView = () => {
    switch (activeView) {
      case 'dashboard':
        return <DashboardView />;
      case 'graph':
        return <GraphView />;
      case 'workspace':
        return <WorkspaceView />;
      case 'task':
        return <TaskView />;
      case 'security':
        return <SecurityView />;
      default:
        return <DashboardView />;
    }
  };

  return (
    <div className="flex h-screen w-screen bg-[#07080c] overflow-hidden">
      <Sidebar />
      <div className="flex-1 flex flex-col h-screen overflow-hidden">
        <Header />
        <main className="flex-1 overflow-hidden bg-gradient-to-b from-[#07080c] to-[#0b0d14]">
          {renderView()}
        </main>
      </div>
    </div>
  );
};

export default App;
