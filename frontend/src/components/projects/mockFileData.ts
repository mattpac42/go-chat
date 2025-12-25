import { FileNode } from '@/types';

/**
 * Mock file data with 2-tier reveal metadata
 * This demonstrates the data structure that will work with the schema
 * being designed for the App Map feature.
 */
export const mockFileTree: FileNode[] = [
  {
    id: 'folder-frontend',
    name: 'frontend',
    path: 'frontend',
    type: 'folder',
    functionalGroup: 'User Interface',
    children: [
      {
        id: 'folder-frontend-src',
        name: 'src',
        path: 'frontend/src',
        type: 'folder',
        children: [
          {
            id: 'folder-frontend-src-components',
            name: 'components',
            path: 'frontend/src/components',
            type: 'folder',
            functionalGroup: 'UI Components',
            children: [
              {
                id: 'file-header',
                name: 'Header.tsx',
                path: 'frontend/src/components/Header.tsx',
                type: 'file',
                language: 'tsx',
                shortDescription: 'Navigation bar with logo and menu links',
                longDescription: 'Displays the main site navigation, handles mobile hamburger menu, and shows user authentication status.',
                functionalGroup: 'Navigation',
                content: `import React from 'react';

export function Header() {
  return (
    <header className="bg-white shadow-sm">
      <nav className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <span className="text-xl font-bold text-teal-600">
              MyApp
            </span>
          </div>
          <div className="flex items-center space-x-4">
            <a href="/about" className="text-gray-600 hover:text-gray-900">
              About
            </a>
            <a href="/contact" className="text-gray-600 hover:text-gray-900">
              Contact
            </a>
          </div>
        </div>
      </nav>
    </header>
  );
}`,
              },
              {
                id: 'file-footer',
                name: 'Footer.tsx',
                path: 'frontend/src/components/Footer.tsx',
                type: 'file',
                language: 'tsx',
                shortDescription: 'Site footer with copyright and links',
                longDescription: 'Contains footer navigation, social media links, and copyright information.',
                functionalGroup: 'Navigation',
                content: `import React from 'react';

export function Footer() {
  return (
    <footer className="bg-gray-100 mt-auto">
      <div className="max-w-7xl mx-auto py-6 px-4">
        <p className="text-center text-gray-500 text-sm">
          &copy; 2025 MyApp. All rights reserved.
        </p>
      </div>
    </footer>
  );
}`,
              },
              {
                id: 'file-contact-form',
                name: 'ContactForm.tsx',
                path: 'frontend/src/components/ContactForm.tsx',
                type: 'file',
                language: 'tsx',
                shortDescription: 'Contact form with email and message fields',
                longDescription: 'Collects user contact information and messages, validates input, and submits to the backend API.',
                functionalGroup: 'Contact Form',
                content: `import React, { useState } from 'react';

export function ContactForm() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      // API call would go here
      await new Promise(resolve => setTimeout(resolve, 1000));
      alert('Message sent!');
      setEmail('');
      setMessage('');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="email" className="block text-sm font-medium">
          Email
        </label>
        <input
          type="email"
          id="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="mt-1 block w-full rounded border-gray-300"
          required
        />
      </div>
      <div>
        <label htmlFor="message" className="block text-sm font-medium">
          Message
        </label>
        <textarea
          id="message"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          rows={4}
          className="mt-1 block w-full rounded border-gray-300"
          required
        />
      </div>
      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full py-2 px-4 bg-teal-600 text-white rounded hover:bg-teal-700"
      >
        {isSubmitting ? 'Sending...' : 'Send Message'}
      </button>
    </form>
  );
}`,
              },
            ],
          },
          {
            id: 'folder-frontend-src-pages',
            name: 'pages',
            path: 'frontend/src/pages',
            type: 'folder',
            functionalGroup: 'Pages',
            children: [
              {
                id: 'file-homepage',
                name: 'index.tsx',
                path: 'frontend/src/pages/index.tsx',
                type: 'file',
                language: 'tsx',
                shortDescription: 'Homepage with welcome message and features',
                longDescription: 'The main landing page visitors see first. Displays hero section, feature highlights, and call-to-action.',
                functionalGroup: 'Homepage',
                content: `import React from 'react';
import { Header } from '../components/Header';
import { Footer } from '../components/Footer';

export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        <div className="max-w-7xl mx-auto py-12 px-4">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Welcome to MyApp
          </h1>
          <p className="text-xl text-gray-600 mb-8">
            Build something amazing today.
          </p>
          <a
            href="/contact"
            className="inline-block bg-teal-600 text-white px-6 py-3 rounded-lg"
          >
            Get Started
          </a>
        </div>
      </main>
      <Footer />
    </div>
  );
}`,
              },
              {
                id: 'file-about',
                name: 'about.tsx',
                path: 'frontend/src/pages/about.tsx',
                type: 'file',
                language: 'tsx',
                shortDescription: 'About page with company information',
                longDescription: 'Tells visitors about the company, its mission, and team members.',
                functionalGroup: 'About Page',
                content: `import React from 'react';
import { Header } from '../components/Header';
import { Footer } from '../components/Footer';

export default function AboutPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        <div className="max-w-7xl mx-auto py-12 px-4">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            About Us
          </h1>
          <p className="text-gray-600">
            We are a team passionate about building great software.
          </p>
        </div>
      </main>
      <Footer />
    </div>
  );
}`,
              },
              {
                id: 'file-contact',
                name: 'contact.tsx',
                path: 'frontend/src/pages/contact.tsx',
                type: 'file',
                language: 'tsx',
                shortDescription: 'Contact page with form for reaching out',
                longDescription: 'Allows visitors to send messages to the team using the contact form.',
                functionalGroup: 'Contact Form',
                content: `import React from 'react';
import { Header } from '../components/Header';
import { Footer } from '../components/Footer';
import { ContactForm } from '../components/ContactForm';

export default function ContactPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        <div className="max-w-md mx-auto py-12 px-4">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Contact Us
          </h1>
          <p className="text-gray-600 mb-8">
            Have a question? We'd love to hear from you.
          </p>
          <ContactForm />
        </div>
      </main>
      <Footer />
    </div>
  );
}`,
              },
            ],
          },
        ],
      },
    ],
  },
  {
    id: 'folder-backend',
    name: 'backend',
    path: 'backend',
    type: 'folder',
    functionalGroup: 'Backend Services',
    children: [
      {
        id: 'file-server',
        name: 'server.go',
        path: 'backend/server.go',
        type: 'file',
        language: 'go',
        shortDescription: 'Main server that handles all requests',
        longDescription: 'The heart of the backend. Sets up routes, middleware, and starts the HTTP server.',
        functionalGroup: 'Backend Services',
        content: `package main

import (
    "log"
    "net/http"
)

func main() {
    // Set up routes
    http.HandleFunc("/api/contact", handleContact)
    http.HandleFunc("/api/health", handleHealth)

    // Start server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func handleContact(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // Handle contact form submission
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Message received"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}`,
      },
      {
        id: 'file-db',
        name: 'database.go',
        path: 'backend/database.go',
        type: 'file',
        language: 'go',
        shortDescription: 'Database connection and queries',
        longDescription: 'Manages the PostgreSQL database connection, runs migrations, and provides query helpers.',
        functionalGroup: 'Database',
        content: `package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB(connStr string) error {
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    if err = db.Ping(); err != nil {
        return err
    }

    log.Println("Database connected")
    return nil
}

func saveContact(email, message string) error {
    _, err := db.Exec(
        "INSERT INTO contacts (email, message) VALUES ($1, $2)",
        email, message,
    )
    return err
}`,
      },
    ],
  },
  {
    id: 'file-package-json',
    name: 'package.json',
    path: 'package.json',
    type: 'file',
    language: 'json',
    shortDescription: 'Project dependencies and npm scripts',
    longDescription: 'Lists all the packages this project needs and defines commands like "npm start" and "npm build".',
    functionalGroup: 'Configuration',
    content: `{
  "name": "my-app",
  "version": "1.0.0",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "eslint ."
  },
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/react": "^18.0.0",
    "eslint": "^8.0.0"
  }
}`,
  },
  {
    id: 'file-readme',
    name: 'README.md',
    path: 'README.md',
    type: 'file',
    language: 'markdown',
    shortDescription: 'Project documentation and setup guide',
    longDescription: 'Explains what this project does, how to set it up, and how to use it.',
    functionalGroup: 'Documentation',
    content: `# MyApp

A simple web application with contact form functionality.

## Getting Started

1. Install dependencies:
   \`\`\`bash
   npm install
   \`\`\`

2. Start the development server:
   \`\`\`bash
   npm run dev
   \`\`\`

3. Open [http://localhost:3000](http://localhost:3000)

## Features

- Homepage with welcome message
- About page
- Contact form with email and message

## Tech Stack

- Next.js 14
- React 18
- Go backend
- PostgreSQL database
`,
  },
];

/**
 * Flatten file tree to get all files (for search/filtering)
 */
export function flattenFileTree(nodes: FileNode[]): FileNode[] {
  const result: FileNode[] = [];

  function traverse(node: FileNode) {
    if (node.type === 'file') {
      result.push(node);
    }
    if (node.children) {
      node.children.forEach(traverse);
    }
  }

  nodes.forEach(traverse);
  return result;
}

/**
 * Get files by functional group
 */
export function getFilesByGroup(nodes: FileNode[]): Record<string, FileNode[]> {
  const files = flattenFileTree(nodes);
  const groups: Record<string, FileNode[]> = {};

  for (const file of files) {
    const group = file.functionalGroup || 'Other';
    if (!groups[group]) {
      groups[group] = [];
    }
    groups[group].push(file);
  }

  return groups;
}
