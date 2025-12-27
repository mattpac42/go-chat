'use client';

// Types for Cost Savings data
export interface CostSavingsData {
  pmMinutes: number;        // Minutes of PM-equivalent work
  devHours: number;         // Hours of dev-equivalent work
  messageCount: number;     // Total messages in session
  filesGenerated: number;   // Files created
  tokensUsed: number;       // Approximate tokens (for AI cost calc)
}

export interface CostSavingsCardProps {
  data: CostSavingsData;
  showDetailed?: boolean;   // Show breakdown or just totals
}

// Rate constants (midpoint values)
const PM_HOURLY_RATE = 80;        // $80/hr for PM consulting
const DEV_HOURLY_RATE = 112.50;   // $112.50/hr for development

// AI cost per 1K tokens
const INPUT_COST_PER_1K = 0.003;  // $0.003 per 1K input tokens
const OUTPUT_COST_PER_1K = 0.015; // $0.015 per 1K output tokens
const INPUT_OUTPUT_RATIO = 0.6;   // 60% input, 40% output estimate

/**
 * Calculate PM consulting value based on minutes spent
 */
function calculatePmValue(minutes: number): number {
  return (minutes / 60) * PM_HOURLY_RATE;
}

/**
 * Calculate development value based on hours spent
 */
function calculateDevValue(hours: number): number {
  return hours * DEV_HOURLY_RATE;
}

/**
 * Calculate estimated AI cost based on token usage
 * Uses 60/40 input/output split estimation
 */
function calculateAiCost(tokens: number): number {
  const inputTokens = tokens * INPUT_OUTPUT_RATIO;
  const outputTokens = tokens * (1 - INPUT_OUTPUT_RATIO);

  const inputCost = (inputTokens / 1000) * INPUT_COST_PER_1K;
  const outputCost = (outputTokens / 1000) * OUTPUT_COST_PER_1K;

  return inputCost + outputCost;
}

/**
 * Format currency for display
 */
function formatCurrency(value: number): string {
  if (value < 1) {
    return `$${value.toFixed(2)}`;
  }
  return `$${Math.round(value).toLocaleString()}`;
}

/**
 * Format time for display
 */
function formatTime(minutes: number, isHours: boolean = false): string {
  if (isHours) {
    const hours = minutes;
    if (hours < 1) {
      return `${Math.round(hours * 60)} min`;
    }
    return `${hours.toFixed(1)} hours`;
  }

  if (minutes < 60) {
    return `${Math.round(minutes)} min`;
  }
  const hours = minutes / 60;
  return `${hours.toFixed(1)} hours`;
}

/**
 * Calculate total hours equivalent for summary
 */
function calculateTotalHours(pmMinutes: number, devHours: number): number {
  return (pmMinutes / 60) + devHours;
}

// Icon components
function UserIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
      />
    </svg>
  );
}

function CodeIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
      />
    </svg>
  );
}

function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <span className="text-xs font-semibold text-gray-500 uppercase tracking-wide">
      {children}
    </span>
  );
}

interface ValueCardProps {
  icon: React.ReactNode;
  label: string;
  timeValue: string;
  monetaryValue: string;
}

function ValueCard({ icon, label, timeValue, monetaryValue }: ValueCardProps) {
  return (
    <div className="bg-gray-50 rounded-lg p-4 flex-1">
      <div className="flex items-center gap-2 mb-2">
        {icon}
        <SectionLabel>{label}</SectionLabel>
      </div>
      <p className="text-lg font-semibold text-gray-900">{timeValue}</p>
      <p className="text-sm text-teal-600 font-medium">~{monetaryValue} value</p>
    </div>
  );
}

/**
 * CostSavingsCard - Displays savings comparison between professional rates and AI costs
 *
 * Shows two-column layout with PM consulting and development time equivalents,
 * plus a bottom panel showing total value delivered vs actual AI cost.
 */
export function CostSavingsCard({
  data,
  showDetailed = true,
}: CostSavingsCardProps) {
  const { pmMinutes, devHours, messageCount, filesGenerated, tokensUsed } = data;

  // Calculate values
  const pmValue = calculatePmValue(pmMinutes);
  const devValue = calculateDevValue(devHours);
  const totalValue = pmValue + devValue;
  const aiCost = calculateAiCost(tokensUsed);
  const totalHours = calculateTotalHours(pmMinutes, devHours);

  return (
    <div className="bg-white rounded-lg border border-gray-200 shadow-sm p-4 md:p-6">
      {/* Header */}
      <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-4">
        Your Savings
      </h3>

      {/* Two-column value cards */}
      <div className="flex flex-col sm:flex-row gap-3 mb-4">
        <ValueCard
          icon={<UserIcon className="w-4 h-4 text-gray-400" />}
          label="PM Time"
          timeValue={formatTime(pmMinutes)}
          monetaryValue={formatCurrency(pmValue)}
        />
        <ValueCard
          icon={<CodeIcon className="w-4 h-4 text-gray-400" />}
          label="Dev Time"
          timeValue={formatTime(devHours, true)}
          monetaryValue={formatCurrency(devValue)}
        />
      </div>

      {/* Total value panel */}
      <div className="bg-teal-50 rounded-lg p-4 border border-teal-100">
        <SectionLabel>Total Value Delivered</SectionLabel>
        <p className="text-2xl font-bold text-teal-700 mt-1">
          {formatCurrency(totalValue)}
        </p>
        <p className="text-sm text-gray-600 mt-1">
          Equivalent to ~{totalHours.toFixed(1)} hours of consulting
        </p>
        <p className="text-sm text-gray-600">
          Your actual AI cost: <span className="font-medium">~{formatCurrency(aiCost)}</span>
        </p>
      </div>

      {/* Detailed breakdown (optional) */}
      {showDetailed && (
        <div className="mt-4 pt-4 border-t border-gray-100 grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">Messages exchanged:</span>
            <span className="ml-2 text-gray-900 font-medium">{messageCount}</span>
          </div>
          <div>
            <span className="text-gray-500">Files generated:</span>
            <span className="ml-2 text-gray-900 font-medium">{filesGenerated}</span>
          </div>
        </div>
      )}

      {/* Disclaimer */}
      <p className="mt-4 text-xs text-gray-400">
        * Estimates based on average industry rates. PM rate: ${PM_HOURLY_RATE}/hr, Dev rate: ${DEV_HOURLY_RATE}/hr.
      </p>
    </div>
  );
}
