import { useMemo } from 'react';
import type { CostSavingsData } from '@/components/savings';
import { useWageSettings, DEFAULT_WAGE_SETTINGS } from './useWageSettings';

// AI cost per 1K tokens
const INPUT_COST_PER_1K = 0.003;  // $0.003 per 1K input tokens
const OUTPUT_COST_PER_1K = 0.015; // $0.015 per 1K output tokens
const INPUT_OUTPUT_RATIO = 0.6;   // 60% input, 40% output estimate

// Estimation factors for deriving values from session metrics
const PM_MINUTES_PER_MESSAGE = 1.5;     // Each message represents ~1.5 min of PM thinking
const DEV_HOURS_PER_FILE = 0.5;         // Each file represents ~30 min of dev work
const DEV_HOURS_PER_100_MESSAGES = 0.25; // Additional dev time for message complexity
const DESIGNER_HOURS_PER_MESSAGE = 0.5; // Each designer message represents ~30 min of design work
const TOKENS_PER_MESSAGE = 500;         // Average tokens per message (both directions)

/**
 * Session data that can be used to estimate cost savings
 */
export interface SessionMetrics {
  messageCount: number;
  filesGenerated: number;
  tokensUsed?: number; // If available, use actual tokens; otherwise estimate
  designerMessageCount?: number; // Count of messages from designer agent
}

/**
 * Calculated savings result
 */
export interface CostSavingsResult {
  data: CostSavingsData;
  pmValue: number;
  devValue: number;
  designerValue: number;
  totalValue: number;
  aiCost: number;
  savingsMultiplier: number; // How many times more value vs cost
}

/**
 * Estimate PM minutes based on message count
 * Assumes each message exchange represents about 1.5 minutes of PM thinking/planning
 */
function estimatePmMinutes(messageCount: number): number {
  return messageCount * PM_MINUTES_PER_MESSAGE;
}

/**
 * Estimate dev hours based on files generated and message complexity
 * Each file represents significant dev work, messages add incremental time
 */
function estimateDevHours(filesGenerated: number, messageCount: number): number {
  const fileBasedHours = filesGenerated * DEV_HOURS_PER_FILE;
  const messageBasedHours = (messageCount / 100) * DEV_HOURS_PER_100_MESSAGES;
  return fileBasedHours + messageBasedHours;
}

/**
 * Estimate designer hours based on designer message count
 * Each substantial designer response represents ~30 min of design work
 */
function estimateDesignerHours(designerMessageCount: number): number {
  return designerMessageCount * DESIGNER_HOURS_PER_MESSAGE;
}

/**
 * Estimate tokens if not provided
 */
function estimateTokens(messageCount: number): number {
  return messageCount * TOKENS_PER_MESSAGE;
}

/**
 * Calculate PM consulting value based on minutes spent
 */
function calculatePmValue(minutes: number, hourlyRate: number): number {
  return (minutes / 60) * hourlyRate;
}

/**
 * Calculate development value based on hours spent
 */
function calculateDevValue(hours: number, hourlyRate: number): number {
  return hours * hourlyRate;
}

/**
 * Calculate designer value based on hours spent
 */
function calculateDesignerValue(hours: number, hourlyRate: number): number {
  return hours * hourlyRate;
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
 * Hook to calculate cost savings from session metrics
 *
 * Uses rates from the useWageSettings hook for calculations.
 *
 * @param metrics - Session metrics (messageCount, filesGenerated, optionally tokensUsed, designerMessageCount)
 * @returns Calculated savings data and computed values
 *
 * @example
 * ```tsx
 * const { data, totalValue, aiCost, savingsMultiplier } = useCostSavings({
 *   messageCount: 25,
 *   filesGenerated: 5,
 *   designerMessageCount: 3,
 * });
 *
 * // Use with CostSavingsCard
 * <CostSavingsCard data={data} />
 * ```
 */
export function useCostSavings(metrics: SessionMetrics): CostSavingsResult {
  const { settings } = useWageSettings();

  return useMemo(() => {
    const { messageCount, filesGenerated, tokensUsed, designerMessageCount = 0 } = metrics;
    const { pmHourlyRate, devHourlyRate, designerHourlyRate } = settings;

    // Estimate or use provided values
    const pmMinutes = estimatePmMinutes(messageCount);
    const devHours = estimateDevHours(filesGenerated, messageCount);
    const designerHours = estimateDesignerHours(designerMessageCount);
    const tokens = tokensUsed ?? estimateTokens(messageCount);

    // Calculate monetary values using configured rates
    const pmValue = calculatePmValue(pmMinutes, pmHourlyRate);
    const devValue = calculateDevValue(devHours, devHourlyRate);
    const designerValue = calculateDesignerValue(designerHours, designerHourlyRate);
    const totalValue = pmValue + devValue + designerValue;
    const aiCost = calculateAiCost(tokens);

    // Calculate savings multiplier (how many times more value vs cost)
    const savingsMultiplier = aiCost > 0 ? totalValue / aiCost : 0;

    const data: CostSavingsData = {
      pmMinutes,
      devHours,
      designerHours,
      messageCount,
      filesGenerated,
      tokensUsed: tokens,
    };

    return {
      data,
      pmValue,
      devValue,
      designerValue,
      totalValue,
      aiCost,
      savingsMultiplier,
    };
  }, [metrics, settings]);
}

// Export calculation functions for direct use if needed
export {
  calculatePmValue,
  calculateDevValue,
  calculateDesignerValue,
  calculateAiCost,
  estimatePmMinutes,
  estimateDevHours,
  estimateDesignerHours,
};

// Export default rates for backward compatibility
export const PM_HOURLY_RATE = DEFAULT_WAGE_SETTINGS.pmHourlyRate;
export const DEV_HOURLY_RATE = DEFAULT_WAGE_SETTINGS.devHourlyRate;
export const DESIGNER_HOURLY_RATE = DEFAULT_WAGE_SETTINGS.designerHourlyRate;
