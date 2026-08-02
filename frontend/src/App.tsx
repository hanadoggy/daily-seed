import { useEffect, useState } from 'react';
import { Leaf, Settings, BarChart2, Pencil, Eye, Unlock, CalendarDays } from 'lucide-react';
import { Routes, Route, Link, useLocation, useNavigate } from 'react-router-dom';
import { AutoMigrationPrompt } from '@/features/progress/AutoMigrationPrompt';
import { ManageModal } from '@/features/manage/ManageModal';
import { useAppStore } from '@/store/useAppStore';
import { todayJST } from '@/lib/dayjs';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { ThemeSwitcher } from '@/components/theme-switcher';
import { MainPage } from '@/pages/MainPage';
import { AnalyticsPage } from '@/pages/AnalyticsPage';

const MODE_CLASS_MAP = {
  Growth: 'mode-growth',
  Rest: 'mode-rest',
  Office: 'mode-office',
  Remote: 'mode-remote',
} as const;

function App() {
  const { currentMode, selectedDate, setDateAndFetch, fetchMasterData, error, isAdminMode, toggleAdminMode } = useAppStore();
  const [initialized, setInitialized] = useState(false);
  const [showManage, setShowManage] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (!initialized) {
      setDateAndFetch(todayJST());
      fetchMasterData();
      setInitialized(true);
    }
  }, [initialized, setDateAndFetch, fetchMasterData]);

  const modeClass = MODE_CLASS_MAP[currentMode];
  const formattedDate = selectedDate
    ? new Date(selectedDate + 'T00:00:00+09:00').toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      })
    : '';

  const isAnalytics = location.pathname === '/analytics';
  const isToday = selectedDate === todayJST();

  return (
    <div className={cn('min-h-screen bg-background transition-colors duration-500', modeClass)}>
      <div className="mx-auto max-w-6xl px-4 py-6 lg:px-8">
        {/* Header */}
        <header className="mb-8">
          <div className="flex items-center gap-3 mb-1">
            <Link to="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-mode-accent-soft">
                <Leaf className="h-5 w-5 text-mode-accent" />
              </div>
              <h1 className="text-xl font-bold tracking-tight">Daily Seed</h1>
            </Link>
            
            {!isAnalytics && (
              <Button
                variant="outline"
                size="sm"
                disabled={isToday}
                onClick={() => setDateAndFetch(todayJST())}
                className={cn(
                  'ml-1 gap-1.5 transition-all',
                  isToday
                    ? 'border-mode-accent/30 text-mode-accent/50 bg-mode-accent/5 opacity-60 cursor-not-allowed'
                    : 'border-mode-accent/60 text-mode-accent bg-mode-accent/10 hover:bg-mode-accent/20 cursor-pointer shadow-sm font-medium'
                )}
              >
                <CalendarDays className="h-3.5 w-3.5" />
                Today
              </Button>
            )}

            <div className="ml-auto flex items-center gap-2">
              {/* Edit / View / Admin Mode Button (Dashboard only) */}
              {!isAnalytics && (
                isToday ? (
                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-1.5 border-primary/30 text-primary bg-primary/5 cursor-default"
                  >
                    <Pencil className="h-3.5 w-3.5" />
                    Edit Mode
                  </Button>
                ) : isAdminMode ? (
                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-1.5 border-destructive/40 text-destructive hover:bg-destructive/10"
                    onClick={toggleAdminMode}
                  >
                    <Unlock className="h-3.5 w-3.5" />
                    Admin Mode
                  </Button>
                ) : (
                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-1.5 border-accent/40 text-accent-foreground hover:bg-accent/10"
                    onClick={toggleAdminMode}
                  >
                    <Eye className="h-3.5 w-3.5" />
                    View Mode
                  </Button>
                )
              )}

              <ThemeSwitcher />
              
              {isAnalytics ? (
                <Button
                  variant="outline"
                  size="sm"
                  className="gap-1.5 border-green-500/30 text-green-600 hover:bg-green-500/10 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300"
                  onClick={() => navigate('/')}
                >
                  <Leaf className="h-3.5 w-3.5" />
                  Dashboard
                </Button>
              ) : (
                <Button
                  variant="outline"
                  size="sm"
                  className="gap-1.5 border-orange-500/30 text-orange-600 hover:bg-orange-500/10 hover:text-orange-700 dark:text-orange-400 dark:hover:text-orange-300"
                  onClick={() => navigate('/analytics')}
                >
                  <BarChart2 className="h-3.5 w-3.5" />
                  Analytics
                </Button>
              )}

              <Button
                id="manage-button"
                variant="outline"
                size="sm"
                className="gap-1.5"
                onClick={() => setShowManage(true)}
              >
                <Settings className="h-3.5 w-3.5" />
                Manage
              </Button>
            </div>
          </div>
          {selectedDate && (
            <p className="text-sm text-muted-foreground ml-12">{formattedDate}</p>
          )}
        </header>

        {/* Error banner */}
        {error && (
          <div className="mb-6 rounded-xl border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
            {error}
          </div>
        )}

        <Routes>
          <Route path="/" element={<MainPage />} />
          <Route path="/analytics" element={<AnalyticsPage />} />
        </Routes>
      </div>

      {/* Management Panel */}
      <ManageModal open={showManage} onOpenChange={setShowManage} />
      <AutoMigrationPrompt />
    </div>
  );
}

export default App;
