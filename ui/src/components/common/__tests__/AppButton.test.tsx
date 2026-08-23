import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import AppButton from '../AppButton';

describe('AppButton', () => {
  it('renders children correctly', () => {
    render(<AppButton>Click Me</AppButton>);
    expect(screen.getByText('Click Me')).toBeInTheDocument();
  });

  it('handles click events', () => {
    const handleClick = vi.fn();
    render(<AppButton onClick={handleClick}>Click Me</AppButton>);
    
    fireEvent.click(screen.getByText('Click Me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('applies variant classes correctly', () => {
    render(<AppButton variant="danger">Delete</AppButton>);
    const button = screen.getByText('Delete');
    expect(button.className).toContain('text-danger-700');
  });
});
