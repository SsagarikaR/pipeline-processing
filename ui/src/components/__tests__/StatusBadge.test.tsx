import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import StatusBadge from '../StatusBadge';

describe('StatusBadge', () => {
  it('renders pending status with correct text', () => {
    render(<StatusBadge status="pending" />);
    const badge = screen.getByText('pending');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('bg-neutral-100');
  });

  it('renders completed status with correct style', () => {
    render(<StatusBadge status="completed" />);
    const badge = screen.getByText('completed');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('bg-success-100');
  });

  it('renders running status with pulse animation', () => {
    render(<StatusBadge status="running" />);
    const badge = screen.getByText('running');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('animate-pulse');
  });
});
