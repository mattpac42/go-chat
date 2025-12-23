#!/usr/bin/env python3
"""
Automated Context Tracking with A/B Testing
Implements both Option 1 (system warning parsing) and Option 2 (env/API check)
Compares results to determine most accurate method
"""

import re
import json
import sys
import os
import subprocess
from datetime import datetime
from pathlib import Path

class ContextTracker:
    """Production-ready context tracker with A/B testing"""

    STATE_FILE = Path('.claude/.context-state.json')
    AB_TEST_LOG = Path('.claude/.context-ab-test.jsonl')

    def __init__(self):
        self.load_state()

    def load_state(self):
        """Load previous context state"""
        if self.STATE_FILE.exists():
            try:
                with open(self.STATE_FILE) as f:
                    state = json.load(f)
                    self.estimated_tokens = state.get('estimated_tokens', 30000)
                    self.last_actual = state.get('last_actual', 30000)
                    self.last_handoff_pct = state.get('last_handoff_pct', 0)
                    self.preferred_method = state.get('preferred_method', 'auto')
            except:
                self._reset_state()
        else:
            self._reset_state()

    def _reset_state(self):
        """Initialize default state"""
        self.estimated_tokens = 30000  # System overhead baseline
        self.last_actual = 30000
        self.last_handoff_pct = 0
        self.preferred_method = 'auto'

    # ==================== OPTION 0: Parse /context CLI Command ====================

    def method_context_command(self, context_output=None):
        """
        Option 0: Parse /context CLI command output directly (highest confidence)
        Pattern: claude-sonnet-4-5-20250929 · 44k/200k tokens (22%)

        Args:
            context_output: Optional string containing /context command output
                           If None, will check CLAUDE_CONTEXT_OUTPUT env var
        """
        # Get context output from parameter, environment, or return None
        output = context_output

        if not output:
            # Check environment variable
            output = os.environ.get('CLAUDE_CONTEXT_OUTPUT')

        if not output:
            # No context data available
            return None

        try:
            # Parse pattern: "44k/200k tokens (22%)"
            match = re.search(r'(\d+)k/(\d+)k tokens \((\d+)%\)', output)

            if match:
                used_k = int(match.group(1))
                total_k = int(match.group(2))
                percentage = int(match.group(3))

                used = used_k * 1000
                total = total_k * 1000

                return {
                    'used': used,
                    'total': total,
                    'remaining': total - used,
                    'percentage': percentage,
                    'method': 'context_command',
                    'confidence': 'high'
                }
        except (ValueError, AttributeError):
            # Parsing failed, return None to try other methods
            pass

        return None

    # ==================== OPTION 1: Parse System Warnings ====================

    def method_system_warning(self, text):
        """
        Option 1: Parse system warning from Claude's response
        Pattern: <system_warning>Token usage: 91498/200000; 108502 remaining</system_warning>
        """
        # Try XML-style system warning first
        match = re.search(r'<system_warning>Token usage: (\d+)/(\d+); (\d+) remaining</system_warning>',
                         text, re.IGNORECASE)

        if match:
            used = int(match.group(1))
            total = int(match.group(2))
            remaining = int(match.group(3))
            return {
                'used': used,
                'total': total,
                'remaining': remaining,
                'percentage': round((used / total) * 100),
                'method': 'system_warning',
                'confidence': 'high'
            }

        # Try plain text format: "Token usage: 91498/200000"
        match = re.search(r'Token usage:\s*(\d+)/(\d+)', text, re.IGNORECASE)
        if match:
            used = int(match.group(1))
            total = int(match.group(2))
            remaining = total - used
            return {
                'used': used,
                'total': total,
                'remaining': remaining,
                'percentage': round((used / total) * 100),
                'method': 'system_warning_plain',
                'confidence': 'high'
            }

        return None

    # ==================== OPTION 2: Environment Variables / API ====================

    def method_environment(self):
        """
        Option 2: Check for Claude Code environment variables or API
        """
        # Check for token usage in environment
        tokens_used = os.environ.get('CLAUDE_TOKENS_USED')
        tokens_total = os.environ.get('CLAUDE_TOKENS_TOTAL')

        if tokens_used and tokens_total:
            used = int(tokens_used)
            total = int(tokens_total)
            return {
                'used': used,
                'total': total,
                'remaining': total - used,
                'percentage': round((used / total) * 100),
                'method': 'environment',
                'confidence': 'high'
            }

        # Check for SSE port and try API query
        sse_port = os.environ.get('CLAUDE_CODE_SSE_PORT')
        if sse_port:
            # Future: Could try querying the SSE endpoint
            # For now, return None (not implemented)
            pass

        return None

    # ==================== OPTION 3: Estimation Fallback ====================

    def method_estimation(self, response_length=5000):
        """
        Fallback: Estimate based on response length
        Conservative estimate: ~5k tokens per substantial exchange
        """
        # Update estimate
        self.estimated_tokens += response_length

        return {
            'used': self.estimated_tokens,
            'total': 200000,
            'remaining': 200000 - self.estimated_tokens,
            'percentage': round((self.estimated_tokens / 200000) * 100),
            'method': 'estimation',
            'confidence': 'low'
        }

    # ==================== A/B Testing Logic ====================

    def get_context_info(self, response_text="", context_output=None):
        """
        Run A/B test: Try all methods and compare results
        Returns: (best_result, all_results)

        Args:
            response_text: Claude's response text for parsing
            context_output: Optional /context command output string
        """
        results = {}

        # Try Option 0: /context CLI command (highest priority)
        result_0 = self.method_context_command(context_output)
        if result_0:
            results['context_command'] = result_0

        # Try Option 1: System warning parsing
        result_1 = self.method_system_warning(response_text)
        if result_1:
            results['system_warning'] = result_1

        # Try Option 2: Environment variables
        result_2 = self.method_environment()
        if result_2:
            results['environment'] = result_2

        # Always have Option 3: Estimation fallback
        result_3 = self.method_estimation(len(response_text) // 4)
        results['estimation'] = result_3

        # Select best result based on confidence and availability
        if self.preferred_method == 'auto':
            # Auto-select: Prefer high-confidence methods in priority order
            if 'context_command' in results:
                best = results['context_command']
            elif 'system_warning' in results:
                best = results['system_warning']
            elif 'environment' in results:
                best = results['environment']
            else:
                best = results['estimation']
        else:
            # Use preferred method if available
            best = results.get(self.preferred_method, results['estimation'])

        # Log A/B test results for analysis
        self._log_ab_test(results, best)

        # Update state with best result
        if best['confidence'] == 'high':
            self.estimated_tokens = best['used']
            self.last_actual = best['used']

        return best, results

    def _log_ab_test(self, all_results, chosen):
        """Log A/B test results to JSONL file for later analysis"""
        log_entry = {
            'timestamp': datetime.now().isoformat(),
            'all_methods': {k: v for k, v in all_results.items()},
            'chosen_method': chosen['method'],
            'chosen_percentage': chosen['percentage'],
            'chosen_confidence': chosen['confidence']
        }

        # Append to JSONL log
        self.AB_TEST_LOG.parent.mkdir(exist_ok=True)
        with open(self.AB_TEST_LOG, 'a') as f:
            f.write(json.dumps(log_entry) + '\n')

    # ==================== Display & Formatting ====================

    def format_display(self, context_info, show_ab_test=False):
        """Generate display string with optional A/B test details"""
        pct = context_info['percentage']
        used = context_info['used']
        total = context_info['total']
        method = context_info['method']
        confidence = context_info['confidence']

        # Calculate blocks (20 total, each = 5%)
        blocks = round(pct / 5)
        viz = []
        for i in range(1, 21):
            if i <= blocks:
                if i <= 12:      # 0-60%
                    viz.append("🟩")
                elif i <= 15:    # 60-75%
                    viz.append("🟨")
                elif i <= 17:    # 75-85%
                    viz.append("🟧")
                else:            # 85-100%
                    viz.append("🟥")
            else:
                viz.append("⬛")

        display = f"Context: {''.join(viz)} {pct}% ({used//1000}k/{total//1000}k)"

        # Add status message
        if pct >= 85:
            display += " 🚨 New session recommended"
        elif pct >= 75:
            display += " 🔄 Session handoff recommended"
        elif pct >= 60:
            display += " ⚠️ Approaching handoff"

        # Add method indicator
        if show_ab_test:
            display += f" [method: {method}, confidence: {confidence}]"
        elif confidence == 'low':
            display += " [estimated]"

        return display

    def format_ab_comparison(self, all_results):
        """Format A/B test comparison for debugging"""
        lines = ["\n📊 A/B Test Results:"]

        for method, result in all_results.items():
            pct = result['percentage']
            used = result['used']
            conf = result['confidence']
            lines.append(f"  • {method}: {pct}% ({used//1000}k) - confidence: {conf}")

        return '\n'.join(lines)

    # ==================== Handoff Detection ====================

    def check_handoff_needed(self, percentage):
        """Check if automatic handoff should trigger at 75%"""
        if percentage >= 75 and self.last_handoff_pct < 75:
            self.last_handoff_pct = 75
            return True, "75% threshold reached"

        return False, None

    # ==================== State Persistence ====================

    def save_state(self):
        """Persist state for next run"""
        state = {
            'estimated_tokens': self.estimated_tokens,
            'last_actual': self.last_actual,
            'last_handoff_pct': self.last_handoff_pct,
            'preferred_method': self.preferred_method,
            'last_updated': datetime.now().isoformat()
        }

        self.STATE_FILE.parent.mkdir(exist_ok=True)
        with open(self.STATE_FILE, 'w') as f:
            json.dump(state, f, indent=2)

    # ==================== Analytics ====================

    @classmethod
    def analyze_ab_tests(cls):
        """Analyze A/B test log to determine most accurate method"""
        if not cls.AB_TEST_LOG.exists():
            return "No A/B test data available yet"

        method_stats = {}

        with open(cls.AB_TEST_LOG) as f:
            for line in f:
                try:
                    entry = json.loads(line)
                    method = entry['chosen_method']

                    if method not in method_stats:
                        method_stats[method] = {
                            'count': 0,
                            'total_confidence_high': 0
                        }

                    method_stats[method]['count'] += 1
                    if entry['chosen_confidence'] == 'high':
                        method_stats[method]['total_confidence_high'] += 1
                except:
                    continue

        # Generate report
        lines = ["📊 A/B Test Analysis Report:"]
        lines.append(f"Total samples: {sum(s['count'] for s in method_stats.values())}")
        lines.append("")

        for method, stats in sorted(method_stats.items(), key=lambda x: x[1]['count'], reverse=True):
            high_pct = (stats['total_confidence_high'] / stats['count'] * 100) if stats['count'] > 0 else 0
            lines.append(f"  {method}:")
            lines.append(f"    - Used: {stats['count']} times")
            lines.append(f"    - High confidence: {stats['total_confidence_high']} ({high_pct:.1f}%)")

        lines.append("")
        lines.append("Recommendation:")

        # Recommend method with most high-confidence results
        best_method = max(method_stats.items(),
                         key=lambda x: x[1]['total_confidence_high'])[0]
        lines.append(f"  Prefer '{best_method}' (most reliable)")

        return '\n'.join(lines)


def main():
    """Main entry point for context tracking"""
    import argparse

    parser = argparse.ArgumentParser(description='Automated context tracking with A/B testing')
    parser.add_argument('--show-ab', action='store_true',
                       help='Show A/B test comparison')
    parser.add_argument('--analyze', action='store_true',
                       help='Analyze A/B test results')
    parser.add_argument('--method', choices=['auto', 'context_command', 'system_warning', 'environment', 'estimation'],
                       default='auto', help='Preferred tracking method')
    parser.add_argument('--context-output', type=str,
                       help='Output from /context command for parsing')

    args = parser.parse_args()

    # Handle analysis mode
    if args.analyze:
        print(ContextTracker.analyze_ab_tests())
        return

    # Create tracker
    tracker = ContextTracker()
    tracker.preferred_method = args.method

    # Get response text (from stdin)
    response_text = sys.stdin.read() if not sys.stdin.isatty() else ""

    # Get context info (runs A/B test)
    best_result, all_results = tracker.get_context_info(response_text, args.context_output)

    # Display formatted context
    print(tracker.format_display(best_result, show_ab_test=args.show_ab))

    # Show A/B comparison if requested
    if args.show_ab:
        print(tracker.format_ab_comparison(all_results))

    # Check for automatic handoff
    should_handoff, reason = tracker.check_handoff_needed(best_result['percentage'])
    if should_handoff:
        print(f"\n🔔 AUTOMATIC HANDOFF TRIGGERED: {reason}")
        print("📝 Main agent should invoke /handoff command now")
        print("   Run: /handoff")

    # Save state for next run
    tracker.save_state()


if __name__ == "__main__":
    main()
