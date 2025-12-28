import { render, screen } from '@testing-library/react';
import {
  PersonaIntroduction,
  isIntroductionMessage,
  isTeamIntroMessage,
} from '@/components/chat/PersonaIntroduction';
import { Message } from '@/types';

describe('PersonaIntroduction', () => {
  const teamIntroMessage: Message = {
    id: 'intro-root-team-project-1',
    projectId: 'project-1',
    role: 'assistant',
    content:
      "Great! Now that we understand your project, let me introduce the team who'll help build it.\n\nMeet Bloom, our designer - she'll craft the visual experience and user interface.\n\nAnd Harvest, our developer - he'll write the code and bring everything to life.",
    timestamp: '2025-12-24T10:00:00Z',
    agentType: 'product_manager',
  };

  const designerIntroMessage: Message = {
    id: 'intro-designer-project-1',
    projectId: 'project-1',
    role: 'assistant',
    content:
      "Hi! I'm excited to design your project. I'll focus on making it intuitive and beautiful.",
    timestamp: '2025-12-24T10:00:01Z',
    agentType: 'designer',
  };

  const developerIntroMessage: Message = {
    id: 'intro-developer-project-1',
    projectId: 'project-1',
    role: 'assistant',
    content:
      "Hey! Ready to build this. I'll handle the technical implementation and make sure everything works smoothly.",
    timestamp: '2025-12-24T10:00:02Z',
    agentType: 'developer',
  };

  const regularMessage: Message = {
    id: 'msg-123',
    projectId: 'project-1',
    role: 'assistant',
    content: 'Let me help you with that.',
    timestamp: '2025-12-24T10:01:00Z',
    agentType: 'designer',
  };

  describe('PersonaIntroduction component', () => {
    it('renders team introduction with special styling', () => {
      render(<PersonaIntroduction message={teamIntroMessage} isTeamIntro />);

      expect(screen.getByTestId('team-introduction')).toBeInTheDocument();
      expect(screen.getByText(/introducing the team/i)).toBeInTheDocument();
      expect(
        screen.getByText(/let me introduce the team/i)
      ).toBeInTheDocument();
    });

    it('renders designer intro with correct styling', () => {
      render(<PersonaIntroduction message={designerIntroMessage} />);

      expect(screen.getByTestId('persona-intro-designer')).toBeInTheDocument();
      expect(
        screen.getByText(/excited to design your project/i)
      ).toBeInTheDocument();
    });

    it('renders developer intro with correct styling', () => {
      render(<PersonaIntroduction message={developerIntroMessage} />);

      expect(screen.getByTestId('persona-intro-developer')).toBeInTheDocument();
      expect(screen.getByText(/Ready to build this/i)).toBeInTheDocument();
    });
  });

  describe('isIntroductionMessage helper', () => {
    it('returns true for introduction messages', () => {
      expect(isIntroductionMessage(teamIntroMessage)).toBe(true);
      expect(isIntroductionMessage(designerIntroMessage)).toBe(true);
      expect(isIntroductionMessage(developerIntroMessage)).toBe(true);
    });

    it('returns false for regular messages', () => {
      expect(isIntroductionMessage(regularMessage)).toBe(false);
    });
  });

  describe('isTeamIntroMessage helper', () => {
    it('returns true for team introduction messages', () => {
      expect(isTeamIntroMessage(teamIntroMessage)).toBe(true);
    });

    it('returns false for individual persona intros', () => {
      expect(isTeamIntroMessage(designerIntroMessage)).toBe(false);
      expect(isTeamIntroMessage(developerIntroMessage)).toBe(false);
    });

    it('returns false for regular messages', () => {
      expect(isTeamIntroMessage(regularMessage)).toBe(false);
    });
  });
});
