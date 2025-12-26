'use client';

export type DiscoveryStage = 'welcome' | 'problem' | 'personas' | 'mvp' | 'summary' | 'complete';

interface DiscoveryProgressProps {
  currentStage: DiscoveryStage;
  onStageClick?: (stage: DiscoveryStage) => void;
  isMobile?: boolean;
}

const STAGES: DiscoveryStage[] = ['welcome', 'problem', 'personas', 'mvp', 'summary'];

const STAGE_LABELS: Record<DiscoveryStage, string> = {
  welcome: 'Welcome',
  problem: 'Problem',
  personas: 'Personas',
  mvp: 'MVP',
  summary: 'Summary',
  complete: 'Complete',
};

function getStageIndex(stage: DiscoveryStage): number {
  if (stage === 'complete') return STAGES.length;
  return STAGES.indexOf(stage);
}

function isStageCompleted(stage: DiscoveryStage, currentStage: DiscoveryStage): boolean {
  const currentIndex = getStageIndex(currentStage);
  const stageIndex = getStageIndex(stage);
  return stageIndex < currentIndex;
}

function isStageCurrent(stage: DiscoveryStage, currentStage: DiscoveryStage): boolean {
  return stage === currentStage && currentStage !== 'complete';
}

export function DiscoveryProgress({
  currentStage,
  onStageClick,
  isMobile = false,
}: DiscoveryProgressProps) {
  const currentIndex = getStageIndex(currentStage);
  const displayNumber = currentStage === 'complete' ? 5 : Math.min(currentIndex + 1, 5);

  const handleClick = () => {
    if (onStageClick) {
      onStageClick(currentStage);
    }
  };

  if (isMobile) {
    return (
      <button
        onClick={handleClick}
        className="flex items-center gap-2 px-2 py-1 rounded-md hover:bg-gray-100 transition-colors"
        aria-label={`Discovery progress: ${displayNumber} of 5, current stage: ${STAGE_LABELS[currentStage]}`}
      >
        <div className="flex items-center gap-1.5">
          {STAGES.map((stage, index) => {
            const completed = isStageCompleted(stage, currentStage);
            const current = isStageCurrent(stage, currentStage);
            const filled = completed || current || currentStage === 'complete';

            return (
              <div
                key={stage}
                className={`w-2.5 h-2.5 rounded-full ${
                  filled
                    ? 'bg-teal-500'
                    : 'border-2 border-gray-300 bg-transparent'
                }`}
                aria-label={`${STAGE_LABELS[stage]}: ${completed ? 'completed' : current ? 'current' : 'upcoming'}`}
              />
            );
          })}
        </div>
        <span className="text-sm text-gray-500">{displayNumber}/5</span>
      </button>
    );
  }

  return (
    <div
      className="flex items-center gap-2"
      role="progressbar"
      aria-valuenow={displayNumber}
      aria-valuemin={1}
      aria-valuemax={5}
      aria-label={`Discovery progress: ${displayNumber} of 5`}
    >
      <div className="flex items-center gap-2">
        {STAGES.map((stage, index) => {
          const completed = isStageCompleted(stage, currentStage);
          const current = isStageCurrent(stage, currentStage);
          const filled = completed || current || currentStage === 'complete';

          return (
            <div key={stage} className="flex items-center gap-2">
              <div
                className={`w-3 h-3 rounded-full ${
                  filled
                    ? 'bg-teal-500'
                    : 'border-2 border-gray-300 bg-transparent'
                }`}
                aria-label={`${STAGE_LABELS[stage]}: ${completed ? 'completed' : current ? 'current' : 'upcoming'}`}
              />
              {current && (
                <span className="text-sm font-medium text-gray-700">
                  {STAGE_LABELS[stage]}
                </span>
              )}
            </div>
          );
        })}
      </div>
      <span className="text-sm text-gray-500 ml-2">
        {displayNumber} of 5
      </span>
    </div>
  );
}
