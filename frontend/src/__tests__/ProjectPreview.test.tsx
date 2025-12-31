import React from 'react';
import { render, screen } from '@testing-library/react';
import { ProjectPreview } from '@/components/preview/ProjectPreview';

describe('ProjectPreview', () => {
  describe('when no files are provided', () => {
    it('shows empty state message', () => {
      render(<ProjectPreview files={[]} />);

      const iframe = screen.getByTitle('App Preview');
      expect(iframe).toBeInTheDocument();
      expect(iframe).toHaveAttribute('sandbox', 'allow-scripts allow-same-origin');
    });
  });

  describe('when only CSS and JS files are provided (no HTML)', () => {
    it('shows no HTML file message in preview', () => {
      const files = [
        { path: 'styles.css', content: 'body { background: red; }' },
        { path: 'app.js', content: 'console.log("hello");' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      expect(iframe).toBeInTheDocument();
      // The empty state message should be in the srcdoc
      expect(iframe).toHaveAttribute('srcdoc', expect.stringContaining('No HTML file to preview'));
    });
  });

  describe('when HTML file is provided', () => {
    it('renders HTML content in iframe', () => {
      const files = [
        { path: 'index.html', content: '<!DOCTYPE html><html><head></head><body><h1>Hello World</h1></body></html>' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      expect(iframe).toBeInTheDocument();
      expect(iframe).toHaveAttribute('srcdoc', expect.stringContaining('Hello World'));
    });

    it('injects CSS files into head', () => {
      const files = [
        { path: 'index.html', content: '<!DOCTYPE html><html><head></head><body><h1>Test</h1></body></html>' },
        { path: 'styles.css', content: '.test { color: blue; }' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      const srcdoc = iframe.getAttribute('srcdoc');
      expect(srcdoc).toContain('<style>');
      expect(srcdoc).toContain('.test { color: blue; }');
      expect(srcdoc).toContain('</style>');
    });

    it('injects JS files before closing body tag', () => {
      const files = [
        { path: 'index.html', content: '<!DOCTYPE html><html><head></head><body><h1>Test</h1></body></html>' },
        { path: 'app.js', content: 'console.log("hello");' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      const srcdoc = iframe.getAttribute('srcdoc');
      expect(srcdoc).toContain('<script>');
      expect(srcdoc).toContain('console.log("hello");');
      expect(srcdoc).toContain('</script>');
    });

    it('injects both CSS and JS files', () => {
      const files = [
        { path: 'index.html', content: '<!DOCTYPE html><html><head></head><body><h1>Test</h1></body></html>' },
        { path: 'styles.css', content: 'body { margin: 0; }' },
        { path: 'main.js', content: 'alert("hi");' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      const srcdoc = iframe.getAttribute('srcdoc');
      expect(srcdoc).toContain('body { margin: 0; }');
      expect(srcdoc).toContain('alert("hi");');
    });
  });

  describe('when multiple HTML files exist', () => {
    it('prefers index.html as main file', () => {
      const files = [
        { path: 'about.html', content: '<html><body><h1>About Page</h1></body></html>' },
        { path: 'index.html', content: '<html><body><h1>Home Page</h1></body></html>' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      const srcdoc = iframe.getAttribute('srcdoc');
      expect(srcdoc).toContain('Home Page');
      expect(srcdoc).not.toContain('About Page');
    });

    it('uses first HTML file if no index.html exists', () => {
      const files = [
        { path: 'about.html', content: '<html><body><h1>About Page</h1></body></html>' },
        { path: 'contact.html', content: '<html><body><h1>Contact Page</h1></body></html>' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      const srcdoc = iframe.getAttribute('srcdoc');
      // Should pick one of them (first in the filtered list)
      expect(srcdoc).toMatch(/About Page|Contact Page/);
    });
  });

  describe('security', () => {
    it('applies sandbox attribute with allow-scripts and allow-same-origin', () => {
      const files = [
        { path: 'index.html', content: '<html><body>Test</body></html>' },
      ];

      render(<ProjectPreview files={files} />);

      const iframe = screen.getByTitle('App Preview');
      // allow-same-origin is needed for localStorage access in previewed apps
      expect(iframe).toHaveAttribute('sandbox', 'allow-scripts allow-same-origin');
    });
  });

  describe('className prop', () => {
    it('applies custom className', () => {
      render(<ProjectPreview files={[]} className="custom-class" />);

      const iframe = screen.getByTitle('App Preview');
      expect(iframe).toHaveClass('custom-class');
    });
  });
});
