import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ProjectCard } from '@/components/projects/ProjectCard';
import { Project } from '@/types';

describe('ProjectCard', () => {
  const mockProject: Project = {
    id: 'test-project-1',
    title: 'Test Project',
    createdAt: '2025-12-24T10:00:00Z',
    updatedAt: '2025-12-24T10:47:00Z',
  };

  const mockOnClick = jest.fn();
  const mockOnRename = jest.fn().mockResolvedValue(undefined);
  const mockOnDelete = jest.fn().mockResolvedValue(undefined);

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('default view', () => {
    it('renders project title and date', () => {
      render(<ProjectCard project={mockProject} />);

      expect(screen.getByText('Test Project')).toBeInTheDocument();
      expect(screen.getByText(/Dec 24/)).toBeInTheDocument();
    });

    it('calls onClick when clicking the card', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onClick={mockOnClick} />);

      await user.click(screen.getByText('Test Project'));

      expect(mockOnClick).toHaveBeenCalledWith(mockProject);
    });

    it('shows pencil icon on hover when onRename is provided', () => {
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      expect(renameButton).toBeInTheDocument();
    });
  });

  describe('edit mode', () => {
    it('enters edit mode when clicking pencil icon', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      expect(screen.getByRole('textbox')).toBeInTheDocument();
      expect(screen.getByRole('textbox')).toHaveValue('Test Project');
    });

    it('shows Save button in edit mode', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      expect(screen.getByRole('button', { name: /save project name/i })).toBeInTheDocument();
    });

    it('shows Cancel button in edit mode', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      // Use the specific aria-label to target the edit mode Cancel button
      expect(screen.getByRole('button', { name: /cancel editing/i })).toBeInTheDocument();
    });

    it('shows trash icon in edit mode when onDelete is provided', async () => {
      const user = userEvent.setup();
      render(
        <ProjectCard
          project={mockProject}
          onRename={mockOnRename}
          onDelete={mockOnDelete}
        />
      );

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      // Look for the delete button in edit mode (trash icon)
      expect(screen.getByRole('button', { name: /delete project/i })).toBeInTheDocument();
    });

    it('saves new title when clicking Save button', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const input = screen.getByRole('textbox');
      await user.clear(input);
      await user.type(input, 'New Project Name');

      const saveButton = screen.getByRole('button', { name: /save project name/i });
      await user.click(saveButton);

      expect(mockOnRename).toHaveBeenCalledWith(mockProject, 'New Project Name');
    });

    it('cancels edit when clicking Cancel button', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const input = screen.getByRole('textbox');
      await user.clear(input);
      await user.type(input, 'Changed Name');

      const cancelButton = screen.getByRole('button', { name: /cancel editing/i });
      await user.click(cancelButton);

      // Should exit edit mode without saving
      expect(mockOnRename).not.toHaveBeenCalled();
      expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
      expect(screen.getByText('Test Project')).toBeInTheDocument();
    });

    it('saves on Enter key press', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const input = screen.getByRole('textbox');
      await user.clear(input);
      await user.type(input, 'Enter Saved Name{enter}');

      expect(mockOnRename).toHaveBeenCalledWith(mockProject, 'Enter Saved Name');
    });

    it('cancels on Escape key press', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const input = screen.getByRole('textbox');
      await user.clear(input);
      await user.type(input, 'Escape Cancel{escape}');

      expect(mockOnRename).not.toHaveBeenCalled();
      expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
    });
  });

  describe('delete flow from edit mode', () => {
    it('shows delete confirmation when clicking trash in edit mode', async () => {
      const user = userEvent.setup();
      render(
        <ProjectCard
          project={mockProject}
          onRename={mockOnRename}
          onDelete={mockOnDelete}
        />
      );

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const deleteButton = screen.getByRole('button', { name: /delete project/i });
      await user.click(deleteButton);

      // The delete confirmation panel should now be visible
      expect(screen.getByText(/delete\?/i)).toBeInTheDocument();
    });

    it('completes delete from edit mode', async () => {
      const user = userEvent.setup();
      render(
        <ProjectCard
          project={mockProject}
          onRename={mockOnRename}
          onDelete={mockOnDelete}
        />
      );

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const deleteButton = screen.getByRole('button', { name: /delete project/i });
      await user.click(deleteButton);

      // Find the confirm delete button in the sliding panel (text is just "Delete")
      const confirmButtons = screen.getAllByRole('button').filter(
        btn => btn.textContent === 'Delete'
      );
      expect(confirmButtons.length).toBeGreaterThan(0);
      await user.click(confirmButtons[0]);

      expect(mockOnDelete).toHaveBeenCalledWith(mockProject);
    });

    it('cancels delete confirmation and returns to view mode', async () => {
      const user = userEvent.setup();
      render(
        <ProjectCard
          project={mockProject}
          onRename={mockOnRename}
          onDelete={mockOnDelete}
        />
      );

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const deleteButton = screen.getByRole('button', { name: /delete project/i });
      await user.click(deleteButton);

      // Cancel in the delete confirmation panel - find the one with "Cancel" text
      // that does not have aria-label (the confirm panel Cancel)
      const cancelButtons = screen.getAllByRole('button').filter(
        btn => btn.textContent === 'Cancel' && !btn.getAttribute('aria-label')
      );
      expect(cancelButtons.length).toBeGreaterThan(0);
      await user.click(cancelButtons[0]);

      expect(mockOnDelete).not.toHaveBeenCalled();
    });
  });

  describe('styling', () => {
    it('applies active styling when isActive is true', () => {
      render(<ProjectCard project={mockProject} isActive={true} />);

      // Card should have teal border when active
      const card = screen.getByText('Test Project').closest('[class*="border-teal"]');
      expect(card).toBeInTheDocument();
    });

    it('Save button has primary teal styling', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const saveButton = screen.getByRole('button', { name: /save project name/i });
      expect(saveButton).toHaveClass('bg-teal-500');
    });

    it('Cancel button has secondary gray styling', async () => {
      const user = userEvent.setup();
      render(<ProjectCard project={mockProject} onRename={mockOnRename} />);

      const renameButton = screen.getByRole('button', { name: /rename/i });
      await user.click(renameButton);

      const cancelButton = screen.getByRole('button', { name: /cancel editing/i });
      expect(cancelButton).toHaveClass('border-gray-300');
    });
  });
});
