import { useEffect } from 'react';
import { Sparkles, Trophy, Award, X } from 'lucide-react';

interface MilestoneCelebrationModalProps {
  isOpen: boolean;
  onClose: () => void;
  habitTitle: string;
  milestone: number;
}

export function MilestoneCelebrationModal({
  isOpen,
  onClose,
  habitTitle,
  milestone,
}: MilestoneCelebrationModalProps) {
  useEffect(() => {
    if (!isOpen) return;

    // Auto-dismiss after 2 seconds
    const timer = setTimeout(() => {
      onClose();
    }, 2500);

    return () => clearTimeout(timer);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-in fade-in duration-300">
      <div className="relative w-full max-w-sm overflow-hidden rounded-2xl border border-amber-500/30 bg-card p-6 shadow-2xl animate-in zoom-in-95 slide-in-from-bottom-4 duration-300">
        {/* Animated Background Glow */}
        <div className="absolute -top-12 -left-12 w-32 h-32 bg-amber-500/20 rounded-full blur-2xl animate-pulse" />
        <div className="absolute -bottom-12 -right-12 w-32 h-32 bg-orange-500/20 rounded-full blur-2xl animate-pulse" />

        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-3 right-3 p-1 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors"
          aria-label="Close modal"
        >
          <X className="w-4 h-4" />
        </button>

        {/* Modal Content */}
        <div className="flex flex-col items-center text-center space-y-4 pt-2">
          <div className="relative">
            <div className="p-4 rounded-full bg-gradient-to-br from-amber-400/20 to-orange-500/20 border border-amber-500/30 text-amber-500 shadow-inner">
              <Trophy className="w-10 h-10 animate-bounce" />
            </div>
            <Sparkles className="w-5 h-5 text-amber-400 absolute -top-1 -right-1 animate-spin duration-1000" />
            <Award className="w-4 h-4 text-orange-400 absolute -bottom-1 -left-1" />
          </div>

          <div className="space-y-1">
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-500 border border-amber-500/20">
              <Sparkles className="w-3.5 h-3.5" />
              Milestone Achieved!
            </span>
            <h3 className="text-xl font-extrabold text-foreground pt-1">
              {milestone}-Day Streak!
            </h3>
            <p className="text-sm text-muted-foreground font-medium">
              {habitTitle}
            </p>
          </div>

          <p className="text-xs text-muted-foreground italic bg-muted/30 px-3 py-1.5 rounded-lg border border-border/50">
            Outstanding consistency! Keep up the incredible momentum.
          </p>
        </div>
      </div>
    </div>
  );
}
