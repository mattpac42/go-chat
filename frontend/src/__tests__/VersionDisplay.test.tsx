import { render, screen } from '@testing-library/react';
import { VersionDisplay } from '@/components/settings/VersionDisplay';

describe('VersionDisplay', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    jest.resetModules();
    process.env = { ...originalEnv };
  });

  afterAll(() => {
    process.env = originalEnv;
  });

  it('renders version and git hash', () => {
    process.env.NEXT_PUBLIC_APP_VERSION = '0.1.0';
    process.env.NEXT_PUBLIC_GIT_HASH = 'fb50661';

    render(<VersionDisplay />);

    const versionElement = screen.getByTestId('version-display');
    expect(versionElement).toBeInTheDocument();
    expect(versionElement).toHaveTextContent('v0.1.0 (fb50661)');
  });

  it('displays default version when env vars are not set', () => {
    delete process.env.NEXT_PUBLIC_APP_VERSION;
    delete process.env.NEXT_PUBLIC_GIT_HASH;

    render(<VersionDisplay />);

    const versionElement = screen.getByTestId('version-display');
    expect(versionElement).toHaveTextContent('v0.0.0 (dev)');
  });

  it('applies subtle styling', () => {
    process.env.NEXT_PUBLIC_APP_VERSION = '0.1.0';
    process.env.NEXT_PUBLIC_GIT_HASH = 'abc1234';

    render(<VersionDisplay />);

    const versionElement = screen.getByTestId('version-display');
    expect(versionElement).toHaveClass('text-xs', 'text-gray-400');
  });
});
